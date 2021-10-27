package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	protoFile    string
	includePaths []string
	messageType  string
)

func protoFlags() *pflag.FlagSet {
	protoFlags := pflag.NewFlagSet("proto_flags", pflag.PanicOnError)
	protoFlags.StringVar(&messageType, "message_type", "", "proto message `type` to use for decoding . Required for `wire_format=proto`")
	protoFlags.StringSliceVar(&includePaths, "include_paths", []string{}, "`include_paths` containing dependencies of proto. Required for `wire_format=proto`")
	protoFlags.StringVar(&protoFile, "proto_file", "", "`proto_file` to be used for decoding kafka message. Required for `wire_format=proto`")
	return protoFlags
}

func lookupFlagValue(cmd *cobra.Command, name string) string {
	return cmd.Flags().Lookup(name).Value.String()
}

func markFlagsRequired(cmd *cobra.Command, flags []string) {
	for _, f := range flags {
		cmd.MarkFlagRequired(f)
	}
}
