package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

const appVersion = "0.1.0"

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print current version of kafkatail",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf(appVersion)
	},
}
