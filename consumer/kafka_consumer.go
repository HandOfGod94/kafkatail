package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"gopkg.in/tomb.v2"
)

type Message = string

type kafkaConsumer struct {
	bootstrapServers []string
	topic            string
	options          Options
}

type Options struct {
	GroupID      string
	Offset       int64
	Partition    int
	FromDateTime time.Time
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

func (kc *kafkaConsumer) Consume(ctx context.Context, tb *tomb.Tomb, decoder wire.Decoder) <-chan Message {
	r, err := kc.initReader(ctx)
	if err != nil {
		log.Fatal("failed to initialize kafka consumer:", err)
	}

	outChan := make(chan string)

	tb.Go(func() error {
		for {
			m, err := r.ReadMessage(ctx)
			if err != nil {
				// TODO: return custom wrapped error with contextual info
				return err
			}

			value, err := decoder.Decode(m.Value)
			if err != nil {
				log.Printf("failed to decode message. error: %v", err)
				return err
			}
			msg := bytes.NewBufferString("")
			fmt.Fprintln(msg, "====================Message====================")
			fmt.Fprintf(msg, "============Partition: %v, Offset: %v==========\n", m.Partition, m.Offset)
			fmt.Fprintln(msg, value)
			outChan <- msg.String()
		}
	})

	return outChan

}

func (kc *kafkaConsumer) initReader(ctx context.Context) (*kafka.Reader, error) {
	log.Printf("Starting consumer with config: %+v", kc)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   kc.bootstrapServers,
		Topic:     kc.topic,
		Partition: kc.options.Partition,
	})

	if !kc.options.FromDateTime.IsZero() {
		err := r.SetOffsetAt(ctx, kc.options.FromDateTime)
		if err != nil {
			return nil, err
		}
		return r, nil
	}

	err := r.SetOffset(kc.options.Offset)
	if err != nil {
		return nil, err
	}

	return r, nil
}
