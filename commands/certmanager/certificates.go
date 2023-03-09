package certmanager

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/services/certmanager/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-cert-manager"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	certCmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

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
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		if resp != nil {
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if cert != nil {
			r.OutputJSON = cert
			r.KeyValue = getCertRows(cert)                          // map header -> rows
			r.Columns = printer.GetHeadersAllDefault(allCols, cols) // headers
		}
	}
	return r
}

func printProperties(value *ionoscloud.CertificateDto, flag bool) string {
	if !flag {
		return *value.Properties.Certificate
	}
	return *value.Properties.CertificateChain
}

type CertPrint struct {
	CertId      string `json:"CertId,omitempty"`
	DisplayName string `json:"DisplayName,omitempty"`
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
		}
		o := structs.Map(certPrint)
		out = append(out, o)
	}
	return out
}

var allCols = structs.Names(CertPrint{})

func CertificatesIds() []string {
	client, _ := config.GetClient()
	svc := resources.NewCertsService(client, context.Background())
	certs, _, _ := svc.List()
	return functional.Map(*certs.GetItems(), func(dto ionoscloud.CertificateDto) string {
		return *dto.GetId()
	})
}
