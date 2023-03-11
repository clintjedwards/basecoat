package service

import (
	"github.com/spf13/cobra"
)

var CmdService = &cobra.Command{
	Use:   "service",
	Short: "Manages service related commands for Basecoat.",
	Long: `Manages service related commands for the Basecoat Service/API.

These commands help with managing and running the Basecoat service.`,
}
