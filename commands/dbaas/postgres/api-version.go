package postgres

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func APIVersionCmd() *core.Command {
	ctx := context.TODO()
	apiversionCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "api-version",
			Aliases:          []string{"api", "info"},
			Short:            "PostgreSQL API Version Operations",
			Long:             "The sub-commands of `ionosctl dbaas postgres api-version` allow you to get information available DBaaS PostgreSQL API Versions.",
			TraverseChildren: true,
		},
	}
	globalFlags := apiversionCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultAPIVersionCols, tabheaders.ColsMessage(defaultAPIVersionCols))
	_ = viper.BindPFlag(core.GetFlagName(apiversionCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = apiversionCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAPIVersionCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, apiversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "api-version",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List DBaaS PostgreSQL API Versions",
		LongDesc:   "Use this command to retrieve all available DBaaS PostgreSQL API versions.",
		Example:    listAPIVersionExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunAPIVersionList,
		InitClient: true,
	})
	_ = list // Actually used - added through "NewCommand" func. TODO: This is confusing!

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, apiversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "api-version",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get current version of DBaaS PostgreSQL API",
		LongDesc:   "Use this command to get the current version of DBaaS PostgreSQL API.",
		Example:    getAPIVersionExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunAPIVersionGet,
		InitClient: true,
	})
	_ = get // Actually used - added through "NewCommand" func. TODO: This is confusing!

	return apiversionCmd
}

func RunAPIVersionList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting all available API Versions..."))

	versionList, _, err := c.CloudApiDbaasPgsqlServices.Infos().List()
	if err != nil {
		return err
	}

	var versionListConverted []map[string]interface{}
	for _, v := range versionList.Versions {
		temp, err := resource2table.ConvertDbaasPostgresAPIVersionToTable(v)
		if err != nil {
			return err
		}

		versionListConverted = append(versionListConverted, temp...)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(versionList.Versions, versionListConverted,
		tabheaders.GetHeadersAllDefault(defaultAPIVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

func RunAPIVersionGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting the current API Version..."))

	apiVersion, _, err := c.CloudApiDbaasPgsqlServices.Infos().Get()
	if err != nil {
		return err
	}

	apiVersionConverted, err := resource2table.ConvertDbaasPostgresAPIVersionToTable(apiVersion.APIVersion)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(apiVersion.APIVersion, apiVersionConverted,
		tabheaders.GetHeadersAllDefault(defaultAPIVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
	return nil
}

// Output Printing

var (
	defaultAPIVersionCols = []string{"Version", "SwaggerUrl"}
)
