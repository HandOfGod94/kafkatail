package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/segmentio/kafka-go"
)

type groupConsumer struct {
	reader *kafka.Reader
}

func NewGroupConsumer(bootstrapServers []string, topic, groupID string, offset int64) (*groupConsumer, error) {
	log.Printf("starting group consumer with config: bootstrapServers %v, topic %s, groupID %s", bootstrapServers, topic, groupID)
	return &groupConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:     bootstrapServers,
			GroupID:     groupID,
			Topic:       topic,
			StartOffset: offset,
		}),
	}, nil

}

func (gc *groupConsumer) Consume(ctx context.Context, decoder wire.Decoder) <-chan Result {
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

func (gc *groupConsumer) Close() error {
	log.Println("closing group consumer")
	return gc.reader.Close()
}
