package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/basecoat/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	"github.com/clintjedwards/basecoat/config"
	"github.com/spf13/cobra"
)

var cmdFormulasList = &cobra.Command{
	Use:   "list",
	Short: "List all formulas",
	Run:   runFormulasListCmd,
}

func runFormulasListCmd(cmd *cobra.Command, args []string) {

	config, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed to read configuration")
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
		log.Fatalf("could not connect to basecoat: %v", err)
	}
	defer conn.Close()

	basecoatClient := api.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+config.CommandLine.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	formulas, err := basecoatClient.ListFormulas(ctx, &api.ListFormulasRequest{})
	if err != nil {
		log.Fatalf("could not get list of formula: %v", err)
	}

	for key, value := range formulas.Formulas {
		fmt.Printf("%s :: %s\n", key, value.String())
	}
}

func init() {
	cmdFormulas.AddCommand(cmdFormulasList)
}
