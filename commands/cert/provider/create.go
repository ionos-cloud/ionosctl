package provider

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/viper"
)

func ProviderPostCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "provider",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create an Provider",
		Example:   "ionosctl certmanager provider create --name NAME --email EMAIL --server SERVER",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagEmail, constants.FlagServer)
			if err != nil {
				return err
			}

			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagKeyId, constants.FlagKeySecret)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {

			input := cert.Provider{ExternalAccountBinding: cert.NewProviderExternalAccountBinding()}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.Name = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagEmail); viper.IsSet(fn) {
				input.Email = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagServer); viper.IsSet(fn) {
				input.Server = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagKeyId); viper.IsSet(fn) {
				input.ExternalAccountBinding.KeyId = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagKeySecret); viper.IsSet(fn) {
				input.ExternalAccountBinding.KeySecret = pointer.From(viper.GetString(fn))
			}

			z, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersPost(context.Background()).
				ProviderCreate(cert.ProviderCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating the Provider: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.CertManagerProvider, z, tabheaders.GetHeadersAllDefault(allCols, cols))
			if err != nil {
				return fmt.Errorf("failed generating the output: %w", err)
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the certificate Provider")
	cmd.AddStringFlag(constants.FlagEmail, "", "", "The email address of the certificate requester")
	cmd.AddStringFlag(constants.FlagServer, "", "", "The URL of the certificate Provider")
	cmd.AddStringFlag(constants.FlagKeyId, "", "", "The key ID of the external account binding")
	cmd.AddStringFlag(constants.FlagKeySecret, "", "", "The key secret of the external account binding")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
