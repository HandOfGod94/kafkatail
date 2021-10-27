// +build integration

package kafkatail_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/hashicorp/go-multierror"
	"github.com/segmentio/kafka-go"
)

const (
	kafkaTestTopic              = "kafkatail-test"
	kafkaTestTopicWithPartition = "kafkatail-test-topic-with-partition"
	kafkaPlainTextTopic         = "kafkatail-plaintext-test-topic"
	kafkaProtoTopic             = "kafkatail-proto-test-topic"
)

func setup() error {
	var result error

	topicConfig := kafka.TopicConfig{Topic: kafkaTestTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		result = multierror.Append(result, err)
	}

	topicConfig = kafka.TopicConfig{Topic: kafkaTestTopicWithPartition, NumPartitions: 2, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		result = multierror.Append(result, err)
	}

	topicConfig = kafka.TopicConfig{Topic: kafkaPlainTextTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		result = multierror.Append(result, err)
	}

	topicConfig = kafka.TopicConfig{Topic: kafkaProtoTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func teardown() error {
	var result error

	if err := kafkatest.DeleteTopic(context.Background(), kafkaTestTopic); err != nil {
		result = multierror.Append(result, err)
	}

	if err := kafkatest.DeleteTopic(context.Background(), kafkaTestTopicWithPartition); err != nil {
		result = multierror.Append(result, err)
	}

	if err := kafkatest.DeleteTopic(context.Background(), kafkaPlainTextTopic); err != nil {
		result = multierror.Append(result, err)
	}

	if err := kafkatest.DeleteTopic(context.Background(), kafkaProtoTopic); err != nil {
		result = multierror.Append(result, err)
	}

	return result
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal("failed to setup test. error:", err)
	}
	statusCode := m.Run()

	if err := teardown(); err != nil {
		log.Fatal("failed to delete topic. error:", err)
	}
	os.Exit(statusCode)
}
