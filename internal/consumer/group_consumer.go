package consumer

import (
	"context"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

var _ ClosableConsumer = (*GroupConsumer)(nil)

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
	return printMsgs(ctx, gc.reader, decoder)
}

func (gc *GroupConsumer) Close() error {
	log.Println("closing group consumer")
	return gc.reader.Close()
}
