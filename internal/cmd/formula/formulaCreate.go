package formula

import (
	"context"
	"fmt"
	"strings"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"

	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdFormulaCreate = &cobra.Command{
	Use:     "create <name>",
	Short:   "Create a new formula",
	Long:    `Create a new formula.`,
	Example: `$ basecoat formula create "Formula Name"`,
	RunE:    formulaCreate,
	Args:    cobra.ExactArgs(1),
}

func init() {
	cmdFormulaCreate.Flags().StringP("number", "u", "", "Special formula number")
	cmdFormulaCreate.Flags().StringP("notes", "o", "", "Notes about the formula")
	cmdFormulaCreate.Flags().StringArrayP("base", "b", []string{}, "Bases to add to the formula. The syntax is <id>:<amount>.")
	cmdFormulaCreate.Flags().StringArrayP("colorant", "c", []string{}, "Colorants to add to the formula. The syntax is <id>:<amount>.")
	CmdFormula.AddCommand(cmdFormulaCreate)
}

func formulaCreate(cmd *cobra.Command, args []string) error {
	name := args[0]

	number, err := cmd.Flags().GetString("number")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	notes, err := cmd.Flags().GetString("notes")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	basesRaw, err := cmd.Flags().GetStringArray("base")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	colorantsRaw, err := cmd.Flags().GetStringArray("colorant")
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	bases := map[string]string{}

	for _, base := range basesRaw {
		id, amount, found := strings.Cut(base, ":")
		if !found {
			cl.State.Fmt.Err(fmt.Sprintf("could not parse base; must be in format <id>:<amount> : %v", err))
			cl.State.Fmt.Finish()
			return err
		}

		bases[id] = amount
	}

	colorants := map[string]string{}

	for _, colorant := range colorantsRaw {
		id, amount, found := strings.Cut(colorant, ":")
		if !found {
			cl.State.Fmt.Err(fmt.Sprintf("could not parse colorant; must be in format <id>:<amount> : %v", err))
			cl.State.Fmt.Finish()
			return err
		}

		colorants[id] = amount
	}

	cl.State.Fmt.Print("Creating formula", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	resp, err := client.CreateFormula(ctx, &proto.CreateFormulaRequest{
		Name:   name,
		Number: number,
		Notes:  notes,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not create formula: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Created formula: [%s] %q", resp.Formula.Id, resp.Formula.Name))

	for base, amount := range bases {
		_, err := client.AssociateBaseWithFormula(ctx, &proto.AssociateBaseWithFormulaRequest{
			Formula: resp.Formula.Id,
			Base:    base,
			Amount:  amount,
		})
		if err != nil {
			cl.State.Fmt.Err(fmt.Sprintf("could not link base to formula: %v", err))
			cl.State.Fmt.Finish()
			return err
		}
		cl.State.Fmt.Success(fmt.Sprintf("Attached Base: %s of %s", amount, base))
	}

	for colorant, amount := range colorants {
		_, err := client.AssociateColorantWithFormula(ctx, &proto.AssociateColorantWithFormulaRequest{
			Formula:  resp.Formula.Id,
			Colorant: colorant,
			Amount:   amount,
		})
		if err != nil {
			cl.State.Fmt.Err(fmt.Sprintf("could not link colorant to formula: %v", err))
			cl.State.Fmt.Finish()
			return err
		}

		cl.State.Fmt.Success(fmt.Sprintf("Attached Colorant: %s of %s", amount, colorant))
	}

	cl.State.Fmt.Finish()
	return nil
}
