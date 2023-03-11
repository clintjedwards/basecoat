package formula

import (
	"context"
	"fmt"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdFormulaUpdate = &cobra.Command{
	Use:     "update <id>",
	Short:   "Update an formula",
	Long:    `Update an formula.`,
	Example: `$ basecoat formula update FyrjxCQ`,
	RunE:    formulaUpdate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	cmdFormulaUpdate.Flags().StringP("name", "n", "", "Human readable formula name")
	cmdFormulaUpdate.Flags().StringP("number", "u", "", "Specialized formula number")
	cmdFormulaUpdate.Flags().StringP("notes", "o", "", "Notes about a specific formula")
	CmdFormula.AddCommand(cmdFormulaUpdate)
}

func formulaUpdate(cmd *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Updating formula", polyfmt.Pretty)

	name, err := cmd.Flags().GetString("name")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	number, err := cmd.Flags().GetString("number")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	notes, err := cmd.Flags().GetString("number")
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

	updateFormulaRequest := &proto.UpdateFormulaRequest{
		Id: id,
	}

	if cmd.Flags().Changed("name") {
		updateFormulaRequest.Name = name
	}

	if cmd.Flags().Changed("number") {
		updateFormulaRequest.Number = number
	}

	if cmd.Flags().Changed("notes") {
		updateFormulaRequest.Notes = notes
	}

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.UpdateFormula(ctx, updateFormulaRequest)
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update formula: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated formula: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
