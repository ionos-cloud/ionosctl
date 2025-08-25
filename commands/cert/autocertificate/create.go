package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/cert/provider"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/viper"
)

func AutocertificatePostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create an AutoCertificate. Create a zone with the same name as the common-name parameter and confirm that the zone is enabled",
		Example:   "ionosctl certmanager autocertificate create --name NAME --provider-id PROVIDER --common-name COMMONNAME --key-algorithm rsa2048",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagProviderID, constants.FlagCommonName, constants.FlagKeyAlgorithm)
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := cert.AutoCertificate{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagProviderID); viper.IsSet(fn) {
				input.Provider = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagCommonName); viper.IsSet(fn) {
				input.CommonName = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagKeyAlgorithm); viper.IsSet(fn) {
				input.KeyAlgorithm = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagSubjectAlternativeNames); viper.IsSet(fn) {
				input.SubjectAlternativeNames = viper.GetStringSlice(fn)
			}

			z, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesPost(context.Background()).
				AutoCertificateCreate(cert.AutoCertificateCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating the AutoCertificate: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerAutocertificate, z, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the AutoCertificate")
	cmd.AddStringFlag(constants.FlagProviderID, constants.FlagIdShort, "", "The certificate provider used to issue the AutoCertificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return provider.ProviderIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddStringFlag(constants.FlagCommonName, "", "", "The common name (DNS) of the certificate to issue")
	cmd.AddStringFlag(constants.FlagKeyAlgorithm, "", "", "The key algorithm used to generate the certificate.", core.RequiredFlagOption(), core.WithCompletion(
		func() []string {
			return []string{"rsa2048", "rsa3072", "rsa4096"}
		}, constants.CertApiRegionalURL, constants.CertLocations))
	cmd.AddStringSliceFlag(constants.FlagSubjectAlternativeNames, "", []string{""}, "Optional additional names to be added to the issued certificate")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
