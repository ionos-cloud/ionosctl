package provider

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func ProviderFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "provider",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a Provider",
		Example:   "ionosctl certmanager provider get --provider-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagProviderID); err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			providerId := viper.GetString(core.GetFlagName(c.NS, constants.FlagProviderID))
			r, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersFindById(context.Background(), providerId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting the Provider: %w", err)
			}

			return c.Printer(allCols).Print(r)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagProviderID, constants.FlagIdShort, "", "The certificate Provider used to issue the certificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return ProviderIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
