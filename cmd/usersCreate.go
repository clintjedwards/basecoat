package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/howeyc/gopass"
	"github.com/spf13/cobra"
)

var cmdUsersCreate = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a single user",
	Long:  `Users act as a bucket for all formula and job keys`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runUsersCreateCmd,
}

func runUsersCreateCmd(cmd *cobra.Command, args []string) {
	name := args[0]

	fmt.Printf("Password: ")
	pass, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatalf("failed retrieve password")
	}

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration: %v", err)
	}

	creds, err := credentials.NewClientTLSFromFile(config.TLSCertPath, "")
	if err != nil {
		log.Fatalf("failed to get certificates: %v", err)
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(creds))

	hostPortTuple := strings.Split(config.URL, ":")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostPortTuple[0], hostPortTuple[1]), opts...)
	if err != nil {
		log.Fatalf("could not connect to server: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+config.CommandLine.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.CreateAccount(ctx, &api.CreateAccountRequest{
		Id:       name,
		Password: string(pass),
	})
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}
}

func init() {
	cmdUsers.AddCommand(cmdUsersCreate)
}
