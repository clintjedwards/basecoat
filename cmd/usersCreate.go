package cmd

import (
	"fmt"
	"log"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/storage"
	"github.com/clintjedwards/toolkit/password"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var cmdUsersCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a single user",
	Long:  `Users act as a bucket for all formula and job keys`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runUsersCreateCmd,
}

func runUsersCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatalf("failed retrieve password")
	}

	hash, err := password.HashPassword(pass)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	storage, err := storage.InitStorage()
	if err != nil {
		log.Fatalf("could not connect to storage: %v", err)
	}

	err = storage.CreateUser(name, &api.User{
		Name: name,
		Hash: string(hash),
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
}

func init() {
	cmdUsers.AddCommand(cmdUsersCreate)
}
