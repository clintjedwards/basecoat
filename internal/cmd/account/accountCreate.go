package account

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdAccountCreate = &cobra.Command{
	Use:     "create <name>",
	Short:   "Create a new account",
	Long:    `Create a new account.`,
	Example: `$ basecoat account create "Account Name"`,
	RunE:    accountCreate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	CmdAccount.AddCommand(cmdAccountCreate)
}

func accountCreate(_ *cobra.Command, args []string) error {
	name := args[0]

	cl.State.Fmt.Print("Creating account", polyfmt.Pretty)

	password1 := cl.State.Fmt.Question("Password: ")
	password2 := cl.State.Fmt.Question("Retype Password: ")

	if password1 != password2 {
		cl.State.Fmt.Err("Passwords do not match")
		cl.State.Fmt.Finish()
		return fmt.Errorf("passwords do not match")
	}

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CreateAccount(ctx, &proto.CreateAccountRequest{
		Name:     name,
		Password: password1,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not create account: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Created account: [%s] %q", resp.Account.Id, resp.Account.Name))
	cl.State.Fmt.Finish()
	return nil
}
