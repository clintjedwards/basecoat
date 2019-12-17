package cmd

import (
	"fmt"
	"log"

	"github.com/clintjedwards/toolkit/version"
	"github.com/spf13/cobra"
)

var appVersion = "0.0.dev_000000_333333"

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Basecoat",
	Run:   runVersionCmd,
}

func runVersionCmd(cmd *cobra.Command, args []string) {
	info, err := version.Parse(appVersion)
	if err != nil {
		log.Fatalf("could not parse version: %v", err)
	}

	fmt.Printf("Basecoat v%s %s %s\n", info.Semver, info.Epoch, info.Hash)
}

func init() {
	RootCmd.AddCommand(cmdVersion)
}
