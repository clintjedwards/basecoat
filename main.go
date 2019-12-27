//go:generate go run frontend/generate.go

// Basecoat is an internal formula management system
package main

import (
	"github.com/clintjedwards/basecoat/cmd"
)

func main() {
	cmd.RootCmd.Execute()
}
