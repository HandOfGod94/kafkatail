// +build integration

package kafkatail_test

import (
	"context"
	"io"
	"os/exec"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/handofgod94/kafkatail/testdata"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var starwarsHuman = testdata.Human{
	Mass:   32.0,
	Height: &testdata.Human_Height{Unit: testdata.LengthUnit_METER, Value: 1.0},
}

func TestKafkatalProto(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			desc: "prints decoded messages based on valid proto provided via args",
			cmd:  `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto`,
			want: "height:",
		},
		{
			desc:    "prints decode error when it fails to decode",
			cmd:     `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=foo.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto`,
			want:    "failed to decode message",
			wantErr: true,
		},
	}
	const topic = "kafkatail-test-proto"
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

			msg, err := proto.Marshal(&starwarsHuman)
			if err != nil {
				t.Log("failed to marshal message:", err)
				t.FailNow()
			}

			kafkatest.SendMessage(t, context.Background(), []string{localBroker}, topic, nil, msg)
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
