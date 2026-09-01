package provider

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
		ShortDesc: "Register an ACME certificate provider",
		LongDesc: `Register an ACME certificate authority (for example Let's Encrypt) that auto-certificates can then use to be issued and auto-renewed.

Required: --name (a label), --email (the ACME account contact address), and --server (the CA's ACME directory URL, e.g. https://acme-v02.api.letsencrypt.org/directory).

External account binding (EAB) is optional and links your IONOS account to a pre-registered account at the CA. If you use it, --key-id and --key-secret must be supplied together.`,
		Example: `# Register Let's Encrypt (production directory)
ionosctl certmanager provider create --name letsencrypt --email admin@example.com --server https://acme-v02.api.letsencrypt.org/directory

# Register a CA that requires external account binding (EAB)
ionosctl certmanager provider create --name my-ca --email admin@example.com --server https://acme.my-ca.example/directory --key-id my-eab-key-id --key-secret my-eab-secret`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagEmail, constants.FlagServer)
			if err != nil {
				return err
			}

			c.Command.Command.MarkFlagsRequiredTogether(constants.FlagKeyId, constants.FlagKeySecret)

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := cert.Provider{}

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
				if input.ExternalAccountBinding == nil {
					input.ExternalAccountBinding = &cert.ProviderExternalAccountBinding{}
				}
				input.ExternalAccountBinding.KeyId = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagKeySecret); viper.IsSet(fn) {
				if input.ExternalAccountBinding == nil {
					input.ExternalAccountBinding = &cert.ProviderExternalAccountBinding{}
				}
				input.ExternalAccountBinding.KeySecret = pointer.From(viper.GetString(fn))
			}

			z, _, err := client.Must().CertManagerClient.ProviderApi.ProvidersPost(context.Background()).
				ProviderCreate(cert.ProviderCreate{
					Properties: input,
				}).Execute()

			if err != nil {
				return fmt.Errorf("failed creating the Provider: %w", err)
			}

			return c.Printer(allCols).Print(z)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "A label to identify this provider (e.g. \"Let's Encrypt\"). Required")
	cmd.AddStringFlag(constants.FlagEmail, "", "", "Contact email registered with the ACME account; the CA uses it for expiry and policy notices. Required")
	cmd.AddStringFlag(constants.FlagServer, "", "", "The CA's ACME directory URL, e.g. https://acme-v02.api.letsencrypt.org/directory. Required")
	cmd.AddStringFlag(constants.FlagKeyId, "", "", "External account binding (EAB) key ID issued by the CA. Optional; must be given together with --key-secret")
	cmd.AddStringFlag(constants.FlagKeySecret, "", "", "External account binding (EAB) HMAC key secret issued by the CA. Optional; must be given together with --key-id. Write-only")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
