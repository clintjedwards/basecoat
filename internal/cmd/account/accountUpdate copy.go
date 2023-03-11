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

var cmdAccountUpdate = &cobra.Command{
	Use:     "update <id>",
	Short:   "Update an account",
	Long:    `Update an account.`,
	Example: `$ basecoat account update FyrjxCQ`,
	RunE:    accountUpdate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	cmdAccountUpdate.Flags().StringP("name", "n", "", "Human readable account name")
	cmdAccountUpdate.Flags().BoolP("password", "p", false, "Reset the account password")
	CmdAccount.AddCommand(cmdAccountUpdate)
}

func accountUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Updating account", polyfmt.Pretty)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	password, err := cmd.Flags().GetBool("password")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	updateAccountRequest := &proto.UpdateAccountRequest{
		Id: id,
	}

	if cmd.Flags().Changed("name") {
		updateAccountRequest.Name = name
	}

	if password {
		password1 := cl.State.Fmt.Question("Password: ")
		password2 := cl.State.Fmt.Question("Retype Password: ")

		if password1 != password2 {
			cl.State.Fmt.Err("Passwords do not match")
			cl.State.Fmt.Finish()
			return fmt.Errorf("passwords do not match")
		}

		updateAccountRequest.Password = password1
	}

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.UpdateAccount(ctx, updateAccountRequest)
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update account: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated account: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
