package dataplatform

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func VersionsCmd() *core.Command {
	ctx := context.TODO()
	versionsCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "versions",
			Aliases:          []string{"v"},
			Short:            "Managed Data Platform API versions",
			Long:             "The sub-commands of `ionosctl dataplatform cluster` allow you to manage the Data Platform Clusters under your account.",
			TraverseChildren: true,
		},
	}

	/*
		List Command
	*/
	list := core.NewCommand(ctx, versionsCmd, core.CommandBuilder{
		Namespace:  "dataplatform",
		Resource:   "versions",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Managed Data Platform API versions",
		LongDesc:   "Use this command to retrieve a list of Managed Data Platform API versions.",
		Example:    listVersionsExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunVersionsList,
		InitClient: true,
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, "When using text output, don't print headers")
	list.AddStringSliceFlag(config.ArgCols, "", defaultVersionsCols, printer.ColsMessage(allVersionsCols))
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allVersionsCols, cobra.ShellCompDirectiveNoFileComp
	})

	return versionsCmd
}

func RunVersionsList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Versions...")
	versions, _, err := c.DataPlatformServices.Versions().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getVersionsPrint(c, getVersions(versions)))
}

// Output Printing

var (
	defaultVersionsCols = []string{"DataPlatformVersion"}
	allVersionsCols     = []string{"DataPlatformVersion"}
)

type VersionsPrint struct {
	DataPlatformVersion string `json:"DataPlatformVersion,omitempty"`
}

func getVersionsPrint(c *core.CommandConfig, versionsList []string) printer.Result {
	r := printer.Result{}
	if c != nil {
		if versionsList != nil {
			r.OutputJSON = versionsList
			r.KeyValue = getVersionsKVMaps(versionsList)
			r.Columns = getVersionsCols(core.GetFlagName(c.NS, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getVersionsCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultVersionsCols
	}

	columnsMap := map[string]string{
		"DataPlatformVersion": "DataPlatformVersion",
	}
	var versionsCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			versionsCols = append(versionsCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return versionsCols
}

func getVersions(versions []string) []string {
	v := make([]string, 0)
	if len(versions) > 0 {
		for _, version := range versions {
			v = append(v, version)
		}
	}
	return v
}

func getVersionsKVMaps(versions []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(versions))
	for _, version := range versions {
		var versionsPrint VersionsPrint
		versionsPrint.DataPlatformVersion = version

		o := structs.Map(versionsPrint)
		out = append(out, o)
	}
	return out
}
