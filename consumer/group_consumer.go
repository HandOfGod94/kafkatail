package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

type GroupConsumer struct {
	reader *kafka.Reader
}

type GroupConsumerOpts struct {
	BootstrapServers []string `validate:"required"`
	GroupID          string   `validate:"required"`
	Topic            string   `validate:"required"`
}

func NewGroupConsumer(opts GroupConsumerOpts) (*GroupConsumer, error) {
	log.Printf("starting group consumer with config: %+v", opts)

	validate := validator.New()
	if err := validate.Struct(opts); err != nil {
		return nil, err
	}

	return &GroupConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: opts.BootstrapServers,
			GroupID: opts.GroupID,
			Topic:   opts.Topic,
		}),
	}, nil

}

func (gc *GroupConsumer) Consume(ctx context.Context, decoder WireDecoder) <-chan Result {
	resultChan := make(chan Result)

	go func(ctx context.Context) {
		for {
			m, err := gc.reader.ReadMessage(ctx)
			if err != nil {
				resultChan <- Result{Err: err}
				return
			}

			value, err := decoder.Decode(m.Value)
			if err != nil {
				log.Printf("failed to decode message. error: %v", err)
				resultChan <- Result{Err: err}
				return
			}
			msg := bytes.NewBufferString("")
			fmt.Fprintln(msg, "====================Message====================")
			fmt.Fprintf(msg, "============Partition: %v, Offset: %v==========\n", m.Partition, m.Offset)
			fmt.Fprintln(msg, value)
			resultChan <- Result{Message: msg.String()}
		}
	}(ctx)

	return resultChan
}

func (gc *GroupConsumer) Close() error {
	log.Println("closing group consumer")
	return gc.reader.Close()
}
