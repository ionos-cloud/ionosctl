/*
This is used for supporting completion in the CLI.
Option: --datacenter-id --server-id --backupunit-id, usually --<RESOURCE_TYPE>-id
*/
package completer

import (
	"context"
	"fmt"
	"io"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

func BackupUnitsIds(outErr io.Writer) []string {
	backupUnitSvc := resources.NewBackupUnitService(client.Must(), context.Background())
	backupUnits, _, err := backupUnitSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	datacenterSvc := resources.NewDataCenterService(client.Must(), context.Background())
	datacenters, _, err := datacenterSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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

func FirewallRulesIds(outErr io.Writer, datacenterId, serverId, nicId string) []string {
	firewallRuleSvc := resources.NewFirewallRuleService(client.Must(), context.Background())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	flowLogSvc := resources.NewFlowLogService(client.Must(), context.Background())
	flowLogs, _, err := flowLogSvc.List(datacenterId, serverId, nicId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	groups, _, err := groupSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	imageSvc := resources.NewImageService(client.Must(), context.Background())
	images, _, err := imageSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	ipBlockSvc := resources.NewIpBlockService(client.Must(), context.Background())
	ipBlocks, _, err := ipBlockSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListClusters(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListVersions()
	clierror.CheckErrorAndDie(err, outErr)
	return k8ss
}

func K8sNodesIds(outErr io.Writer, clusterId, nodepoolId string) []string {
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	k8sSvc := resources.NewK8sService(client.Must(), context.Background())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	lanSvc := resources.NewLanService(client.Must(), context.Background())
	lans, _, err := lanSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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

func LoadbalancersIds(outErr io.Writer, datacenterId string) []string {
	loadbalancerSvc := resources.NewLoadbalancerService(client.Must(), context.Background())
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	locationSvc := resources.NewLocationService(client.Must(), context.Background())
	locations, _, err := locationSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natgateways, _, err := natgatewaySvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natFlowLogs, _, err := natgatewaySvc.ListFlowLogs(datacenterId, natgatewayId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	natgatewaySvc := resources.NewNatGatewayService(client.Must(), context.Background())
	natgateways, _, err := natgatewaySvc.ListRules(datacenterId, natgatewayId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	networkloadbalancers, _, err := networkloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	networkloadbalancerSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	natFlowLogs, _, err := networkloadbalancerSvc.ListFlowLogs(datacenterId, networkloadbalancerId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	nlbSvc := resources.NewNetworkLoadBalancerService(client.Must(), context.Background())
	natForwardingRules, _, err := nlbSvc.ListForwardingRules(datacenterId, nlbId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	nicSvc := resources.NewNicService(client.Must(), context.Background())
	nics, _, err := nicSvc.List(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	nicSvc := resources.NewLoadbalancerService(client.Must(), context.Background())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	pccSvc := resources.NewPrivateCrossConnectService(client.Must(), context.Background())
	pccs, _, err := pccSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	reqSvc := resources.NewRequestService(client.Must(), context.Background())
	requests, _, err := reqSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	userSvc := resources.NewUserService(client.Must(), context.Background())
	res, _, err := userSvc.ListResources()
	clierror.CheckErrorAndDie(err, outErr)
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
	S3KeySvc := resources.NewS3KeyService(client.Must(), context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	servers, _, err := serverSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	res, _, err := groupSvc.ListResources(groupId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	snapshotSvc := resources.NewSnapshotService(client.Must(), context.Background())
	snapshots, _, err := snapshotSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	tplSvc := resources.NewTemplateService(client.Must(), context.Background())
	tpls, _, err := tplSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	userSvc := resources.NewUserService(client.Must(), context.Background())
	users, _, err := userSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	groupSvc := resources.NewGroupService(client.Must(), context.Background())
	users, _, err := groupSvc.ListUsers(groupId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	volumeSvc := resources.NewVolumeService(client.Must(), context.Background())
	volumes, _, err := volumeSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	serverSvc := resources.NewServerService(client.Must(), context.Background())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	applicationloadbalancers, _, err := applicationloadbalancerSvc.List(datacenterId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	applicationloadbalancerSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	natFlowLogs, _, err := applicationloadbalancerSvc.ListFlowLogs(datacenterId, applicationloadbalancerId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	albSvc := resources.NewApplicationLoadBalancerService(client.Must(), context.Background())
	natForwardingRules, _, err := albSvc.ListForwardingRules(datacenterId, albId, resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
	targetGroupSvc := resources.NewTargetGroupService(client.Must(), context.Background())
	targetGroups, _, err := targetGroupSvc.List(resources.ListQueryParams{})
	clierror.CheckErrorAndDie(err, outErr)
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
