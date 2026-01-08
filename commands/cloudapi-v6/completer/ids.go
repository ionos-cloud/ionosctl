package completer

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
)

func BackupUnitsIds() []string {
	backupUnitSvc := resources.NewBackupUnitService(client.Must(), context.Background())
	backupUnits, _, err := backupUnitSvc.List()
	if err != nil {
		return nil
	}
	backupUnitsIds := make([]string, 0)
	if items, ok := backupUnits.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				backupUnitsIds = append(backupUnitsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return backupUnitsIds
}

func AttachedCdromsIds(datacenterId, serverId string) []string {
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId)
	if err != nil {
		return nil
	}
	attachedCdromsIds := make([]string, 0)
	if items, ok := cdroms.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedCdromsIds = append(attachedCdromsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedCdromsIds
}

func DatacenterIdsFilterLocation(loc string) []string {
	req := client.Must().CloudClient.DataCentersApi.DatacentersGet(context.Background())
	if loc != "" {
		req = req.Filter("location", loc)
	}
	dcs, _, err := req.Execute()
	if err != nil {
		return nil
	}

	var dcIds []string
	for _, item := range *dcs.Items {
		var completion string
		if item.Id == nil {
			continue
		}
		completion = *item.Id
		if props, ok := item.GetPropertiesOk(); ok {
			if name, ok := props.GetNameOk(); ok {
				completion = fmt.Sprintf("%s\t%s", completion, *name)
			}
			if location, ok := props.GetLocationOk(); ok {
				completion = fmt.Sprintf("%s - %s", completion, *location)
			}
		}

		dcIds = append(dcIds, completion)
	}
	return dcIds
}

func DataCentersIds(filters ...func(datacenter ionoscloud.Datacenter) bool) []string {
	datacenterSvc := resources.NewDataCenterService(client.Must(), context.Background())
	datacenters, _, err := datacenterSvc.List()
	if err != nil {
		return nil
	}
	dcIds := make([]string, 0)
	if items, ok := datacenters.GetItemsOk(); ok {
		for _, item := range *items {
			var completion string
			if item.Id == nil {
				continue
			}

			skip := false
			for _, filter := range filters {
				if !filter(item) {
					skip = true
					break
				}
			}

			if skip {
				continue
			}

			completion = *item.Id
			if props, ok := item.GetPropertiesOk(); ok {
				if name, ok := props.GetNameOk(); ok {
					// Here is where the completion descriptions start
					completion = fmt.Sprintf("%s\t%s", completion, *name)
				}
				if location, ok := props.GetLocationOk(); ok {
					completion = fmt.Sprintf("%s - %s", completion, *location)
				}
			}

			dcIds = append(dcIds, completion)
		}
	} else {
		return nil
	}
	return dcIds
}

func FirewallRulesIds(datacenterId, serverId, nicId string) []string {
	firewallRuleSvc := resources.NewFirewallRuleService(client.Must(), context.Background())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId)
	if err != nil {
		return nil
	}
	firewallRulesIds := make([]string, 0)
	if items, ok := firewallRules.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				firewallRulesIds = append(firewallRulesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return firewallRulesIds
}

func FlowLogsIds(datacenterId, serverId, nicId string) []string {
	flowLogSvc := resources.NewFlowLogService(client.Must(), context.Background())
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId)
	if err != nil {
		return nil
	}
	flowLogsIds := make([]string, 0)
	if items, ok := flowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				flowLogsIds = append(flowLogsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return flowLogsIds
}

func GroupsIds() []string {
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	groups, _, err := groupSvc.List()
	if err != nil {
		return nil
	}

	items, ok := groups.GetItemsOk()
	if !ok || items == nil || len(*items) == 0 {
		return nil
	}

	out := make([]string, 0, len(*items))
	var b strings.Builder

	for _, item := range *items {
		b.Reset()

		// id + tab
		if id, ok := item.GetIdOk(); ok && id != nil && *id != "" {
			b.WriteString(*id)
			b.WriteByte('\t')
		} else {
			continue
		}

		// name
		name := ""
		if item.Properties != nil && item.Properties.Name != nil {
			name = *item.Properties.Name
		}
		if name != "" {
			b.WriteString(name)
		} else {
			b.WriteString("(no name)")
		}

		// users count, if available
		countText := ""
		if item.Entities != nil && item.Entities.Users != nil {
			if groupMembers, ok := item.Entities.GetUsersOk(); ok && groupMembers != nil {
				usersItems, ok := groupMembers.GetItemsOk()
				if !ok || usersItems == nil {
					continue
				}
				countText = fmt.Sprintf("%d", len(*usersItems))
			}
		}
		if countText != "" && countText != "0" {
			b.WriteString("; users: ")
			b.WriteString(countText)
		}

		out = append(out, b.String())
	}

	return out
}

func ImageIds(customFilters ...func(ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest) []string {
	req := client.Must().CloudClient.ImagesApi.ImagesGet(context.Background()).
		Depth(1).
		OrderBy("public")

	for _, cf := range customFilters {
		req = cf(req)
	}

	ls, _, err := req.Execute()
	if err != nil || ls.Items == nil {
		return nil
	}

	var completions []string
	for _, image := range *ls.Items {
		completion := *image.Id + "\t"

		if props := image.Properties; props == nil {
			continue
		}

		if license := image.Properties.LicenceType; license != nil {
			completion = fmt.Sprintf("%s %s", completion, *license)
		}

		if imgType := image.Properties.ImageType; imgType != nil {
			completion = fmt.Sprintf("%s %s", completion, *imgType)
		}

		if public := image.Properties.Public; public != nil {
			if *public {
				completion = fmt.Sprintf("%s public", completion)
			} else {
				completion = fmt.Sprintf("%s private", completion)
			}
		}

		if aliases := image.Properties.ImageAliases; aliases != nil && len(*aliases) > 0 && (*aliases)[0] != "" {
			completion = fmt.Sprintf("%s [%s]", completion, strings.Join(*aliases, ","))
		}

		if name := image.Properties.Name; name != nil {
			completion = fmt.Sprintf("%s (%s)", completion, *name)
		}

		if loc := image.Properties.Location; loc != nil {
			completion = fmt.Sprintf("%s from %s", completion, *loc)
		}

		completions = append(completions, completion)
	}
	return completions
}

func IpBlocksIds() []string {
	ipBlockSvc := resources.NewIpBlockService(client.Must(), context.Background())
	ipBlocks, _, err := ipBlockSvc.List()
	if err != nil {
		return nil
	}

	items, ok := ipBlocks.GetItemsOk()
	if !ok || items == nil || len(*items) == 0 {
		return nil
	}

	out := make([]string, 0, len(*items))
	var b strings.Builder

	for _, item := range *items {
		b.Reset()

		// id + tab
		if id, ok := item.GetIdOk(); ok && id != nil && *id != "" {
			b.WriteString(*id)
			b.WriteByte('\t')
		} else {
			continue
		}

		// name
		name := ""
		if item.Properties != nil && item.Properties.Name != nil {
			name = *item.Properties.Name
		}
		if name != "" {
			b.WriteString(name)
		} else {
			b.WriteString("(no name)")
		}

		// location
		location := ""
		if item.Properties != nil && item.Properties.Location != nil {
			location = *item.Properties.Location
		}
		if location != "" {
			b.WriteString("; location: ")
			b.WriteString(location)
		}

		// ips list/count
		if item.Properties != nil && item.Properties.Ips != nil {
			ips := *item.Properties.Ips
			if len(ips) > 0 {
				b.WriteString("; ips: [")
				b.WriteString(strings.Join(ips, ", "))
				b.WriteString("]")
			}
		} else {
			b.WriteString("; ips: none")
		}

		out = append(out, b.String())
	}

	return out
}

func K8sClustersIds() []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListClusters()
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}

func K8sVersionsIds() []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListVersions()
	if err != nil {
		return nil
	}
	return k8ss
}

func K8sNodesIds(clusterId, nodepoolId string) []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId)
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}

func K8sNodePoolsIds(clusterId string) []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId)
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				k8ssIds = append(k8ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return k8ssIds
}

func LansIds(datacenterId string) []string {
	lanSvc := resources.NewLanService(client.Must(), context.Background())
	lans, _, err := lanSvc.List(datacenterId)
	if err != nil {
		return nil
	}
	lansIds := make([]string, 0)
	if items, ok := lans.GetItemsOk(); ok {
		for _, item := range *items {
			var completion string
			if item.Id == nil {
				continue
			}
			completion = *item.Id
			if props, ok := item.GetPropertiesOk(); ok {
				if name, ok := props.GetNameOk(); ok {
					// Here is where the completion descriptions start
					completion = fmt.Sprintf("%s\t%s", completion, *name)
				}
				if public, ok := props.GetPublicOk(); ok {
					if *public {
						completion = fmt.Sprintf("%s (public)", completion)
					} else {
						completion = fmt.Sprintf("%s (private)", completion)
					}
				}
			}

			lansIds = append(lansIds, completion)
		}
	} else {
		return nil
	}
	return lansIds
}

func LoadbalancersIds(datacenterId string) []string {
	loadbalancerSvc := resources.NewLoadbalancerService(client.Must(), context.Background())
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId)
	if err != nil {
		return nil
	}
	loadbalancersIds := make([]string, 0)
	if items, ok := loadbalancers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				loadbalancersIds = append(loadbalancersIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return loadbalancersIds
}

func LocationIds() []string {
	locationSvc := resources.NewLocationService(client.Must(), context.Background())
	locations, _, err := locationSvc.List()
	if err != nil {
		return nil
	}

	items, ok := locations.GetItemsOk()
	if !ok || items == nil || len(*items) == 0 {
		return nil
	}

	out := make([]string, 0, len(*items))
	var b strings.Builder

	for _, item := range *items {
		// reset builder for this item
		b.Reset()

		if id, ok := item.GetIdOk(); ok && id != nil && *id != "" {
			b.WriteString(*id)
			b.WriteByte('\t')
		}

		if item.Properties != nil && item.Properties.Name != nil {
			b.WriteString(*item.Properties.Name)
		}
		b.WriteString("; supports CPUs: ")

		if item.Properties == nil || item.Properties.CpuArchitecture == nil || len(*item.Properties.CpuArchitecture) == 0 {
			b.WriteString("none")
		} else {
			first := true
			for _, cpu := range *item.Properties.CpuArchitecture {
				if cpu.CpuFamily == nil {
					continue
				}
				if !first {
					b.WriteString(", ")
				}
				first = false
				b.WriteString(*cpu.CpuFamily)
			}
			if first {
				b.WriteString("none")
			}
		}

		out = append(out, b.String())
	}

	return out
}

func NatGatewaysIds(datacenterId string) []string {
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natgateways, _, err := natgatewaySvc.List(datacenterId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natgateways.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func NatGatewayFlowLogsIds(datacenterId, natgatewayId string) []string {
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natFlowLogs, _, err := natgatewaySvc.ListFlowLogs(datacenterId, natgatewayId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func NatGatewayRulesIds(datacenterId, natgatewayId string) []string {
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natgateways.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func NetworkLoadBalancersIds(datacenterId string) []string {
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := networkloadbalancers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func NetworkLoadBalancerFlowLogsIds(datacenterId, networkloadbalancerId string) []string {
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func ForwardingRulesIds(datacenterId, nlbId string) []string {
	nlbSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func NicsIds(datacenterId, serverId string) []string {
	nicSvc := resources.NewNicService(client.Must(), context.Background())
	nics, _, err := nicSvc.List(datacenterId, serverId)
	if err != nil {
		return nil
	}
	nicsIds := make([]string, 0)
	if items, ok := nics.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				nicsIds = append(nicsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return nicsIds
}

func AttachedNicsIds(datacenterId, loadbalancerId string) []string {
	nicSvc := resources.NewLoadbalancerService(client.Must(), context.Background())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId)
	if err != nil {
		return nil
	}
	attachedNicsIds := make([]string, 0)
	if items, ok := nics.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedNicsIds = append(attachedNicsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedNicsIds
}

func PccsIds() []string {
	pccSvc := resources.NewPrivateCrossConnectService(client.Must(), context.Background())
	pccs, _, err := pccSvc.List()
	if err != nil {
		return nil
	}
	pccsIds := make([]string, 0)
	if items, ok := pccs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				pccsIds = append(pccsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return pccsIds
}

func RequestsIds() []string {
	reqSvc := resources.NewRequestService(client.Must(), context.Background())
	requests, _, err := reqSvc.List()
	if err != nil {
		return nil
	}
	reqIds := make([]string, 0)
	if items, ok := requests.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				reqIds = append(reqIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return reqIds
}

func ResourcesIds() []string {
	userSvc := resources.NewUserService(client.Must(), context.Background())
	res, _, err := userSvc.ListResources()
	if err != nil {
		return nil
	}
	resIds := make([]string, 0)
	if items, ok := res.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				resIds = append(resIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return resIds
}

func S3KeyIds(userId string) []string {
	S3KeySvc := resources.NewS3KeyService(client.Must(), context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := S3Keys.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func ServersIds(datacenterId string) []string {
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	servers, _, err := serverSvc.List(datacenterId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := servers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func GroupResourcesIds(groupId string) []string {
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	res, _, err := groupSvc.ListResources(groupId)
	if err != nil {
		return nil
	}
	resIds := make([]string, 0)
	if items, ok := res.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				resIds = append(resIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return resIds
}

func SnapshotIds() []string {
	snapshotSvc := resources.NewSnapshotService(client.Must(), context.Background())
	snapshots, _, err := snapshotSvc.List()
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := snapshots.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func TemplatesIds() []string {
	tplSvc := resources.NewTemplateService(client.Must(), context.Background())
	templates, _, err := tplSvc.List()
	if err != nil {
		return nil
	}

	items, ok := templates.GetItemsOk()
	if !ok || items == nil {
		return nil
	}

	tplIds := make([]string, 0, len(*items))
	for _, item := range *items {
		if item.Id == nil {
			continue
		}

		completion := *item.Id

		if props, ok := item.GetPropertiesOk(); ok {
			// Example: "Basic Cube M | 4 cores, 16 GB RAM, 240 GB (Basic Templates)"
			parts := make([]string, 0)

			if name, ok := props.GetNameOk(); ok {
				parts = append(parts, *name)
			}

			if cores, ok := props.GetCoresOk(); ok {
				parts = append(parts, fmt.Sprintf("%.0f cores", *cores))
			}

			if ram, ok := props.GetRamOk(); ok {
				parts = append(parts, fmt.Sprintf("%.0f GB RAM", *ram/1024))
			}

			if storage, ok := props.GetStorageSizeOk(); ok {
				parts = append(parts, fmt.Sprintf("%.0f GB", *storage))
			}

			if category, ok := props.GetCategoryOk(); ok {
				parts = append(parts, fmt.Sprintf("(%s)", *category))
			}

			if gpus, ok := props.GetGpusOk(); ok && gpus != nil && len(*gpus) > 0 {
				gpuInfo := resource2table.FormatGPUs(*gpus)
				parts = append(parts, gpuInfo)
			}

			if len(parts) > 0 {
				completion = fmt.Sprintf("%s\t%s", completion, strings.Join(parts, " | "))
			}
		}

		tplIds = append(tplIds, completion)
	}

	return tplIds
}

func UsersIds() []string {
	userSvc := resources.NewUserService(client.Must(), context.Background())
	users, _, err := userSvc.List()
	if err != nil {
		return nil
	}

	items, ok := users.GetItemsOk()
	if !ok || items == nil || len(*items) == 0 {
		return nil
	}

	out := make([]string, 0, len(*items))
	var b strings.Builder

	for _, item := range *items {
		b.Reset()

		// id + tab
		if id, ok := item.GetIdOk(); ok && id != nil && *id != "" {
			b.WriteString(*id)
			b.WriteByte('\t')
		} else {
			continue
		}

		// build user helper info: email, full name, admin
		email := ""
		if item.Properties != nil && item.Properties.Email != nil {
			email = *item.Properties.Email
		}

		first := ""
		if item.Properties != nil && item.Properties.Firstname != nil {
			first = *item.Properties.Firstname
		}

		last := ""
		if item.Properties != nil && item.Properties.Lastname != nil {
			last = *item.Properties.Lastname
		}

		admin := "false"
		if item.Properties != nil && item.Properties.Administrator != nil && *item.Properties.Administrator {
			admin = "true"
		}

		// format: email — First Last; admin: true
		if email != "" {
			b.WriteString(email)
			b.WriteString(" ")
		}
		if first != "" || last != "" {
			b.WriteString("(")
			if first != "" {
				b.WriteString(first)
			}
			if first != "" && last != "" {
				b.WriteString(" ")
			}
			if last != "" {
				b.WriteString(last)
			}
			b.WriteString(") ")
		}
		b.WriteString("; admin: ")
		b.WriteString(admin)

		out = append(out, b.String())
	}

	return out
}

func GroupUsersIds(groupId string) []string {
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	users, _, err := groupSvc.ListUsers(groupId)
	if err != nil {
		return nil
	}
	usersIds := make([]string, 0)
	if items, ok := users.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				usersIds = append(usersIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return usersIds
}

func VolumesIds(datacenterId string) []string {
	volumeSvc := resources.NewVolumeService(client.Must(), context.Background())
	volumes, _, err := volumeSvc.List(datacenterId)
	if err != nil {
		return nil
	}

	items, ok := volumes.GetItemsOk()
	if !ok || items == nil || len(*items) == 0 {
		return nil
	}

	out := make([]string, 0, len(*items))
	var b strings.Builder

	for _, item := range *items {
		b.Reset()

		// id + tab
		if id, ok := item.GetIdOk(); ok && id != nil && *id != "" {
			b.WriteString(*id)
			b.WriteByte('\t')
		} else {
			continue
		}

		// name
		name := ""
		if item.Properties != nil && item.Properties.Name != nil {
			name = *item.Properties.Name
		}
		if name != "" {
			b.WriteString(name)
		} else {
			b.WriteString("(no name)")
		}

		// image alias
		image := ""
		if item.Properties != nil && item.Properties.ImageAlias != nil {
			image = *item.Properties.ImageAlias
		}
		b.WriteString("; image: ")
		if image != "" {
			b.WriteString(image)
		} else {
			b.WriteString("none")
		}

		out = append(out, b.String())
	}

	return out
}

func AttachedVolumesIds(datacenterId, serverId string) []string {
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId)
	if err != nil {
		return nil
	}
	attachedVolumesIds := make([]string, 0)
	if items, ok := volumes.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				attachedVolumesIds = append(attachedVolumesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return attachedVolumesIds
}

func ApplicationLoadBalancersIds(datacenterId string) []string {
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	applicationloadbalancers, _, err := applicationloadbalancerSvc.List(datacenterId)
	if err != nil {
		return nil
	}
	albIds := make([]string, 0)
	if items, ok := applicationloadbalancers.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				albIds = append(albIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return albIds
}

func ApplicationLoadBalancerFlowLogsIds(datacenterId, applicationloadbalancerId string) []string {
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	natFlowLogs, _, err := applicationloadbalancerSvc.ListFlowLogs(datacenterId, applicationloadbalancerId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func AlbForwardingRulesIds(datacenterId, albId string) []string {
	albSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId)
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func TargetGroupIds() []string {
	targetGroupSvc := resources.NewTargetGroupService(client.Must(), context.Background())
	targetGroups, _, err := targetGroupSvc.List()
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := targetGroups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				ssIds = append(ssIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return ssIds
}

func GpusIds(datacenterId, serverId string) []string {
	if datacenterId == "" || serverId == "" {
		return nil
	}

	gpus, _, err := client.Must().CloudClient.GraphicsProcessingUnitCardsApi.DatacentersServersGPUsGet(context.Background(), datacenterId, serverId).Execute()
	if err != nil {
		return nil
	}

	if gpus.GetItems() == nil {
		return nil
	}

	ids := make([]string, 0)
	for _, gpu := range *gpus.GetItems() {
		if gpu.GetId() != nil {
			ids = append(ids, *gpu.GetId())
		}
	}

	return ids
}
