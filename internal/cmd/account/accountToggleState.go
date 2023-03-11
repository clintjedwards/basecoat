package account

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/internal/cmd/format"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdAccountToggleState = &cobra.Command{
	Use:     "toggle-state <id>",
	Short:   "Toggle the state of an account",
	Long:    `Toggle the state an account.`,
	Example: `$ basecoat account toggle-state FyrjxCQ`,
	RunE:    accountToggleState,
	Args:    cobra.ExactArgs(1),
}

func init() {
	CmdAccount.AddCommand(cmdAccountToggleState)
}

func accountToggleState(_ *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Updating account", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	state, err := client.ToggleAccountState(ctx, &proto.ToggleAccountStateRequest{Id: id})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not toggle account state: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Toggled account state: %q to %s", id,
		format.ColorizeAccountState(format.NormalizeEnumValue(state.State.String(), "Unknown"))))
	cl.State.Fmt.Finish()
	return nil
}
