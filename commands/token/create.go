package token

import (
	"context"
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/utils"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	"github.com/spf13/viper"
)

func TokenPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
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
		},
	)
	cmd.AddIntFlag(
		authservice.ArgContractNo, "", 0,
		"Users with multiple contracts can provide the contract number, for which the token is generated",
	)
	cmd.AddSetFlag(
		constants.FlagTtl, "", "1Y", []string{"1h", "4h", "1D", "7D", "1M", "2M", "3M", "6M", "1Y"},
		"Token Time to Live in seconds",
	)

	return cmd
}

func runTokenCreate(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Generating new token.."))

	var contractNumber int32
	if viper.IsSet(core.GetFlagName(c.NS, authservice.ArgContractNo)) {
		fmt.Fprintf(
			c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				contractNumberMessage,
				viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo)),
			),
		)
		contractNumber = viper.GetInt32(core.GetFlagName(c.NS, authservice.ArgContractNo))
	}

	var ttl int
	var err error
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagTtl)) {
		fmt.Fprintf(
			c.Command.Command.ErrOrStderr(),
			jsontabwriter.GenerateVerboseOutput(
				"Token TTL: %v",
				viper.GetString(core.GetFlagName(c.NS, constants.FlagTtl)),
			),
		)

		ttl, err = utils.ConvertTime(viper.GetString(core.GetFlagName(c.NS, constants.FlagTtl)), utils.Seconds)
		if err != nil {
			return err
		}
	}
	newJwt, _, err := c.AuthV1Services.Tokens().Create(contractNumber, int32(ttl))
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
