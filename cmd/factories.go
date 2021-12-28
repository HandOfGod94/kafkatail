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
	Consume(ctx context.Context, decoder wire.Decoder) <-chan consumer.Result
}

func consumerFactory(topic, groupID string, fromDateTime time.Time) (closerConsumer, error) {
	if groupID != "" {
		return consumer.NewGroupConsumer(bootstrapServers, topic, groupID)
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

func decoderFactory(wireFormat wire.Format) wire.Decoder {
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
