package consumer_test

import (
	"testing"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/stretchr/testify/assert"
)

type factoryargs struct {
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

func TestCreateConsumer(t *testing.T) {
	testCases := []struct {
		desc    string
		factory factoryargs
		want    consumer.ClosableConsumer
	}{
		{
			desc: "creates GroupConsumeer when groupID is present",
			factory: factoryargs{
				bootstrapServers: []string{"localhost:8081"},
				groupID:          "groupID",
				topic:            "test-topic",
				partition:        0,
				offset:           10,
				fromDateTime:     time.Time{},
				wireFormat:       wire.PlainText,
				protoFile:        "foo.proto",
				messageType:      "Bar",
				includePaths:     []string{"./hello"},
			},
			want: &consumer.GroupConsumer{},
		},
		{
			desc: "creates partitions consumer when groupID is absent",
			factory: factoryargs{
				bootstrapServers: []string{"localhost:8081"},
				topic:            "test-topic",
				partition:        0,
				offset:           10,
				fromDateTime:     time.Time{},
				wireFormat:       wire.PlainText,
				protoFile:        "foo.proto",
				messageType:      "Bar",
				includePaths:     []string{"./hello"},
			},
			want: &consumer.PartitionConsumer{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cf := consumer.NewFactory(
				tc.factory.bootstrapServers,
				tc.factory.groupID,
				tc.factory.topic,
				tc.factory.partition,
				tc.factory.offset,
				tc.factory.fromDateTime,
				tc.factory.wireFormat,
				tc.factory.protoFile,
				tc.factory.messageType,
				tc.factory.includePaths,
			)

			got, err := cf.CreateConsumer()
			assert.NoError(t, err)
			assert.IsType(t, tc.want, got)
		})
	}
}

func TestWireDecoder(t *testing.T) {
	testCases := []struct {
		desc    string
		factory factoryargs
		want    consumer.WireDecoder
	}{
		{
			desc: "creates plain text decoder for plaintext wireFormat",
			factory: factoryargs{
				bootstrapServers: []string{"localhost:8081"},
				groupID:          "groupID",
				topic:            "test-topic",
				partition:        0,
				offset:           10,
				fromDateTime:     time.Time{},
				wireFormat:       wire.PlainText,
				protoFile:        "foo.proto",
				messageType:      "Bar",
				includePaths:     []string{"./hello"},
			},
			want: &wire.PlaintextDecoder{},
		},
		{
			desc: "creates protodecoder for proto wireformat",
			factory: factoryargs{
				bootstrapServers: []string{"localhost:8081"},
				groupID:          "groupID",
				topic:            "test-topic",
				partition:        0,
				offset:           10,
				fromDateTime:     time.Time{},
				wireFormat:       wire.Proto,
				protoFile:        "foo.proto",
				messageType:      "Bar",
				includePaths:     []string{"./hello"},
			},
			want: &wire.ProtoDecoder{},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			cf := consumer.NewFactory(
				tc.factory.bootstrapServers,
				tc.factory.groupID,
				tc.factory.topic,
				tc.factory.partition,
				tc.factory.offset,
				tc.factory.fromDateTime,
				tc.factory.wireFormat,
				tc.factory.protoFile,
				tc.factory.messageType,
				tc.factory.includePaths,
			)

			got, err := cf.CreateDecoder()
			assert.NoError(t, err)
			assert.IsType(t, tc.want, got)
		})
	}
}
