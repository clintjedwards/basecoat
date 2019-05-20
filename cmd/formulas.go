package cmd

import (
	"github.com/spf13/cobra"
)

var cmdFormulas = &cobra.Command{
	Use:   "formulas",
	Short: "Controls operations that can be performed on formulas",
}

func init() {
	RootCmd.AddCommand(cmdFormulas)
}
