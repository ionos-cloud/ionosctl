package secondary_zones

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns/secondary-zones/transfer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var allCols = []string{"Id", "Name", "Description", "PrimaryIPs", "State"}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "secondary-zone",
			Aliases:          []string{"secondary-zones", "sz"},
			Short:            "All commands related to secondary zones",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, allCols, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.AddCommand(transfer.Root())
	cmd.AddCommand(createCmd())
	cmd.AddCommand(deleteCmd())
	cmd.AddCommand(listCmd())
	cmd.AddCommand(getCmd())
	cmd.AddCommand(updateCmd())

	return cmd
}

func resolve(nameOrID string) (string, error) {
	if _, err := uuid.FromString(nameOrID); err == nil {
		return nameOrID, nil
	}

	secZones, _, err := client.Must().DnsClient.SecondaryZonesApi.SecondaryzonesGet(context.Background()).FilterZoneName(nameOrID).Execute()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve zones by name %s: %w", nameOrID, err)
	}
	if secZones.Items == nil || len(*secZones.Items) < 1 {
		return "", fmt.Errorf("no zones found with name %s", nameOrID)
	}

	return *(*secZones.Items)[0].Id, nil
}
