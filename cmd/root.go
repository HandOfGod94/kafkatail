/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
	"github.com/thediveo/enumflag"
	"gopkg.in/tomb.v2"
)

type WireFormat enumflag.Flag

const (
	PlainTextFormat WireFormat = iota
	ProtoFormat
	AvroFormat
)

var (
	bootstrapServers []string
	topic            string
	groupID          string
	wireForamt       WireFormat
)

var WireFormatIDs = map[WireFormat][]string{
	PlainTextFormat: {"plaintext"},
	ProtoFormat:     {"proto"},
	AvroFormat:      {"avro"},
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkatail",
	Short: "Tail kafka logs of any wire format",
	Long: `Print kafka messages from any topic, of any wire format (avro, plaintext, protobuf)
on console`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: make timeouts configurable
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tm, tctx := tomb.WithContext(ctx)
		stopChan := make(chan os.Signal, 2)
		signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
		msgChan :=
			consumer.Options{
				GroupID: groupID,
			}.New(bootstrapServers, topic).
				Consume(tm, tctx)

		for {
			select {
			case <-tm.Dead():
				log.Fatalf("failed to read messages from kafka. error: %v", tm.Err())
			case msg := <-msgChan:
				fmt.Println(msg)
			case <-stopChan:
				log.Printf("Stopping application")
				os.Exit(0)
			}
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringSliceVar(&bootstrapServers, "bootstrap_servers", []string{}, "list of kafka `bootstrap_servers` separated by comma")
	rootCmd.Flags().StringVar(&topic, "topic", "", "`topic` whose message you want to tail")
	rootCmd.Flags().StringVar(&groupID, "group_id", "", "kafka consumer `group_id` to be used for subscribing to topic")
	rootCmd.Flags().Var(enumflag.New(&wireForamt, "wire_format", WireFormatIDs, enumflag.EnumCaseSensitive),
		"wire_format",
		"Wire format of messages in topic",
	)

	rootCmd.Flags().Lookup("wire_format").NoOptDefVal = "plaintext"

	rootCmd.MarkFlagRequired("bootstrap_servers")
	rootCmd.MarkFlagRequired("topic")
}
