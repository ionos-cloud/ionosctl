package certmanager

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func CertGetApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "api-version",
		Aliases:    []string{"api", "info"},
		ShortDesc:  "Get Certificate Manager API Version",
		LongDesc:   "Use this command to retrieve API Version.",
		Example:    "ionosctl certificate-manager api-version",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdGetApiVersion,
		InitClient: true,
	})

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allAPICols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddBoolFlag(constants.ArgNoHeaders, "n", false, "Response delete all certificates")

	return cmd
}

func CmdGetApiVersion(c *core.CommandConfig) error {
	c.Printer.Verbose("Getting Api Version...")
	headers, err := c.Command.Command.Flags().GetBool(constants.ArgNoHeaders)
	if err != nil {
		return err
	}
	APIVersion, _, err := c.CertificateManagerServices.Certs().GetApiVersion()
	if err != nil {
		return err
	}
	return c.Printer.Print(getApiPrint(nil, c, &[]sdkgo.ApiInfoDto{APIVersion}, headers))
}

func getApiPrint(resp *sdkgo.APIResponse, c *core.CommandConfig, cert *[]sdkgo.ApiInfoDto, headers bool) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if cert != nil {
			r.OutputJSON = cert
			r.KeyValue = getApiRows(cert)                              // map header -> rows
			r.Columns = printer.GetHeadersAllDefault(allAPICols, cols) // headers
		}
	}
	return r
}

func getApiRows(apis *[]sdkgo.ApiInfoDto) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*apis))
	for _, api := range *apis {
		var apiPrint ApiPrint
		if nameOk, ok := api.GetNameOk(); ok && nameOk != nil {
			apiPrint.Name = *nameOk
		}
		if hrefOk, ok := api.GetHrefOk(); ok && hrefOk != nil {
			apiPrint.Href = *hrefOk
		}
		if versionOk, ok := api.GetVersionOk(); ok && versionOk != nil {
			apiPrint.Version = *versionOk
		}
		o := structs.Map(apiPrint)
		out = append(out, o)
	}
	return out
}

var allAPICols = structs.Names(ApiPrint{})

func getAPIHeaders(customColumns []string) []string {
	if customColumns == nil {
		return allAPICols[0:3]
	}

	return customColumns
}
