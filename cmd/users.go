package cmd

import (
	"github.com/spf13/cobra"
)

var cmdUsers = &cobra.Command{
	Use:   "users",
	Short: "Controls operations that can be performed on users. Admin route.",
}

func init() {
	RootCmd.AddCommand(cmdUsers)
}
