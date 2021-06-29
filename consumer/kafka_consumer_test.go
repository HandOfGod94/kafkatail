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
		GroupID:   "foo",
		Offset:    10,
		Partition: 2,
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

	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("world"))

	<-ctx.Done()
	got := w.String()
	assert.Contains(t, got, "hello\n")
	assert.Contains(t, got, "world\n")
}

func TestConsume_WithParition(t *testing.T) {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	topic := kafka.TopicConfig{
		Topic:             "kafkatail-test-topic",
		NumPartitions:     2,
		ReplicationFactor: 1,
	}
	createTopicWithConfig(t, context.Background(), topic)
	defer deleteTopic(t, context.Background(), "kafkatail-test-topic")

	c := consumer.Options{Partition: 0}.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	go func(ctx context.Context) {
		c.Consume(ctx, w, wire.NewPlaintextDecoder())
	}(ctx)

	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("bar"), []byte("world"))

	<-ctx.Done()
	assert.Contains(t, w.String(), "Partition: 0")
}

func TestConsumer_WithOffeset(t *testing.T) {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	createTopic(t, context.Background(), "kafkatail-test-topic")
	defer deleteTopic(t, context.Background(), "kafkatail-test-topic")

	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(t, ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("world"))

	c := consumer.Options{Offset: kafka.FirstOffset}.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	c.Consume(ctx, w, wire.NewPlaintextDecoder())

	assert.Equal(t, "hello\nworld\n", w.String())
}
