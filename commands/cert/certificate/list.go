package certificate

import (
	"context"

	cert "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func CertListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace:  "certmanager",
		Resource:   "certificates",
		Verb:       "list",
		Aliases:    []string{"l"},
		ShortDesc:  "List all Certificates",
		LongDesc:   "Use this command to retrieve all Certificates.",
		Example:    "ionosctl certificate-manager list",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdList,
		InitClient: true,
	})

	return cmd
}

func CmdList(c *core.CommandConfig) error {
	c.Verbose("Getting Certificates...")

	return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
		certClient := cert.NewAPIClient(cfg)
		ls, _, err := certClient.CertificateApi.CertificatesGet(context.Background()).Execute()
		return ls, err
	})
}
