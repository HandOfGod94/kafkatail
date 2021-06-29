package cmd

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/handofgod94/kafkatail/consumer"
	"github.com/handofgod94/kafkatail/wire"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
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
	Use:   "kafkatail [flags] topic",
	Short: "Tail kafka logs of any wire format",
	Long: `Print kafka messages from any topic, of any wire format (avro, plaintext, protobuf)
on console`,
	Version: appVersion,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		topic := args[0]

		parsedDT, err := time.Parse(time.RFC3339, fromDateTime)
		if err != nil {
			log.Fatal("invalid datetime provided:", err)
		}

		err = consumer.Options{
			GroupID:      groupID,
			Offset:       offset,
			Partition:    partition,
			FromDateTime: parsedDT,
		}.New(bootstrapServers, topic).Consume(context.Background(), os.Stdout, decoderFactory(wireForamt))

		log.Fatal("error while consuming messages:", err)
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
	rootCmd.Flags().StringVar(&fromDateTime, "from_datetime", zeroTime, "time from which you want to tail in RFC3339 datetime format")
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
