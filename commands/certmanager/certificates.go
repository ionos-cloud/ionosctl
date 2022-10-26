package certmanager

import (
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	return certCmd
}

func getCertPrint(resp *ionoscloud.APIResponse, c *core.CommandConfig, dcs *ionoscloud.CertificateDto) printer.Result {
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

func getCertRows(Cert *ionoscloud.CertificateDto) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, 1)

	var CertPrint CertPrint
	if idOk, ok := Cert.GetIdOk(); ok && idOk != nil {
		CertPrint.CertId = *idOk
	}
	if propertiesOk, ok := Cert.GetPropertiesOk(); ok && propertiesOk != nil {
		if displayNameOk, ok := propertiesOk.GetNameOk(); ok && displayNameOk != nil {
			CertPrint.DisplayName = *displayNameOk
		}
		// if displayNameOk, ok := propertiesOk.get; ok && displayNameOk != nil {
		// 	CertPrint.DisplayName = *displayNameOk
		// }
	}
	if metadataOk, ok := Cert.GetMetadataOk(); ok && metadataOk != nil {
		if stateOk, ok := metadataOk.GetStateOk(); ok && stateOk != nil {
			CertPrint.State = string(*stateOk)
		}
	}
	o := structs.Map(CertPrint)
	out = append(out, o)
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
