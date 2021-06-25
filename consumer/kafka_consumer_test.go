package consumer_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tomb.v2"
)

const defaultTimeout = 5 * time.Second

var kafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   defaultTimeout,
	Transport: nil,
}

func TestConsume_Errors(t *testing.T) {
	testCases := []struct {
		desc             string
		bootstrapServers []string
		topic            string
		expectedErr      string
	}{
		{"with invalid brokers config", []string{"foo:9093"}, "test_topic", "lookup foo: no such host"},
		{"with invalid topic", []string{"localhost:9093"}, "nonexistent_topic", context.DeadlineExceeded.Error()},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			c := consumer.New(tc.bootstrapServers, tc.topic)
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			tmb, tctx := tomb.WithContext(ctx)

			c.Consume(tmb, tctx)

			<-tmb.Dead()

			cancel()
			assert.False(t, tmb.Alive())
			assert.Contains(t, tmb.Err().Error(), tc.expectedErr)
		})
	}
}

func TestReadTimeout(t *testing.T) {
	broker := []string{"localhost:9093"}
	topic := "test_topic"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	tmb, tctx := tomb.WithContext(ctx)

	c := consumer.New(broker, topic)
	c.Consume(tmb, tctx)

	<-tmb.Dead()

	assert.False(t, tmb.Alive())
	assert.ErrorIs(t, tmb.Err(), context.DeadlineExceeded)
}

func TestConsume_Success(t *testing.T) {
	broker := []string{"localhost:9093"}
	topic := "test_topic"
	sent := "hello world"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	tmb, tctx := tomb.WithContext(ctx)

	c := consumer.New(broker, topic)
	msgChan := c.Consume(tmb, tctx)

	sendMessage(ctx, broker, topic, sent)
	received := <-msgChan

	assert.Equal(t, sent, received)
}

func TestConsume_WithMultipleMessages(t *testing.T) {
	broker := []string{"localhost:9093"}
	topic := "test_topic"
	sent1 := "hello"
	sent2 := "world"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	tmb, tctx := tomb.WithContext(ctx)

	c := consumer.New(broker, topic)
	msgChan := c.Consume(tmb, tctx)

	sendMessage(ctx, broker, topic, sent1)
	sendMessage(ctx, broker, topic, sent2)
	received1 := <-msgChan
	received2 := <-msgChan

	assert.Equal(t, sent1, received1)
	assert.Equal(t, sent2, received2)
}

func sendMessage(ctx context.Context, brokers []string, topic, message string) {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("foo"),
		Value: []byte(message),
	}); err != nil {
		log.Fatal("failed to write messages: ", err)
	}

	if err := w.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
}

func createTopic(ctx context.Context, topic string) error {
	_, err := kafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{{Topic: topic, NumPartitions: 1}},
	})

	if err != nil {
		return err
	}
	return nil
}

func deleteTopic(ctx context.Context, topic string) error {
	_, err := kafkaRawClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})

	if err != nil {
		return err
	}
	return nil
}

func TestMain(m *testing.M) {
	topic := "test_topic"
	if err := createTopic(context.Background(), topic); err != nil {
		log.Fatal("failed to create topic:", err)
	}

	code := m.Run()

	if err := deleteTopic(context.Background(), topic); err != nil {
		log.Fatal("failed to delete topic:", err)
	}
	os.Exit(code)
}
