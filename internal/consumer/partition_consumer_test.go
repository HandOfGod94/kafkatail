package consumer_test

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/internal/consumer"
	"github.com/handofgod94/kafkatail/internal/wire"
	"github.com/handofgod94/kafkatail/kafkatest"
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
	//lint:ignore SA5001 this test is for checking err is nil
	defer c.Close()
	assert.NoError(suite.T(), err)
}

func (suite *PartitionConsumerTestSuite) TestConsume() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c, _ := consumer.NewPartitionConsumer(ctx, suite.opts)
	defer c.Close()
	outChan := c.Consume(ctx, wire.NewPlaintextDecoder())

	kafkatest.SendMessage(suite.T(), suite.opts.BootstrapServers, suite.opts.Topic, nil, []byte("hello"))
	kafkatest.SendMessage(suite.T(), suite.opts.BootstrapServers, suite.opts.Topic, nil, []byte("world"))

	got := kafkatest.ReadChanMessages(ctx, outChan)
	assert.Contains(suite.T(), got[0].Message, "hello")
	assert.Contains(suite.T(), got[1].Message, "world")
}

func TestPartitionConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(PartitionConsumerTestSuite))
}
