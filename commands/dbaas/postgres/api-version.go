package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	shared "github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/cobra"
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
			Short:            "Show the DBaaS PostgreSQL REST API version",
			Long:             "Report the version of the DBaaS PostgreSQL REST API itself (the management API this CLI talks to), not the PostgreSQL engine version. For engine versions see 'dbaas postgres version'. Use `list` for all API versions the service exposes or `get` for the one currently in use, each with a link to its OpenAPI/Swagger document.",
			TraverseChildren: true,
		},
	}
	apiversionCmd.AddColsFlag(allAPIVersionCols)

	/*
		List Command
	*/
	list := core.NewCommand(ctx, apiversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-postgres",
		Resource:   "api-version",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List all DBaaS PostgreSQL API versions",
		LongDesc:   "Retrieve every version of the DBaaS PostgreSQL REST API that the service exposes, each with a link to its OpenAPI/Swagger document.",
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
		ShortDesc:  "Get the current DBaaS PostgreSQL API version",
		LongDesc:   "Retrieve the version of the DBaaS PostgreSQL REST API currently in use, along with a link to its OpenAPI/Swagger document.",
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

	return c.Printer(allAPIVersionCols).Print(versionList)
}

func RunAPIVersionGet(c *core.CommandConfig) error {
	c.Verbose("Getting the current API Version...")

	apiVersion, _, err := client.Must().PostgresClient.MetadataApi.InfosVersionGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	return c.Printer(allAPIVersionCols).Print(apiVersion)
}
