package account

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

var cmdAccountList = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `List all accounts.

A short listing of all currently registered accounts.`,
	Example: `$ basecoat account list`,
	RunE:    accountList,
}

func init() {
	CmdAccount.AddCommand(cmdAccountList)
}

func accountList(_ *cobra.Command, _ []string) error {
	cl.State.Fmt.Print("Retrieving accounts", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)
	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := client.ListAccounts(ctx, &proto.ListAccountsRequest{})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not list accounts: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	if len(resp.Accounts) == 0 {
		cl.State.Fmt.Println("No accounts found")
		cl.State.Fmt.Finish()
		return err
	}

	data := [][]string{}
	for _, account := range resp.Accounts {
		data = append(data, []string{
			account.Id,
			account.Name,
			format.ColorizeAccountState(format.NormalizeEnumValue(account.State.String(), "Unknown")),
			format.UnixMilli(account.Created, "Never", cl.State.Config.Detail),
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

	table.SetHeader([]string{"ID", "Name", "State", "Created"})
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
