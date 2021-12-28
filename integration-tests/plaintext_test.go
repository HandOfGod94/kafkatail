// +build integration

package kafkatail_test

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	. "github.com/handofgod94/kafkatail/kafkatest"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const kafkaPlainTextTopic = "kafkatail-plaintext-test-topic"

type PlainTextTestSuite struct {
	suite.Suite
}

func (suite *PlainTextTestSuite) SetupSuite() {
	topicConfig := kafka.TopicConfig{Topic: kafkaPlainTextTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		suite.FailNow("failed to create test topic:", err)
	}
}

func (suite *PlainTextTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), kafkaPlainTextTopic)
}

func (suite *PlainTextTestSuite) TestKafkatailPlaintext() {
	testCases := []struct {
		desc    string
		cmd     string
		message string
		want    string
		wantErr bool
	}{
		{
			desc:    "when topic doesn't exist",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 foobar",
			message: "hello world",
			want:    "",
		},
		{
			desc:    "when messages are present in topic",
			cmd:     "kafkatail --bootstrap_servers=localhost:9093 --offset=-2 kafkatail-plaintext-test-topic",
			message: "hello world",
			want:    "hello world",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMessage(t, []string{LocalBroker}, kafkaPlainTextTopic, nil, []byte(tc.message))
			got := cmd.GetOutput()
			assert.Contains(t, string(got), tc.want)
		})
	}
}

func TestPlainTextTestSuite(t *testing.T) {
	suite.Run(t, new(PlainTextTestSuite))
}
