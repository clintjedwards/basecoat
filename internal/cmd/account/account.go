package account

import (
	"github.com/spf13/cobra"
)

var CmdAccount = &cobra.Command{
	Use:   "account",
	Short: "Manage accounts",
	Long:  `Manage accounts`,
}
