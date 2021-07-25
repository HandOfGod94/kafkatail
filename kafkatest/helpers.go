package kafkatest

import (
	"context"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
)

var DefaultTimeout = 10 * time.Second

var KafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   DefaultTimeout,
	Transport: nil,
}

func SendMessage(t *testing.T, ctx context.Context, brokers []string, topic string, key, message []byte) {
	w := &kafka.Writer{
		Addr:  kafka.TCP(brokers...),
		Topic: topic,
	}

	if err := w.WriteMessages(ctx, kafka.Message{
		Key:   key,
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

func SendMessageToPartition(t *testing.T, ctx context.Context, brokers []string, topic string, parition int, key, message []byte) {
	record := kafka.Record{
		Value: kafka.NewBytes(message),
	}
	produceReq := kafka.ProduceRequest{Topic: topic, Partition: parition, RequiredAcks: kafka.RequireOne, Records: kafka.NewRecordReader(record)}
	resp, err := KafkaRawClient.Produce(ctx, &produceReq)

	if err != nil {
		t.Log("failed to send produce request to broker:", err)
		t.FailNow()
	}

	if resp.Error != nil {
		t.Log("failed to produce message: ", resp.Error)
		t.FailNow()
	}
}

func CreateTopicWithConfig(t *testing.T, ctx context.Context, topic kafka.TopicConfig) {
	resp, err := KafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{topic},
	})

	if err != nil {
		t.Log("failed to create topic:", err)
		t.FailNow()
	}

	if resp.Errors[topic.Topic] != nil {
		t.Logf("failed to create topic. errors: %+v", resp.Errors)
		t.FailNow()
	}
}

func DeleteTopic(t *testing.T, ctx context.Context, topic string) {
	resp, err := KafkaRawClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
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
