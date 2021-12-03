package pg

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	pgsqlresources "github.com/ionos-cloud/ionosctl/services/dbaas-pg/resources"
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
			Long:             "The sub-commands of `ionosctl pg api-version` allow you to get information available DBaaS PostgreSQL API Versions.",
			TraverseChildren: true,
		},
	}
	globalFlags := apiversionCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultAPIVersionCols, printer.ColsMessage(defaultAPIVersionCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(apiversionCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = apiversionCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultAPIVersionCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, apiversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
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

	/*
		Get Command
	*/
	core.NewCommand(ctx, apiversionCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
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

	return apiversionCmd
}

func RunAPIVersionList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting all available API Versions...")
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Infos().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getAPIVersionPrint(c, getAPIVersion(&versionList)))
}

func RunAPIVersionGet(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting the current API Version...")
	apiVersion, _, err := c.CloudApiDbaasPgsqlServices.Infos().Get()
	if err != nil {
		return err
	}
	return c.Printer.Print(getAPIVersionPrint(c, &[]pgsqlresources.APIVersion{apiVersion}))
}

// Output Printing

var defaultAPIVersionCols = []string{"Version", "SwaggerUrl"}

type APIVersionPrint struct {
	Version    string `json:"Version,omitempty"`
	SwaggerUrl string `json:"SwaggerUrl,omitempty"`
}

func getAPIVersion(a *pgsqlresources.APIVersionList) *[]pgsqlresources.APIVersion {
	u := make([]pgsqlresources.APIVersion, 0)
	if a != nil {
		for _, item := range a.Versions {
			u = append(u, pgsqlresources.APIVersion{APIVersion: item})
		}
	}
	return &u
}

func getAPIVersionPrint(c *core.CommandConfig, postgresVersionList *[]pgsqlresources.APIVersion) printer.Result {
	r := printer.Result{}
	if c != nil {
		if postgresVersionList != nil {
			r.OutputJSON = postgresVersionList
			r.KeyValue = getAPIVersionsKVMaps(postgresVersionList)
			r.Columns = getAPIVersionCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getAPIVersionCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pgsqlVersionCols []string
		columnsMap := map[string]string{
			"Version":    "Version",
			"SwaggerUrl": "SwaggerUrl",
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
		return defaultAPIVersionCols
	}
}

func getAPIVersionsKVMaps(apiVersions *[]pgsqlresources.APIVersion) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)
	if apiVersions != nil {
		for _, apiVersion := range *apiVersions {
			var uPrint APIVersionPrint
			if versionOk, ok := apiVersion.GetVersionOk(); ok && versionOk != nil {
				uPrint.Version = *versionOk
			}
			if swaggerUrlOk, ok := apiVersion.GetSwaggerUrlOk(); ok && swaggerUrlOk != nil {
				uPrint.SwaggerUrl = *swaggerUrlOk
			}
			o := structs.Map(uPrint)
			out = append(out, o)
		}
	}
	return out
}
