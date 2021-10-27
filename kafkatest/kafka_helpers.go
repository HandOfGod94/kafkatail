package kafkatest

import (
	"context"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"gopkg.in/tomb.v2"
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

func SendMultipleMessagesToParition(t *testing.T, ctx context.Context, brokers []string, topic string, msgs map[int]string) {
	for partition, msg := range msgs {
		SendMessageToPartition(t, ctx, brokers, topic, partition, nil, []byte(msg))
	}
}

func SendMessageToPartition(t *testing.T, ctx context.Context, brokers []string, topic string, parition int, key, message []byte) {
	conn, err := net.Dial("tcp", brokers[0])
	if err != nil {
		t.Log("Failed to connect to kafka broker. %w", err)
		t.FailNow()
	}
	connConfig := kafka.ConnConfig{
		ClientID:  "kafkatail-test-client",
		Topic:     topic,
		Partition: parition,
	}

	kconn := kafka.NewConnWith(conn, connConfig)
	msg := kafka.Message{
		Key:   key,
		Value: message,
	}

	_, err = kconn.WriteMessages(msg)
	if err != nil {
		t.Log("failed to send write message to broker:", err)
		t.FailNow()
	}
}

func CreateTopicWithConfig(ctx context.Context, topic kafka.TopicConfig) error {
	resp, err := KafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{topic},
	})

	if err != nil {
		return err
	}

	return resp.Errors[topic.Topic]
}

func DeleteTopic(ctx context.Context, topic string) error {
	resp, err := KafkaRawClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})

	if err != nil {
		return err
	}

	return resp.Errors[topic]
}

func ReadChanMessages(tb *tomb.Tomb, c <-chan string) string {
	var got strings.Builder

	loop := true
	for loop {
		select {
		case msg := <-c:
			got.WriteString(msg)
		case <-tb.Dead():
			loop = false
		}
	}
	return got.String()
}
