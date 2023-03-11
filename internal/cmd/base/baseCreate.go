package base

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdBaseCreate = &cobra.Command{
	Use:     "create <label> <manufacturer>",
	Short:   "Create a new base",
	Long:    `Create a new base.`,
	Example: `$ basecoat base create "Off White" "Benjamin Moore"`,
	RunE:    baseCreate,
	Args:    cobra.ExactArgs(2),
}

func init() {
	CmdBase.AddCommand(cmdBaseCreate)
}

func baseCreate(_ *cobra.Command, args []string) error {
	label := args[0]
	manufacturer := args[1]

	cl.State.Fmt.Print("Creating base", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CreateBase(ctx, &proto.CreateBaseRequest{
		Label:        label,
		Manufacturer: manufacturer,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not create base: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Created base: [%s] %q", resp.Base.Id, resp.Base.Manufacturer+" - "+resp.Base.Label))

	cl.State.Fmt.Finish()
	return nil
}
