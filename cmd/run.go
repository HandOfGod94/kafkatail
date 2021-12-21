package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/spf13/cobra"
	"gopkg.in/tomb.v2"
)

func runKafkaTail(cmd *cobra.Command, args []string) error {
	topic := args[0]
	parsedDT, err := time.Parse(time.RFC3339, fromDateTime)
	if err != nil {
		return fmt.Errorf("invalid datetime provided: %w", err)
	}

	tb, ctx := tomb.WithContext(context.Background())
	c := consumer.Options{
		GroupID:      groupID,
		Offset:       offset,
		Partition:    partition,
		FromDateTime: parsedDT,
	}.New(bootstrapServers, topic)

	kr, err := c.InitReader(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize reader: %w", err)
	}

	outChan := c.Consume(ctx, tb, decoderFactory(wireForamt), kr)
	exitCode := receiveMessages(tb, outChan)

	kr.Close()
	log.Printf("stopping application, with exitcode: %d", exitCode)
	os.Exit(exitCode)

	return nil
}

func receiveMessages(tb *tomb.Tomb, outChan <-chan string) int {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	loop := true
	exitCode := 0
	for loop {
		select {
		case <-tb.Dead():
			err := tb.Err()
			log.Print("error while consuming messages:", err)
			loop = false
			exitCode = 1
		case msg := <-outChan:
			fmt.Println(msg)
		case <-sigs:
			loop = false
		}
	}
	return exitCode
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
