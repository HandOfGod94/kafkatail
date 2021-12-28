package cmd

import (
	"time"

	"github.com/handofgod94/kafkatail/wire"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag"
)

var (
	bootstrapServers []string
	groupID          string
	wireForamt       wire.Format
	offset           int64
	partition        int
	fromDateTime     string
	protoFile        string
	includePaths     []string
	messageType      string
)

const appVersion = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kafkatail [flags] topic",
	Short:   "Tail kafka logs of any wire format",
	Long:    `Tail kafka messages from any topic, of any wire format on console (plaintext, protobuf)`,
	Version: appVersion,
	Args:    cobra.MinimumNArgs(1),
	Example: `
	# tail messages from a topic
	kafkatail --bootstrap_servers=localhost:9093 kafkatail-test

	# tail proto messages from a topic
	kafkatail --bootstrap_servers=localhost:9093 --wire_format=proto --proto_file=starwars.proto --include_paths="../testdata" --message_type=Human kafkatail-test-proto

	# tail messages from an offset. Default: -1 (latest). For earliest, use offset=-2
	kafkatail --bootstrap_servers=localhost:9093 --offset=12 kafkatail-test-base

	# tail messages from specific time
	kafkatail --bootstrap_servers=localhost:9093 --from_datetime=2021-06-28T15:04:23Z kafkatail-test-base

	# tail messages from specific partition. Default: 0
	kafkatail --bootstrap_servers=localhost:9093 --partition=5 kafkatail-test-base

	# tail from multiple partitions, using group_id
	kafkatail --bootstrap_servers=localhost:9093 --group_id=myfoo kafka-consume-gorup-id-int-test
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		formatFlag := cmd.Flags().Lookup("wire_format")
		if formatFlag.Value.String() == "proto" {
			cmd.MarkFlagRequired("proto_file")
			cmd.MarkFlagRequired("include_paths")
			cmd.MarkFlagRequired("message_type")
		}
	},
	RunE: runKafkaTail,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	zeroTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Format(time.RFC3339)

	rootCmd.Flags().StringSliceVarP(&bootstrapServers, "bootstrap_servers", "b", []string{}, "list of kafka `bootstrap_servers` separated by comma")
	rootCmd.Flags().StringVar(&groupID, "group_id", "", "[Optional] kafka consumer `group_id` to be used for subscribing to topic")
	rootCmd.Flags().Int64Var(&offset, "offset", -1, "kafka offset to start consuming from. Possible Values: -1=latest, -2=earliest, n=nth offset")
	rootCmd.Flags().IntVar(&partition, "partition", 0, "kafka partition to consume from")
	rootCmd.Flags().StringVar(&fromDateTime, "from_datetime", zeroTime, "tail from specific past datetime in RFC3339 format")
	rootCmd.Flags().StringVar(&messageType, "message_type", "", "proto message `type` to use for decoding . Required for `wire_format=proto`")
	rootCmd.Flags().StringSliceVar(&includePaths, "include_paths", []string{}, "`include_paths` containing dependencies of proto. Required for `wire_format=proto`")
	rootCmd.Flags().StringVar(&protoFile, "proto_file", "", "`proto_file` to be used for decoding kafka message. Required for `wire_format=proto`")
	rootCmd.Flags().Var(enumflag.New(&wireForamt, "wire_format", wire.FormatIDs, enumflag.EnumCaseSensitive),
		"wire_format",
		"Wire format of messages in topic",
	)

	rootCmd.Flags().Lookup("wire_format").NoOptDefVal = "plaintext"

	rootCmd.MarkFlagRequired("bootstrap_servers")
}
