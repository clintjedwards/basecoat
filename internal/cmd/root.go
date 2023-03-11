package cmd

import (
	"fmt"
	"strings"

	"github.com/clintjedwards/basecoat/internal/cmd/account"
	"github.com/clintjedwards/basecoat/internal/cmd/base"
	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/internal/cmd/colorant"
	"github.com/clintjedwards/basecoat/internal/cmd/formula"
	"github.com/clintjedwards/basecoat/internal/cmd/service"
	"github.com/spf13/cobra"
)

var appVersion = "0.0.dev_000000"

// RootCmd is the base command for the basecoat cli
var RootCmd = &cobra.Command{
	Use:     "basecoat",
	Short:   "Basecoat is a formula tracking and search tool",
	Version: " ", // We leave this added but empty so that the rootcmd will supply the -v flag
	PersistentPreRun: func(cmd *cobra.Command, _ []string) {
		cl.InitState(cmd)
	},
}

func init() {
	RootCmd.SetVersionTemplate(humanizeVersion(appVersion))
	RootCmd.AddCommand(service.CmdService)
	RootCmd.AddCommand(account.CmdAccount)
	RootCmd.AddCommand(formula.CmdFormula)
	RootCmd.AddCommand(base.CmdBase)
	RootCmd.AddCommand(colorant.CmdColorant)

	RootCmd.PersistentFlags().String("config", "", "configuration file path")
	RootCmd.PersistentFlags().Bool("detail", false, "show extra detail for some commands (ex. Exact time instead of humanized)")
	RootCmd.PersistentFlags().String("format", "", "output format; accepted values are 'pretty', 'json', 'silent'")
	RootCmd.PersistentFlags().Bool("no-color", false, "disable color output")
	RootCmd.PersistentFlags().String("host", "", "specify the URL of the server to communicate to")
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return RootCmd.Execute()
}

func humanizeVersion(version string) string {
	semver, hash, err := strings.Cut(version, "_")
	if !err {
		return ""
	}
	return fmt.Sprintf("basecoat %s [%s]\n", semver, hash)
}
