/*
This is used for supporting completion in the CLI.
Option: --datacenter-id --server-id --backupunit-id, usually --<RESOURCE_TYPE>-id
*/
package completer

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	compute "github.com/ionos-cloud/sdk-go/v6"
)

func BackupUnitsIds() []string {
	backupUnitSvc := resources.NewBackupUnitService(client.Must(), context.Background())
	backupUnits, _, err := backupUnitSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	backupUnitsIds := make([]string, 0)
	if items, ok := backupUnits.BackupUnits.GetItemsOk(); ok && items != nil {
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
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId, resources.ListQueryParams{})
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

func DataCentersIds(filters ...func(datacenter compute.Datacenter) bool) []string {
	datacenterSvc := resources.NewDataCenterService(client.Must(), context.Background())
	datacenters, _, err := datacenterSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	dcIds := make([]string, 0)
	if items, ok := datacenters.Datacenters.GetItemsOk(); ok {
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
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	firewallRulesIds := make([]string, 0)
	if items, ok := firewallRules.FirewallRules.GetItemsOk(); ok && items != nil {
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
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	flowLogsIds := make([]string, 0)
	if items, ok := flowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
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
	groups, _, err := groupSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	groupsIds := make([]string, 0)
	if items, ok := groups.Groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				groupsIds = append(groupsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return groupsIds
}

func ImageIds(customFilters ...func(compute.ApiImagesGetRequest) compute.ApiImagesGetRequest) []string {
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

		if public, ok := image.Properties.GetPublicOk(); ok {
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
	ipBlocks, _, err := ipBlockSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := ipBlocks.IpBlocks.GetItemsOk(); ok && items != nil {
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

func K8sClustersIds() []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListClusters(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesClusters.GetItemsOk(); ok && items != nil {
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
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesNodes.GetItemsOk(); ok && items != nil {
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
	k8ss, _, err := k8sSvc.ListNodePools(clusterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	k8ssIds := make([]string, 0)
	if items, ok := k8ss.KubernetesNodePools.GetItemsOk(); ok && items != nil {
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
	lans, _, err := lanSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	lansIds := make([]string, 0)
	if items, ok := lans.Lans.GetItemsOk(); ok {
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
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	loadbalancersIds := make([]string, 0)
	if items, ok := loadbalancers.Loadbalancers.GetItemsOk(); ok && items != nil {
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
	locations, _, err := locationSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	lcIds := make([]string, 0)
	if items, ok := locations.Locations.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				lcIds = append(lcIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return lcIds
}

func NatGatewaysIds(datacenterId string) []string {
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natgateways, _, err := natgatewaySvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natgateways.NatGateways.GetItemsOk(); ok && items != nil {
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
	natFlowLogs, _, err := natgatewaySvc.ListFlowLogs(datacenterId, natgatewayId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
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
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natgateways.NatGatewayRules.GetItemsOk(); ok && items != nil {
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
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := networkloadbalancers.NetworkLoadBalancers.GetItemsOk(); ok && items != nil {
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
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
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
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.NetworkLoadBalancerForwardingRules.GetItemsOk(); ok && items != nil {
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
	nics, _, err := nicSvc.List(datacenterId, serverId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	nicsIds := make([]string, 0)
	if items, ok := nics.Nics.GetItemsOk(); ok && items != nil {
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
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	attachedNicsIds := make([]string, 0)
	if items, ok := nics.BalancedNics.GetItemsOk(); ok && items != nil {
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
	pccs, _, err := pccSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	pccsIds := make([]string, 0)
	if items, ok := pccs.PrivateCrossConnects.GetItemsOk(); ok && items != nil {
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
	requests, _, err := reqSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	reqIds := make([]string, 0)
	if items, ok := requests.Requests.GetItemsOk(); ok && items != nil {
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
	if items, ok := res.Resources.GetItemsOk(); ok && items != nil {
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
	S3Keys, _, err := S3KeySvc.List(userId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := S3Keys.S3Keys.GetItemsOk(); ok && items != nil {
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
	servers, _, err := serverSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := servers.Servers.GetItemsOk(); ok && items != nil {
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
	res, _, err := groupSvc.ListResources(groupId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	resIds := make([]string, 0)
	if items, ok := res.ResourceGroups.GetItemsOk(); ok && items != nil {
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
	snapshots, _, err := snapshotSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := snapshots.Snapshots.GetItemsOk(); ok && items != nil {
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
	tpls, _, err := tplSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	tplsIds := make([]string, 0)
	if items, ok := tpls.Templates.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				tplsIds = append(tplsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return tplsIds
}

func UsersIds() []string {
	userSvc := resources.NewUserService(client.Must(), context.Background())
	users, _, err := userSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	usersIds := make([]string, 0)
	if items, ok := users.Users.GetItemsOk(); ok && items != nil {
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

func GroupUsersIds(groupId string) []string {
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	users, _, err := groupSvc.ListUsers(groupId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	usersIds := make([]string, 0)
	if items, ok := users.GroupMembers.GetItemsOk(); ok && items != nil {
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
	volumes, _, err := volumeSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	volumesIds := make([]string, 0)
	if items, ok := volumes.Volumes.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				volumesIds = append(volumesIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return volumesIds
}

func AttachedVolumesIds(datacenterId, serverId string) []string {
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	attachedVolumesIds := make([]string, 0)
	if items, ok := volumes.AttachedVolumes.GetItemsOk(); ok && items != nil {
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
	applicationloadbalancers, _, err := applicationloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	albIds := make([]string, 0)
	if items, ok := applicationloadbalancers.ApplicationLoadBalancers.GetItemsOk(); ok && items != nil {
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
	natFlowLogs, _, err := applicationloadbalancerSvc.ListFlowLogs(datacenterId, applicationloadbalancerId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natFlowLogs.FlowLogs.GetItemsOk(); ok && items != nil {
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
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId, resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := natForwardingRules.ApplicationLoadBalancerForwardingRules.GetItemsOk(); ok && items != nil {
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
	targetGroups, _, err := targetGroupSvc.List(resources.ListQueryParams{})
	if err != nil {
		return nil
	}
	ssIds := make([]string, 0)
	if items, ok := targetGroups.TargetGroups.GetItemsOk(); ok && items != nil {
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
