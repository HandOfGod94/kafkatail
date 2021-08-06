// +build integration

package kafkatail_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestKafkatailBase(t *testing.T) {
	testCases := []struct {
		desc    string
		cmd     string
		want    string
		msg     string
		wantErr bool
	}{
		{
			desc: "check version",
			cmd:  "kafkatail --version",
			want: "dev",
		},
		{
			desc:    "print error for missing required args",
			cmd:     "kafkatail",
			want:    `requires at least 1 arg(s)`,
			wantErr: true,
		},
		{
			desc:    "print error for invalid wire_fomrat",
			cmd:     "kafkatail --bootstrap_servers=1.1.1.1:9093 --wire_format=foo test",
			want:    "must be 'avro', 'plaintext', 'proto'",
			wantErr: true,
		},
		{
			desc: "prints messages for valid args",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with offset option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "Offset: 0",
		},
		{
			desc: "with partition option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --partition=0 kafkatail-test-base",
			msg:  "hello world",
			want: "Partition: 0",
		},
		{
			desc: "with `from_datetime` option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with shorthand flag",
			cmd:  "kafkatail -b=localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with spaces instead of `=` in command",
			cmd:  "kafkatail -b localhost:9093 --offset=-2 kafkatail-test-base",
			msg:  "hello world",
			want: "hello world",
		},
	}

	const topic = "kafkatail-test-base"
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
			cmd := command{t: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.execute(ctx)

			kafkatest.SendMessage(t, context.Background(), []string{localBroker}, topic, nil, []byte(tc.msg))
			got := cmd.getOutput()

			assert.Contains(t, got, tc.want)
		})
	}
}

func TestTailForMultipleParitions(t *testing.T) {
	topic := "kafka-consume-gorup-id-int-test"
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	}
	kafkatest.CreateTopicWithConfig(t, context.Background(), topicConfig)
	defer kafkatest.DeleteTopic(t, context.Background(), topic)

	testCases := []struct {
		desc         string
		cmd          string
		messages     map[int]string
		wantMessages []string
		wantErr      bool
	}{
		{
			desc: "tail with group_id flag",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --group_id=myfoo kafka-consume-gorup-id-int-test",
			messages: map[int]string{
				0: "hello",
				1: "world",
			},
			wantMessages: []string{
				`
				====================Message====================
				============Partition: 0, Offset: 0==========
				hello
				`,
				`
				====================Message====================
				============Partition: 1, Offset: 0==========
				world
				`,
			},
		},
		{
			desc: "tail withput group_id flag",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-1 kafka-consume-gorup-id-int-test",
			messages: map[int]string{
				0: "hello",
				1: "world",
			},
			wantMessages: []string{
				`
				====================Message====================
				============Partition: 0, Offset: 0==========
				hello
				`,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := command{t: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.execute(ctx)

			kafkatest.SendMultipleMessagesToParition(t, context.Background(), []string{localBroker}, topic, tc.messages)

			got := cmd.getOutput()
			actual := sanitizeString(string(got))

			for _, wt := range tc.wantMessages {
				assert.Contains(t, actual, sanitizeString(wt))
			}
			assert.Equal(t, len(tc.wantMessages), strings.Count(actual, "Message"))
		})
	}
}
