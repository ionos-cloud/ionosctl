package token

import (
	"context"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenScopesCmd() *core.Command {
	scopesCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "scope",
			Aliases: []string{"s", "scopes"},
			Short:   "Registry Tokens' Scopes Operations",
			Long: "Manage token scopes for a registry. You can list, add, and remove scopes from a token. " +
				"Scopes are used to grant access to a registry. " +
				"Each token can have multiple scopes. ",
			TraverseChildren: true,
		},
	}

	scopesCmd.AddCommand(TokenScopesListCmd())
	return scopesCmd
}

func TokenScopesListCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "token",
			Resource:   "scope",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "Get a token scopes",
			LongDesc:   "Use this command to get a token scopes of a container registry.",
			Example:    "ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID]",
			PreCmdRun:  PreCmdTokenScopesList,
			CmdRun:     CmdGetTokenScopesList,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("registry-id", "r", "", "Registry ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag("token-id", "t", "", "Token ID")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"token-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all scopes of all tokens of a registry.")

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func CmdGetTokenScopesList(c *core.CommandConfig) error {
	reg_id := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
	token_id := viper.GetString(core.GetFlagName(c.NS, "token-id"))
	token, _, err := c.ContainerRegistryServices.Token().Get(token_id, reg_id)
	if err != nil {
		return err
	}

	return c.Printer.Print(getTokenScopePrint(nil, c, &token, true))
}

func PreCmdTokenScopesList(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", "registry-id")
	if err != nil {
		return err
	}
	return nil
}

var (
	defaultScopeCols = []string{"DisplayName", "Type", "Actions"}
	allScopeCols     = []string{"TokenId", "DisplayName", "Type", "Actions"}
)

func getTokenScopePrint(
	resp *ionoscloud.APIResponse, c *core.CommandConfig, response *ionoscloud.TokenResponse,
	first bool,
) printer.Result {
	r := printer.Result{}

	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(
				core.GetFlagName(
					c.NS, constants.ArgWaitForRequest,
				),
			) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if response != nil {
			r.OutputJSON = response
			r.Columns = printer.GetHeadersListAll(allColsScopes, defaultScopeCols, "TokenId", allScopeCols, viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))) // headers
			r.KeyValue = getTokensScopeRows(response) // map header -> rows
		}
	}
	return r
}

type TokenScopePrint struct {
	DisplayName string `json:"DisplayName,omitempty"`
	TokenId     string `json:"TokenId,omitempty"`
	Type        string `json:"Type,omitempty"`
	Actions     string `json:"Actions,omitempty"`
}

func getTokensScopeRows(token *ionoscloud.TokenResponse) []map[string]interface{} {
	scopes := token.Properties.GetScopes()
	out := make([]map[string]interface{}, 0, len(*scopes))
	for _, scope := range *scopes {
		var tokenScopePrint TokenScopePrint
		if nameOk, ok := scope.GetNameOk(); ok && nameOk != nil {
			tokenScopePrint.DisplayName = *nameOk
		}
		if typeOk, ok := scope.GetTypeOk(); ok && typeOk != nil {
			tokenScopePrint.Type = *typeOk
		}
		if actionsOk, ok := scope.GetActionsOk(); ok && actionsOk != nil {
			for _, action := range *actionsOk {
				tokenScopePrint.Actions += string(action) + ", "
			}
			tokenScopePrint.Actions = tokenScopePrint.Actions[:len(tokenScopePrint.Actions)-2]
		}
		tokenScopePrint.TokenId = *token.Id
		o := structs.Map(tokenScopePrint)
		out = append(out, o)
	}
	return out
}

var allColsScopes = structs.Names(TokenScopePrint{})
