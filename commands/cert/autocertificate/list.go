package autocertificate

import (
	"context"

	cert "github.com/ionos-cloud/sdk-go-bundle/products/cert/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func AutocertificateListCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "autocertificate",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "Retrieve AutoCertificate list",
		Example:   "ionosctl certmanager autocertificate list",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				certClient := cert.NewAPIClient(cfg)
				req := certClient.AutoCertificateApi.AutoCertificatesGet(context.Background())

				if fn := core.GetFlagName(c.NS, constants.FlagCommonName); viper.IsSet(fn) {
					req = req.FilterCommonName(viper.GetString(fn))
				}

				ls, _, err := req.Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCommonName, "", "", "Filter by the common name (DNS)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
