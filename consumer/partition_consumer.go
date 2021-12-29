package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/segmentio/kafka-go"
)

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
	resultChan := make(chan Result)

	go func(ctx context.Context) {
		for {
			m, err := pc.reader.ReadMessage(ctx)
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

func (pc *PartitionConsumer) Close() error {
	return pc.reader.Close()
}

func seekToOffset(ctx context.Context, reader *kafka.Reader, offset int64, fromDateTime time.Time) error {
	if !fromDateTime.IsZero() {
		return reader.SetOffsetAt(ctx, fromDateTime)
	}

	return reader.SetOffset(offset)
}
