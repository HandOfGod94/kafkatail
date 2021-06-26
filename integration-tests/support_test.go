// +build integration

package kafktail_test

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

var kafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   5 * time.Second,
	Transport: nil,
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

func createTopic(ctx context.Context, topic string) {
	resp, err := kafkaRawClient.CreateTopics(ctx, &kafka.CreateTopicsRequest{
		Topics: []kafka.TopicConfig{{Topic: topic, NumPartitions: 1, ReplicationFactor: 1}},
	})

	if err != nil {
		log.Fatal("failed to create topic:", err)
	}

	if resp.Errors[topic] != nil {
		log.Fatalf("failed to create topic. errors: %+v", resp.Errors)
	}
}

func deleteTopic(ctx context.Context, topic string) {
	resp, err := kafkaRawClient.DeleteTopics(ctx, &kafka.DeleteTopicsRequest{
		Topics: []string{topic},
	})

	if err != nil {
		log.Fatal("failed to delete topic:", err)
	}

	if resp.Errors[topic] != nil {
		log.Fatalf("failed to create topic. errors: %+v", resp.Errors)
	}
}

func appNameAndArgs(cmd string) (appName string, args []string) {
	tokens := strings.Split(cmd, " ")
	appName = tokens[0]
	args = tokens[1:]
	return
}
