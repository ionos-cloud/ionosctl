package token

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/token/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenGetCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Token",
		LongDesc:   "Use this command to retrieve details about a Token by using its ID.\n\nRequired values to run command:\n\n* Token Id",
		Example:    getTokenExample,
		PreCmdRun:  preRunTokenId,
		CmdRun:     runTokenGet,
		InitClient: true,
	})

	cmd.AddUUIDFlag(constants.FlagTokenId, constants.FlagIdShort, "", constants.DescTokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", constants.DescToken, core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagContract, "", 0, "Users with multiple contracts must provide the contract number, for which the token information is displayed")

	return cmd
}

func preRunTokenId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagTokenId}, []string{constants.ArgToken})
}

func runTokenGet(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagTokenId)) {
		return runTokenGetById(c)
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return runTokenGetByToken(c)
	}

	return nil
}

func runTokenGetById(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Getting Token with ID: %v...", viper.GetString(core.GetFlagName(c.NS, constants.FlagTokenId))))

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagContract)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract))))
	}

	req := client.Must().AuthClient.TokensApi.TokensFindById(context.Background(), viper.GetString(core.GetFlagName(c.NS, constants.FlagTokenId)))
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract)) != 0 {
		req = req.XContractNumber(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract)))
	}
	sdkToken, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.AuthToken, sdkToken,
		tabheaders.GetHeaders(allTokenCols, defaultTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func runTokenGetByToken(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, constants.ArgToken))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Token content is: %s", token))

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Token with ID: %v...", tokenId))
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagContract)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract))))
	}

	req := client.Must().AuthClient.TokensApi.TokensFindById(context.Background(), fmt.Sprintf("%v", tokenId))
	if viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract)) != 0 {
		req = req.XContractNumber(viper.GetInt32(core.GetFlagName(c.NS, constants.FlagContract)))
	}
	tokenObj, _, err := req.Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.AuthToken, tokenObj,
		tabheaders.GetHeaders(allTokenCols, defaultTokenCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}
