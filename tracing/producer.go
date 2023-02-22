package tracing

import (
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

import (
	"context"
)

func SendMessages(ctx context.Context, writer *kafka.Writer, messages ...kafka.Message) error {
	for _, message := range messages {
		span := startProducerSpan(ctx, &message)
		err := writer.WriteMessages(ctx, message)
		finishProducerSpan(span, err)
		return err
	}
	return nil
}

func startProducerSpan(ctx context.Context, msg *kafka.Message) trace.Span {
	carrier := NewMessageCarrier(msg)

	attrs := []attribute.KeyValue{
		semconv.MessagingSystem("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationName(msg.Topic),
		semconv.MessagingOperationPublish,
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindProducer),
	}
	ctx, span := otel.Tracer("github.com/andser91/otel-kafka-go").
		Start(ctx, fmt.Sprintf("%s publish", msg.Topic), opts...)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return span
}

func finishProducerSpan(span trace.Span, err error) {
	if err != nil {
		span.SetStatus(codes.Error, err.Error())
	}
	span.End()
}
