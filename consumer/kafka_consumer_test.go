package consumer_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

const defaultTimeout = 5 * time.Second

func TestCreateConsumerWithOptions(t *testing.T) {
	opts := consumer.Options{
		GroupID: "foo",
		Offset:  10,
	}

	c := opts.New([]string{"localhost:9093"}, "footest")

	assert.NotNil(t, c)
}

func TestConsume_Errors(t *testing.T) {
	testCases := []struct {
		desc             string
		bootstrapServers []string
		topic            string
		expectedErr      string
	}{
		{"with invalid brokers config", []string{"foo:9093"}, "test_topic", "no such host"},
		{"with invalid topic", []string{"localhost:9093"}, "nonexistent_topic", "deadline exceeded"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			w := bytes.NewBufferString("")
			c := consumer.New(tc.bootstrapServers, tc.topic)
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			err := c.Consume(ctx, w, wire.NewPlaintextDecoder())

			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestConsume_Success(t *testing.T) {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	createTopic(t, context.Background(), "kafkatail-test-topic")
	defer deleteTopic(t, context.Background(), "kafkatail-test-topic")

	c := consumer.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	go func(ctx context.Context) {
		c.Consume(ctx, w, wire.NewPlaintextDecoder())
	}(ctx)

	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("hello"))
	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("world"))

	<-ctx.Done()
	assert.Equal(t, "hello\nworld\n", w.String())
}

func TestConsumer_WithOffeset(t *testing.T) {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	createTopic(t, context.Background(), "kafkatail-test-topic")
	defer deleteTopic(t, context.Background(), "kafkatail-test-topic")

	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("hello"))
	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("world"))

	c := consumer.Options{Offset: kafka.FirstOffset}.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	c.Consume(ctx, w, wire.NewPlaintextDecoder())

	assert.Equal(t, "hello\nworld\n", w.String())
}
