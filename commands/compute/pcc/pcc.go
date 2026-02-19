package pcc

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allPccJSONPaths = map[string]string{
		"PccId":       "id",
		"Name":        "properties.name",
		"Description": "properties.description",
		"State":       "metadata.state",
	}

	allPccPeerJSONPaths = map[string]string{
		"LanId":          "id",
		"LanName":        "name",
		"DatacenterId":   "datacenterId",
		"DatacenterName": "datacenterName",
		"Location":       "location",
	}

	defaultPccCols = []string{"PccId", "Name", "Description", "State"}
)

func PccCmd() *core.Command {
	pccCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "pcc",
			Aliases:          []string{"cc"},
			Short:            "Cross-Connect Operations",
			Long:             "The sub-commands of `ionosctl compute pcc` allow you to list, get, create, update, delete Cross-Connect. To add Cross-Connect to a Lan, check the `ionosctl compute lan update` command.",
			TraverseChildren: true,
		},
	}
	globalFlags := pccCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultPccCols, tabheaders.ColsMessage(defaultPccCols))
	_ = viper.BindPFlag(core.GetFlagName(pccCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = pccCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPccCols, cobra.ShellCompDirectiveNoFileComp
	})

	pccCmd.AddCommand(PccListCmd())
	pccCmd.AddCommand(PccGetCmd())
	pccCmd.AddCommand(PccCreateCmd())
	pccCmd.AddCommand(PccUpdateCmd())
	pccCmd.AddCommand(PccDeleteCmd())
	pccCmd.AddCommand(PeersCmd())

	return core.WithConfigOverride(pccCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
