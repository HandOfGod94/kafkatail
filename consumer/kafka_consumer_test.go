package consumer_test

import (
	"context"
	"testing"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/kafkatest"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func TestConsumeSuccses(t *testing.T) {
	type fields struct {
		bootstrapServers []string
		topic            string
	}

	testCases := []struct {
		desc      string
		fields    fields
		opts      consumer.Options
		msgToSend []string
		want      string
	}{
		{
			desc: "with default values",
			fields: fields{
				bootstrapServers: []string{"localhost:9093"},
				topic:            "kafkatail-test-topic",
			},
			msgToSend: []string{"hello", "world"},
			want:      "hello",
		},
		{
			desc: "with partition option",
			fields: fields{
				bootstrapServers: []string{"localhost:9093"},
				topic:            "kafkatail-test-topic",
			},
			msgToSend: []string{"hello", "world"},
			opts:      consumer.Options{Partition: 0},
			want:      "hello",
		},
		{
			desc: "with offset option",
			fields: fields{
				bootstrapServers: []string{"localhost:9093"},
				topic:            "kafkatail-test-topic",
			},
			msgToSend: []string{"hello", "world"},
			opts:      consumer.Options{Offset: kafka.LastOffset},
			want:      "hello",
		},
	}

	kafkatest.CreateTopicWithConfig(context.Background(), kafka.TopicConfig{Topic: "kafkatail-test-topic", NumPartitions: 2, ReplicationFactor: 1})
	defer kafkatest.DeleteTopic(context.Background(), "kafkatail-test-topic")
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), kafkatest.DefaultTimeout)
			defer cancel()

			c := consumer.New(tc.fields.bootstrapServers, tc.fields.topic)
			kr, _ := c.InitReader(ctx)
			outChan := c.Consume(ctx, wire.NewPlaintextDecoder(), kr)

			kafkatest.SendMessage(t, ctx, tc.fields.bootstrapServers, tc.fields.topic, nil, []byte("hello"))
			kafkatest.SendMessage(t, ctx, tc.fields.bootstrapServers, tc.fields.topic, nil, []byte("world"))

			got := kafkatest.ReadChanMessages(ctx, outChan)
			assert.Contains(t, got[0].Message, tc.want)
		})
	}
}

func TestCreateConsumerWithOptions(t *testing.T) {
	opts := consumer.Options{
		GroupID:   "foo",
		Offset:    10,
		Partition: 2,
	}

	c := opts.New([]string{"localhost:9093"}, "footest")

	assert.NotNil(t, c)
}

func TestConsume_Errors(t *testing.T) {
	testCases := []struct {
		desc             string
		bootstrapServers []string
		topic            string
		expectedErr      string
	}{
		{"with invalid brokers config", []string{"foo:9093"}, "test_topic", "lookup foo:"},
		{"with invalid topic", []string{"localhost:9093"}, "nonexistent_topic", "deadline exceeded"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), kafkatest.DefaultTimeout)
			defer cancel()

			c := consumer.New(tc.bootstrapServers, tc.topic)
			kr, _ := c.InitReader(ctx)
			outChan := c.Consume(ctx, wire.NewPlaintextDecoder(), kr)

			got := kafkatest.ReadChanMessages(ctx, outChan)

			assert.Contains(t, got[0].Err.Error(), tc.expectedErr)
		})
	}
}
