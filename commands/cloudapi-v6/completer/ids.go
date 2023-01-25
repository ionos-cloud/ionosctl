/*
This is used for supporting completion in the CLI.
Option: --datacenter-id --server-id --backupunit-id, usually --<RESOURCE_TYPE>-id
*/
package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
)

func BackupUnitsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	backupUnitSvc := resources.NewBackupUnitService(client, context.TODO())
	backupUnits, _, err := backupUnitSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func AttachedCdromsIds(outErr io.Writer, datacenterId, serverId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func DataCentersIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	datacenterSvc := resources.NewDataCenterService(client, context.TODO())
	datacenters, _, err := datacenterSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
	dcIds := make([]string, 0)
	if items, ok := datacenters.Datacenters.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				dcIds = append(dcIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return dcIds
}

func FirewallRulesIds(outErr io.Writer, datacenterId, serverId, nicId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	firewallRuleSvc := resources.NewFirewallRuleService(client, context.TODO())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func FlowLogsIds(outErr io.Writer, datacenterId, serverId, nicId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	flowLogSvc := resources.NewFlowLogService(client, context.TODO())
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func GroupsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	groups, _, err := groupSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ImageIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	imageSvc := resources.NewImageService(client, context.TODO())
	images, _, err := imageSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func IpBlocksIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	ipBlockSvc := resources.NewIpBlockService(client, context.TODO())
	ipBlocks, _, err := ipBlockSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func K8sClustersIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListClusters(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func K8sVersionsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListVersions()
	clierror.CheckError(err, outErr)
	return k8ss
}

func K8sNodesIds(outErr io.Writer, clusterId, nodepoolId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func K8sNodePoolsIds(outErr io.Writer, clusterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func LansIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	lanSvc := resources.NewLanService(client, context.TODO())
	lans, _, err := lanSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
	lansIds := make([]string, 0)
	if items, ok := lans.Lans.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				lansIds = append(lansIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return lansIds
}

func LoadbalancersIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	loadbalancerSvc := resources.NewLoadbalancerService(client, context.TODO())
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func LocationIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	locationSvc := resources.NewLocationService(client, context.TODO())
	locations, _, err := locationSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NatGatewaysIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natgateways, _, err := natgatewaySvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NatGatewayFlowLogsIds(outErr io.Writer, datacenterId, natgatewayId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natFlowLogs, _, err := natgatewaySvc.ListFlowLogs(datacenterId, natgatewayId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NatGatewayRulesIds(outErr io.Writer, datacenterId, natgatewayId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	natgatewaySvc := resources.NewNatGatewayService(client, context.TODO())
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NetworkLoadBalancersIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NetworkLoadBalancerFlowLogsIds(outErr io.Writer, datacenterId, networkloadbalancerId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ForwardingRulesIds(outErr io.Writer, datacenterId, nlbId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	nlbSvc := resources.NewNetworkLoadBalancerService(client, context.TODO())
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func NicsIds(outErr io.Writer, datacenterId, serverId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	nicSvc := resources.NewNicService(client, context.TODO())
	nics, _, err := nicSvc.List(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func AttachedNicsIds(outErr io.Writer, datacenterId, loadbalancerId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	nicSvc := resources.NewLoadbalancerService(client, context.TODO())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func PccsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	pccSvc := resources.NewPrivateCrossConnectService(client, context.TODO())
	pccs, _, err := pccSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func RequestsIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	reqSvc := resources.NewRequestService(client, context.TODO())
	requests, _, err := reqSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ResourcesIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	userSvc := resources.NewUserService(client, context.TODO())
	res, _, err := userSvc.ListResources()
	clierror.CheckError(err, outErr)
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

func S3KeyIds(outErr io.Writer, userId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	S3KeySvc := resources.NewS3KeyService(client, context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ServersIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	servers, _, err := serverSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func GroupResourcesIds(outErr io.Writer, groupId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	res, _, err := groupSvc.ListResources(groupId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func SnapshotIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	snapshotSvc := resources.NewSnapshotService(client, context.TODO())
	snapshots, _, err := snapshotSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func TemplatesIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	tplSvc := resources.NewTemplateService(client, context.TODO())
	tpls, _, err := tplSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func UsersIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	userSvc := resources.NewUserService(client, context.TODO())
	users, _, err := userSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func GroupUsersIds(outErr io.Writer, groupId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	users, _, err := groupSvc.ListUsers(groupId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func VolumesIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	volumeSvc := resources.NewVolumeService(client, context.TODO())
	volumes, _, err := volumeSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func AttachedVolumesIds(outErr io.Writer, datacenterId, serverId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ApplicationLoadBalancersIds(outErr io.Writer, datacenterId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	applicationloadbalancers, _, err := applicationloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func ApplicationLoadBalancerFlowLogsIds(outErr io.Writer, datacenterId, applicationloadbalancerId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	natFlowLogs, _, err := applicationloadbalancerSvc.ListFlowLogs(datacenterId, applicationloadbalancerId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func AlbForwardingRulesIds(outErr io.Writer, datacenterId, albId string) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	albSvc := resources.NewApplicationLoadBalancerService(client, context.TODO())
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId, resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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

func TargetGroupIds(outErr io.Writer) []string {
	client, err := config.GetClient()
	clierror.CheckError(err, outErr)
	targetGroupSvc := resources.NewTargetGroupService(client, context.TODO())
	targetGroups, _, err := targetGroupSvc.List(resources.ListQueryParams{})
	clierror.CheckError(err, outErr)
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
