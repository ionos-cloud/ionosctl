package certificate

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

func CertDeleteCmd() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Namespace: "certmanager",
			Resource:  "certificates",
			Verb:      "delete",
			Aliases:   []string{"d"},
			ShortDesc: "Delete a certificate by ID, or all certificates",
			LongDesc: `Delete an uploaded certificate. Pass --certificate-id to delete a single certificate, or --all to delete every certificate in the account.

Make sure no product (ALB, CDN) still references the certificate before deleting it, otherwise those services will lose their TLS material.`,
			Example: "ionosctl certmanager certificate delete " +
				core.FlagsUsage(constants.FlagCertId) + "\n" +
				"ionosctl certmanager certificate delete --all",
			PreCmdRun:  PreCmdDelete,
			CmdRun:     CmdDelete,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(constants.FlagCertId, constants.FlagIdShort, "", "The ID (UUID) of the certificate to delete. Required unless --all is set", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return CertificatesIds()
		}, constants.CertApiRegionalURL, constants.CertLocations),
	)
	cmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete every certificate in the account. Use instead of --certificate-id")

	return cmd
}

func CmdDelete(c *core.CommandConfig) error {
	allFlag, err := c.Command.Command.Flags().GetBool(constants.ArgAll)
	if err != nil {
		return err
	}

	if allFlag {
		return core.DeleteAll(c, core.DeleteAllOptions[cert.CertificateRead]{
			Resource: "Certificate",
			List: func() ([]cert.CertificateRead, error) {
				certs, _, err := client.Must().CertManagerClient.CertificateApi.CertificatesGet(context.Background()).Execute()
				if err != nil {
					return nil, err
				}
				return certs.Items, nil
			},
			Summary: func(z cert.CertificateRead) string {
				return fmt.Sprintf("name: %s, id: %s", z.Properties.Name, z.Id)
			},
			ID: func(z cert.CertificateRead) string {
				return z.Id
			},
			Delete: func(z cert.CertificateRead) error {
				_, err := client.Must().CertManagerClient.CertificateApi.CertificatesDelete(context.Background(), z.Id).Execute()
				return err
			},
		})
	} else {
		id, err := c.Command.Command.Flags().GetString(constants.FlagCertId)
		if err != nil {
			return err
		}

		msg := fmt.Sprintf("delete Certificate ID: %s", id)
		if !confirm.FAsk(c.Command.Command.InOrStdin(), msg, viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		_, err = client.Must().CertManagerClient.CertificateApi.CertificatesDelete(context.Background(), id).Execute()

		return err
	}
}

func PreCmdDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS,
		[]string{constants.FlagCertId},
		[]string{constants.ArgAll},
	)
}
