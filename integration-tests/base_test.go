// +build integration

package kafkatail_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	. "github.com/handofgod94/kafkatail/kafkatest"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	kafkaTestTopic              = "kafkatail-test"
	kafkaTestTopicWithPartition = "kafkatail-test-topic-with-partition"
)

type BaseTestSuite struct {
	suite.Suite
}

func (suite *BaseTestSuite) SetupSuite() {
	topicConfig := kafka.TopicConfig{Topic: kafkaTestTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		suite.FailNow("failed to create test topic:", err)
	}

	topicConfig = kafka.TopicConfig{Topic: kafkaTestTopicWithPartition, NumPartitions: 2, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		suite.FailNow("failed to create test-topic-with-partition topic:", err)
	}
}

func (suite *BaseTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), kafkaTestTopic)
	kafkatest.DeleteTopic(context.Background(), kafkaTestTopicWithPartition)
}

func (suite *BaseTestSuite) TestKafkatailBase() {
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
			want: "0.1.4",
		},
		{
			desc:    "print error for missing required args",
			cmd:     "kafkatail",
			want:    `requires at least 1 arg(s)`,
			wantErr: true,
		},
		{
			desc:    "print error for invalid wire_format",
			cmd:     "kafkatail --bootstrap_servers=1.1.1.1:9093 --wire_format=foo test",
			want:    "must be 'plaintext', 'proto'",
			wantErr: true,
		},
		{
			desc: "prints messages for valid args",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-test",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with offset and partition option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --partition=0 --offset=-2 kafkatail-test",
			msg:  "hello world",
			want: "Partition: 0, Offset: 0",
		},
		{
			desc: "with `from_datetime` option",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with shorthand flag",
			cmd:  "kafkatail -b=localhost:9093 --offset=-2 kafkatail-test",
			msg:  "hello world",
			want: "hello world",
		},
		{
			desc: "with spaces instead of `=` in command",
			cmd:  "kafkatail -b localhost:9093 --offset=-2 kafkatail-test",
			msg:  "hello world",
			want: "hello world",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMessage(t, []string{LocalBroker}, kafkaTestTopic, nil, []byte(tc.msg))
			got := cmd.GetOutput()

			assert.Contains(t, got, tc.want)
		})
	}
}

func (suite *BaseTestSuite) TestTailForMultiplePartitions() {
	testCases := []struct {
		desc         string
		cmd          string
		messages     map[int]string
		wantMessages []string
		wantErr      bool
	}{
		{
			desc: "tail with group_id flag",
			cmd:  "kafkatail --bootstrap_servers=localhost:9093 --group_id=myfoo --offset=-2 kafkatail-test-topic-with-partition",
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
	}
	for _, tc := range testCases {
		suite.T().Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMultipleMessagesToPartition(t, []string{LocalBroker}, kafkaTestTopicWithPartition, tc.messages)

			got := cmd.GetOutput()

			for _, want := range tc.wantMessages {
				assert.Contains(t, MinifyString(got), MinifyString(want))
			}
			assert.Equal(t, len(tc.wantMessages), strings.Count(got, "Message"))
		})
	}
}

func TestBaseTestSuite(t *testing.T) {
	suite.Run(t, new(BaseTestSuite))
}
