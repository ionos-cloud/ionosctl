package snapshot

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultSnapshotCols = []string{"SnapshotId", "Name", "LicenceType", "Size", "State"}
)

func SnapshotCmd() *core.Command {
	snapshotCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "snapshot",
			Aliases:          []string{"ss", "snap"},
			Short:            "Snapshot Operations",
			Long:             "The sub-commands of `ionosctl compute snapshot` allow you to see information, to create, update, delete Snapshots.",
			TraverseChildren: true,
		},
	}
	globalFlags := snapshotCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultSnapshotCols, tabheaders.ColsMessage(defaultSnapshotCols))
	_ = viper.BindPFlag(core.GetFlagName(snapshotCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = snapshotCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultSnapshotCols, cobra.ShellCompDirectiveNoFileComp
	})

	snapshotCmd.AddCommand(SnapshotListCmd())
	snapshotCmd.AddCommand(SnapshotGetCmd())
	snapshotCmd.AddCommand(SnapshotCreateCmd())
	snapshotCmd.AddCommand(SnapshotUpdateCmd())
	snapshotCmd.AddCommand(SnapshotRestoreCmd())
	snapshotCmd.AddCommand(SnapshotDeleteCmd())

	return core.WithConfigOverride(snapshotCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
