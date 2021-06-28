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

func (ao *AppOptions) Start() {
	err :=
		consumer.Options{
			GroupID: ao.GroupID,
		}.New(ao.BootstrapServers, ao.Topic).
			Consume(context.Background(), wire.NewPlaintextDecoder(), "")

	log.Fatal("error while consuming messages:", err)
}
