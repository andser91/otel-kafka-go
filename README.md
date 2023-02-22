# otel-kafka-go

Library for instrumenting "github.com/segmentio/kafka-go" kafka client. 

## How to use

### Producer:

Send the message wrapping kafka.Writer in a function that manage new Span. 
The function needs current context, the writer and the message/messages to send

SendMessages(ctx context.Context, writer *kafka.Writer, messages ...kafka.Message)

### Consumer 
The consumer needs to retrieve the instrumented context and create a new Span.
The function "Extract" take current context and the message as parameters and return an instrumented context and a new active Span.

Extract(ctx context.Context, message *kafka.Message) (context.Context, trace.Span)

