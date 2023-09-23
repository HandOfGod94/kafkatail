package consumer

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

type Result struct {
	Message string
	Err     error
}

func printMsgs(ctx context.Context, reader *kafka.Reader, decoder WireDecoder) <-chan Result {
	resultChan := make(chan Result)
	go func(ctx context.Context) {
		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				resultChan <- Result{Err: err}
				return
			}

			msg := bytes.NewBufferString("")
			fmt.Fprintln(msg, "====================Message====================")
			fmt.Fprintf(msg, "============Partition: %v, Offset: %v==========\n", m.Partition, m.Offset)
			fmt.Fprintln(msg, "====================Header====================")
			for _, header := range m.Headers {
				fmt.Fprintf(msg, "%s: %s\n", header.Key, string(header.Value))
			}

			value, err := decoder.Decode(m.Value)
			if err != nil {
				log.Printf("failed to decode message. error: %v", err)
				resultChan <- Result{Err: err}
				return
			}
			fmt.Fprintln(msg, "====================Payload====================")
			fmt.Fprintln(msg, value)
			resultChan <- Result{Message: msg.String()}
		}
	}(ctx)

	return resultChan
}
