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

var cmdBaseUpdate = &cobra.Command{
	Use:     "update <id>",
	Short:   "Update an base",
	Long:    `Update an base.`,
	Example: `$ basecoat base update FyrjxCQ`,
	RunE:    baseUpdate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	cmdBaseUpdate.Flags().StringP("label", "l", "", "Human readable base name")
	cmdBaseUpdate.Flags().StringP("manufacturer", "m", "", "Manufacturer of the base")
	CmdBase.AddCommand(cmdBaseUpdate)
}

func baseUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Updating base", polyfmt.Pretty)

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

	updateBaseRequest := &proto.UpdateBaseRequest{
		Id: id,
	}

	if cmd.Flags().Changed("label") {
		updateBaseRequest.Label = label
	}

	if cmd.Flags().Changed("manufacturer") {
		updateBaseRequest.Manufacturer = manufacturer
	}

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.UpdateBase(ctx, updateBaseRequest)
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update base: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated base: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
