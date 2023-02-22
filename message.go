package tracing

import "github.com/segmentio/kafka-go"

type MessageCarrier struct {
	msg *kafka.Message
}

func NewMessageCarrier(msg *kafka.Message) MessageCarrier {
	return MessageCarrier{msg: msg}
}

func (c MessageCarrier) Get(key string) string {
	for _, h := range c.msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c MessageCarrier) Set(key, value string) {
	for i := 0; i < len(c.msg.Headers); i++ {
		if c.msg.Headers[i].Key == key {
			c.msg.Headers = append(c.msg.Headers[:i], c.msg.Headers[i+1:]...)
		}
	}

	c.msg.Headers = append(c.msg.Headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c MessageCarrier) Keys() []string {
	var keys []string
	for i, _ := range c.msg.Headers {
		keys = append(keys, c.msg.Headers[i].Key)
	}
	return keys
}
