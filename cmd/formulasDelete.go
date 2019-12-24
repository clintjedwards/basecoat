package cmd

import (
	"fmt"
	"log"
	"strings"

	"context"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var cmdFormulasDelete = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a single formula",
	Args:  cobra.MinimumNArgs(1),
	Run:   runFormulasDeleteCmd,
}

func runFormulasDeleteCmd(cmd *cobra.Command, args []string) {
	id := args[0]

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration")
	}

	hostPortTuple := strings.Split(config.URL, ":")

	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", hostPortTuple[0], hostPortTuple[1]))
	if err != nil {
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+config.CommandLine.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = basecoatClient.DeleteFormula(ctx, &api.DeleteFormulaRequest{
		Id: id,
	})
	if err != nil {
		log.Fatalf("could not delete formula: %v", err)
	}
}

func init() {
	cmdFormulas.AddCommand(cmdFormulasDelete)
}
