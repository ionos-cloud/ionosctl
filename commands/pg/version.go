package pg

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	pgsqlcompleter "github.com/ionos-cloud/ionosctl/commands/pg/completer"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/services/dbaas-pg"
	pgsqlresources "github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PgsqlVersionCmd() *core.Command {
	ctx := context.TODO()
	pgsqlversionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "DBaaS Postgres Version Operations",
			Long:             "The sub-commands of `ionosctl dbaas-pgsql version` allow you to get information about available DBaaS Postgres Versions.",
			TraverseChildren: true,
		},
	}
	globalFlags := pgsqlversionCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultPgsqlVersionCols, printer.ColsMessage(defaultPgsqlVersionCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(pgsqlversionCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = pgsqlversionCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPgsqlVersionCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, pgsqlversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
		Resource:   "version",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List DBaaS Postgres Versions",
		LongDesc:   "Use this command to retrieve all available DBaaS Postgres versions.",
		Example:    listVersionExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunPgsqlVersionList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, pgsqlversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
		Resource:   "version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get DBaaS Postgres Versions for a Cluster",
		LongDesc:   "Use this command to retrieve a list of all PostgreSQL versions available for a specified Cluster.\n\nRequired values to run command:\n\n* Cluster Id",
		Example:    getVersionExample,
		PreCmdRun:  PreRunClusterId,
		CmdRun:     RunPgsqlVersionGet,
		InitClient: true,
	})
	get.AddStringFlag(cloudapidbaaspgsql.ArgClusterId, cloudapidbaaspgsql.ArgIdShort, "", cloudapidbaaspgsql.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapidbaaspgsql.ArgClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return pgsqlcompleter.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return pgsqlversionCmd
}

func RunPgsqlVersionList(c *core.CommandConfig) error {
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Versions().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getPgsqlVersionPrint(c, &versionList))
}

func RunPgsqlVersionGet(c *core.CommandConfig) error {
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Versions().Get(
		viper.GetString(core.GetFlagName(c.NS, cloudapidbaaspgsql.ArgClusterId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getPgsqlVersionPrint(c, &versionList))
}

// Output Printing

var defaultPgsqlVersionCols = []string{"Name"}

type PgsqlVersionPrint struct {
	Name string `json:"Name,omitempty"`
}

func getPgsqlVersionPrint(c *core.CommandConfig, postgresVersionList *pgsqlresources.PostgresVersionList) printer.Result {
	r := printer.Result{}
	if c != nil {
		if postgresVersionList != nil {
			r.OutputJSON = postgresVersionList
			r.KeyValue = getPgsqlVersionsKVMaps(postgresVersionList)
			r.Columns = getPgsqlVersionCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getPgsqlVersionCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pgsqlVersionCols []string
		columnsMap := map[string]string{
			"Name": "Name",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				pgsqlVersionCols = append(pgsqlVersionCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return pgsqlVersionCols
	} else {
		return defaultPgsqlVersionCols
	}
}

func getPgsqlVersionsKVMaps(postgresVersionList *pgsqlresources.PostgresVersionList) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)
	if postgresVersionList != nil {
		if dataOk, ok := postgresVersionList.GetDataOk(); ok && dataOk != nil {
			for _, data := range *dataOk {
				var uPrint PgsqlVersionPrint
				if nameOk, ok := data.GetNameOk(); ok && nameOk != nil {
					uPrint.Name = *nameOk
				}
				o := structs.Map(uPrint)
				out = append(out, o)
			}
		}
	}
	return out
}
