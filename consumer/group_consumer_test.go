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
	bootstrapServer []string
	topic           string
	groupID         string
}

func (suite *GroupConsumerTestSuite) SetupSuite() {
	suite.bootstrapServer = []string{"localhost:9093"}
	suite.topic = "kafka-consume-group-id-test"
	suite.groupID = "foo-test-id"
	kafkatest.CreateTopicWithConfig(context.Background(), kafka.TopicConfig{
		Topic:             suite.topic,
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
}

func (gcts *GroupConsumerTestSuite) TearDownSuite() {
	kafkatest.DeleteTopic(context.Background(), gcts.topic)
}

func (suite *GroupConsumerTestSuite) TestNewGroupConsumer_ShouldNotReturnError() {
	c, err := consumer.NewGroupConsumer(suite.bootstrapServer, suite.topic, suite.groupID)
	//lint:ignore SA5001 this test is for checking err is nil
	defer c.Close()
	assert.NoError(suite.T(), err)
}

func (suite *GroupConsumerTestSuite) TestConsume_ReturnsMessagesAcrossAllThePartitions() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	gc, _ := consumer.NewGroupConsumer(suite.bootstrapServer, suite.topic, suite.groupID)
	defer gc.Close()
	resultChan := gc.Consume(ctx, wire.NewPlaintextDecoder())

	kafkatest.SendMultipleMessagesToPartition(suite.T(), suite.bootstrapServer, suite.topic, map[partition]string{
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
