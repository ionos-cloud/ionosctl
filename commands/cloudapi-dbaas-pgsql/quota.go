package cloudapi_dbaas_pgsql

import (
	"context"
	"errors"
	"io"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	pgsqlresources "github.com/ionos-cloud/ionosctl/services/cloudapi-dbaas-pgsql/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func QuotaCmd() *core.Command {
	ctx := context.TODO()
	quotaCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "quota",
			Aliases:          []string{"q"},
			Short:            "DBaaS Postgres Quota Operations",
			Long:             "The sub-commands of `ionosctl dbaas-pgsql quota` allow you to get information about limits on your account.",
			TraverseChildren: true,
		},
	}
	globalFlags := quotaCmd.GlobalFlags()
	globalFlags.StringSliceP(config.ArgCols, "", defaultQuotaCols, printer.ColsMessage(defaultQuotaCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(quotaCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = quotaCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultQuotaCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, quotaCmd, core.CommandBuilder{
		Namespace:  "dbaas-pgsql",
		Resource:   "quota",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Quota Limits and Usage",
		LongDesc:   "Use this command to get the current quota limits and usage.",
		Example:    getQuotaExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunQuotaList,
		InitClient: true,
	})

	return quotaCmd
}

func RunQuotaList(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting the current quota limits and usage...")
	versionList, _, err := c.CloudApiDbaasPgsqlServices.Quotas().Get()
	if err != nil {
		return err
	}
	return c.Printer.Print(getQuotaPrint(c, &versionList))
}

// Output Printing

var defaultQuotaCols = []string{"QuotaUsage", "QuotaLimits"}

type QuotaPrint struct {
	QuotaUsage  map[string]interface{} `json:"QuotaUsage,omitempty"`
	QuotaLimits map[string]interface{} `json:"QuotaLimits,omitempty"`
}

func getQuotaPrint(c *core.CommandConfig, postgresVersionList *pgsqlresources.QuotaList) printer.Result {
	r := printer.Result{}
	if c != nil {
		if postgresVersionList != nil {
			r.OutputJSON = postgresVersionList
			r.KeyValue = getQuotasKVMaps(postgresVersionList)
			r.Columns = getQuotaCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getQuotaCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var pgsqlVersionCols []string
		columnsMap := map[string]string{
			"QuotaUsage":  "QuotaUsage",
			"QuotaLimits": "QuotaLimits",
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
		return defaultQuotaCols
	}
}

func getQuotasKVMaps(postgresVersionList *pgsqlresources.QuotaList) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 0)
	if postgresVersionList != nil {
		var uPrint QuotaPrint
		if nameOk, ok := postgresVersionList.GetQuotaLimitsOk(); ok && nameOk != nil {
			uPrint.QuotaLimits = *nameOk
		}
		if nameOk, ok := postgresVersionList.GetQuotaUsageOk(); ok && nameOk != nil {
			uPrint.QuotaUsage = *nameOk
		}
		o := structs.Map(uPrint)
		out = append(out, o)
	}
	return out
}
