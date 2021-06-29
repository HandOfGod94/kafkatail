package consumer

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
)

type kafkaConsumer struct {
	bootstrapServers []string
	topic            string
	options          Options
}

type Options struct {
	GroupID   string
	Offset    int64
	Partition int
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
		options: Options{
			Offset: kafka.LastOffset,
		},
	}
}

func (kc *kafkaConsumer) Consume(ctx context.Context, w io.Writer, decoder wire.Decoder) error {
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

		value, err := decoder.Decode(m.Value)
		if err != nil {
			log.Printf("failed to decode message. error: %v", err)
			continue
		}
		fmt.Fprintln(w, "====================Message====================")
		fmt.Fprintf(w, "============Partition: %v, Offset: %v==========\n", m.Partition, m.Offset)
		fmt.Fprintln(w, value)
	}

}

func (kc *kafkaConsumer) initReader() (*kafka.Reader, error) {
	log.Printf("Starting consumer with config: %+v", kc)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   kc.bootstrapServers,
		Topic:     kc.topic,
		Partition: kc.options.Partition,
	})

	err := r.SetOffset(kc.options.Offset)
	if err != nil {
		return nil, err
	}

	return r, nil
}
