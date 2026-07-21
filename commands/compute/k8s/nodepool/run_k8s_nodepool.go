package nodepool

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	k8scluster "github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/cluster"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func RunK8sNodePoolCreateFromJSON(c *core.CommandConfig, propertiesFromJson ionoscloud.KubernetesNodePoolForPost) error {
	np, _, err := client.Must().CloudClient.KubernetesApi.K8sNodepoolsPost(context.Background(),
		c.Flags().String(constants.FlagClusterId)).KubernetesNodePool(propertiesFromJson).Execute()
	if err != nil {
		return fmt.Errorf("failed creating nodepool from json properties: %w", err)
	}

	return handleApiResponseK8sNodepoolCreate(c, np)
}

func RunK8sNodePoolCreate(c *core.CommandConfig) error {
	newNodePool, err := getNewK8sNodePool(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating K8s NodePool in K8s Cluster with ID: %v", c.Flags().String(constants.FlagClusterId))

	u, resp, err := c.CloudApiV6Services.K8s().CreateNodePool(c.Flags().String(constants.FlagClusterId), *newNodePool)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return handleApiResponseK8sNodepoolCreate(c, u.KubernetesNodePool)
}

func handleApiResponseK8sNodepoolCreate(c *core.CommandConfig, pool ionoscloud.KubernetesNodePool) error {
	return c.Printer(allK8sNodePoolCols).Print(pool)
}

func PreRunK8sNodePoolsList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId},
		[]string{cloudapiv6.ArgAll},
	); err != nil {
		return err
	}
	return nil
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId)
}

func PreRunK8sClusterNodePoolDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId, constants.FlagNodepoolId},
		[]string{constants.FlagClusterId, cloudapiv6.ArgAll},
	)
}

func RunK8sNodePoolListAll(c *core.CommandConfig) error {
	clusters, _, err := c.CloudApiV6Services.K8s().ListClusters()
	if err != nil {
		return err
	}

	var allNodePools []ionoscloud.KubernetesNodePool
	totalTime := time.Duration(0)

	for _, cluster := range k8scluster.GetK8sClusters(clusters) {
		nodePools, resp, err := c.CloudApiV6Services.K8s().ListNodePools(*cluster.GetId())
		if err != nil {
			return err
		}

		items, ok := nodePools.GetItemsOk()
		if !ok || items == nil {
			continue
		}

		allNodePools = append(allNodePools, *items...)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		c.Verbose(constants.MessageRequestTime, totalTime)
	}

	return c.Printer(allK8sNodePoolCols).Print(allNodePools)
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	if c.Flags().Bool(cloudapiv6.ArgAll) {
		return RunK8sNodePoolListAll(c)
	}

	c.Verbose("Getting K8s NodePools from K8s Cluster with ID: %v", c.Flags().String(constants.FlagClusterId))

	k8ss, resp, err := c.CloudApiV6Services.K8s().ListNodePools(c.Flags().String(constants.FlagClusterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sNodePoolCols).Prefix("items").Print(k8ss.KubernetesNodePools)
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	k8sNodePoolId := c.Flags().String(constants.FlagNodepoolId)
	k8sClusterId := c.Flags().String(constants.FlagClusterId)

	c.Verbose("K8s node pool with id: %v from K8s Cluster with id: %v is getting...", k8sNodePoolId, k8sClusterId)

	u, resp, err := c.CloudApiV6Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sNodePoolCols).Print(u.KubernetesNodePool)
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	oldNodePool, _, err := c.CloudApiV6Services.K8s().GetNodePool(c.Flags().String(constants.FlagClusterId),
		c.Flags().String(constants.FlagNodepoolId))
	if err != nil {
		return err
	}

	newNodePool, err := getNewK8sNodePoolUpdated(oldNodePool, c)
	if err != nil {
		return err
	}
	newNodePoolUpdated, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(c.Flags().String(constants.FlagClusterId),
		c.Flags().String(constants.FlagNodepoolId), newNodePool)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sNodePoolCols).Print(newNodePoolUpdated.KubernetesNodePool)
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	k8sClusterId := c.Flags().String(constants.FlagClusterId)
	k8sNodePoolId := c.Flags().String(constants.FlagNodepoolId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllK8sNodepools(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete k8s node pool", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting K8s node pool with id: %v from K8s Cluster with id: %v...", k8sNodePoolId, k8sClusterId)

	resp, err := c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Kubernetes Nodepool successfully deleted")
	return nil
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePoolForPost, error) {
	var k8sversion string

	if c.Flags().Changed(cloudapiv6.ArgK8sVersion) {
		k8sversion = c.Flags().String(cloudapiv6.ArgK8sVersion)
	} else {
		clusterId := c.Flags().String(constants.FlagClusterId)

		k8sCluster, _, err := c.CloudApiV6Services.K8s().GetCluster(clusterId)
		if err != nil {
			return nil, fmt.Errorf("failed to get k8s cluster to fetch default version: %w", err)
		}

		k8sVerPtr := k8sCluster.GetProperties().GetK8sVersion()
		if k8sVerPtr == nil {
			// Fallback: fetch the default k8s version from the API
			k8sversion, err = k8scluster.GetK8sVersion(c)
			if err != nil {
				return nil, fmt.Errorf("k8s version is not set on the cluster, and failed to fetch default version: %w", err)
			}
		} else {
			k8sversion = *k8sVerPtr
		}
	}

	ramSize, err := utils2.ConvertSize(c.Flags().String(constants.FlagRam), utils2.MegaBytes)
	if err != nil {
		return nil, err
	}

	storageSize, err := utils2.ConvertSize(c.Flags().String(constants.FlagStorageSize), utils2.GigaBytes)
	if err != nil {
		return nil, err
	}

	name := c.Flags().String(cloudapiv6.ArgName)
	nodeCount := c.Flags().Int32(constants.FlagNodeCount)
	dcId := c.Flags().String(cloudapiv6.ArgDataCenterId)
	cpuFamily := c.Flags().String(constants.FlagCpuFamily)
	cores := c.Flags().Int32(constants.FlagCores)
	availabilityZone := c.Flags().String(constants.FlagAvailabilityZone)
	storageType := c.Flags().String(constants.FlagStorageType)

	// Set Properties
	nodePoolProperties := ionoscloud.KubernetesNodePoolPropertiesForPost{}
	nodePoolProperties.SetName(name)
	c.Verbose("Property Name set: %v", name)

	nodePoolProperties.SetK8sVersion(k8sversion)
	c.Verbose("Property K8sVersion set: %v", k8sversion)

	nodePoolProperties.SetNodeCount(nodeCount)
	c.Verbose("Property NodeCount set: %v", nodeCount)

	nodePoolProperties.SetDatacenterId(dcId)
	c.Verbose("Property DatacenterId set: %v", dcId)

	if c.Flags().Changed(constants.FlagCpuFamily) {
		nodePoolProperties.SetCpuFamily(c.Flags().String(constants.FlagCpuFamily))
		c.Verbose("Property CPU Family set: %v", cpuFamily)
	}

	if c.Flags().Changed(constants.FlagServerType) {
		nodePoolProperties.SetServerType(ionoscloud.KubernetesNodePoolServerType(c.Flags().String(constants.FlagServerType)))
	}

	nodePoolProperties.SetCoresCount(cores)
	c.Verbose("Property CoresCount set: %v", cores)

	if ramSize < 0 || ramSize > math.MaxInt32 {
		return nil, fmt.Errorf("RAM size %d is out of allowed int32 range [0-%d]", ramSize, math.MaxInt32)
	}
	nodePoolProperties.SetRamSize(int32(ramSize))
	c.Verbose("Property RAM Size set: %vMB", int32(ramSize))

	nodePoolProperties.SetAvailabilityZone(availabilityZone)
	c.Verbose("Property Availability Zone set: %v", availabilityZone)

	if storageSize < 0 || storageSize > math.MaxInt32 {
		return nil, fmt.Errorf("storage size %d is out of allowed int32 range [0-%d]", storageSize, math.MaxInt32)
	}
	nodePoolProperties.SetStorageSize(int32(storageSize))
	c.Verbose("Property Storage Size set: %vGB", int32(storageSize))

	nodePoolProperties.SetStorageType(storageType)
	c.Verbose("Property Storage Type set: %v", storageType)

	if c.Flags().Changed(constants.FlagLabels) {
		keyValueMapLabels := c.Flags().StringToString(constants.FlagLabels)
		nodePoolProperties.SetLabels(keyValueMapLabels)

		c.Verbose("Property Labels set: %v", keyValueMapLabels)
	}

	if c.Flags().Changed(constants.FlagAnnotations) {
		keyValueMapAnnotations := c.Flags().StringToString(constants.FlagAnnotations)
		nodePoolProperties.SetAnnotations(keyValueMapAnnotations)

		c.Verbose("Property Annotations set: %v", keyValueMapAnnotations)
	}

	// Add LANs
	if c.Flags().Changed(cloudapiv6.ArgLanIds) {
		newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		lanIds := c.Flags().IntSlice(cloudapiv6.ArgLanIds)
		dhcp := c.Flags().Bool(cloudapiv6.ArgDhcp)

		for _, lanId := range lanIds {
			if lanId < math.MinInt32 || lanId > math.MaxInt32 {
				return nil, fmt.Errorf("LAN ID %d is out of allowed int32 range [%d-%d]", lanId, math.MinInt32, math.MaxInt32)
			}
			id := int32(lanId)

			c.Verbose("Property Lan ID set: %v", id)
			c.Verbose("Property Dhcp set: %v", dhcp)

			newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
				Id:   &id,
				Dhcp: &dhcp,
			})
		}

		nodePoolProperties.SetLans(newLans)
	}

	return &resources.K8sNodePoolForPost{
		KubernetesNodePoolForPost: ionoscloud.KubernetesNodePoolForPost{
			Properties: &nodePoolProperties,
		},
	}, nil
}

func getNewK8sNodePoolUpdated(oldNodePool *resources.K8sNodePool, c *core.CommandConfig) (resources.K8sNodePoolForPut, error) {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}

	if properties, ok := oldNodePool.GetPropertiesOk(); ok && properties != nil {
		if c.Flags().Changed(cloudapiv6.ArgK8sVersion) {
			vers := c.Flags().String(cloudapiv6.ArgK8sVersion)
			propertiesUpdated.SetK8sVersion(vers)

			c.Verbose("Property K8sVersion set: %v", vers)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}

		if c.Flags().Changed(constants.FlagNodeCount) {
			nodeCount := c.Flags().Int32(constants.FlagNodeCount)
			propertiesUpdated.SetNodeCount(nodeCount)

			c.Verbose("Property NodeCount set: %v", nodeCount)
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgK8sMinNodeCount) ||
			c.Flags().Changed(cloudapiv6.ArgK8sMaxNodeCount) {
			var minCount, maxCount int32

			autoScaling := properties.GetAutoScaling()
			if c.Flags().Changed(cloudapiv6.ArgK8sMinNodeCount) {
				minCount = c.Flags().Int32(cloudapiv6.ArgK8sMinNodeCount)

				c.Verbose("Property MinNodeCount set: %v", minCount)
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}

			if c.Flags().Changed(cloudapiv6.ArgK8sMaxNodeCount) {
				maxCount = c.Flags().Int32(cloudapiv6.ArgK8sMaxNodeCount)

				c.Verbose("Property MaxNodeCount set: %v", maxCount)
			} else {
				if m, ok := autoScaling.GetMaxNodeCountOk(); ok && m != nil {
					maxCount = *m
				}
			}

			if minCount == 0 && maxCount == 0 {
				propertiesUpdated.SetAutoScaling(ionoscloud.KubernetesAutoScaling{})
			} else {
				propertiesUpdated.SetAutoScaling(ionoscloud.KubernetesAutoScaling{
					MinNodeCount: &minCount,
					MaxNodeCount: &maxCount,
				})
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceDay) ||
			c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceTime) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := k8scluster.GetMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})

				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgK8sAnnotationKey) &&
			c.Flags().Changed(cloudapiv6.ArgK8sAnnotationValue) {
			key := c.Flags().String(cloudapiv6.ArgK8sAnnotationKey)
			value := c.Flags().String(cloudapiv6.ArgK8sAnnotationValue)
			propertiesUpdated.SetAnnotations(map[string]string{
				key: value,
			})

			c.Verbose("Property Annotations set: key: %v, value: %v", key, value)
		}

		if c.Flags().Changed(cloudapiv6.ArgLabelKey) &&
			c.Flags().Changed(cloudapiv6.ArgLabelValue) {
			key := c.Flags().String(cloudapiv6.ArgLabelKey)
			value := c.Flags().String(cloudapiv6.ArgLabelValue)
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})

			c.Verbose("Property Labels set: key: %v, value: %v", key, value)
		}

		if c.Flags().Changed(constants.FlagLabels) {
			keyValueMapLabels := c.Flags().StringToString(constants.FlagLabels)
			propertiesUpdated.SetLabels(keyValueMapLabels)

			c.Verbose("Property Labels set: %v", keyValueMapLabels)
		}

		if c.Flags().Changed(constants.FlagAnnotations) {
			keyValueMapAnnotations := c.Flags().StringToString(constants.FlagAnnotations)
			propertiesUpdated.SetAnnotations(keyValueMapAnnotations)
			c.Verbose("Property Annotations set: %v", keyValueMapAnnotations)
		}

		if c.Flags().Changed(cloudapiv6.ArgLanIds) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)

			// Append existing LANs
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}

			// Add new LANs
			lanIds := c.Flags().IntSlice(cloudapiv6.ArgLanIds)
			dhcp := c.Flags().Bool(cloudapiv6.ArgDhcp)
			for _, lanId := range lanIds {
				if lanId < math.MinInt32 || lanId > math.MaxInt32 {
					return resources.K8sNodePoolForPut{}, fmt.Errorf("LAN ID %d is out of allowed int32 range [%d-%d]", lanId, math.MinInt32, math.MaxInt32)
				}
				id := int32(lanId)
				newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
					Id:   &id,
					Dhcp: &dhcp,
				})

				c.Verbose("Property Lans set: %v", id)
			}

			propertiesUpdated.SetLans(newLans)
		}

		if c.Flags().Changed(cloudapiv6.ArgPublicIps) {
			publicIps := c.Flags().StringSlice(cloudapiv6.ArgPublicIps)
			propertiesUpdated.SetPublicIps(publicIps)

			c.Verbose("Property PublicIps set: %v", publicIps)
		}

		// serverType
		if c.Flags().Changed(constants.FlagServerType) {
			serverType := c.Flags().String(constants.FlagServerType)
			propertiesUpdated.SetServerType(ionoscloud.KubernetesNodePoolServerType(serverType))

			c.Verbose("Property ServerType set: %v", serverType)
		}
	}

	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}, nil
}

func DeleteAllK8sNodepools(c *core.CommandConfig) error {
	k8sClusterId := c.Flags().String(constants.FlagClusterId)

	c.Verbose("K8sCluster ID: %v", k8sClusterId)
	c.Verbose("Getting K8sNodePools...")

	k8sNodePools, resp, err := c.CloudApiV6Services.K8s().ListNodePools(k8sClusterId)
	if err != nil {
		return err
	}

	k8sNodePoolsItems, ok := k8sNodePools.GetItemsOk()
	if !ok || k8sNodePoolsItems == nil {
		return fmt.Errorf("could not get items of Kubernetes Nodepools")
	}

	if len(*k8sNodePoolsItems) <= 0 {
		return fmt.Errorf("no Kubernetes Nodepools found")
	}

	c.Msg("K8sNodePools to be deleted:")

	var multiErr error
	for _, dc := range *k8sNodePoolsItems {
		id := dc.GetId()
		name := dc.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete K8sNodePool with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, *id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
		}

		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
