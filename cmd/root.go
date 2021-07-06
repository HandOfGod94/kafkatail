package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
	"gopkg.in/tomb.v2"
)

var (
	bootstrapServers []string
	groupID          string
	wireForamt       wire.Format
	protoFile        string
	includePaths     []string
	messageType      string
	offset           int64
	partition        int
	fromDateTime     string
)

const appVersion = "0.1.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kafkatail [flags] topic",
	Short:   "Tail kafka logs of any wire format",
	Long:    `Tail kafka messages from any topic, of any wire format on console (plaintext, protobuf)`,
	Version: appVersion,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]

		parsedDT, err := time.Parse(time.RFC3339, fromDateTime)
		if err != nil {
			log.Fatal("invalid datetime provided:", err)
		}

		tb, ctx := tomb.WithContext(context.Background())

		outChan := consumer.Options{
			GroupID:      groupID,
			Offset:       offset,
			Partition:    partition,
			FromDateTime: parsedDT,
		}.New(bootstrapServers, topic).Consume(ctx, tb, decoderFactory(wireForamt))

		for {
			select {
			case <-tb.Dead():
				err := tb.Err()
				log.Fatal("Stopping Application. error while consuming messages:", err)
			case msg := <-outChan:
				fmt.Println(msg)
			}
		}
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	zeroTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	rootCmd.Flags().StringSliceVarP(&bootstrapServers, "bootstrap_servers", "b", []string{}, "list of kafka `bootstrap_servers` separated by comma")
	rootCmd.Flags().StringVar(&groupID, "group_id", "", "[Optional] kafka consumer `group_id` to be used for subscribing to topic")
	rootCmd.Flags().StringVar(&protoFile, "proto_file", "", "`proto_file` to be used for decoding kafka message. Required for `wire_format=proto`")
	rootCmd.Flags().StringSliceVar(&includePaths, "include_paths", []string{}, "`include_paths` containing dependencies of proto. Required for `wire_format=proto`")
	rootCmd.Flags().StringVar(&messageType, "message_type", "", "proto message `type` to use for decoding . Required for `wire_format=proto`")
	rootCmd.Flags().Int64Var(&offset, "offset", -1, "kafka offset to start consuming from. Possible Values: -1=latest, -2=earliest, n=nth offset")
	rootCmd.Flags().IntVar(&partition, "partition", 0, "kafka partition to consume from")
	rootCmd.Flags().StringVar(&fromDateTime, "from_datetime", zeroTime, "tail from specific past datetime in RFC3339 format")
	rootCmd.Flags().Var(enumflag.New(&wireForamt, "wire_format", wire.FormatIDs, enumflag.EnumCaseSensitive),
		"wire_format",
		"Wire format of messages in topic",
	)

	rootCmd.Flags().Lookup("wire_format").NoOptDefVal = "plaintext"

	rootCmd.MarkFlagRequired("bootstrap_servers")
	rootCmd.MarkFlagRequired("topic")
}

func decoderFactory(wireFormat wire.Format) wire.Decoder {
	if wireFormat == wire.PlainText {
		return wire.NewPlaintextDecoder()
	} else if wireFormat == wire.Proto {
		return wire.NewProtoDecoder(protoFile, messageType, includePaths)
	} else {
		log.Fatalf("unsupported message type. received: %v, supported: %v", messageType, "plaintext, proto")
		return nil
	}
}
