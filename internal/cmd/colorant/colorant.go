package colorant

import (
	"github.com/spf13/cobra"
)

var CmdColorant = &cobra.Command{
	Use:   "colorant",
	Short: "Manage colorants",
	Long:  `Manage colorants`,
}
