package zone

import (
	"context"
	"fmt"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/tabheaders"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	allZoneJSONPaths = map[string]string{
		"Id":          "id",
		"Name":        "properties.zoneName",
		"Description": "properties.description",
		"NameServers": "metadata.nameServers",
		"Enabled":     "properties.enabled",
		"State":       "metadata.state",
	}

	allCols = []string{"Id", "Name", "Description", "NameServers", "Enabled", "State"}
)

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "zone",
			Aliases:          []string{"z", "zones"},
			Short:            "The sub-commands of `ionosctl dns zone` allow you to manage DNS zones. A DNS zone serves as an authoritative source of information about which IP addresses belong to which domains",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, tabheaders.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())

	return cmd
}

// Zones returns all zones matching the given filters
func Zones(fs ...Filter) (dns.ZoneReadList, error) {
	// Hack to enforce the dns-level flag default for API URL on the completions too
	if url := config.GetServerUrl(); url == constants.DefaultApiURL {
		viper.Set(constants.ArgServerUrl, "")
	}

	req := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background())

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return dns.ZoneReadList{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return dns.ZoneReadList{}, err
	}
	return ls, nil
}

func ZonesProperty[V any](f func(dns.ZoneRead) V, fs ...Filter) []V {
	recs, err := Zones(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request dns.ApiZonesGetRequest) (dns.ApiZonesGetRequest, error)

// Resolve resolves nameOrId (the name of a zone, or the ID of a zone) - to the ID of the zone.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func Resolve(nameOrId string) (string, error) {
	uid, errParseUuid := uuid.FromString(nameOrId)
	zId := uid.String()
	if errParseUuid != nil {
		// nameOrId is a name
		ls, _, errFindZoneByName := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).FilterZoneName(nameOrId).Limit(1).Execute()
		if errFindZoneByName != nil {
			return "", fmt.Errorf("failed finding a zone by name: %w", errFindZoneByName)
		}
		if len(*ls.Items) < 1 {
			return "", fmt.Errorf("could not find zone by name %s: got %d zones", nameOrId, len(*ls.Items))
		}
		zId = *(*ls.Items)[0].Id
	}
	return zId, nil
}
