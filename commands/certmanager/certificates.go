package certmanager

import (
	"github.com/fatih/structs"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go-cert-manager"
)

func CertCmd() *core.Command {
	certCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "certificate-manager",
			Short:            "Certificate Manager Operations",
			Long:             "The sub-commands of `ionosctl certificate-manager` allows you to manage the SSL Certificates under your account.",
			TraverseChildren: true,
		},
	}
	certCmd.AddCommand(CertGetCmd())
	certCmd.AddCommand(CertCreateCmd())
	certCmd.AddCommand(CertListCmd())
	certCmd.AddCommand(CertDeleteCmd())
	certCmd.AddCommand(CertGetApiVersionCmd())
	certCmd.AddCommand(CertUpdateCmd())
	return certCmd
}

func getCertPrint(resp *ionoscloud.APIResponse, c *core.CommandConfig, cert *[]ionoscloud.CertificateDto) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, FlagArgWaitForState)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if cert != nil {
			r.OutputJSON = cert
			r.KeyValue = getCertRows(cert)                                                        // map header -> rows
			r.Columns = getCertHeaders(viper.GetStringSlice(core.GetFlagName(c.NS, FlagArgCols))) // headers
		}
	}
	return r
}

func printProperties(value *ionoscloud.CertificateDto, c *core.CommandConfig, flags []bool) string {
	var printString string
	if flags[0] {
		printString += *value.Properties.Certificate
	}
	if flags[1] {
		printString += *value.Properties.CertificateChain
	}
	return printString
}

type CertPrint struct {
	CertId      string `json:"CertId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
	State       string `json:"State,omitempty"`
}

type ApiPrint struct {
	Name    string `json:"Name,omitempty"`
	Href    string `json:"Href,omitempty"`
	Version string `json:"Version,omitempty"`
}

func getCertRows(certs *[]ionoscloud.CertificateDto) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(*certs))
	for _, cert := range *certs {
		var certPrint CertPrint
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
