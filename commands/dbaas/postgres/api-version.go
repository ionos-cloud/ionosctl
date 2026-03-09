package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var allAPIVersionCols = []table.Column{
	{Name: "Version", JSONPath: "name", Default: true},
	{Name: "SwaggerUrl", Default: true, Format: func(item map[string]any) any {
		swaggerUrl, _ := item["swaggerUrl"].(string)
		if swaggerUrl == "" {
			return ""
		}
		if strings.HasPrefix(swaggerUrl, "appserver:8181/postgresql") {
			swaggerUrl = strings.TrimPrefix(swaggerUrl, "appserver:8181/postgresql")
		}
		if !strings.HasPrefix(swaggerUrl, shared.DefaultIonosServerUrl) {
			swaggerUrl = fmt.Sprintf("%s%s", shared.DefaultIonosServerUrl, swaggerUrl)
		}
		return swaggerUrl
	}},
}

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
	globalFlags.StringSliceP(constants.ArgCols, "", nil, table.ColsMessage(allAPIVersionCols))
	_ = viper.BindPFlag(core.GetFlagName(apiversionCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = apiversionCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return table.AllCols(allAPIVersionCols), cobra.ShellCompDirectiveNoFileComp
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
	c.Verbose("Getting all available API Versions...")

	versionList, _, err := client.Must().PostgresClient.MetadataApi.InfosVersionsGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	return c.Out(table.Sprint(allAPIVersionCols, versionList, cols))
}

func RunAPIVersionGet(c *core.CommandConfig) error {
	c.Verbose("Getting the current API Version...")

	apiVersion, _, err := client.Must().PostgresClient.MetadataApi.InfosVersionGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	return c.Out(table.Sprint(allAPIVersionCols, apiVersion, cols))
}
