package cmd

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
)

type closerConsumer interface {
	io.Closer
	Consume(ctx context.Context, decoder consumer.WireDecoder) <-chan consumer.Result
}

func consumerFactory(topic, groupID string, fromDateTime time.Time) (closerConsumer, error) {
	if groupID != "" {
		return consumer.NewGroupConsumer(consumer.GroupConsumerOpts{
			BootstrapServers: bootstrapServers,
			GroupID:          groupID,
			Topic:            topic,
		})
	} else {
		return consumer.NewPartitionConsumer(context.Background(), consumer.PartitionConsumerOpts{
			BootstrapServers: bootstrapServers,
			Topic:            topic,
			Partition:        partition,
			Offset:           offset,
			FromDateTime:     fromDateTime,
		})
	}
}

func decoderFactory(wireFormat wire.Format) consumer.WireDecoder {
	switch wireFormat {
	case wire.PlainText:
		return wire.NewPlaintextDecoder()
	case wire.Proto:
		return wire.NewProtoDecoder(protoFile, messageType, includePaths)
	default:
		log.Fatalf("unsupported message type. received: %v, supported: %v", messageType, "plaintext, proto")
		return nil
	}
}
