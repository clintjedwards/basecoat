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

var cmdColorantUpdate = &cobra.Command{
	Use:     "update <id>",
	Short:   "Update an colorant",
	Long:    `Update an colorant.`,
	Example: `$ basecoat colorant update FyrjxCQ`,
	RunE:    colorantUpdate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	cmdColorantUpdate.Flags().StringP("label", "l", "", "Human readable colorant name")
	cmdColorantUpdate.Flags().StringP("manufacturer", "m", "", "Manufacturer of the colorant")
	CmdColorant.AddCommand(cmdColorantUpdate)
}

func colorantUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Updating colorant", polyfmt.Pretty)

	label, err := cmd.Flags().GetString("label")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	manufacturer, err := cmd.Flags().GetString("manufacturer")
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

	updateColorantRequest := &proto.UpdateColorantRequest{
		Id: id,
	}

	if cmd.Flags().Changed("label") {
		updateColorantRequest.Label = label
	}

	if cmd.Flags().Changed("manufacturer") {
		updateColorantRequest.Manufacturer = manufacturer
	}

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.UpdateColorant(ctx, updateColorantRequest)
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update colorant: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated colorant: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
