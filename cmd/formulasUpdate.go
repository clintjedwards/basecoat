package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/clintjedwards/basecoat/api"
	"github.com/clintjedwards/basecoat/config"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var cmdFormulasUpdate = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a single formula",
	Args:  cobra.MinimumNArgs(1),
	Run:   runFormulasUpdateCmd,
}

func runFormulasUpdateCmd(cmd *cobra.Command, args []string) {
	id := args[0]
	name, _ := cmd.Flags().GetString("name")
	number, _ := cmd.Flags().GetString("number")
	notes, _ := cmd.Flags().GetString("notes")
	jobsraw, _ := cmd.Flags().GetStringSlice("jobs")

	var jobs []string
	for _, id := range jobsraw {
		jobs = append(jobs, id)
	}

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

	_, err = basecoatClient.UpdateFormula(ctx, &api.UpdateFormulaRequest{
		Id:     id,
		Name:   name,
		Number: number,
		Notes:  notes,
		Jobs:   jobs,
	})
	if err != nil {
		log.Fatalf("could not update formula: %v", err)
	}
}

func init() {
	var name string
	cmdFormulasUpdate.Flags().StringVarP(&name, "name", "n", "", "formula human readable name")

	var number string
	cmdFormulasUpdate.Flags().StringVarP(&number, "number", "u", "", "formula number used internally")

	var notes string
	cmdFormulasUpdate.Flags().StringVarP(&notes, "notes", "o", "", "any additional notes for the formula")

	var jobs []string
	cmdFormulasUpdate.Flags().StringSliceVarP(&jobs, "jobs", "j", []string{}, "comma separated list of jobs by id in which this formulas has been used")

	cmdFormulas.AddCommand(cmdFormulasUpdate)
}
