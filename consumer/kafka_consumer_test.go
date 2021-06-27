package consumer_test

import (
	"context"
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

const defaultTimeout = 5 * time.Second

var kafkaRawClient = kafka.Client{
	Addr:      kafka.TCP("localhost:9093"),
	Timeout:   defaultTimeout,
	Transport: nil,
}

func TestConsume_Errors(t *testing.T) {
	testCases := []struct {
		desc             string
		bootstrapServers []string
		topic            string
		expectedErr      string
	}{
		{"with invalid brokers config", []string{"foo:9093"}, "test_topic", "no such host"},
		{"with invalid topic", []string{"localhost:9093"}, "nonexistent_topic", "deadline exceeded"},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			c := consumer.New(tc.bootstrapServers, tc.topic)
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()
			err := c.Consume(ctx)

			assert.Contains(t, err.Error(), tc.expectedErr)
		})
	}
}
