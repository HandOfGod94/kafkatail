package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"gopkg.in/tomb.v2"
)

type AppOptions struct {
	BootstrapServers []string
	Topic            string
	GroupID          string
	WireForamt       wire.Format
}

func (ao *AppOptions) Start() {
	tumb, tmCtx := tomb.WithContext(context.Background())
	msgChan :=
		consumer.Options{
			GroupID: ao.GroupID,
		}.New(ao.BootstrapServers, ao.Topic).
			Consume(tumb, tmCtx)

	stopChan := make(chan os.Signal, 2)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	for {
		select {
		case <-tumb.Dead():
			log.Fatalf("failed to read messages from kafka. error: %v", tumb.Err())
		case msg := <-msgChan:
			fmt.Println(msg)
		case <-stopChan:
			log.Printf("Stopping application")
			os.Exit(0)
		}
	}
}
