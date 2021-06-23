package consumer_test

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

const defaultTimeout = 2 * time.Second

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
		{"with invalid brokers config", []string{"foo:9093"}, "test_topic", "no such host"},
		{"with invalid topic", []string{"localhost:9093"}, "nonexistent_topic", "deadline exceeded"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			c := consumer.New(tc.bootstrapServers, tc.topic)
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			_, errorChan := c.Consume(ctx)
			err := <-errorChan

			cancel()
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}

func TestReadTimeout(t *testing.T) {
	broker := []string{"localhost:9093"}
	topic := "test_topic"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	c := consumer.New(broker, topic)
	_, errChan := c.Consume(ctx)
	err := <-errChan

	assert.Error(t, err)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestConsume_Success(t *testing.T) {
	broker := []string{"localhost:9093"}
	topic := "test_topic"
	sent := "hello world"
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()
	err := createTopic(ctx, topic)
	if err != nil {
		log.Fatal("failed to create topic:", err)
	}

	c := consumer.New(broker, topic)
	msgChan, _ := c.Consume(ctx)

	sendMessage(ctx, broker, topic, sent)
	received := <-msgChan

	deleteTopic(ctx, topic)
	assert.Equal(t, sent, received)
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
