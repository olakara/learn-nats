# Go OTel Integration Notes

These are the three places instrumentation actually has to be written in
your Go code — the compose file gives you nowhere to send data without this.

## 1. NATS JetStream — trace context through message headers

NATS does **not** propagate trace context automatically like some Kafka
client libraries do. You inject/extract manually using message headers.
This is the piece that determines whether JetStream redeliveries show up
as a coherent retry chain or as orphaned traces.

```go
import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

// natsHeaderCarrier adapts nats.Header to the TextMapCarrier interface
// OTel's propagator expects.
type natsHeaderCarrier nats.Header

func (c natsHeaderCarrier) Get(key string) string {
	values := nats.Header(c)[key]
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c natsHeaderCarrier) Set(key, value string) {
	nats.Header(c).Set(key, value)
}

func (c natsHeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(c))
	for k := range c {
		keys = append(keys, k)
	}
	return keys
}

// Publisher side: inject the current span context into headers before publish.
func PublishWithTrace(ctx context.Context, js nats.JetStreamContext, subject string, data []byte) error {
	msg := &nats.Msg{Subject: subject, Data: data, Header: nats.Header{}}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier(msg.Header))
	_, err := js.PublishMsg(msg)
	return err
}

// Consumer side: extract the parent span context before processing.
// Do this for EVERY delivery attempt, including redeliveries — that's
// what keeps retries attached to the original trace instead of starting
// a new, disconnected one.
func HandleWithTrace(msg *nats.Msg, tracerName string) (context.Context, trace.Span) {
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), natsHeaderCarrier(msg.Header))
	tracer := otel.Tracer(tracerName)
	return tracer.Start(ctx, "nats.consume."+msg.Subject)
}
```

Also record `nats.Msg.Header.Get(nats.MsgIdHdr)` and the JetStream
`NumDelivered()` metadata as span attributes — that's how you'll tell
"first attempt" from "3rd redelivery" when you're staring at Tempo later.

## 2. Postgres — otelsql

Wraps `database/sql` so every query becomes a span, without you writing
manual `tracer.Start()` calls around every query.

```go
import (
	"github.com/XSAM/otelsql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

db, err := otelsql.Open("pgx", dsn, otelsql.WithAttributes(
	semconv.DBSystemPostgreSQL,
))
if err != nil {
	log.Fatal(err)
}
// Optional: register DB pool stats (open conns, idle conns, wait count)
// as metrics too, not just traces.
otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(
	semconv.DBSystemPostgreSQL,
))
```

## 3. Redis — redisotel

```go
import (
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

rdb := redis.NewClient(&redis.Options{Addr: "redis:6379"})
if err := redisotel.InstrumentTracing(rdb); err != nil {
	log.Fatal(err)
}
if err := redisotel.InstrumentMetrics(rdb); err != nil {
	log.Fatal(err)
}
```

## 4. Sampler — AlwaysSample for now

At <10 events/sec there's no cost justification for probabilistic or
tail-based sampling yet. Every trace, every time:

```go
import "go.opentelemetry.io/otel/sdk/trace"

tp := trace.NewTracerProvider(
	trace.WithSampler(trace.AlwaysSample()),
	// ... exporter, resource, etc.
)
```

**When to revisit this:** once you're consistently above ~100 events/sec
per service, or once Tempo storage/query latency becomes noticeable.
At that point switch to a tail-based sampling policy in the Collector
(sample 100% of errors and slow traces, a small percentage of the rest) —
not head-based random sampling, which would blind you to rare failures
just as easily as it saves storage.

## 5. Shared setup — don't let each service configure OTel independently

With multiple Go services, hand-rolling SDK setup per service is how you
end up with inconsistent resource attributes, mismatched service names, or
one service silently missing a signal type. Put the SDK bootstrap
(tracer provider, meter provider, logger provider, resource attributes,
OTLP exporter config, shutdown hooks) in one internal shared package
(`internal/otel` or similar) that every service imports and calls once
from `main()`. Config differences (service name, environment) go through
env vars, not duplicated setup code.
