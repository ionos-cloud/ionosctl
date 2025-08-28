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
	cmd.AddUUIDFlag(constants.FlagTokenId, constants.FlagIdShort, "", constants.DescTokenId, core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTokenId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.TokensIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", constants.DescToken, core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.FlagCurrent, constants.FlagCurrentShortAuth, false, "Delete the Token that is currently used. This requires a token to be set for authentication via environment variable IONOS_TOKEN or via config file", core.RequiredFlagOption())
	cmd.Command.Flags().MarkShorthandDeprecated(constants.FlagCurrent, "it will be removed in a future release.")
	cmd.AddBoolFlag(constants.FlagExpired, constants.FlagExpiredShortAuth, false, "Delete the Tokens that are currently expired", core.RequiredFlagOption())
	cmd.Command.Flags().MarkShorthandDeprecated(constants.FlagExpired, "it will be removed in a future release.")

	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts must provide the contract number, for which the tokens are deleted")

	// Deprecated: use -a instead of -A
	cmd.AddBoolFlag(constants.ArgAllAddedAsHidden, constants.ArgAllShortDeprecated, false, "Delete the Tokens under your account", core.RequiredFlagOption())
	cmd.Command.Flags().MarkShorthandDeprecated(constants.ArgAllAddedAsHidden, "please use -a instead")
	_ = cmd.Command.Flags().MarkHidden(constants.ArgAllAddedAsHidden)
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete the Tokens under your account", core.RequiredFlagOption())

	return cmd
}

func preRunTokenDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.FlagTokenId}, []string{constants.FlagCurrent}, []string{constants.FlagExpired}, []string{constants.ArgAll}, []string{constants.ArgToken})
}

func runTokenDelete(c *core.CommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(contractNumberMessage,
			viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagTokenId)) {
		return runTokenDeleteById(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagCurrent)) && viper.GetBool(core.GetFlagName(c.NS, constants.FlagCurrent)) {
		return runTokenDeleteCurrent(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagExpired)) && viper.GetBool(core.GetFlagName(c.NS, constants.FlagExpired)) {
		return runTokenDeleteExpired(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgAll)) && viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)) {
		return runTokenDeleteAll(c)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return runTokenDeleteByToken(c)
	}
	return nil
}

func runTokenDeleteAll(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("delete all tokens"), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("ALL", viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Status: all tokens have been successfully deleted"))
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting expired tokens..."))

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByCriteria("EXPIRED", viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Status: expired tokens have been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting expired tokens")
}

func runTokenDeleteCurrent(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
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

	if !confirm.FAsk(c.Command.Command.InOrStdin(), ask.String(), viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	tokenResponse, _, err := c.AuthV1Services.Tokens().DeleteByID(fmt.Sprintf("%v", tokenId), viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if tokenResponse != nil {
		if success, ok := tokenResponse.GetSuccessOk(); ok && success != nil {
			if *success {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Status: token has been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting token")
}

func runTokenDeleteById(c *core.CommandConfig) error {
	tokenId := viper.GetString(core.GetFlagName(c.NS, constants.FlagTokenId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Token ID: %s", tokenId))
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
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Status: token has been successfully deleted"))
				return nil
			}
		}
	}
	return errors.New("error deleting token")
}

func runTokenDeleteByToken(c *core.CommandConfig) error {
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
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Status: token has been successfully deleted"))
				return nil
			}
		}
	}

	return errors.New("error deleting token")
}
