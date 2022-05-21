package consumer

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/handofgod94/kafkatail/wire"
)

type ClosableConsumer interface {
	io.Closer
	Consume(ctx context.Context, decoder WireDecoder) <-chan Result
}

type WireDecoder interface {
	Decode([]byte) (string, error)
}

type factory struct {
	bootstrapServers []string
	groupID          string
	topic            string
	partition        int
	offset           int64
	fromDateTime     time.Time
	wireFormat       wire.Format
	protoFile        string
	messageType      string
	includePaths     []string
}

func NewFactory(
	bootstrapServers []string,
	groupID string,
	topic string,
	partition int,
	offset int64,
	fromDateTime time.Time,
	wireFormat wire.Format,
	protoFile string,
	messageType string,
	includePaths []string,
) *factory {
	return &factory{
		bootstrapServers: bootstrapServers,
		groupID:          groupID,
		topic:            topic,
		partition:        partition,
		offset:           offset,
		fromDateTime:     fromDateTime,
		wireFormat:       wireFormat,
		protoFile:        protoFile,
		messageType:      messageType,
		includePaths:     includePaths,
	}
}

func (o *factory) CreateConsumer() (ClosableConsumer, error) {
	if o.groupID != "" {
		return NewGroupConsumer(GroupConsumerOpts{
			BootstrapServers: o.bootstrapServers,
			GroupID:          o.groupID,
			Topic:            o.topic,
		})
	} else {
		return NewPartitionConsumer(context.Background(), PartitionConsumerOpts{
			BootstrapServers: o.bootstrapServers,
			Topic:            o.topic,
			Partition:        o.partition,
			Offset:           o.offset,
			FromDateTime:     o.fromDateTime,
		})
	}
}

func (o *factory) CreateDecoder() (WireDecoder, error) {
	switch o.wireFormat {
	case wire.PlainText:
		return wire.NewPlaintextDecoder(), nil
	case wire.Proto:
		return wire.NewProtoDecoder(o.protoFile, o.messageType, o.includePaths), nil
	default:
		return nil, fmt.Errorf("unsupported message type. received: %v, supported: %v", o.messageType, "plaintext, proto")
	}
}
