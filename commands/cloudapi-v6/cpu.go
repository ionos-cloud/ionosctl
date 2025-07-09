package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	defaultCpuCols = []string{"CpuFamily", "MaxCores", "MaxRam", "Vendor"}
)

func CpuCmd() *core.Command {
	ctx := context.TODO()
	cpuCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cpu",
			Short:            "Location CPU Architecture Operations",
			Long:             "The sub-command of `ionosctl location cpu` allows you to see information about available CPU Architectures in different Locations.",
			TraverseChildren: true,
		},
	}
	globalFlags := cpuCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultCpuCols, tabheaders.ColsMessage(defaultCpuCols))
	_ = viper.BindPFlag(core.GetFlagName(cpuCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = cpuCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultCpuCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, cpuCmd, core.CommandBuilder{
		Namespace:  "location",
		Resource:   "cpu",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List available CPU Architecture from a Location",
		LongDesc:   "Use this command to get information about available CPU Architectures from a specific Location.\n\nRequired values to run command:\n\n* Location Id",
		Example:    listLocationCpuExample,
		PreCmdRun:  PreRunLocationId,
		CmdRun:     RunLocationCpuList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgLocationId, "", "", cloudapiv6.LocationId, core.RequiredFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocationId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)

	return core.WithConfigOverride(cpuCmd, "compute", "")
}

func RunLocationCpuList(c *core.CommandConfig) error {
	locId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId))

	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}

	loc, resp, err := c.CloudApiV6Services.Locations().GetByRegionAndLocationId(ids[0], ids[1], resources.QueryParams{})
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := loc.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting location properties")
	}

	cpus, ok := properties.GetCpuArchitectureOk()
	if !ok || cpus == nil {
		return fmt.Errorf("error getting cpu architectures")
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Cpu, *cpus, tabheaders.GetHeadersAllDefault(defaultCpuCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}
