package service

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdServiceCreateAPIToken = &cobra.Command{
	Use:   "create-api-token <id>",
	Short: "Create a new API Token for the given account",
	RunE:  serviceCreateAPIToken,
}

func init() {
	cmdServiceCreateAPIToken.Flags().IntP("duration", "d", 86400, "The duration of the api token")
	CmdService.AddCommand(cmdServiceCreateAPIToken)
}

func serviceCreateAPIToken(_ *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Creating api token", polyfmt.Pretty)

	password1 := cl.State.Fmt.Question("Password: ")

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CreateAPIToken(ctx, &proto.CreateAPITokenRequest{
		Account:  id,
		Password: password1,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not create account: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Created token: %q", resp.Key))
	cl.State.Fmt.Warning("Please remember to save token as it won't be shown again")
	cl.State.Fmt.Finish()
	return nil
}
