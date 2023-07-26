package authv1

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/v6/commands/auth-v1/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/confirm"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	authv1 "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "token",
		Resource:  "token",
		Verb:      "delete",
		Aliases:   []string{"d"},
		ShortDesc: "Delete one or multiple Tokens",
		LongDesc: `Use this command to delete a specified Token by token Id or multiple Tokens (based on a criteria: CURRENT, EXPIRED, ALL) from your account. With parameter values ALL and EXPIRED, 'Basic Authentication' or 'Token Authentication' tokens with valid credentials must be encapsulated in the header. With value CURRENT, only the 'Token Authentication' with valid credentials is required.

Required values to run command:

* Token Id/Token/CURRENT/EXPIRED/ALL`,
		Example:    deleteTokenExample,
		PreCmdRun:  preRunTokenDelete,
		CmdRun:     runTokenDelete,
		InitClient: true,
	})
	cmd.AddUUIDFlag(authv1.ArgTokenId, authv1.ArgIdShort, "", authv1.TokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(authv1.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(authv1.ArgToken, authv1.ArgTokenShort, "", authv1.Token, core.RequiredFlagOption())
	cmd.AddBoolFlag(authv1.ArgCurrent, authv1.ArgCurrentShort, false, "Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file", core.RequiredFlagOption())
	cmd.AddBoolFlag(authv1.ArgExpired, authv1.ArgExpiredShort, false, "Delete the Tokens that are currently expired", core.RequiredFlagOption())
	cmd.AddBoolFlag(authv1.ArgAll, authv1.ArgAllShort, false, "Delete the Tokens under your account", core.RequiredFlagOption())
	cmd.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are deleted")

	return cmd
}

func preRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authv1.ArgTokenId}, []string{authv1.ArgCurrent}, []string{authv1.ArgExpired}, []string{authv1.ArgAll}, []string{authv1.ArgToken})
}

func runTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgTokenId)) {
		return runTokenDeleteById(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgCurrent)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgCurrent)) {
		return runTokenDeleteCurrent(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgExpired)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgExpired)) {
		return runTokenDeleteExpired(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgAll)) && viper.GetBool(core.GetFlagName(c.NS, authv1.ArgAll)) {
		return runTokenDeleteAll(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgToken)) {
		return runTokenDeleteByToken(c)
	}
	return nil
}

func runTokenDeleteAll(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete all tokens")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting all tokens...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Warn("Status: all tokens have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting all tokens")
}

func runTokenDeleteExpired(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete expired tokens")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting expired tokens...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Warn("Status: expired tokens have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting expired tokens")
}

func runTokenDeleteCurrent(c *core.CommandConfig) error {
	c.Printer.Verbose("Note: This operation is based on Authorization Header for Bearer Token")
	if viper.GetString(constants.Token) == "" {
		return errors.New(fmt.Sprintf("no token found. Please make sure you have exported the %s environment variable or you have token set in the config file",
			sdkgoauth.IonosTokenEnvVar))
	}
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete CURRENT token")); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting current token...")
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("CURRENT", viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Warn("Status: current token have been successfully deleted")
			}
		}
	}
	return errors.New("error deleting current token")
}

func runTokenDeleteById(c *core.CommandConfig) error {
	tokenId := viper.GetString(core.GetFlagName(c.NS, authv1.ArgTokenId))
	c.Printer.Verbose("Token ID: %s", tokenId)
	if err := utils.AskForConfirm(c.Stdin, c.Printer, fmt.Sprintf("delete token with ID: %s", tokenId)); err != nil {
		return err
	}
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(tokenId, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Warn("Status: token has been successfully deleted")
			}
		}
	}
	return errors.New("error deleting token")
}

func runTokenDeleteByToken(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, authv1.ArgToken))
	c.Printer.Verbose("Token content is: %s", token)

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	if !confirm.Ask(fmt.Sprintf("delete token with ID: %s", tokenId), viper.GetBool(constants.ArgForce)) {
		return nil
	}

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				return c.Printer.Warn("Status: token has been successfully deleted")
			}
		}
	}

	return errors.New("error deleting token")
}
