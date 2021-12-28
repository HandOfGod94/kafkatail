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
	"github.com/spf13/cobra"
)

type status = int

func runKafkaTail(cmd *cobra.Command, args []string) error {
	topic := args[0]
	parsedDT, err := time.Parse(time.RFC3339, fromDateTime)
	if err != nil {
		return fmt.Errorf("invalid datetime provided: %w", err)
	}

	konsumer, err := consumerFactory(topic, groupID, parsedDT)
	if err != nil {
		return err
	}

	resultChan := konsumer.Consume(context.Background(), decoderFactory(wireForamt))
	exitCode := <-receiveMessages(resultChan)

	konsumer.Close()
	log.Printf("stopping application, with exitcode: %d", exitCode)
	os.Exit(exitCode)

	return nil
}

func receiveMessages(resultChan <-chan consumer.Result) <-chan status {
	sigs := make(chan os.Signal, 1)
	exitCode := make(chan status)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		loop := true
		for loop {
			select {
			case result := <-resultChan:
				if result.Err != nil {
					loop = false
					log.Print("error while consuming messages:", result.Err)
					exitCode <- 1
				}
				fmt.Println(result.Message)
			case <-sigs:
				loop = false
				exitCode <- 0
			}
		}
	}()

	return exitCode
}
