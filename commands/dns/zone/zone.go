package zone

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	dns "github.com/ionos-cloud/sdk-go-dns"
	"github.com/spf13/cobra"
)

func ZoneCommand() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "zone",
			Short:            "DNS zones",
			Aliases:          []string{"z", "zones"},
			Long:             "The sub-commands of `ionosctl dns zone` allow you to perform operations on DNS zones",
			TraverseChildren: true,
		},
	}

	cmd.Command.PersistentFlags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allCols, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.Command.PersistentFlags().Bool(constants.ArgNoHeaders, false, "When using text output, don't print headers")

	cmd.AddCommand(ZonesGetCmd())
	cmd.AddCommand(ZonesDeleteCmd())
	cmd.AddCommand(ZonesPostCmd())
	cmd.AddCommand(ZonesPutCmd())
	cmd.AddCommand(ZonesFindByIdCmd())

	return cmd
}

// Helper functions for printing zone

func getZonesPrint(c *core.CommandConfig, data dns.ZoneReadList) printer.Result {
	r := printer.Result{}
	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	if c != nil {
		// TODO for r.OutputJSON: This loses all kinds of information in `-o json`, like `limit`, `offset`, etc. See https://github.com/ionos-cloud/ionosctl/issues/249
		// But we are forced to do this otherwise we'd have this JSON output:
		// {
		//  "items": {
		//    "items": [
		// ...
		r.OutputJSON = data.Items // TODO: See above comment. Remove `.Items` once JSON marshalling works as one would expect
		r.KeyValue = makeZonePrintObj(*data.Items...)
		r.Columns = printer.GetHeadersAllDefault(allCols, cols)
	}
	return r
}

func getZonePrint(c *core.CommandConfig, data dns.ZoneRead) printer.Result {
	return getZonesPrint(c, dns.ZoneReadList{Items: &[]dns.ZoneRead{data}})
}

type zonePrint struct {
	Id          string `json:"ID,omitempty"`
	Name        string `json:"Name,omitempty"`
	Description string `json:"Content,omitempty"`
	NameServers string `json:"NameServers,omitempty"`
	Enabled     bool   `json:"Enabled,omitempty"`
	State       string `json:"State,omitempty"`
}

var allCols = structs.Names(zonePrint{})

func makeZonePrintObj(data ...dns.ZoneRead) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(data))

	for _, item := range data {
		var printObj zonePrint
		printObj.Id = *item.GetId()

		if p, ok := item.GetPropertiesOk(); ok {
			printObj.Enabled = *p.Enabled
			printObj.Description = *p.Description
			printObj.Name = *p.ZoneName
		}
		if m, ok := item.GetMetadataOk(); ok && m != nil {
			printObj.State = string(*m.State)
			printObj.NameServers = strings.Join(*m.Nameservers, ", ")
		}

		o := structs.Map(printObj)
		out = append(out, o)
	}
	return out
}

func ZoneNames() []string {
	ls, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), func(t dns.ZoneRead) string {
		return *t.Properties.ZoneName
	})
}

func Zones(f func(dns.ZoneRead) string) []string {
	ls, _, err := client.Must().DnsClient.ZonesApi.ZonesGet(context.Background()).Execute()
	if err != nil {
		return nil
	}
	return functional.Map(*ls.GetItems(), f)
}

// Resolve resolves nameOrId (the name of a zone, or the ID of a zone) - to the ID of the zone.
// If it's an ID, it's returned as is. If it's not, then it's a name, and we try to resolve it
func Resolve(nameOrId string) (string, error) {
	uid, errParseUuid := uuid.Parse(nameOrId)
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
