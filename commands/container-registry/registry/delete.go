package registry

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/containerregistry/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RegDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace: "container-registry",
			Resource:  "registry",
			Verb:      "delete",
			Aliases:   []string{"d"},
			ShortDesc: "Delete a registry",
			LongDesc:  "Delete a registry and all of its contents (repositories, artifacts and tokens). This is irreversible: pushed images and issued credentials are permanently removed. Pass --all to delete every registry in the contract.",
			Example: `# Delete one registry
ionosctl container-registry registry delete --id REGISTRY_ID

# Delete every registry in the contract
ionosctl container-registry registry delete --all`,
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagRegistryId, "i", "", "The unique ID of the registry to delete", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.FlagRegistryId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete every registry in the contract instead of a single one (--id is then ignored)")

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	allFlag, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if err != nil {
		return err
	}

	if allFlag {
		return core.DeleteAll(c, core.DeleteAllOptions[containerregistry.RegistryResponse]{
			Resource: "Container Registry",
			List: func() ([]containerregistry.RegistryResponse, error) {
				regs, _, err := client.Must().RegistryClient.RegistriesApi.RegistriesGet(context.Background()).Execute()
				if err != nil {
					return nil, err
				}
				return regs.Items, nil
			},
			Summary: func(reg containerregistry.RegistryResponse) string {
				return fmt.Sprintf("name: %s, id: %s", reg.Properties.Name, *reg.Id)
			},
			ID: func(reg containerregistry.RegistryResponse) string {
				return *reg.Id
			},
			Delete: func(reg containerregistry.RegistryResponse) error {
				_, err := client.Must().RegistryClient.RegistriesApi.RegistriesDelete(context.Background(), *reg.Id).Execute()
				return err
			},
		})
	} else {
		id, err := c.Command.Command.Flags().GetString(constants.FlagRegistryId)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("delete Container Registry: %s", id)

		if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = client.Must().RegistryClient.RegistriesApi.RegistriesDelete(context.Background(), id).Execute()
		if err != nil {
			return err
		}
	}

	return nil
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagRegistryId},
		[]string{constants.ArgAll},
	)
}
