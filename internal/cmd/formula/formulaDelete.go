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

var cmdFormulaDelete = &cobra.Command{
	Use:     "delete <id>",
	Short:   "Delete  an formula",
	Long:    `Delete an formula.`,
	Example: `$ basecoat formula delete FyrjxCQ`,
	RunE:    formulaDelete,
	Args:    cobra.ExactArgs(1),
}

func init() {
	CmdFormula.AddCommand(cmdFormulaDelete)
}

func formulaDelete(_ *cobra.Command, args []string) error {
	id := args[0]

	cl.State.Fmt.Print("Deleting formula", polyfmt.Pretty)

	conn, err := cl.State.Connect()
	if err != nil {
		cl.State.Fmt.Err(err)
		cl.State.Fmt.Finish()
		return err
	}

	client := proto.NewBasecoatClient(conn)

	md := metadata.Pairs("Authorization", "Bearer "+cl.State.Config.Token)
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	_, err = client.DeleteFormula(ctx, &proto.DeleteFormulaRequest{
		Id: id,
	})
	if err != nil {
		cl.State.Fmt.Err(fmt.Sprintf("could not update formula: %v", err))
		cl.State.Fmt.Finish()
		return err
	}

	cl.State.Fmt.Success(fmt.Sprintf("Updated formula: %q", id))
	cl.State.Fmt.Finish()
	return nil
}
