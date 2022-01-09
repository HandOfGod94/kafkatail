package cmd

import (
	"bytes"
	"context"
	"testing"
	"time"

	. "github.com/handofgod94/kafkatail/kafkatest"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestValidations_Errors(t *testing.T) {
	testCases := []struct {
		desc    string
		cmdArgs []string
		want    string
	}{
		{
			desc:    "topic is required args",
			cmdArgs: []string{},
			want:    `requires at least 1 arg(s), only received 0`,
		},
		{
			desc:    "bootstrap_servers is required flag",
			cmdArgs: []string{"kafka-test"},
			want:    `required flag(s) "bootstrap_servers" not set`,
		},
		{
			desc:    "wire_format can only be one of proto, plaintext or avro",
			cmdArgs: []string{"--bootstrap_servers=localhost:9093", "--wire_format=foo", "kafka-test"},
			want:    `invalid argument "foo" for "--wire_format"`,
		},
		{
			desc:    "schema_file is required if wire_format is 'avro'",
			cmdArgs: []string{"--bootstrap_servers=localhost:9093", "--wire_format=avro", "kafka-test"},
			want:    `required flag(s) "schema_file" not set`,
		},
		{
			desc:    "proto_file is required if wire_format is 'proto'",
			cmdArgs: []string{"--bootstrap_servers=localhost:9093", "--wire_format=proto", "kafka-test"},
			want:    `required flag(s) "include_paths", "message_type", "proto_file" not set`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			b := bytes.NewBufferString("")
			rootCmd.SetOut(b)
			rootCmd.SetArgs(tc.cmdArgs)
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
			defer cancel()
			err := rootCmd.ExecuteContext(ctx)
			assert.NotNil(t, err)
			assert.Contains(t, err.Error(), tc.want)
		})
	}
}

func TestValidations_Success(t *testing.T) {
	testCases := []struct {
		desc    string
		cmdArgs []string
		msgSent string
		want    string
	}{
		{
			desc:    "with valid args for plaintext messages",
			cmdArgs: []string{"--bootstrap_servers=localhost:9093", "kafka-test"},
		},
	}

	const topic = "kafka-test"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	}
	if err := CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		t.Log("failed to create topic", err)
		t.FailNow()
	}
	t.Cleanup(func() {
		if err := DeleteTopic(context.Background(), topic); err != nil {
			t.Log("failed to delete topic", err)
			t.FailNow()
		}
	})
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			b := bytes.NewBufferString("")
			rootCmd.SetOut(b)
			rootCmd.SetArgs(tc.cmdArgs)
			ctx, cancel := context.WithTimeout(rootCmd.Context(), 1*time.Second)
			defer cancel()
			err := rootCmd.ExecuteContext(ctx)
			SendMessage(t, []string{LocalBroker}, topic, nil, []byte(tc.msgSent))

			assert.NotNil(t, err)
			got := b.String()

			assert.Contains(t, got, tc.want)
		})
	}
}
