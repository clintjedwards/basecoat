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

var cmdColorantDelete = &cobra.Command{
	Use:     "delete <id>",
	Short:   "Delete  an colorant",
	Long:    `Delete an colorant.`,
	Example: `$ basecoat colorant delete FyrjxCQ`,
	RunE:    colorantDelete,
	Args:    cobra.ExactArgs(1),
}

func init() {
	CmdColorant.AddCommand(cmdColorantDelete)
}

func colorantDelete(_ *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Deleting colorant", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.DeleteColorant(ctx, &proto.DeleteColorantRequest{
		Id: id,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update colorant: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated colorant: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
