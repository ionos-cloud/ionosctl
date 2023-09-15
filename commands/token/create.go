package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
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
	cmd.AddIntFlag(authservice.ArgContractNo, "", 0, "Users with multiple contracts can provide the contract number, for which the token is generated")

	return cmd
}

func runTokenCreate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Generating new token.."))

	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(contractNumberMessage,
			viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))))
	}

	newJwt, _, err := c.AuthV1Services.Tokens().Create(viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)))
	if err != nil {
		return err
	}

	if newJwt != nil {
		if token, ok := newJwt.GetTokenOk(); ok && token != nil {
			fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateRawOutput(*token))

			return nil
		} else {
			return errors.New("error getting generated token")
		}
	} else {
		return errors.New("error getting generated JWT")
	}
}
