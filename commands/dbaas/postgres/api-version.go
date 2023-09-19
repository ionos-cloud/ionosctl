package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/json2table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	sdkgo "github.com/ionos-cloud/sdk-go-dbaas-postgres"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultAPIVersionCols, printer.ColsMessage(defaultAPIVersionCols))
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
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

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
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, "When using text output, don't print headers")

	return apiversionCmd
}

func RunAPIVersionList(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting all available API Versions..."))

	versionList, _, err := c.CloudApiDbaasPgsqlServices.Infos().List()
	if err != nil {
		return err
	}

	var versionListConverted []map[string]interface{}
	for _, v := range versionList.Versions {
		temp, err := convertAPIVersionToTable(v)
		if err != nil {
			return err
		}

		versionListConverted = append(versionListConverted, temp...)
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(versionList.Versions, versionListConverted,
		tabheaders.GetHeadersAllDefault(defaultAPIVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

func RunAPIVersionGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Getting the current API Version..."))

	apiVersion, _, err := c.CloudApiDbaasPgsqlServices.Infos().Get()
	if err != nil {
		return err
	}

	apiVersionConverted, err := convertAPIVersionToTable(apiVersion.APIVersion)
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutputPreconverted(apiVersion.APIVersion, apiVersionConverted,
		tabheaders.GetHeadersAllDefault(defaultAPIVersionCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)
	return nil
}

// Output Printing

var (
	allAPIVersionJSONPaths = map[string]string{
		"Version": "name",
	}

	defaultAPIVersionCols = []string{"Version", "SwaggerUrl"}
)

func convertAPIVersionToTable(apiVersion sdkgo.APIVersion) ([]map[string]interface{}, error) {
	swaggerUrlOk, ok := apiVersion.GetSwaggerUrlOk()
	if !ok || swaggerUrlOk == nil {
		return nil, fmt.Errorf("could not retrieve PostgreSQL API Version swagger URL")
	}

	if strings.HasPrefix(*swaggerUrlOk, "appserver:8181/postgresql") {
		*swaggerUrlOk = strings.TrimPrefix(*swaggerUrlOk, "appserver:8181/postgresql")
	}
	if !strings.HasPrefix(*swaggerUrlOk, sdkgo.DefaultIonosServerUrl) {
		*swaggerUrlOk = fmt.Sprintf("%s%s", sdkgo.DefaultIonosServerUrl, *swaggerUrlOk)
	}

	temp, err := json2table.ConvertJSONToTable("", allAPIVersionJSONPaths, apiVersion)
	if err != nil {
		return nil, fmt.Errorf("could not convert from JSON to Table format: %w", err)
	}

	temp[0]["SwaggerUrl"] = *swaggerUrlOk

	return temp, nil
}
