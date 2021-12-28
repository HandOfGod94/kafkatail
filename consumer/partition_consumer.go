package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
	"gopkg.in/validator.v2"
)

type partitionConsumer struct {
	reader *kafka.Reader
}

type PartitionConsumerOpts struct {
	BootstrapServers []string `validate:"min=1"`
	Topic            string   `validate:"min=1"`
	Partition        int      `validate:"min=0"`
	Offset           int64
	FromDateTime     time.Time
}

func NewPartitionConsumer(ctx context.Context, opts PartitionConsumerOpts) (*partitionConsumer, error) {
	log.Printf("starting partition consumer with config: %+v", opts)
	if errs := validator.Validate(opts); errs != nil {
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

	return &partitionConsumer{reader}, nil
}

func (pc *partitionConsumer) Consume(ctx context.Context, decoder wire.Decoder) <-chan Result {
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

func (pc *partitionConsumer) Close() error {
	return pc.reader.Close()
}

func seekToOffset(ctx context.Context, reader *kafka.Reader, offset int64, fromDateTime time.Time) error {
	if err := reader.SetOffset(offset); err != nil {
		return err
	}

	if !fromDateTime.IsZero() {
		err := reader.SetOffsetAt(ctx, fromDateTime)
		if err != nil {
			return fmt.Errorf("failed to initialize partition consumer. error: %w", err)
		}
	}

	return nil
}
