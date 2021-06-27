package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
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

func (kc *kafkaConsumer) Consume(ctx context.Context) error {
	r, err := kc.initReader()
	if err != nil {
		log.Fatal("failed to initialize kafka consumer:", err)
	}

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			// TODO: return custom wrapped error with contextual info
			return err
		}
		fmt.Println(string(m.Value))
	}

}

func (kc *kafkaConsumer) initReader() (*kafka.Reader, error) {
	log.Printf("Starting consumer with config: %+v", kc)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: kc.bootstrapServers,
		Topic:   kc.topic,
	})

	// TODO: move offset to options
	err := r.SetOffset(kafka.LastOffset)
	if err != nil {
		return nil, err
	}

	return r, nil
}
