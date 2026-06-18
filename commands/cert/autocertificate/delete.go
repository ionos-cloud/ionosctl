package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/spf13/viper"
)

func AutocertificateDeleteCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "delete",
		Aliases:   []string{"del", "d"},
		ShortDesc: "Delete an AutoCertificate",
		Example:   "ionosctl certmanager autocertificate delete --autocertificate-id ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{constants.ArgAll}, []string{constants.FlagAutocertificateID})
		},
		CmdRun: func(c *core.CommandConfig) error {
			if all := viper.GetBool(core.GetFlagName(c.NS, constants.ArgAll)); all {
				return deleteAll(c)
			}

			autocertificateId := viper.GetString(core.GetFlagName(c.NS, constants.FlagAutocertificateID))
			z, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesFindById(context.Background(), autocertificateId).Execute()
			if err != nil {
				return fmt.Errorf("failed getting AutoCertificate by id %s: %w", autocertificateId, err)
			}
			yes := confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Are you sure you want to delete AutoCertificate with name: %s, id: %s ", z.Properties.Name, z.Id),
				viper.GetBool(constants.ArgForce))
			if !yes {
				return fmt.Errorf(confirm.UserDenied)
			}

			_, err = client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesDelete(context.Background(), autocertificateId).Execute()
			if err != nil {
				return fmt.Errorf("failed to delete the AutoCertificate: %w", err)
			}

			return err
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagAutocertificateID, constants.FlagIdShort, "", "Provide the specified AutoCertificate", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return AutocertificateIDs()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)

	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all AutoCertificates")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

func deleteAll(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[cert.AutoCertificateRead]{
		Resource: "AutoCertificate",
		List: func() ([]cert.AutoCertificateRead, error) {
			xs, _, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesGet(context.Background()).Execute()
			if err != nil {
				return nil, fmt.Errorf("failed getting the AutoCertificates: %w", err)
			}
			return xs.GetItems(), nil
		},
		Summary: func(z cert.AutoCertificateRead) string {
			return fmt.Sprintf("name: %s, id: %s", z.Properties.Name, z.Id)
		},
		ID: func(z cert.AutoCertificateRead) string {
			return z.Id
		},
		Delete: func(z cert.AutoCertificateRead) error {
			_, err := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesDelete(context.Background(), z.Id).Execute()
			return err
		},
	})
}
