package consumer

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"gopkg.in/tomb.v2"
)

type kafkaConsumer struct {
	bootstrapServers []string
	topic            string
	options          Options
}

type Options struct {
	GroupID string
}

func (o Options) New(bootstrapServers []string, topic string) *kafkaConsumer {
	c := New(bootstrapServers, topic)
	c.options = o

	return c
}

func New(bootstrapServers []string, topic string) *kafkaConsumer {
	return &kafkaConsumer{
		bootstrapServers: bootstrapServers,
		topic:            topic,
		options:          Options{},
	}
}

func (kc *kafkaConsumer) Consume(ctx context.Context) (<-chan string, <-chan error) {
	outChan := make(chan string)
	errorChan := make(chan error)

	log.Printf("Starting consumer with config: %+v", kc)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kc.bootstrapServers,
		Topic:   kc.topic,
	})

	t, tCtx := tomb.WithContext(ctx)
	t.Go(func() error {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				errorChan <- err
				break
			}
			select {
			case <-tCtx.Done():
				break
			default:
				outChan <- string(m.Value)
			}
		}
		if err := r.Close(); err != nil {
			errorChan <- err
		}
		close(outChan)
		close(errorChan)
		return nil
	})

	return outChan, errorChan
}
