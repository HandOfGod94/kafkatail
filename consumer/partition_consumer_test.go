package consumer_test

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PartitionConsumerTestSuite struct {
	suite.Suite
	opts consumer.PartitionConsumerOpts
}

func (suite *PartitionConsumerTestSuite) SetupSuite() {
	suite.opts = consumer.PartitionConsumerOpts{
		BootstrapServers: []string{"localhost:9093"},
		Topic:            "kafkatail-test-topic",
		Offset:           kafka.FirstOffset,
		Partition:        0,
		FromDateTime:     time.Time{},
	}

	kafkatest.CreateTopicWithConfig(context.Background(), kafka.TopicConfig{
		Topic:             suite.opts.Topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
}

func (suite *PartitionConsumerTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), suite.opts.Topic)
}

func (suite *PartitionConsumerTestSuite) TestNewPartitionConsumer_SuccessWithValidValues() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c, err := consumer.NewPartitionConsumer(ctx, suite.opts)
	defer c.Close()
	assert.NoError(suite.T(), err)
}

func (suite *PartitionConsumerTestSuite) TestConsume() {
	testCases := []struct {
		desc      string
		opts      consumer.PartitionConsumerOpts
		msgToSend []string
		want      string
	}{
		{
			desc:      "should consume messages with valid partition consumer opts",
			opts:      suite.opts,
			msgToSend: []string{"hello", "world"},
			want:      "hello",
		},
		{
			desc:      "should consume only from specific partition when partition number opts is provided",
			opts:      suite.opts.WithPartition(0),
			msgToSend: []string{"hello", "world"},
			want:      "hello",
		},
		{
			desc:      "should consume messages from offset when offset opts is provided",
			opts:      suite.opts.WithOffset(kafka.LastOffset),
			msgToSend: []string{"hello", "world"},
			want:      "world",
		},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			c, _ := consumer.NewPartitionConsumer(ctx, tc.opts)
			defer c.Close()
			outChan := c.Consume(ctx, wire.NewPlaintextDecoder())

			kafkatest.SendMessage(t, suite.opts.BootstrapServers, suite.opts.Topic, nil, []byte("hello"))
			kafkatest.SendMessage(t, suite.opts.BootstrapServers, suite.opts.Topic, nil, []byte("world"))

			got := kafkatest.ReadChanMessages(ctx, outChan)
			assert.Contains(t, got[0].Message, tc.want)
		})
	}
}

func TestPartitionConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(PartitionConsumerTestSuite))
}
