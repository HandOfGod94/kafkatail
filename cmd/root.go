package cmd

import (
	"github.com/handofgod94/kafkatail/app"
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
		opts := app.AppOptions{
			BootstrapServers: bootstrapServers,
			Topic:            args[0],
			GroupID:          groupID,
			WireForamt:       wireForamt,
			ProtoFile:        protoFile,
			Includes:         includePaths,
			MessageType:      messageType,
			Offset:           offset,
			Partition:        partition,
		}

		opts.Start()
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringSliceVarP(&bootstrapServers, "bootstrap_servers", "b", []string{}, "list of kafka `bootstrap_servers` separated by comma")
	rootCmd.Flags().StringVar(&groupID, "group_id", "", "[Optional] kafka consumer `group_id` to be used for subscribing to topic")
	rootCmd.Flags().StringVar(&protoFile, "proto_file", "", "`proto_file` to be used for decoding kafka message. Required for `wire_format=proto`")
	rootCmd.Flags().StringSliceVar(&includePaths, "include_paths", []string{}, "`include_paths` containing dependencies of proto. Required for `wire_format=proto`")
	rootCmd.Flags().StringVar(&messageType, "message_type", "", "proto message `type` to use for decoding . Required for `wire_format=proto`")
	rootCmd.Flags().Int64Var(&offset, "offset", -1, "kafka offset to start consuming from. Possible Values: -1=latest, -2=earliest, n=nth offset")
	rootCmd.Flags().IntVar(&partition, "partition", 0, "kafka partition to consume from")
	rootCmd.Flags().Var(enumflag.New(&wireForamt, "wire_format", wire.FormatIDs, enumflag.EnumCaseSensitive),
		"wire_format",
		"Wire format of messages in topic",
	)

	rootCmd.Flags().Lookup("wire_format").NoOptDefVal = "plaintext"

	rootCmd.MarkFlagRequired("bootstrap_servers")
	rootCmd.MarkFlagRequired("topic")
}
