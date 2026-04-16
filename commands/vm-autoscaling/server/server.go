package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling/group"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
)

// allCols defines server columns. The first column is the autoscaling group server ID.
// Remaining columns come from the CloudAPI server lookup and are populated via
// enrichAutoscalingServer.
var allCols = []table.Column{
	{Name: "GroupServerId", JSONPath: "id", Default: true},
	{Name: "ServerId", JSONPath: "cloudapi.id", Default: true},
	{Name: "DatacenterId", Default: true, Format: func(item map[string]any) any {
		return table.Navigate(item, "cloudapi.href")
	}},
	{Name: "Name", JSONPath: "cloudapi.properties.name", Default: true},
	{Name: "AvailabilityZone", JSONPath: "cloudapi.properties.availabilityZone", Default: true},
	{Name: "Cores", JSONPath: "cloudapi.properties.cores", Default: true},
	{Name: "RAM", JSONPath: "cloudapi.properties.ram", Default: true},
	{Name: "CpuFamily", JSONPath: "cloudapi.properties.cpuFamily", Default: true},
	{Name: "VmState", JSONPath: "cloudapi.properties.vmState", Default: true},
	{Name: "State", JSONPath: "cloudapi.metadata.state", Default: true},
	{Name: "TemplateId", JSONPath: "cloudapi.properties.templateUuid", Default: true},
	{Name: "Type", JSONPath: "cloudapi.properties.type", Default: true},
	{Name: "BootCdromId", JSONPath: "cloudapi.properties.bootCdrom.id", Default: true},
	{Name: "BootVolumeId", JSONPath: "cloudapi.properties.bootVolume.id", Default: true},
	{Name: "NicMultiQueue", JSONPath: "cloudapi.properties.nicMultiQueue", Default: true},
}

func Root() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "server",
			Aliases:          []string{"s", "sv", "vm", "vms", "servers"},
			Short:            "Autoscaling Servers Operations",
			Long:             "The sub-commands of `ionosctl autoscaling server` allow you to manage the Autoscaling Servers under your account.",
			TraverseChildren: true,
		},
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Get())

	cmd.AddColsFlag(allCols)

	return cmd
}

// enrichAutoscalingServer takes an autoscaling Server, looks up the corresponding
// CloudAPI server, and returns a merged map with the autoscaling server data at the
// top level and CloudAPI server data nested under "cloudapi".
func enrichAutoscalingServer(sv vmasc.Server) (map[string]any, error) {
	if sv.Properties == nil || sv.Properties.DatacenterServer == nil ||
		sv.Properties.DatacenterServer.Id == nil || sv.Properties.DatacenterServer.Href == nil {
		return nil, fmt.Errorf("server properties are incomplete: %+v", sv)
	}

	hrefFields := strings.FieldsFunc(*sv.Properties.DatacenterServer.Href, func(r rune) bool { return r == '/' })
	dcId := hrefFields[len(hrefFields)-3]
	cloudApiId := *sv.Properties.DatacenterServer.Id

	cloudApiServer, _, err := client.Must().CloudClient.ServersApi.DatacentersServersFindById(
		context.Background(), dcId, cloudApiId).Execute()
	if err != nil {
		return nil, fmt.Errorf("could not find server %s in datacenter %s via CloudAPI: %w", cloudApiId, dcId, err)
	}

	// Build a merged structure: autoscaling server ID at top level, CloudAPI data nested
	return map[string]any{
		"id":       sv.Id,
		"cloudapi": cloudApiServer,
	}, nil
}

// enrichAutoscalingServers enriches all servers in a collection via CloudAPI lookups.
func enrichAutoscalingServers(sc vmasc.ServerCollection) ([]map[string]any, error) {
	if sc.Items == nil {
		return nil, fmt.Errorf("could not retrieve items")
	}

	var result []map[string]any
	for _, sv := range *sc.Items {
		enriched, err := enrichAutoscalingServer(sv)
		if err != nil {
			return nil, err
		}
		result = append(result, enriched)
	}
	return result, nil
}

func Servers(fs ...Filter) (vmasc.ServerCollection, error) {
	groupIds := group.GroupsProperty(func(r vmasc.Group) string {
		if r.Id == nil {
			return ""
		}
		return *r.Id
	})

	var ls vmasc.ServerCollection
	ls.Items = pointer.From(make([]vmasc.Server, 0))
	for _, groupId := range groupIds {
		actions, err := GroupServers(groupId, fs...)
		if err != nil {
			return vmasc.ServerCollection{}, err
		}
		ls.Items = pointer.From(append(*ls.Items, *actions.Items...))
	}

	return ls, nil
}

func GroupServers(groupId string, fs ...Filter) (vmasc.ServerCollection, error) {
	req := client.Must().VMAscClient.GroupsServersGet(context.Background(), groupId)

	for _, f := range fs {
		var err error
		req, err = f(req)
		if err != nil {
			return vmasc.ServerCollection{}, err
		}
	}

	ls, _, err := req.Execute()
	if err != nil {
		return vmasc.ServerCollection{}, err
	}
	return ls, nil
}

func ServersProperty[V any](f func(resource vmasc.Server) V, fs ...Filter) []V {
	recs, err := Servers(fs...)
	if err != nil {
		return nil
	}
	return functional.Map(*recs.Items, f)
}

type Filter func(request vmasc.ApiGroupsServersGetRequest) (vmasc.ApiGroupsServersGetRequest, error)
