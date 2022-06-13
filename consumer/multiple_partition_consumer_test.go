package consumer_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type MultipleParitionTestSuite struct {
	suite.Suite
	opts consumer.MultiplePartitionConsumerOpts
}

func (suite *MultipleParitionTestSuite) SetupSuite() {
	suite.opts = consumer.MultiplePartitionConsumerOpts{
		BootstrapServers: []string{"localhost:9093"},
		Topic:            "kafkatail-test-topic",
		Offset:           kafka.FirstOffset,
	}

	kafkatest.CreateTopicWithConfig(context.Background(), kafka.TopicConfig{
		Topic:             suite.opts.Topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
}

func (suite *MultipleParitionTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), suite.opts.Topic)
}

func (suite *MultipleParitionTestSuite) TestListPartitions_FetchesAllAvailablePartitions() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	partitions, err := suite.opts.Paritions(ctx)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, len(partitions))
}

func (suite *MultipleParitionTestSuite) TestListPartitions_RetrunsErrorForInvalidConfig() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	opts := consumer.MultiplePartitionConsumerOpts{
		BootstrapServers: []string{"foobar:9093"},
		Topic:            "kafkatail-test-topic",
		Offset:           kafka.FirstOffset,
		FromDateTime:     time.Time{},
	}

	partitions, err := opts.Paritions(ctx)
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), partitions)
}

func (suite *MultipleParitionTestSuite) TestConsume() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gc, _ := consumer.NewMultiplePartitionConsumer(ctx, suite.opts)
	defer gc.Close()

	resultChan := gc.Consume(ctx, wire.NewPlaintextDecoder())

	kafkatest.SendMultipleMessagesToPartition(suite.T(), suite.opts.BootstrapServers, suite.opts.Topic, map[partition]string{
		0: "foo",
		1: "bar",
	})

	results := kafkatest.ReadChanMessages(ctx, resultChan)
	var got strings.Builder
	for _, res := range results[:len(results)-1] {
		got.WriteString(res.Message)
	}

	assert.Contains(suite.T(), got.String(), "Partition: 0")
	assert.Contains(suite.T(), got.String(), "Partition: 1")
	assert.Contains(suite.T(), got.String(), "foo")
	assert.Contains(suite.T(), got.String(), "bar")
}

func TestListPartitionSuite(t *testing.T) {
	suite.Run(t, new(MultipleParitionTestSuite))
}
