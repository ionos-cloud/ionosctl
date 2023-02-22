package token

import (
	"context"
	"fmt"
	"strconv"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
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
	scopesCmd.AddCommand(TokenScopesAddCmd())
	scopesCmd.AddCommand(TokenScopesDeleteCmd())
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

	// TODO: add --all flag
	// cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "List all scopes of all tokens of a registry.")

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

func TokenScopesAddCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "token",
			Resource:  "scope",
			Verb:      "add",
			Aliases:   []string{"a", "ad"},
			ShortDesc: "Add scopes to a token",
			LongDesc:  "Use this command to add scopes to a token of a container registry.",
			Example: "ionosctl container-registry token scope list --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]" +
				" --actions [SCOPE-ACTIONS], --type [SCOPE-TYPE]",
			PreCmdRun:  PreCmdTokenScopesAdd,
			CmdRun:     CmdTokenScopesAdd,
			InitClient: true,
		},
	)

	cmd.AddStringFlag("registry-id", "r", "", "Registry ID", core.RequiredFlagOption())
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

	cmd.AddStringFlag("name", "n", "", "Scope name", core.RequiredFlagOption())
	cmd.AddStringFlag("type", "y", "", "Scope type", core.RequiredFlagOption())
	cmd.AddStringSliceFlag("actions", "a", []string{}, "Scope actions", core.RequiredFlagOption())

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	return cmd
}

func PreCmdTokenScopesAdd(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, "token-id", "registry-id", "name", "actions", "type")
	if err != nil {
		return err
	}
	return nil
}

func CmdTokenScopesAdd(c *core.CommandConfig) error {
	var scope sdkgo.Scope
	var err error

	reg_id, err := c.Command.Command.Flags().GetString("registry-id")
	if err != nil {
		return err
	}

	token_id, err := c.Command.Command.Flags().GetString("token-id")
	if err != nil {
		return err
	}

	name, err := c.Command.Command.Flags().GetString("name")
	if err != nil {
		return err
	}

	actions, err := c.Command.Command.Flags().GetStringSlice("actions")
	if err != nil {
		return err
	}

	scope_type, err := c.Command.Command.Flags().GetString("type")
	if err != nil {
		return err
	}

	scope.SetName(name)
	scope.SetActions(actions)
	scope.SetType(scope_type)

	token, _, err := c.ContainerRegistryServices.Token().Get(token_id, reg_id)
	if err != nil {
		return err
	}

	updateToken := sdkgo.NewPatchTokenInput()
	updateToken.SetExpiryDate(*token.Properties.GetExpiryDate())
	updateToken.SetStatus(*token.Properties.GetStatus())
	scopes := *token.Properties.GetScopes()
	scopes = append(scopes, scope)
	updateToken.SetScopes(scopes)

	tokenUp, _, err := c.ContainerRegistryServices.Token().Patch(token_id, *updateToken, reg_id)
	if err != nil {
		return err
	}

	return c.Printer.Print(getTokenScopePrint(nil, c, &tokenUp, true))
}

func TokenScopesDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "token",
			Resource:  "scope",
			Verb:      "delete",
			Aliases:   []string{"d", "rm", "remove"},
			ShortDesc: "Delete a token scope",
			LongDesc: "Use this command to delete a token scope of a container registry. If a name is provided, the first scope with that" +
				" name will be deleted. It is possible to delete all scopes by providing the --all flag.",
			Example:    "ionosctl container-registry token scope delete --registry-id [REGISTRY-ID], --token-id [TOKEN-ID] --name [SCOPE-NAME]",
			PreCmdRun:  PreCmdTokenScopesDelete,
			CmdRun:     CmdGetTokenScopesDelete,
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

	cmd.AddIntFlag("id", "n", -1, "Scope id")
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

func CmdGetTokenScopesDelete(c *core.CommandConfig) error {
	reg_id := viper.GetString(core.GetFlagName(c.NS, "registry-id"))
	token_id := viper.GetString(core.GetFlagName(c.NS, "token-id"))
	token, _, err := c.ContainerRegistryServices.Token().Get(token_id, reg_id)
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		updateToken := sdkgo.NewPutTokenInputWithDefaults()
		updateProp := sdkgo.NewPostTokenPropertiesWithDefaults()

		updateProp.SetExpiryDate(*token.Properties.GetExpiryDate())
		updateProp.SetStatus(*token.Properties.GetStatus())
		updateProp.SetName(*token.Properties.GetName())
		updateToken.SetProperties(*updateProp)
		msg := fmt.Sprintf("delete all scopes from Token: %s", *token.Id)
		if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
			return err
		}
		_, err = c.ContainerRegistryServices.Token().Delete(token_id, reg_id)
		if err != nil {
			return err
		}
		_, _, err = c.ContainerRegistryServices.Token().Put(token_id, *updateToken, reg_id)
		if err != nil {
			return err
		}

		return nil
	}

	id, err := c.Command.Command.Flags().GetInt("id")
	if err != nil {
		return err
	}
	id--

	updateToken := sdkgo.NewPutTokenInputWithDefaults()
	updateProp := sdkgo.NewPostTokenPropertiesWithDefaults()

	scopes := *token.Properties.GetScopes()
	scopes = append(scopes[:id], scopes[id+1:]...)
	updateProp.SetExpiryDate(*token.Properties.GetExpiryDate())
	updateProp.SetStatus(*token.Properties.GetStatus())
	updateProp.SetName(*token.Properties.GetName())
	updateProp.SetScopes(scopes)
	updateToken.SetProperties(*updateProp)

	msg := fmt.Sprintf("delete scope %d from Token: %s", id + 1, *token.Id)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, msg); err != nil {
		return err
	}
	_, err = c.ContainerRegistryServices.Token().Delete(token_id, reg_id)
	if err != nil {
		return err
	}
	_, _, err = c.ContainerRegistryServices.Token().Put(token_id, *updateToken, reg_id)
	if err != nil {
		return err
	}

	return nil
}

func PreCmdTokenScopesDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{"registry-id", "token-id", "id"},
		[]string{"registry-id", "token-id", constants.ArgAll},
	)
}

var (
	defaultScopeCols = []string{"ScopeId", "DisplayName", "Type", "Actions"}
	allScopeCols     = []string{"ScopeId", "TokenId", "DisplayName", "Type", "Actions"}
)

func getTokenScopePrint(
	resp *sdkgo.APIResponse, c *core.CommandConfig, response *sdkgo.TokenResponse,
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
			r.Columns = printer.GetHeadersListAll(allColsScopes, defaultScopeCols, "TokenId", nil, viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll))) // headers
			r.KeyValue = getTokensScopeRows(response)                                                                                                       // map header -> rows
		}
	}
	return r
}

type TokenScopePrint struct {
	ScopeId     string `json:"ScopeId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
	TokenId     string `json:"TokenId,omitempty"`
	Type        string `json:"Type,omitempty"`
	Actions     string `json:"Actions,omitempty"`
}

func getTokensScopeRows(token *sdkgo.TokenResponse) []map[string]interface{} {
	scopes := token.Properties.GetScopes()
	out := make([]map[string]interface{}, 0, len(*scopes))
	for i, scope := range *scopes {
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
		tokenScopePrint.ScopeId = strconv.Itoa(i + 1)
		o := structs.Map(tokenScopePrint)
		out = append(out, o)
	}
	return out
}

var allColsScopes = structs.Names(TokenScopePrint{})
