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
			Short:            "List available PostgreSQL engine versions",
			Long:             "List the PostgreSQL engine (major) versions offered by DBaaS. Use `list` for every version currently available for new clusters (the values valid for 'cluster create --version'), or `get --cluster-id` for the versions a specific existing cluster can be upgraded to.",
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
		ShortDesc:  "List all available PostgreSQL versions",
		LongDesc:   "Retrieve every PostgreSQL engine version currently offered by DBaaS. These are the values accepted by 'cluster create --version'.",
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
		ShortDesc:  "List the PostgreSQL versions available to a cluster",
		LongDesc:   "Retrieve the PostgreSQL engine versions available for a specific cluster, i.e. the versions it can be upgraded to via 'cluster update --version'.\n\nRequired values to run command:\n\n* Cluster Id",
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
