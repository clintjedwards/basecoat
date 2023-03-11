package colorant

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdColorantCreate = &cobra.Command{
	Use:     "create <label> <manufacturer>",
	Short:   "Create a new colorant",
	Long:    `Create a new colorant.`,
	Example: `$ basecoat colorant create "Off White" "Benjamin Moore"`,
	RunE:    colorantCreate,
	Args:    cobra.ExactArgs(2),
}

func init() {
	CmdColorant.AddCommand(cmdColorantCreate)
}

func colorantCreate(_ *cobra.Command, args []string) error {
	label := args[0]
	manufacturer := args[1]

	cl.State.Fmt.Print("Creating colorant", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CreateColorant(ctx, &proto.CreateColorantRequest{
		Label:        label,
		Manufacturer: manufacturer,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not create colorant: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Created colorant: [%s] %q", resp.Colorant.Id, resp.Colorant.Manufacturer+" - "+resp.Colorant.Label))

	cl.State.Fmt.Finish()
	return nil
}
