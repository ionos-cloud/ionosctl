package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/token/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
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
	cmd.AddUUIDFlag(authservice.ArgTokenId, authservice.ArgIdShort, "", authservice.TokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(authservice.ArgTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(authservice.ArgToken, authservice.ArgTokenShort, "", authservice.Token, core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.ArgCurrent, authservice.ArgCurrentShort, false, "Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file", core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.ArgExpired, authservice.ArgExpiredShort, false, "Delete the Tokens that are currently expired", core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.ArgAll, authservice.ArgAllShort, false, "Delete the Tokens under your account", core.RequiredFlagOption())
	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are deleted")

	return cmd
}

func preRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authservice.ArgTokenId}, []string{authservice.ArgCurrent}, []string{authservice.ArgExpired}, []string{authservice.ArgAll}, []string{authservice.ArgToken})
}

func runTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(contractNumberMessage,
			viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgTokenId)) {
		return runTokenDeleteById(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgCurrent)) && viper.GetBool(core.GetFlagName(c.NS, authservice.ArgCurrent)) {
		return runTokenDeleteCurrent(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgExpired)) && viper.GetBool(core.GetFlagName(c.NS, authservice.ArgExpired)) {
		return runTokenDeleteExpired(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgAll)) && viper.GetBool(core.GetFlagName(c.NS, authservice.ArgAll)) {
		return runTokenDeleteAll(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgToken)) {
		return runTokenDeleteByToken(c)
	}
	return nil
}

func runTokenDeleteAll(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete all tokens"), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL", viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Status: all tokens have been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting all tokens")
}

func runTokenDeleteExpired(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete expired tokens"), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting expired tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED", viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Status: expired tokens have been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting expired tokens")
}

func runTokenDeleteCurrent(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Note: This operation is based on Authorization Header for Bearer Token"))

	if viper.GetString(constants.CfgToken) == "" {
		return errors.New(fmt.Sprintf(
			"no token found. Please make sure you have exported the %s environment variable or you have token set in the config file",
			constants.EnvToken),
		)
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete CURRENT token"), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting current token..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("CURRENT", viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Status: current token have been successfully deleted"))
				return nil
			}
		}
	}
	return errors.New("error deleting current token")
}

func runTokenDeleteById(c *core.CommandConfig) error {
	tokenId := viper.GetString(core.GetFlagName(c.NS, authservice.ArgTokenId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Token ID: %s", tokenId))
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete token with ID: %s", tokenId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(tokenId, viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}
	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Status: token has been successfully deleted"))
				return nil
			}
		}
	}
	return errors.New("error deleting token")
}

func runTokenDeleteByToken(c *core.CommandConfig) error {
	token := viper.GetString(core.GetFlagName(c.NS, authservice.ArgToken))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Token content is: %s", token))

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete token with ID: %s", tokenId), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Status: token has been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting token")
}
