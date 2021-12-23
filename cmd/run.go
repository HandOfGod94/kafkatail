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
)

func runKafkaTail(cmd *cobra.Command, args []string) error {
	topic := args[0]
	parsedDT, err := time.Parse(time.RFC3339, fromDateTime)
	if err != nil {
		return fmt.Errorf("invalid datetime provided: %w", err)
	}

	c := consumer.Options{
		GroupID:      groupID,
		Offset:       offset,
		Partition:    partition,
		FromDateTime: parsedDT,
	}.New(bootstrapServers, topic)

	kr, err := c.InitReader(context.Background())
	if err != nil {
		return fmt.Errorf("failed to initialize reader: %w", err)
	}

	outChan := c.Consume(context.Background(), decoderFactory(wireForamt), kr)
	exitCode := receiveMessages(outChan)

	kr.Close()
	log.Printf("stopping application, with exitcode: %d", exitCode)
	os.Exit(exitCode)

	return nil
}

func receiveMessages(outChan <-chan consumer.Result) int {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	loop := true
	exitCode := 0
	for loop {
		select {
		case result := <-outChan:
			if result.Err != nil {
				log.Print("error while consuming messages:", result.Err)
				loop = false
				exitCode = 1
			}
			fmt.Println(result.Message)
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
