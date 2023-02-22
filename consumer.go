package tracing

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"strconv"
)

func Extract(ctx context.Context, message *kafka.Message) (context.Context, trace.Span) {
	carrier := NewMessageCarrier(message)
	parentSpanContext := otel.GetTextMapPropagator().Extract(ctx, carrier)

	// Create a span.
	attrs := []attribute.KeyValue{
		semconv.MessagingSystem("kafka"),
		semconv.MessagingDestinationKindTopic,
		semconv.MessagingDestinationName(message.Topic),
		semconv.MessagingOperationReceive,
		semconv.MessagingMessageID(strconv.FormatInt(message.Offset, 10)),
		semconv.MessagingKafkaSourcePartition(message.Partition),
	}
	opts := []trace.SpanStartOption{
		trace.WithAttributes(attrs...),
		trace.WithSpanKind(trace.SpanKindConsumer),
	}
	newCtx, span := otel.Tracer("github.com/andser91/otel-kafka-go").
		Start(parentSpanContext, fmt.Sprintf("%s receive", message.Topic), opts...)
	otel.GetTextMapPropagator().Inject(newCtx, carrier)
	return newCtx, span
}
