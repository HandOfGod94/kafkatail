// +build integration

package kafkatail_test

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/kafkatest"
	. "github.com/handofgod94/kafkatail/kafkatest"
	"github.com/handofgod94/kafkatail/testdata"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/proto"
)

const kafkaProtoTopic = "kafkatail-proto-test-topic"

type ProtoTestSuite struct {
	suite.Suite
	testHuman testdata.Human
}

func (suite *ProtoTestSuite) SetupSuite() {
	suite.testHuman = testdata.Human{
		Mass:   32.0,
		Height: &testdata.Human_Height{Unit: testdata.LengthUnit_METER, Value: 1.0},
	}

	topicConfig := kafka.TopicConfig{Topic: kafkaProtoTopic, NumPartitions: 1, ReplicationFactor: 1}
	if err := kafkatest.CreateTopicWithConfig(context.Background(), topicConfig); err != nil {
		suite.FailNow("failed to create test topic:", err)
	}
}

func (suite *ProtoTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), kafkaProtoTopic)
}

func (suite *ProtoTestSuite) TestKafkatalProto() {
	testCases := []struct {
		desc    string
		cmd     string
		want    string
		wantErr bool
	}{
		{
			desc: "prints decoded messages based on valid proto provided via args",
			cmd:  `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --offset=-2 --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-proto-test-topic`,
			want: "height:",
		},
		{
			desc:    "prints decode error when it fails to decode",
			cmd:     `kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --offset=-2 --proto_file=foo.proto --include_paths="../testdata" --message_type=Human kafkatail-proto-test-topic`,
			want:    "failed to decode message",
			wantErr: true,
		},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			msg, err := proto.Marshal(&suite.testHuman)
			if err != nil {
				t.Log("failed to marshal message:", err)
				t.FailNow()
			}

			cmd := Command{T: t, Cmd: tc.cmd, WantErr: tc.wantErr}
			cmd.Execute(ctx)

			SendMessage(t, []string{LocalBroker}, kafkaProtoTopic, nil, msg)

			got := cmd.GetOutput()
			assert.Contains(t, got, tc.want)
		})
	}
}

func TestProtoTestSuite(t *testing.T) {
	suite.Run(t, new(ProtoTestSuite))
}
