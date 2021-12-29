package kafkatest

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/segmentio/kafka-go"
)

var DefaultTimeout = 3 * time.Second
var KafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   DefaultTimeout,
	Transport: nil,
}

func SendMessage(t *testing.T, brokers []string, topic string, key, message []byte) {
	conn, err := net.Dial("tcp", brokers[0])
	if err != nil {
		t.Log("failed to connect to kafka broker:", err)
		t.FailNow()
	}
	defer conn.Close()

	kconn := kafka.NewConnWith(conn, kafka.ConnConfig{ClientID: "kafkatail-test-client", Topic: topic})
	defer kconn.Close()
	_, err = kconn.WriteMessages(kafka.Message{Key: key, Value: message})
	if err != nil {
		t.Log("failed to write messages: ", err)
		t.FailNow()
	}
}

func SendMultipleMessagesToPartition(t *testing.T, brokers []string, topic string, msgs map[int]string) {
	for partition, msg := range msgs {
		SendMessageToPartition(t, brokers, topic, partition, nil, []byte(msg))
	}
}

func SendMessageToPartition(t *testing.T, brokers []string, topic string, partition int, key, message []byte) {
	conn, err := net.Dial("tcp", brokers[0])
	if err != nil {
		t.Log("Failed to connect to kafka broker. %w", err)
		t.FailNow()
	}
	defer conn.Close()

	kconn := kafka.NewConnWith(conn, kafka.ConnConfig{
		ClientID:  "kafkatail-test-client",
		Topic:     topic,
		Partition: partition,
	})
	defer kconn.Close()

	_, err = kconn.WriteMessages(kafka.Message{Key: key, Value: message})
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

func ReadChanMessages(ctx context.Context, c <-chan consumer.Result) []consumer.Result {
	got := make([]consumer.Result, 0)

	loop := true
	for loop {
		select {
		case <-ctx.Done():
			loop = false
			got = append(got, consumer.Result{Err: ctx.Err()})
		case msg := <-c:
			got = append(got, msg)
		}
	}

	return got
}
