package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

var _ ClosableConsumer = (*PartitionConsumer)(nil)

type PartitionConsumer struct {
	reader *kafka.Reader
}

type PartitionConsumerOpts struct {
	BootstrapServers []string `validate:"required"`
	Topic            string   `validate:"required"`
	Partition        int      `validate:"gte=0"`
	Offset           int64
	FromDateTime     time.Time
}

func NewPartitionConsumer(ctx context.Context, opts PartitionConsumerOpts) (*PartitionConsumer, error) {
	log.Printf("starting partition consumer with config: %+v", opts)
	validate := validator.New()
	if errs := validate.Struct(opts); errs != nil {
		return nil, errs
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   opts.BootstrapServers,
		Topic:     opts.Topic,
		Partition: opts.Partition,
	})

	if err := seekToOffset(ctx, reader, opts.Offset, opts.FromDateTime); err != nil {
		return nil, fmt.Errorf("failed to initialize parition consumer: %w", err)
	}

	return &PartitionConsumer{reader}, nil
}

func (pc *PartitionConsumer) Consume(ctx context.Context, decoder WireDecoder) <-chan Result {
	return printMsgs(ctx, pc.reader, decoder)
}

func (pc *PartitionConsumer) Close() error {
	return pc.reader.Close()
}

func seekToOffset(ctx context.Context, reader *kafka.Reader, offset int64, fromDateTime time.Time) error {
	if !fromDateTime.IsZero() {
		return reader.SetOffsetAt(ctx, fromDateTime)
	}

	return reader.SetOffset(offset)
}
