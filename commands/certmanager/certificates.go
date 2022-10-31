package certmanager

import (
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-cert-manager"
)

func CertCmd() *core.Command {
	certCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certificate-manager",
			Short:            "Certificate Manager Operations",
			Long:             "The sub-commands of `ionosctl certificate-manager` allow you to manage the SSL Certificates under your account.",
			TraverseChildren: true,
		},
	}
	certCmd.AddCommand(CertGetCmd())
	certCmd.AddCommand(CertPostCmd())
	certCmd.AddCommand(CertListCmd())
	return certCmd
}

func getCertPrint(resp *ionoscloud.APIResponse, c *core.CommandConfig, dcs *[]ionoscloud.CertificateDto) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForState)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if dcs != nil {
			r.OutputJSON = dcs
			r.KeyValue = getCertRows(dcs)                                                            // map header -> rows
			r.Columns = getCertHeaders(viper.GetStringSlice(core.GetFlagName(c.NS, config.ArgCols))) // headers
		}
	}
	return r
}

type CertPrint struct {
	CertId      string `json:"CertId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
	State       string `json:"State,omitempty"`
}

func getCertRows(certs *[]ionoscloud.CertificateDto) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*certs))
	for _, cert := range *certs {
		var certPrint CertPrint
		if idOk, ok := cert.GetIdOk(); ok && idOk != nil {
			certPrint.CertId = *idOk
		}
		if idOk, ok := cert.GetIdOk(); ok && idOk != nil {
			certPrint.CertId = *idOk
		}
		if propertiesOk, ok := cert.GetPropertiesOk(); ok && propertiesOk != nil {
			if displayNameOk, ok := propertiesOk.GetNameOk(); ok && displayNameOk != nil {
				certPrint.DisplayName = *displayNameOk
			}
			// if displayNameOk, ok := propertiesOk.get; ok && displayNameOk != nil {
			// 	CertPrint.DisplayName = *displayNameOk
			// }
		}
		if metadataOk, ok := cert.GetMetadataOk(); ok && metadataOk != nil {
			if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
				certPrint.State = string(*stateOk)
			}
		}
		o := structs.Map(certPrint)
		out = append(out, o)
	}
	return out
}

var allCols = structs.Names(CertPrint{})

func getCertHeaders(customColumns []string) []string {
	if customColumns == nil {
		return allCols[0:2]
	}
	//for _, c := customColumns {
	//	if slices.Contains(allCols, c) {}
	//}
	return customColumns
}
