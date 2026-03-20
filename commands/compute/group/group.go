package group

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultGroupCols = []string{"GroupId", "Name", "CreateDataCenter", "CreateSnapshot", "CreatePcc", "CreateBackupUnit",
		"CreateInternetAccess", "CreateK8s", "ReserveIp"}
	allGroupCols = []string{"GroupId", "Name", "CreateDataCenter", "CreateSnapshot", "ReserveIp", "AccessActivityLog",
		"CreatePcc", "S3Privilege", "CreateBackupUnit", "CreateInternetAccess", "CreateK8s", "CreateFlowLog",
		"AccessAndManageMonitoring", "AccessAndManageCertificates", "AccessAndManageDns", "ManageDBaaS", "ManageRegistry",
	}
)

func GroupCmd() *core.Command {
	groupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "group",
			Aliases:          []string{"g"},
			Short:            "Group Operations",
			Long:             "The sub-commands of `ionosctl compute group` allow you to list, get, create, update, delete Groups, but also operations: add/remove/list/update User from the Group.",
			TraverseChildren: true,
		},
	}
	globalFlags := groupCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultGroupCols, tabheaders.ColsMessage(allGroupCols))
	_ = viper.BindPFlag(core.GetFlagName(groupCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = groupCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	groupCmd.AddCommand(GroupListCmd())
	groupCmd.AddCommand(GroupGetCmd())
	groupCmd.AddCommand(GroupCreateCmd())
	groupCmd.AddCommand(GroupUpdateCmd())
	groupCmd.AddCommand(GroupDeleteCmd())
	groupCmd.AddCommand(GroupResourceCmd())
	groupCmd.AddCommand(GroupUserCmd())

	return core.WithConfigOverride(groupCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
