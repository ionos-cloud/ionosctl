package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultCpuCols, printer.ColsMessage(defaultCpuCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(cpuCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = cpuCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return cpuCmd
}

func RunLocationCpuList(c *core.CommandConfig) error {
	locId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocationId))
	ids := strings.Split(locId, "/")
	if len(ids) != 2 {
		return errors.New("error getting location id & region id")
	}
	loc, resp, err := c.CloudApiV6Services.Locations().GetByRegionAndLocationId(ids[0], ids[1])
	if resp != nil {
		c.Printer.Verbose(cloudapiv6.RequestTimeMessage, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if properties, ok := loc.GetPropertiesOk(); ok && properties != nil {
		if cpus, ok := properties.GetCpuArchitectureOk(); ok && cpus != nil {
			return c.Printer.Print(printer.Result{
				OutputJSON: cpus,
				KeyValue:   getCpusKVMaps(getCpus(cpus)),
				Columns:    getCpuCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr()),
			})
		} else {
			return errors.New("error getting cpu architectures")
		}
	} else {
		return errors.New("error getting location properties")
	}
}

// Output Printing

var defaultCpuCols = []string{"CpuFamily", "MaxCores", "MaxRam", "Vendor"}

type CpuPrint struct {
	CpuFamily string `json:"CpuFamily,omitempty"`
	MaxCores  int32  `json:"MaxCores,omitempty"`
	MaxRam    string `json:"MaxRam,omitempty"`
	Vendor    string `json:"Vendor,omitempty"`
}

func getCpuCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultCpuCols
	}

	columnsMap := map[string]string{
		"CpuFamily": "CpuFamily",
		"MaxCores":  "MaxCores",
		"MaxRam":    "MaxRam",
		"Vendor":    "Vendor",
	}
	var cpusCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			cpusCols = append(cpusCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return cpusCols
}

func getCpus(cpus *[]ionoscloud.CpuArchitectureProperties) []resources.CpuArchitectureProperties {
	cs := make([]resources.CpuArchitectureProperties, 0)
	if cpus != nil {
		for _, cpuItem := range *cpus {
			cs = append(cs, resources.CpuArchitectureProperties{CpuArchitectureProperties: cpuItem})
		}
	}
	return cs
}

func getCpusKVMaps(cs []resources.CpuArchitectureProperties) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(cs))
	for _, cpuItem := range cs {
		var cpuPrint CpuPrint
		if cpuFamily, ok := cpuItem.GetCpuFamilyOk(); ok && cpuFamily != nil {
			cpuPrint.CpuFamily = *cpuFamily
		}
		if cpuCores, ok := cpuItem.GetMaxCoresOk(); ok && cpuCores != nil {
			cpuPrint.MaxCores = *cpuCores
		}
		if cpuRam, ok := cpuItem.GetMaxRamOk(); ok && cpuRam != nil {
			cpuPrint.MaxRam = fmt.Sprintf("%vMB", *cpuRam)
		}
		if cpuVendor, ok := cpuItem.GetVendorOk(); ok && cpuVendor != nil {
			cpuPrint.Vendor = *cpuVendor
		}
		o := structs.Map(cpuPrint)
		out = append(out, o)
	}
	return out
}
