// +build integration

package kafkatail_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestKafkatailPlaintext(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		message string
		want    string
		wantErr bool
	}{
		{
			desc:    "when topic doesn't exist",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 foobar",
			message: "hello world",
			want:    "",
		},
		{
			desc:    "when messages are present in topic",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 kafkatail-test",
			message: "hello world",
			want:    "hello world",
		},
	}

	const topic = "kafkatail-test"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	}
	kafkatest.CreateTopicWithConfig(t, context.Background(), topicConfig)
	defer kafkatest.DeleteTopic(t, context.Background(), topic)
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			appName, args := appNameAndArgs(tc.cmd)
			cmd := exec.CommandContext(ctx, appName, args...)
			out, err := getOutput(cmd, tc.wantErr)
			if err != nil {
				t.Log("failed to create output pipe:", err)
				t.FailNow()
			}

			err = cmd.Start()
			if err != nil {
				t.Logf("failed to start command: '%v'. error: %v", tc.cmd, err)
				t.FailNow()
			}
			kafkatest.SendMessage(t, context.Background(), []string{localBroker}, topic, nil, []byte(tc.message))
			got, err := io.ReadAll(out)
			if err != nil {
				t.Log("failed to read stdout:", err)
				t.FailNow()
			}

			assert.Contains(t, string(got), tc.want)

			cmd.Wait()
		})
	}
}
