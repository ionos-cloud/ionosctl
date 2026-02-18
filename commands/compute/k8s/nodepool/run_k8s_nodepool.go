package nodepool

import (
	"context"
	"errors"
	"fmt"
	"time"

	k8scluster "github.com/ionos-cloud/ionosctl/v6/commands/compute/k8s/cluster"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func RunK8sNodePoolCreateFromJSON(c *core.CommandConfig, propertiesFromJson ionoscloud.KubernetesNodePoolForPost) error {
	np, _, err := client.Must().CloudClient.KubernetesApi.K8sNodepoolsPost(context.Background(),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))).KubernetesNodePool(propertiesFromJson).Execute()
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Creating K8s NodePool in K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	u, resp, err := c.CloudApiV6Services.K8s().CreateNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), *newNodePool)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	return handleApiResponseK8sNodepoolCreate(c, u.KubernetesNodePool)
}

func handleApiResponseK8sNodepoolCreate(c *core.CommandConfig, pool ionoscloud.KubernetesNodePool) error {
	uConverted, err := resource2table.ConvertK8sNodepoolToTable(pool)
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if id, ok := pool.GetIdOk(); ok && id != nil {
			if err = waitfor.WaitForState(c, waiter.K8sNodePoolStateInterrogator, *id); err != nil {
				return err
			}
			if pool, _, err = client.Must().CloudClient.KubernetesApi.K8sNodepoolsFindById(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), *id).
				Execute(); err != nil {
				return err
			}
		} else {
			return errors.New("error getting new K8s Node Pool id")
		}
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(pool, uConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
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

	var allNodePools []ionoscloud.KubernetesNodePools
	var allNodePoolsConverted []map[string]interface{}
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

		clusterId, ok := cluster.GetIdOk()
		if !ok || clusterId == nil {
			continue
		}

		for _, node := range *items {
			temp, err := resource2table.ConvertK8sNodepoolToTable(node)
			if err != nil {
				return fmt.Errorf("failed to convert from JSON to Table format: %w", err)
			}

			allNodePoolsConverted = append(allNodePoolsConverted, temp[0])
		}

		allNodePools = append(allNodePools, nodePools.KubernetesNodePools)
		totalTime += resp.RequestTime
	}

	if totalTime != time.Duration(0) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, totalTime))
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(allNodePools, allNodePoolsConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolList(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		return RunK8sNodePoolListAll(c)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Getting K8s NodePools from K8s Cluster with ID: %v", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	k8ss, resp, err := c.CloudApiV6Services.K8s().ListNodePools(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	k8ssConverted, err := resource2table.ConvertK8sNodepoolsToTable(k8ss.KubernetesNodePools)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(k8ss.KubernetesNodePools, k8ssConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolGet(c *core.CommandConfig) error {
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	if err := waitfor.WaitForState(c, waiter.K8sNodePoolStateInterrogator, k8sNodePoolId); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"K8s node pool with id: %v from K8s Cluster with id: %v is getting...", k8sNodePoolId, k8sClusterId))

	u, resp, err := c.CloudApiV6Services.K8s().GetNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	uConverted, err := resource2table.ConvertK8sNodepoolToTable(u.KubernetesNodePool)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(u.KubernetesNodePool, uConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolUpdate(c *core.CommandConfig) error {
	oldNodePool, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)))
	if err != nil {
		return err
	}

	newNodePool := getNewK8sNodePoolUpdated(oldNodePool, c)
	_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)), newNodePool)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForState(c, waiter.K8sNodePoolStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))); err != nil {
		return err
	}

	newNodePoolUpdated, _, err := c.CloudApiV6Services.K8s().GetNodePool(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)))

	newNodePoolUpdatedConverted, err := resource2table.ConvertK8sNodepoolToTable(newNodePoolUpdated.KubernetesNodePool)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(newNodePoolUpdated.KubernetesNodePool, newNodePoolUpdatedConverted,
		tabheaders.GetHeaders(allK8sNodePoolCols, defaultK8sNodePoolCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolDelete(c *core.CommandConfig) error {
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllK8sNodepools(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete k8s node pool", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Starting deleting K8s node pool with id: %v from K8s Cluster with id: %v...", k8sNodePoolId, k8sClusterId))

	resp, err := c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, k8sNodePoolId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Kubernetes Nodepool successfully deleted"))
	return nil
}

func getNewK8sNodePool(c *core.CommandConfig) (*resources.K8sNodePoolForPost, error) {
	var k8sversion string

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))
	} else {
		clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

		k8sCluster, _, err := c.CloudApiV6Services.K8s().GetCluster(clusterId)
		if err != nil {
			return nil, fmt.Errorf("failed to get k8s cluster to fetch default version: %w", err)
		}

		k8sVerPtr := k8sCluster.GetProperties().GetK8sVersion()
		if k8sVerPtr == nil {
			return nil, errors.New("k8s version is not set on the cluster")
		}

		k8sversion = *k8sVerPtr
	}

	ramSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagRam)), utils2.MegaBytes)
	if err != nil {
		return nil, err
	}

	storageSize, err := utils2.ConvertSize(viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageSize)), utils2.GigaBytes)
	if err != nil {
		return nil, err
	}

	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	nodeCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagNodeCount))
	dcId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgDataCenterId))
	cpuFamily := viper.GetString(core.GetFlagName(c.NS, constants.FlagCpuFamily))
	cores := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagCores))
	availabilityZone := viper.GetString(core.GetFlagName(c.NS, constants.FlagAvailabilityZone))
	storageType := viper.GetString(core.GetFlagName(c.NS, constants.FlagStorageType))

	// Set Properties
	nodePoolProperties := ionoscloud.KubernetesNodePoolPropertiesForPost{}
	nodePoolProperties.SetName(name)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Name set: %v", name))

	nodePoolProperties.SetK8sVersion(k8sversion)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", k8sversion))

	nodePoolProperties.SetNodeCount(nodeCount)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NodeCount set: %v", nodeCount))

	nodePoolProperties.SetDatacenterId(dcId)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DatacenterId set: %v", dcId))

	if fn := core.GetFlagName(c.NS, constants.FlagCpuFamily); viper.IsSet(fn) {
		nodePoolProperties.SetCpuFamily(viper.GetString(fn))
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property CPU Family set: %v", cpuFamily))
	}

	if fn := core.GetFlagName(c.NS, constants.FlagServerType); viper.IsSet(fn) {
		nodePoolProperties.SetServerType(ionoscloud.KubernetesNodePoolServerType(viper.GetString(fn)))
	}

	nodePoolProperties.SetCoresCount(cores)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property CoresCount set: %v", cores))

	nodePoolProperties.SetRamSize(int32(ramSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property RAM Size set: %vMB", int32(ramSize)))

	nodePoolProperties.SetAvailabilityZone(availabilityZone)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Availability Zone set: %v", availabilityZone))

	nodePoolProperties.SetStorageSize(int32(storageSize))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Storage Size set: %vGB", int32(storageSize)))

	nodePoolProperties.SetStorageType(storageType)
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Storage Type set: %v", storageType))

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLabels)) {
		keyValueMapLabels := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagLabels))
		nodePoolProperties.SetLabels(keyValueMapLabels)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Labels set: %v", keyValueMapLabels))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAnnotations)) {
		keyValueMapAnnotations := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagAnnotations))
		nodePoolProperties.SetAnnotations(keyValueMapAnnotations)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Annotations set: %v", keyValueMapAnnotations))
	}

	// Add LANs
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds)) {
		newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
		lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds))
		dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))

		for _, lanId := range lanIds {
			id := int32(lanId)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Lan ID set: %v", id))
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Dhcp set: %v", dhcp))

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

func getNewK8sNodePoolUpdated(oldNodePool *resources.K8sNodePool, c *core.CommandConfig) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}

	if properties, ok := oldNodePool.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
			vers := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(vers)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", vers))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagNodeCount)) {
			nodeCount := viper.GetInt32(core.GetFlagName(c.NS, constants.FlagNodeCount))
			propertiesUpdated.SetNodeCount(nodeCount)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property NodeCount set: %v", nodeCount))
		} else {
			if n, ok := properties.GetNodeCountOk(); ok && n != nil {
				propertiesUpdated.SetNodeCount(*n)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount)) {
			var minCount, maxCount int32

			autoScaling := properties.GetAutoScaling()
			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount)) {
				minCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMinNodeCount))

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property MinNodeCount set: %v", minCount))
			} else {
				if m, ok := autoScaling.GetMinNodeCountOk(); ok && m != nil {
					minCount = *m
				}
			}

			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount)) {
				maxCount = viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaxNodeCount))

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property MaxNodeCount set: %v", maxCount))
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

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceTime)) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := k8scluster.GetMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})

				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sAnnotationValue))
			propertiesUpdated.SetAnnotations(map[string]string{
				key: value,
			})

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Annotations set: key: %v, value: %v", key, value))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey)) &&
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue)) {
			key := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelKey))
			value := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLabelValue))
			propertiesUpdated.SetLabels(map[string]string{
				key: value,
			})

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Labels set: key: %v, value: %v", key, value))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagLabels)) {
			keyValueMapLabels := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagLabels))
			propertiesUpdated.SetLabels(keyValueMapLabels)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Labels set: %v", keyValueMapLabels))
		}

		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagAnnotations)) {
			keyValueMapAnnotations := viper.GetStringMapString(core.GetFlagName(c.NS, constants.FlagAnnotations))
			propertiesUpdated.SetAnnotations(keyValueMapAnnotations)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Annotations set: %v", keyValueMapAnnotations))
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)

			// Append existing LANs
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}

			// Add new LANs
			lanIds := viper.GetIntSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLanIds))
			dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
			for _, lanId := range lanIds {
				id := int32(lanId)
				newLans = append(newLans, ionoscloud.KubernetesNodePoolLan{
					Id:   &id,
					Dhcp: &dhcp,
				})

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Lans set: %v", id))
			}

			propertiesUpdated.SetLans(newLans)
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps)) {
			publicIps := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgPublicIps))
			propertiesUpdated.SetPublicIps(publicIps)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property PublicIps set: %v", publicIps))
		}

		// serverType
		if viper.IsSet(core.GetFlagName(c.NS, constants.FlagServerType)) {
			serverType := viper.GetString(core.GetFlagName(c.NS, constants.FlagServerType))
			propertiesUpdated.SetServerType(ionoscloud.KubernetesNodePoolServerType(serverType))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property ServerType set: %v", serverType))
		}
	}

	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func DeleteAllK8sNodepools(c *core.CommandConfig) error {
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("K8sCluster ID: %v", k8sClusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting K8sNodePools..."))

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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("K8sNodePools to be deleted:"))

	var multiErr error
	for _, dc := range *k8sNodePoolsItems {
		id := dc.GetId()
		name := dc.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete K8sNodePool with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.K8s().DeleteNodePool(k8sClusterId, *id)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}

		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrWaitDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}
