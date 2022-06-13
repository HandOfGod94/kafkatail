package consumer

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/samber/mo"
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
	partition        mo.Either[int, string]
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
	partition mo.Either[int, string],
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
	}

	if o.isMultiPartition() {
		return NewMultiplePartitionConsumer(context.Background(), MultiplePartitionConsumerOpts{
			BootstrapServers: o.bootstrapServers,
			Topic:            o.topic,
			Offset:           o.offset,
			FromDateTime:     o.fromDateTime,
		})
	}

	return NewPartitionConsumer(context.Background(), PartitionConsumerOpts{
		BootstrapServers: o.bootstrapServers,
		Topic:            o.topic,
		Partition:        o.partition.MustLeft(),
		Offset:           o.offset,
		FromDateTime:     o.fromDateTime,
	})
}

func (o *factory) isMultiPartition() bool {
	return o.partition.IsRight()
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
