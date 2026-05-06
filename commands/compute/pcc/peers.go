package pcc

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allPccPeerCols = []table.Column{
	{Name: "LanId", JSONPath: "id", Default: true},
	{Name: "LanName", JSONPath: "name", Default: true},
	{Name: "DatacenterId", JSONPath: "datacenterId", Default: true},
	{Name: "DatacenterName", JSONPath: "datacenterName", Default: true},
	{Name: "Location", JSONPath: "location", Default: true},
}

func PeersCmd() *core.Command {
	peerCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "peers",
			Short:            "Cross-Connect Peers Operations",
			Long:             "The sub-command of `ionosctl compute pcc peers` allows you to get a list of Peers from a Cross-Connect.",
			TraverseChildren: true,
		},
	}
	peerCmd.AddColsFlag(allPccPeerCols)

	/*
		List Command
	*/
	listPeers := core.NewCommand(context.TODO(), peerCmd, core.CommandBuilder{
		Namespace:  "pcc",
		Resource:   "peers",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "Get a list of Peers from a Cross-Connect",
		LongDesc:   "Use this command to get a list of Peers from a Cross-Connect.\n\nRequired values to run command:\n\n* Pcc Id",
		Example:    `ionosctl compute pcc peers list --pcc-id PCC_ID`,
		PreCmdRun:  PreRunPccId,
		CmdRun:     RunPccPeersList,
		InitClient: true,
	})
	listPeers.AddUUIDFlag(cloudapiv6.ArgPccId, "", "", cloudapiv6.PccId, core.RequiredFlagOption())
	_ = listPeers.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgPccId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.PccsIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return core.WithConfigOverride(peerCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
