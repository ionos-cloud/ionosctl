package postgres

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	pgsqlcompleter "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	pgsqlresources "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func PgsqlVersionCmd() *core.Command {
	ctx := context.TODO()
	pgsqlversionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "version",
			Aliases:          []string{"v"},
			Short:            "PostgreSQL Version Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres version` allow you to get information about available DBaaS PostgreSQL Versions.",
			TraverseChildren: true,
		},
	}
	globalFlags := pgsqlversionCmd.GlobalFlags()
	globalFlags.StringSliceP(
		constants.ArgCols, "", defaultPgsqlVersionCols, printer.ColsMessage(defaultPgsqlVersionCols),
	)
	_ = viper.BindPFlag(
		core.GetFlagName(pgsqlversionCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols),
	)
	_ = pgsqlversionCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return defaultPgsqlVersionCols, cobra.ShellCompDirectiveNoFileComp
		},
	)

	/*
		List Command
	*/
	list := core.NewCommand(
		ctx, pgsqlversionCmd, core.CommandBuilder{
			Namespace:  "dbaas-postgres",
			Resource:   "version",
			Verb:       "list",
			Aliases:    []string{"l", "ls"},
			ShortDesc:  "List DBaaS PostgreSQL Versions",
			LongDesc:   "Use this command to retrieve all available DBaaS PostgreSQL versions.",
			Example:    listVersionExample,
			PreCmdRun:  core.NoPreRun,
			CmdRun:     RunPgsqlVersionList,
			InitClient: true,
		},
	)
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

	/*
		Get Command
	*/
	get := core.NewCommand(
		ctx, pgsqlversionCmd, core.CommandBuilder{
			Namespace:  "dbaas-postgres",
			Resource:   "version",
			Verb:       "get",
			Aliases:    []string{"g"},
			ShortDesc:  "Get DBaaS PostgreSQLVersions for a Cluster",
			LongDesc:   "Use this command to retrieve a list of all PostgreSQL versions available for a specified Cluster.\n\nRequired values to run command:\n\n* Cluster Id",
			Example:    getVersionExample,
			PreCmdRun:  PreRunClusterId,
			CmdRun:     RunPgsqlVersionGet,
			InitClient: true,
		},
	)
	get.AddUUIDFlag(constants.FlagClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(
		constants.FlagClusterId,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return pgsqlcompleter.ClustersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
		},
	)
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
	)
	if err != nil {
		return err
	}
	return c.Printer.Print(getPgsqlVersionPrint(c, &versionList))
}

// Output Printing

var defaultPgsqlVersionCols = []string{"PostgresVersions"}

type PgsqlVersionPrint struct {
	PostgresVersions []string `json:"PostgresVersions,omitempty"`
}

func getPgsqlVersionPrint(
	c *core.CommandConfig, postgresVersionList *pgsqlresources.PostgresVersionList,
) printer.Result {
	r := printer.Result{}
	if c != nil {
		if postgresVersionList != nil {
			r.OutputJSON = postgresVersionList
			r.KeyValue = getPgsqlVersionsKVMaps(postgresVersionList)
			r.Columns = getPgsqlVersionCols(core.GetFlagName(c.Resource, constants.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getPgsqlVersionCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pgsqlVersionCols []string
		columnsMap := map[string]string{
			"PostgresVersions": "PostgresVersions",
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
			var uPrint PgsqlVersionPrint
			for _, data := range *dataOk {
				if nameOk, ok := data.GetNameOk(); ok && nameOk != nil {
					uPrint.PostgresVersions = append(uPrint.PostgresVersions, *nameOk)
				}
			}
			o := structs.Map(uPrint)
			out = append(out, o)
		}
	}
	return out
}
