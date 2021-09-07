package completer

import (
	"context"
	"io"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5/resources"
	"github.com/spf13/viper"
)

func BackupUnitsIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	backupUnitSvc := resources.NewBackupUnitService(client, context.TODO())
	backupUnits, _, err := backupUnitSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	cdroms, _, err := serverSvc.ListCdroms(datacenterId, serverId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	datacenterSvc := resources.NewDataCenterService(client, context.TODO())
	datacenters, _, err := datacenterSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	firewallRuleSvc := resources.NewFirewallRuleService(client, context.TODO())
	firewallRules, _, err := firewallRuleSvc.List(datacenterId, serverId, nicId)
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

func GroupsIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	groups, _, err := groupSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	imageSvc := resources.NewImageService(client, context.TODO())
	images, _, err := imageSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	ipBlockSvc := resources.NewIpBlockService(client, context.TODO())
	ipBlocks, _, err := ipBlockSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListClusters()
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

func K8sNodesIds(outErr io.Writer, clusterId, nodepoolId string) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodes(clusterId, nodepoolId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	k8sSvc := resources.NewK8sService(client, context.TODO())
	k8ss, _, err := k8sSvc.ListNodePools(clusterId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	lanSvc := resources.NewLanService(client, context.TODO())
	lans, _, err := lanSvc.List(datacenterId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	loadbalancerSvc := resources.NewLoadbalancerService(client, context.TODO())
	loadbalancers, _, err := loadbalancerSvc.List(datacenterId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	locationSvc := resources.NewLocationService(client, context.TODO())
	locations, _, err := locationSvc.List()
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

func NicsIds(outErr io.Writer, datacenterId, serverId string) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	nicSvc := resources.NewNicService(client, context.TODO())
	nics, _, err := nicSvc.List(datacenterId, serverId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	nicSvc := resources.NewLoadbalancerService(client, context.TODO())
	nics, _, err := nicSvc.ListNics(datacenterId, loadbalancerId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	pccSvc := resources.NewPrivateCrossConnectService(client, context.TODO())
	pccs, _, err := pccSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	reqSvc := resources.NewRequestService(client, context.TODO())
	requests, _, err := reqSvc.List()
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
	client, err := getClient()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	S3KeySvc := resources.NewS3KeyService(client, context.TODO())
	S3Keys, _, err := S3KeySvc.List(userId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	servers, _, err := serverSvc.List(datacenterId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	res, _, err := groupSvc.ListResources(groupId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	snapshotSvc := resources.NewSnapshotService(client, context.TODO())
	snapshots, _, err := snapshotSvc.List()
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

func UsersIds(outErr io.Writer) []string {
	client, err := getClient()
	clierror.CheckError(err, outErr)
	userSvc := resources.NewUserService(client, context.TODO())
	users, _, err := userSvc.List()
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(client, context.TODO())
	users, _, err := groupSvc.ListUsers(groupId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	volumeSvc := resources.NewVolumeService(client, context.TODO())
	volumes, _, err := volumeSvc.List(datacenterId)
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
	client, err := getClient()
	clierror.CheckError(err, outErr)
	serverSvc := resources.NewServerService(client, context.TODO())
	volumes, _, err := serverSvc.ListVolumes(datacenterId, serverId)
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

// Get Client for Completion Functions
func getClient() (*resources.Client, error) {
	if err := config.Load(); err != nil {
		return nil, err
	}
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	if err != nil {
		return nil, err
	}
	return clientSvc.Get(), nil
}
