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
	"github.com/stretchr/testify/suite"
)

const defaultTimeout = 5 * time.Second

type KafkaConsumerTestSuite struct {
	suite.Suite
}

func TestKafkaConsumerSuite(t *testing.T) {
	suite.Run(t, new(KafkaConsumerTestSuite))
}

func (kct *KafkaConsumerTestSuite) SetupSuite() {
	createTopicWithConfig(kct.T(), context.Background(), kafka.TopicConfig{
		Topic:             "kafkatail-test-topic",
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
}

func (kct *KafkaConsumerTestSuite) TearDownSuite() {
	deleteTopic(kct.T(), context.Background(), "kafkatail-test-topic")
}

func (s *KafkaConsumerTestSuite) TestConsume_Success() {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	c := consumer.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	go func(ctx context.Context) {
		c.Consume(ctx, w, wire.NewPlaintextDecoder())
	}(ctx)

	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("world"))

	<-ctx.Done()
	got := w.String()
	assert.Contains(s.T(), got, "hello\n")
	assert.Contains(s.T(), got, "world\n")
}

func (s *KafkaConsumerTestSuite) TestConsume_WithParition() {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	c := consumer.Options{Partition: 0}.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	go func(ctx context.Context) {
		c.Consume(ctx, w, wire.NewPlaintextDecoder())
	}(ctx)

	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("bar"), []byte("world"))

	<-ctx.Done()
	assert.Contains(s.T(), w.String(), "Partition: 0")
}

func (s *KafkaConsumerTestSuite) TestConsumer_WithOffeset() {
	w := bytes.NewBufferString("")
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("hello"))
	sendMessage(s.T(), ctx, []string{"localhost:9093"}, "kafkatail-test-topic", []byte("foo"), []byte("world"))

	c := consumer.Options{Offset: kafka.FirstOffset}.New([]string{"localhost:9093"}, "kafkatail-test-topic")
	c.Consume(ctx, w, wire.NewPlaintextDecoder())

	got := w.String()
	assert.Contains(s.T(), got, "hello\n")
	assert.Contains(s.T(), got, "hello\n")
}

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
