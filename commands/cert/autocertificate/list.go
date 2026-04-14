package autocertificate

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
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
			req := client.Must().CertManagerClient.AutoCertificateApi.AutoCertificatesGet(context.Background())

			if fn := core.GetFlagName(c.NS, constants.FlagCommonName); viper.IsSet(fn) {
				req = req.FilterCommonName(viper.GetString(fn))
			}

			ls, _, err := req.Execute()
			if err != nil {
				return fmt.Errorf("failed listing the AutoCertificates: %w", err)
			}

			return c.Printer(allCols).Prefix("items").Print(ls)
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagCommonName, "", "", "Filter by the common name (DNS)")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}
