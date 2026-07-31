package token

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func TokenGetCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "get",
			Aliases:    []string{"g"},
			ShortDesc:  "Get a token's properties",
			LongDesc:   "Retrieve a single token of a registry, showing its name, expiry date, status and username. Note: the token password is not returned here - it is only shown once at create/replace time.",
			Example:    "ionosctl container-registry token get --registry-id REGISTRY_ID --token-id TOKEN_ID",
			PreCmdRun:  PreCmdGetToken,
			CmdRun:     CmdGetToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, constants.FlagRegistryIdShort, "", "The unique ID of the registry that owns the token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddStringFlag(FlagTokenId, "", "", "The unique ID of the token to retrieve")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagTokenId,
		func(cobracmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return TokensIds(
				viper.GetString(
					core.GetFlagName(
						cmd.NS, constants.FlagRegistryId,
					),
				),
			), cobra.ShellCompDirectiveNoFileComp
		},
	)

	return cmd
}

func CmdGetToken(c *core.CommandConfig) error {
	regId := viper.GetString(core.GetFlagName(c.NS, constants.FlagRegistryId))
	tokenId := viper.GetString(core.GetFlagName(c.NS, FlagTokenId))

	token, _, err := client.Must().RegistryClient.TokensApi.RegistriesTokensFindById(context.Background(), regId, tokenId).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allCols).Print(token)
}

func PreCmdGetToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagTokenId, constants.FlagRegistryId)
	if err != nil {
		return err
	}

	return nil
}
