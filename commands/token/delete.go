package token

import (
	"context"
	"errors"
	"fmt"
	"strings"

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
	cmd.AddUUIDFlag(authservice.FlagTokenId, authservice.FlagIdShort, "", authservice.TokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(authservice.FlagTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(authservice.FlagToken, authservice.FlagTokenShort, "", authservice.Token, core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.FlagCurrent, authservice.FlagCurrentShort, false, "Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file", core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.FlagExpired, authservice.FlagExpiredShort, false, "Delete the Tokens that are currently expired", core.RequiredFlagOption())
	cmd.AddBoolFlag(authservice.FlagAll, authservice.FlagAllShort, false, "Delete the Tokens under your account", core.RequiredFlagOption())
	cmd.AddIntFlag(authservice.FlagContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are deleted")

	return cmd
}

func preRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{authservice.FlagTokenId}, []string{authservice.FlagCurrent}, []string{authservice.FlagExpired}, []string{authservice.FlagAll}, []string{authservice.FlagToken})
}

func runTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(contractNumberMessage,
			viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagTokenId)) {
		return runTokenDeleteById(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagCurrent)) && viper.GetBool(core.GetFlagName(c.NS, authservice.FlagCurrent)) {
		return runTokenDeleteCurrent(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagExpired)) && viper.GetBool(core.GetFlagName(c.NS, authservice.FlagExpired)) {
		return runTokenDeleteExpired(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagAll)) && viper.GetBool(core.GetFlagName(c.NS, authservice.FlagAll)) {
		return runTokenDeleteAll(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, authservice.FlagToken)) {
		return runTokenDeleteByToken(c)
	}
	return nil
}

func runTokenDeleteAll(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete all tokens"), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL", viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo)))
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
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete expired tokens"), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting expired tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED", viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo)))
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

	tokenToDelete := viper.GetString(constants.EnvToken)
	usedEnv := true // simply used for asking the user for confirmation
	if tokenToDelete == "" {
		usedEnv = false
		tokenToDelete = viper.GetString(constants.CfgToken)
	}
	if tokenToDelete == "" {
		return errors.New(fmt.Sprintf("no token found in environment variable %s or config file", constants.EnvToken))
	}

	ask := strings.Builder{}
	ask.WriteString("delete currently used token in")
	if usedEnv {
		ask.WriteString(" IONOS_TOKEN")
	} else {
		ask.WriteString(" config file")
	}
	username, errUsername := jwt.Username(tokenToDelete)
	if errUsername == nil {
		ask.WriteString(fmt.Sprintf(" of user '%s'", username))
	}
	claims, errExtraInfo := jwt.Claims(tokenToDelete)
	if errExtraInfo == nil {
		role, errRole := jwt.Role(claims)
		if errRole == nil {
			ask.WriteString(fmt.Sprintf(" with role '%s'", role))
		}
		number, errContractNumber := jwt.ContractNumber(claims)
		if errContractNumber == nil {
			ask.WriteString(fmt.Sprintf(" contract '%d'", number))
		}
	}

	headers, err := jwt.Headers(tokenToDelete)
	if err != nil {
		return fmt.Errorf("failed to get headers from token: %w", err)
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return fmt.Errorf("failed to get token ID from token: %w", err)
	}
	ask.WriteString(fmt.Sprintf(" token id '%s'", tokenId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), ask.String(), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo)))
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

func runTokenDeleteById(c *core.CommandConfig) error {
	tokenId := viper.GetString(core.GetFlagName(c.NS, authservice.FlagTokenId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Token ID: %s", tokenId))
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete token with ID: %s", tokenId), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}
	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(tokenId, viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo)))
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
	token := viper.GetString(core.GetFlagName(c.NS, authservice.FlagToken))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Token content is: %s", token))

	headers, err := jwt.Headers(token)
	if err != nil {
		return err
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		return err
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete token with ID: %s", tokenId), viper.GetBool(constants.FlagForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.FlagContractNo)))
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
