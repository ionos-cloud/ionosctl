/*
This is used for supporting completion in the CLI.
Option: --datacenter-id --server-id --backupunit-id, usually --<RESOURCE_TYPE>-id
*/
package completer

import (
	"context"
	"fmt"
	"io"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func BackupUnitsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	backupUnitSvc := resources.NewBackupUnitService(client, context.TODO())
	backupUnits, _, err := backupUnitSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func AttachedCdromsIds(_ io.Writer, datacenterId, serverId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	serverSvc := resources.NewServerService(client, context.TODO())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func DataCentersIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	datacenterSvc := resources.NewDataCenterService(client, context.TODO())
	datacenters, _, err := datacenterSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
	}

	dcIds := make([]string, 0)
	if items, ok := datacenters.Datacenters.GetItemsOk(); ok {
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

func FirewallRulesIds(_ io.Writer, datacenterId, serverId, nicId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	firewallRuleSvc := resources.NewFirewallRuleService(client, context.TODO())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func FlowLogsIds(_ io.Writer, datacenterId, serverId, nicId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	flowLogSvc := resources.NewFlowLogService(client, context.TODO())
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func GroupsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	groupSvc := resources.NewGroupService(client, context.TODO())
	groups, _, err := groupSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ImageIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	imageSvc := resources.NewImageService(client, context.TODO())
	images, _, err := imageSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
	}

	imgsIds := make([]string, 0)
	if items, ok := images.Images.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				imgsIds = append(imgsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return imgsIds
}

func IpBlocksIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	ipBlockSvc := resources.NewIpBlockService(client, context.TODO())
	ipBlocks, _, err := ipBlockSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func K8sClustersIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListClusters(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func K8sVersionsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListVersions()
	if err != nil {
		die.Die(err.Error())
	}

	return k8ss
}

func K8sNodesIds(_ io.Writer, clusterId, nodepoolId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func K8sNodePoolsIds(_ io.Writer, clusterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func LansIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	lanSvc := resources.NewLanService(client, context.TODO())
	lans, _, err := lanSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func LoadbalancersIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	loadbalancerSvc := resources.NewLoadbalancerService(client, context.TODO())
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func LocationIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	locationSvc := resources.NewLocationService(client, context.TODO())
	locations, _, err := locationSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NatGatewaysIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natgateways, _, err := natgatewaySvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NatGatewayFlowLogsIds(_ io.Writer, datacenterId, natgatewayId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natFlowLogs, _, err := natgatewaySvc.ListFlowLogs(datacenterId, natgatewayId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NatGatewayRulesIds(_ io.Writer, datacenterId, natgatewayId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NetworkLoadBalancersIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NetworkLoadBalancerFlowLogsIds(_ io.Writer, datacenterId, networkloadbalancerId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ForwardingRulesIds(_ io.Writer, datacenterId, nlbId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	nlbSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func NicsIds(_ io.Writer, datacenterId, serverId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	nicSvc := resources.NewNicService(client, context.TODO())
	nics, _, err := nicSvc.List(datacenterId, serverId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func AttachedNicsIds(_ io.Writer, datacenterId, loadbalancerId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	nicSvc := resources.NewLoadbalancerService(client, context.TODO())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func PccsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	pccSvc := resources.NewPrivateCrossConnectService(client, context.TODO())
	pccs, _, err := pccSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func RequestsIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	reqSvc := resources.NewRequestService(client, context.TODO())
	requests, _, err := reqSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ResourcesIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	userSvc := resources.NewUserService(client, context.TODO())
	res, _, err := userSvc.ListResources()
	if err != nil {
		die.Die(err.Error())
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

func S3KeyIds(_ io.Writer, userId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	S3KeySvc := resources.NewS3KeyService(client, context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ServersIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	serverSvc := resources.NewServerService(client, context.TODO())
	servers, _, err := serverSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func GroupResourcesIds(_ io.Writer, groupId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	groupSvc := resources.NewGroupService(client, context.TODO())
	res, _, err := groupSvc.ListResources(groupId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func SnapshotIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	snapshotSvc := resources.NewSnapshotService(client, context.TODO())
	snapshots, _, err := snapshotSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func TemplatesIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	tplSvc := resources.NewTemplateService(client, context.TODO())
	tpls, _, err := tplSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func UsersIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	userSvc := resources.NewUserService(client, context.TODO())
	users, _, err := userSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func GroupUsersIds(_ io.Writer, groupId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	groupSvc := resources.NewGroupService(client, context.TODO())
	users, _, err := groupSvc.ListUsers(groupId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func VolumesIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	volumeSvc := resources.NewVolumeService(client, context.TODO())
	volumes, _, err := volumeSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func AttachedVolumesIds(_ io.Writer, datacenterId, serverId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	serverSvc := resources.NewServerService(client, context.TODO())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ApplicationLoadBalancersIds(_ io.Writer, datacenterId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	applicationloadbalancers, _, err := applicationloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func ApplicationLoadBalancerFlowLogsIds(_ io.Writer, datacenterId, applicationloadbalancerId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	natFlowLogs, _, err := applicationloadbalancerSvc.ListFlowLogs(datacenterId, applicationloadbalancerId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func AlbForwardingRulesIds(_ io.Writer, datacenterId, albId string) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	albSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId, resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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

func TargetGroupIds(_ io.Writer) []string {
	client, err := client2.Get()
	if err != nil {
		die.Die(err.Error())
	}

	targetGroupSvc := resources.NewTargetGroupService(client, context.TODO())
	targetGroups, _, err := targetGroupSvc.List(resources.ListQueryParams{})
	if err != nil {
		die.Die(err.Error())
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
