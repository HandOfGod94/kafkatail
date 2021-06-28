package app

import (
	"context"
	"log"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
)

type AppOptions struct {
	BootstrapServers []string
	Topic            string
	GroupID          string
	WireForamt       wire.Format
	ProtoFile        string
	Includes         []string
	MessageType      string
}

func (ao *AppOptions) withDecoder() wire.Decoder {
	if ao.WireForamt == wire.PlainText {
		return wire.NewPlaintextDecoder()
	} else if ao.WireForamt == wire.Proto {
		return wire.NewProtoDecoder(ao.ProtoFile, ao.MessageType, ao.Includes)
	} else {
		log.Fatalf("unsupported message type. received: %v, supported: %v", ao.MessageType, "plaintext, proto")
		return nil
	}
}

func (ao *AppOptions) Start() {
	err :=
		consumer.Options{
			GroupID: ao.GroupID,
		}.New(ao.BootstrapServers, ao.Topic).
			Consume(context.Background(), ao.withDecoder())

	log.Fatal("error while consuming messages:", err)
}
