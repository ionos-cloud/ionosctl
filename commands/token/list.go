package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/viper"
)

func TokenListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Tokens",
		LongDesc:   "Use this command to retrieve a complete list of Tokens under your account, to list the details of all tokens, generated by the user.",
		Example:    listTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     runTokenList,
		InitClient: true,
	})
	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are listed")

	return cmd
}

func runTokenList(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}

	tokens, _, err := c.AuthV1Services.Tokens().List(viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("tokens", jsonpaths.AuthToken, tokens.Tokens,
		tabheaders.GetHeaders(allTokenCols, defaultTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}
