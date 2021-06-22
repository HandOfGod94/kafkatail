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
	"fmt"
	"github.com/spf13/cobra"
)

var (
	bootstrapServers []string
	topic            string
	groupID          string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kafkatail",
	Short: "Tail kafka logs of any wire format",
	Long: `Print kafka messages from any topic, of any wire format (avro, plaintext, protobuf)
on console`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Hello World")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringSliceVar(&bootstrapServers, "bootstrap_servers", []string{}, "list of kafka `bootstrap_servers` separated by comma")
	rootCmd.Flags().StringVar(&topic, "topic", "", "`topic` whose message you want to tail")
	rootCmd.Flags().StringVar(&groupID, "group_id", "", "kafka consumer `group_id` to be used for subscribing to topic")

	rootCmd.MarkFlagRequired("bootstrap_servers")
	rootCmd.MarkFlagRequired("topic")
}
