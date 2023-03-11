package colorant

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

var cmdColorantList = &cobra.Command{
	Use:   "list",
	Short: "List all colorants",
	Long: `List all colorants.

A short listing of all currently registered colorants.`,
	Example: `$ basecoat colorant list`,
	RunE:    colorantList,
}

func init() {
	CmdColorant.AddCommand(cmdColorantList)
}

func colorantList(_ *cobra.Command, _ []string) error {
	cl.State.Fmt.Print("Retrieving colorants", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)
	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.ListColorants(ctx, &proto.ListColorantsRequest{})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not list colorants: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	if len(resp.Colorants) == 0 {
		cl.State.Fmt.Println("No colorants found")
		cl.State.Fmt.Finish()
		return err
	}

	data := [][]string{}
	for _, colorant := range resp.Colorants {
		data = append(data, []string{
			colorant.Id,
			colorant.Manufacturer,
			colorant.Label,
			format.UnixMilli(colorant.Created, "Never", cl.State.Config.Detail),
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

	table.SetHeader([]string{"ID", "Manufacturer", "Label", "Created"})
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
