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

// PublishWithTrace starts a child span for the publish, injects its
// context into the message headers, and publishes via JetStream, so a
// downstream consumer (e.g. accumulator) can extract it and continue the
// same trace.
func PublishWithTrace(ctx context.Context, js nats.JetStreamContext, subject string, data []byte) error {
	tracer := otel.Tracer("maker")
	ctx, span := tracer.Start(ctx, "nats.publish."+subject, trace.WithSpanKind(trace.SpanKindProducer))
	defer span.End()

	msg := &nats.Msg{Subject: subject, Data: data, Header: nats.Header{}}
	otel.GetTextMapPropagator().Inject(ctx, natsHeaderCarrier(msg.Header))
	_, err := js.PublishMsg(msg)
	if err != nil {
		span.RecordError(err)
	}
	return err
}
