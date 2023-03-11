package base

import (
	"github.com/spf13/cobra"
)

var CmdBase = &cobra.Command{
	Use:   "base",
	Short: "Manage bases",
	Long:  `Manage bases`,
}
