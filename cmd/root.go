package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

//RootCmd is the base command for the basecoat cli
var RootCmd = &cobra.Command{
	Use:   "basecoat",
	Short: "Basecoat is a formula tracking and search tool",
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
