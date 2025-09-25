package autocertificate

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/completions"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/spf13/cobra"
)

var (
	allCols = []string{"Id", "Provider", "CommonName", "KeyAlgorithm", "Name", "AlternativeNames", "State"}
)

func AutocertificateCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "autocertificate",
			Aliases:          []string{"autocert", "autocerts", "auto", "autocertificates"},
			Short:            "Auto certificates create new certificates based on a certificate provider",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddCommand(AutocertificatePostCmd())
	cmd.AddCommand(AutocertificateDeleteCmd())
	cmd.AddCommand(AutocertificateFindByIdCmd())
	cmd.AddCommand(AutocertificateListCmd())
	cmd.AddCommand(AutocertificatePutCmd())

	return cmd
}

func AutocertificateIDs() []string {
	autocertificate, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	autocertificateConverted, err := json2table.ConvertJSONToTable("items", jsonpaths.CertManagerAutocertificate, autocertificate)
	if err != nil {
		return nil
	}
	return completions.NewCompleter(autocertificateConverted, "Id").AddInfo("Name").AddInfo("CommonName").AddInfo("KeyAlgorithm").AddInfo("Provider").AddInfo("AlternativeNames").AddInfo("State").ToString()
}
