package consumer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

var _ ClosableConsumer = (*MultiplePartitionConsumer)(nil)

type MultiplePartitionConsumer struct {
	readers []*kafka.Reader
}

type MultiplePartitionConsumerOpts struct {
	BootstrapServers []string `validate:"required"`
	Topic            string   `validate:"required"`
	Offset           int64
	FromDateTime     time.Time
}

func NewMultiplePartitionConsumer(ctx context.Context, opts MultiplePartitionConsumerOpts) (*MultiplePartitionConsumer, error) {
	log.Printf("starting multiple partition consumer with config: %+v", opts)
	validate := validator.New()
	if errs := validate.Struct(opts); errs != nil {
		return nil, errs
	}

	partitions, err := opts.Paritions(ctx)
	if err != nil {
		return nil, err
	}

	readers := make([]*kafka.Reader, len(partitions))
	for i, partition := range partitions {
		readers[i] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:   opts.BootstrapServers,
			Topic:     opts.Topic,
			Partition: partition.ID,
		})

		if err := seekToOffset(ctx, readers[i], opts.Offset, opts.FromDateTime); err != nil {
			return nil, fmt.Errorf("failed to seek offset for partition %d. error: %w", partition.ID, err)
		}
	}

	return &MultiplePartitionConsumer{readers}, nil
}

func (opts *MultiplePartitionConsumerOpts) Paritions(ctx context.Context) ([]kafka.Partition, error) {
	conn, err := kafka.DialContext(ctx, "tcp", opts.BootstrapServers[0])
	if err != nil {
		return nil, err
	}

	return conn.ReadPartitions(opts.Topic)
}

func (mpc *MultiplePartitionConsumer) Consume(ctx context.Context, decoder WireDecoder) <-chan Result {
	out := make(chan Result)

	for _, reader := range mpc.readers {
		go func(partitionOut <-chan Result) {
			for res := range partitionOut {
				out <- res
			}
		}(printMsgs(ctx, reader, decoder))
	}

	return out
}

func (mpc *MultiplePartitionConsumer) Close() error {
	return nil
}
