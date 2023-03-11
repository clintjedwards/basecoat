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

var cmdBaseDelete = &cobra.Command{
	Use:     "delete <id>",
	Short:   "Delete  an base",
	Long:    `Delete an base.`,
	Example: `$ basecoat base delete FyrjxCQ`,
	RunE:    baseDelete,
	Args:    cobra.ExactArgs(1),
}

func init() {
	CmdBase.AddCommand(cmdBaseDelete)
}

func baseDelete(_ *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Deleting base", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.DeleteBase(ctx, &proto.DeleteBaseRequest{
		Id: id,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update base: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated base: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
