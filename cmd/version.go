package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "v0.0.dev"

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Basecoat",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Basecoat %s\n", appVersion)
	},
}

func init() {
	RootCmd.AddCommand(cmdVersion)
}
