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

type partition = int

type GroupConsumerTestSuite struct {
	suite.Suite
	opts consumer.GroupConsumerOpts
}

func (suite *GroupConsumerTestSuite) SetupSuite() {
	suite.opts = consumer.GroupConsumerOpts{
		BootstrapServers: []string{"localhost:9093"},
		Topic:            "kafka-consume-group-id-test",
		GroupID:          "foo-test-id",
	}
	kafkatest.CreateTopicWithConfig(context.Background(), kafka.TopicConfig{
		Topic:             suite.opts.Topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
}

func (suite *GroupConsumerTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), suite.opts.Topic)
}

func (suite *GroupConsumerTestSuite) TestNewGroupConsumer_ShouldNotReturnError() {
	c, err := consumer.NewGroupConsumer(suite.opts)
	//lint:ignore SA5001 this test is for checking err is nil
	defer c.Close()
	assert.NoError(suite.T(), err)
}

func (suite *GroupConsumerTestSuite) TestConsume_ReturnsMessagesAcrossAllThePartitions() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	gc, _ := consumer.NewGroupConsumer(suite.opts)
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

func TestGroupConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(GroupConsumerTestSuite))
}
