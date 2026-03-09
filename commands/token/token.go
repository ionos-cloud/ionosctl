package token

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenCmd() *core.Command {
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "token",
			Short:            "Token Operations",
			Long:             "The sub-commands of `ionosctl token` allow you to generate, list, get, delete Tokens.",
			TraverseChildren: true,
		},
	}
	globalFlags := tokenCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", nil, table.ColsMessage(allTokenCols))
	_ = viper.BindPFlag(core.GetFlagName(tokenCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = tokenCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allTokenCols), cobra.ShellCompDirectiveNoFileComp
	})

	tokenCmd.AddCommand(TokenPostCmd())
	tokenCmd.AddCommand(TokenListCmd())
	tokenCmd.AddCommand(TokenGetCmd())
	tokenCmd.AddCommand(TokenDeleteCmd())
	tokenCmd.AddCommand(TokenParseCmd())

	return core.WithConfigOverride(tokenCmd, []string{"auth"}, constants.DefaultApiURL+"/auth/v1")
}

const contractNumberMessage = "Contract Number: %v"

// Output Printing

var allTokenCols = []table.Column{
	{Name: "TokenId", JSONPath: "id", Default: true},
	{Name: "CreatedDate", JSONPath: "createdDate", Default: true},
	{Name: "ExpirationDate", JSONPath: "expirationDate", Default: true},
	{Name: "Href", JSONPath: "href"},
}
