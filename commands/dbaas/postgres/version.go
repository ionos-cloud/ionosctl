package postgres

import (
	"context"
	"fmt"

	pgsqlcompleter "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	dbaaspg "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultPgsqlVersionCols, tabheaders.ColsMessage(defaultPgsqlVersionCols))
	_ = viper.BindPFlag(core.GetFlagName(pgsqlversionCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = pgsqlversionCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultPgsqlVersionCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, pgsqlversionCmd, core.CommandBuilder{
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
	})
	_ = list // Actually used - added through "NewCommand" func. TODO: This is confusing!

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, pgsqlversionCmd, core.CommandBuilder{
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
	})
	get.AddUUIDFlag(constants.FlagClusterId, dbaaspg.ArgIdShort, "", dbaaspg.ClusterId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return pgsqlcompleter.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return pgsqlversionCmd
}

func RunPgsqlVersionList(c *core.CommandConfig) error {
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Versions().List()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allPgsqlVersionJSONPaths, versionList.PostgresVersionList,
		tabheaders.GetHeadersAllDefault(defaultPgsqlVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunPgsqlVersionGet(c *core.CommandConfig) error {
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Versions().Get(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
	)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", allPgsqlVersionJSONPaths, versionList.PostgresVersionList,
		tabheaders.GetHeadersAllDefault(defaultPgsqlVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

// Output Printing

var (
	allPgsqlVersionJSONPaths = map[string]string{
		"PostgresVersions": "data.*.name",
	}

	defaultPgsqlVersionCols = []string{"PostgresVersions"}
)
