package tracing

import (
	"context"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// natsHeaderCarrier adapts nats.Header to OTel's TextMapCarrier interface.
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

// ExtractContext pulls the trace context injected by an upstream publisher
// (e.g. maker) out of the NATS message headers, so the consumer span below
// continues the same trace instead of starting a new one.
func ExtractContext(ctx context.Context, msg *nats.Msg) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, natsHeaderCarrier(msg.Header))
}

// StartConsumerSpan starts a span for handling a message received on subject.
func StartConsumerSpan(ctx context.Context, subject string) (context.Context, trace.Span) {
	tracer := otel.Tracer("accumulator")
	return tracer.Start(ctx, "nats.consume."+subject, trace.WithSpanKind(trace.SpanKindConsumer))
}
