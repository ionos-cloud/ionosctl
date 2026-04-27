package postgres

import (
	"context"

	pgsqlcompleter "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/postgres/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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
	pgsqlversionCmd.AddColsFlag(allPgsqlVersionCols)

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
	get.AddUUIDFlag(constants.FlagClusterId, constants.FlagIdShort, "", constants.DescCluster, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(constants.FlagClusterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return pgsqlcompleter.ClustersIds(), cobra.ShellCompDirectiveNoFileComp
	})

	return pgsqlversionCmd
}

func RunPgsqlVersionList(c *core.CommandConfig) error {
	versionList, _, err := client.Must().PostgresClient.ClustersApi.PostgresVersionsGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allPgsqlVersionCols).Prefix("data").Print(versionList)
}

func RunPgsqlVersionGet(c *core.CommandConfig) error {
	versionList, _, err := client.Must().PostgresClient.ClustersApi.ClusterPostgresVersionsGet(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allPgsqlVersionCols).Prefix("data").Print(versionList)
}

// Output Printing

var allPgsqlVersionCols = []table.Column{
	{Name: "PostgresVersions", JSONPath: "name", Default: true},
}
