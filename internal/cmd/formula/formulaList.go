package formula

import (
	"context"
	"fmt"
	"strings"

	"github.com/clintjedwards/basecoat/internal/cmd/cl"
	"github.com/clintjedwards/basecoat/internal/cmd/format"
	"github.com/clintjedwards/basecoat/proto"
	"github.com/clintjedwards/polyfmt/v2"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/metadata"
)

var cmdFormulaList = &cobra.Command{
	Use:   "list",
	Short: "List all formulas",
	Long: `List all formulas.

A short listing of all currently registered formulas.`,
	Example: `$ basecoat formula list`,
	RunE:    formulaList,
}

func init() {
	cmdFormulaList.Flags().StringP("filter", "f", "", "Fuzzy search for formulas")
	CmdFormula.AddCommand(cmdFormulaList)
}

func formulaList(cmd *cobra.Command, _ []string) error {
	cl.State.Fmt.Print("Retrieving formulas", polyfmt.Pretty)

	query, err := cmd.Flags().GetString("filter")
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
	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.ListFormulas(ctx, &proto.ListFormulasRequest{
		Filter: query,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not list formulas: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	if len(resp.Formulas) == 0 {
		cl.State.Fmt.Println("No formulas found")
		cl.State.Fmt.Finish()
		return err
	}

	data := [][]string{}
	for _, formula := range resp.Formulas {
		data = append(data, []string{
			formula.Id,
			formula.Name,
			formula.Number,
			format.UnixMilli(formula.Created, "Never", cl.State.Config.Detail),
		})
	}

	table := formatTable(data, !cl.State.Config.NoColor)

	cl.State.Fmt.Println(table)
	cl.State.Fmt.Finish()

	return nil
}

func formatTable(data [][]string, color bool) string {
	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)

	table.SetHeader([]string{"ID", "Name", "Number", "Created"})
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetHeaderLine(true)
	table.SetBorder(false)
	table.SetAutoFormatHeaders(false)
	table.SetRowSeparator("―")
	table.SetRowLine(false)
	table.SetColumnSeparator("")
	table.SetCenterSeparator("")

	if color {
		table.SetHeaderColor(
			tablewriter.Color(tablewriter.FgBlueColor),
			tablewriter.Color(tablewriter.FgBlueColor),
			tablewriter.Color(tablewriter.FgBlueColor),
			tablewriter.Color(tablewriter.FgBlueColor),
		)
		table.SetColumnColor(
			tablewriter.Color(tablewriter.FgYellowColor),
			tablewriter.Color(0),
			tablewriter.Color(0),
			tablewriter.Color(0),
		)
	}

	table.AppendBulk(data)

	table.Render()
	return tableString.String()
}
