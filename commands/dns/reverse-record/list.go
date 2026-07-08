package reverse_record

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	dns "github.com/ionos-cloud/sdk-go-bundle/products/dns/v2"
	"github.com/ionos-cloud/sdk-go-bundle/shared"
	"github.com/spf13/viper"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "list",
		Aliases:   []string{"ls", "l"},
		ShortDesc: "Retrieve all reverse records",
		Example:   "ionosctl dns rr list",
		CmdRun: func(c *core.CommandConfig) error {
			return c.ListAllLocations(allCols, func(cfg *shared.Configuration) (any, error) {
				dnsClient := dns.NewAPIClient(cfg)
				req := dnsClient.ReverseRecordsApi.ReverserecordsGet(context.Background())

				if fn := core.GetFlagName(c.NS, constants.FlagIps); viper.IsSet(fn) {
					req = req.FilterRecordIp(viper.GetStringSlice(fn))
				}

				ls, _, err := req.Execute()
				return ls, err
			})
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagIps, "i", "", "Optional filter for the IP address of the reverse record",
		core.WithCompletion(func() []string {
			return RecordsProperty(func(t dns.ReverseRecordRead) string {
				return t.Properties.Ip
			})
		}, constants.DNSApiRegionalURL, constants.DNSLocations),
	)

	return cmd
}
