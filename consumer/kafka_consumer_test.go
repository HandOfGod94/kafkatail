package consumer_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"gopkg.in/tomb.v2"
)

const defaultTimeout = 5 * time.Second

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
			desc: "with offsett option",
			fields: fields{
				bootstrapServers: []string{"localhost:9093"},
				topic:            "kafkatail-test-topic",
			},
			msgToSend: []string{"hello", "world"},
			opts:      consumer.Options{Offset: kafka.LastOffset},
			want:      "hello",
		},
	}

	createTopicWithConfig(t, context.Background(), kafka.TopicConfig{
		Topic:             "kafkatail-test-topic",
		NumPartitions:     2,
		ReplicationFactor: 1,
	})
	defer deleteTopic(t, context.Background(), "kafkatail-test-topic")
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {

			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			tb, tctx := tomb.WithContext(ctx)

			c := consumer.New(tc.fields.bootstrapServers, tc.fields.topic)
			outChan := c.Consume(tctx, tb, wire.NewPlaintextDecoder())

			sendMessage(t, ctx, tc.fields.bootstrapServers, tc.fields.topic, nil, []byte("hello"))
			sendMessage(t, ctx, tc.fields.bootstrapServers, tc.fields.topic, nil, []byte("world"))

			got := bytes.NewBufferString("")
			select {
			case msg := <-outChan:
				fmt.Fprint(got, msg)
			case <-ctx.Done():
				// no op
			}

			assert.Contains(t, got.String(), tc.want)
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
			c := consumer.New(tc.bootstrapServers, tc.topic)
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			tb, tctx := tomb.WithContext(ctx)
			_ = c.Consume(tctx, tb, wire.NewPlaintextDecoder())

			<-tb.Dead()
			assert.Contains(t, tb.Err().Error(), tc.expectedErr)
		})
	}
}
