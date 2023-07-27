package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	authv1 "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/viper"
)

func TokenPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "token",
		Resource:   "token",
		Verb:       "generate",
		Aliases:    []string{"create"},
		ShortDesc:  "Create a new Token",
		LongDesc:   "Use this command to generate a new Token. Only the JSON Web Token, associated with user credentials, will be displayed.",
		Example:    generateTokenExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     runTokenCreate,
		InitClient: true,
	})
	cmd.AddIntFlag(authv1.ArgContractNo, "", 0, "Users with multiple contracts can provide the contract number, for which the token is generated")

	return cmd
}

func runTokenCreate(c *core.CommandConfig) error {
	c.Printer.Verbose("Generating new token..")
	if viper.IsSet(core.GetFlagName(c.NS, authv1.ArgContractNo)) {
		c.Printer.Verbose(contractNumberMessage, viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	}
	newJwt, _, err := c.AuthV1Services.Tokens().Create(viper.GetInt32(core.GetFlagName(c.NS, authv1.ArgContractNo)))
	if err != nil {
		return err
	}
	if newJwt != nil {
		if token, ok := newJwt.GetTokenOk(); ok && token != nil {
			_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), *token)
			return err
		} else {
			return errors.New("error getting generated token")
		}
	} else {
		return errors.New("error getting generated JWT")
	}
}
