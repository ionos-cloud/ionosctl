package certmanager

import (
	"context"
	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/certmanager/resources"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"

	"github.com/ionos-cloud/ionosctl/pkg/constants"
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

	certCmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = certCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

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
			r.WaitForState = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest)) // this boolean is duplicated everywhere just to do an append of `& wait` to a verbose message
		}
		if cert != nil {
			r.OutputJSON = cert
			r.KeyValue = getCertRows(cert)                                                              // map header -> rows
			r.Columns = getCertHeaders(viper.GetStringSlice(core.GetFlagName(c.NS, constants.ArgCols))) // headers
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

func CertificatesIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	certSvc := resources.NewCertsService(client, context.Background())
	certs, _, err := certSvc.List()
	clierror.CheckError(err, outErr)
	certIds := make([]string, 0)
	if items, ok := certs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				certIds = append(certIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return certIds
}

// Get Client for Completion Functions
func getClient() (*resources.Client, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}
