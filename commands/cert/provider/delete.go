package provider

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/viper"
)

func ProviderDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "provider",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete an Provider",
		Example:   "ionosctl certmanager provider delete --provider-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagProviderID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if c.Command.Command.Flags().Changed(constants.ArgAll) {
				all, _ := c.Command.Command.Flags().GetBool(constants.ArgAll)
				if all {
					return deleteAll(c)
				}
			}

			providerId, _ := c.Command.Command.Flags().GetString(constants.FlagProviderID)
			z, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersFindById(context.Background(), providerId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting Provider by id %s: %w", providerId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete Provider with name: %s, id: %s ", z.Properties.Name, z.Id),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().CertManagerClient.ProviderApi.ProvidersDelete(context.Background(), providerId).Execute()
			if err != nil {
				return fmt.Errorf("failed to delete the Provider: %w", err)
			}

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagProviderID, constants.FlagIdShort, "", "Provide the specified Provider", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return ProviderIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, fmt.Sprintf("Delete all Providers. Required or -%s", constants.FlagGatewayShort))

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Deleting all providers!"))
	xs, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersGet(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("failed getting the Providers: %w", err)
	}
	err = functional.ApplyAndAggregateErrors(xs.GetItems(), func(z cert.ProviderRead) error {
		yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete Provider with name: %s, id: %s ", z.Properties.Name, z.Id),
			viper.GetBool(constants.ArgForce))
		if yes {
			_, delErr := client.Must().CertManagerClient.ProviderApi.ProvidersDelete(context.Background(), z.Id).Execute()
			if delErr != nil {
				return fmt.Errorf("failed deleting %s (name: %s): %w", z.Id, z.Properties.Name, delErr)
			}
		}
		return nil
	})

	return err
}
