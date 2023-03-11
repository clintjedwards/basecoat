package formula

import (
	"github.com/spf13/cobra"
)

var CmdFormula = &cobra.Command{
	Use:   "formula",
	Short: "Manage formulas",
	Long:  `Manage formulas`,
}
