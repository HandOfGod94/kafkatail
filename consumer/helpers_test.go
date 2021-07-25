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

func sendMessage(t *testing.T, ctx context.Context, brokers []string, topic string, key, message []byte) {
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

func sendMessageToPartition(t *testing.T, ctx context.Context, brokers []string, topic string, parition int, key, message []byte) {
	record := kafka.Record{
		Value: kafka.NewBytes(message),
	}
	produceReq := kafka.ProduceRequest{Topic: topic, Partition: parition, RequiredAcks: kafka.RequireOne, Records: kafka.NewRecordReader(record)}
	resp, err := kafkaRawClient.Produce(ctx, &produceReq)

	if err != nil {
		t.Log("failed to send produce request to broker:", err)
		t.FailNow()
	}

	if resp.Error != nil {
		t.Log("failed to produce message: ", resp.Error)
		t.FailNow()
	}
}

func createTopicWithConfig(t *testing.T, ctx context.Context, topic kafka.TopicConfig) {
	resp, err := kafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
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
