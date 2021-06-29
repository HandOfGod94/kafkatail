package consumer_test

import (
	"context"
	"testing"

	"github.com/segmentio/kafka-go"
)

var kafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   defaultTimeout,
	Transport: nil,
}

func sendMessage(t *testing.T, ctx context.Context, brokers []string, topic string, message []byte) {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   []byte("foo"),
		Value: message,
	}); err != nil {
		t.Log("failed to write messages: ", err)
		t.FailNow()
	}

	if err := w.Close(); err != nil {
		t.Log("failed to close writer:", err)
		t.FailNow()
	}
}

func createTopic(t *testing.T, ctx context.Context, topic string) {
	resp, err := kafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}},
	})

	if err != nil {
		t.Log("failed to create topic:", err)
		t.FailNow()
	}

	if resp.Errors[topic] != nil {
		t.Logf("failed to create topic. errors: %+v", resp.Errors)
		t.FailNow()
	}
}

func deleteTopic(t *testing.T, ctx context.Context, topic string) {
	resp, err := kafkaRawClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})

	if err != nil {
		t.Log("failed to delete topic:", err)
		t.FailNow()
	}

	if resp.Errors[topic] != nil {
		t.Logf("failed to create topic. errors: %+v", resp.Errors)
		t.FailNow()
	}
}
