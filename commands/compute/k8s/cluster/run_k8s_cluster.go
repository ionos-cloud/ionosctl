package cluster

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunK8sClusterList(c *core.PreCommandConfig) error {
	return nil
}

func PreRunK8sClusterId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId)
}

func PreRunK8sClusterDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId},
		[]string{cloudapiv6.ArgAll},
	)
}

func RunK8sClusterList(c *core.CommandConfig) error {

	k8ss, resp, err := c.CloudApiV6Services.K8s().ListClusters()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sClusterCols).Prefix("items").Print(k8ss.KubernetesClusters)
}

func RunK8sClusterGet(c *core.CommandConfig) error {
	c.Verbose("K8s cluster with id: %v is getting...", c.Flags().String(constants.FlagClusterId))

	u, resp, err := c.CloudApiV6Services.K8s().GetCluster(c.Flags().String(constants.FlagClusterId))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sClusterCols).Print(u.KubernetesCluster)
}

func RunK8sClusterCreate(c *core.CommandConfig) error {
	newCluster, err := getNewK8sCluster(c)
	if err != nil {
		return err
	}

	c.Verbose("Creating K8s Cluster...")

	u, resp, err := c.CloudApiV6Services.K8s().CreateCluster(*newCluster)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sClusterCols).Print(u.KubernetesCluster)
}

func RunK8sClusterUpdate(c *core.CommandConfig) error {
	oldCluster, _, err := c.CloudApiV6Services.K8s().GetCluster(c.Flags().String(constants.FlagClusterId))
	if err != nil {
		return err
	}

	newCluster := getK8sClusterInfo(oldCluster, c)

	c.Verbose("Updating K8s cluster with ID: %v...", c.Flags().String(constants.FlagClusterId))

	k8sUpd, resp, err := c.CloudApiV6Services.K8s().UpdateCluster(c.Flags().String(constants.FlagClusterId), newCluster)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sClusterCols).Print(k8sUpd.KubernetesCluster)
}

func RunK8sClusterDelete(c *core.CommandConfig) error {
	k8sClusterId := c.Flags().String(constants.FlagClusterId)

	if c.Flags().Bool(cloudapiv6.ArgAll) {
		if err := DeleteAllK8sClusters(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete k8s cluster", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting K8s cluster with id: %v...", k8sClusterId)

	resp, err := c.CloudApiV6Services.K8s().DeleteCluster(k8sClusterId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Kubernetes Cluster successfully deleted")
	return nil
}

func getNewK8sCluster(c *core.CommandConfig) (*resources.K8sClusterForPost, error) {
	var (
		k8sversion string
		err        error
	)

	proper := resources.K8sClusterPropertiesForPost{}
	proper.SetName(c.Flags().String(cloudapiv6.ArgName))

	c.Verbose("Property Name set: %v", c.Flags().String(cloudapiv6.ArgName))

	if c.Flags().Changed(cloudapiv6.ArgK8sVersion) {
		k8sversion = c.Flags().String(cloudapiv6.ArgK8sVersion)

		c.Verbose("Property K8sVersion set: %v", k8sversion)
	} else {
		if k8sversion, err = GetK8sVersion(c); err != nil {
			return nil, err
		}
	}

	proper.SetK8sVersion(k8sversion)

	if c.Flags().Changed(cloudapiv6.ArgS3Bucket) {
		s3buckets := make([]ionoscloud.S3Bucket, 0)
		name := c.Flags().String(cloudapiv6.ArgS3Bucket)
		s3buckets = append(s3buckets, ionoscloud.S3Bucket{
			Name: &name,
		})
		proper.SetS3Buckets(s3buckets)

		c.Verbose("Property S3Buckets set: %v", s3buckets)
	}

	if c.Flags().Changed(cloudapiv6.ArgApiSubnets) {
		proper.SetApiSubnetAllowList(c.Flags().StringSlice(cloudapiv6.ArgApiSubnets))

		c.Verbose("Property ApiSubnetAllowList set: %v", c.Flags().StringSlice(cloudapiv6.ArgApiSubnets))
	}

	if c.Flags().Changed(cloudapiv6.ArgPublic) {
		proper.SetPublic(c.Flags().Bool(cloudapiv6.ArgPublic))
	}

	if c.Flags().Changed(cloudapiv6.ArgLocation) {
		proper.SetLocation(c.Flags().String(cloudapiv6.ArgLocation))
	}

	if c.Flags().Changed(cloudapiv6.ArgNatGatewayIp) {
		proper.SetNatGatewayIp(c.Flags().String(cloudapiv6.ArgNatGatewayIp))
	}

	if c.Flags().Changed(constants.FlagNodeSubnet) {
		proper.SetNodeSubnet(c.Flags().String(constants.FlagNodeSubnet))
	}

	return &resources.K8sClusterForPost{
		KubernetesClusterForPost: ionoscloud.KubernetesClusterForPost{
			Properties: &proper.KubernetesClusterPropertiesForPost,
		},
	}, nil
}

func getK8sClusterInfo(oldUser *resources.K8sCluster, c *core.CommandConfig) resources.K8sClusterForPut {
	propertiesUpdated := resources.K8sClusterPropertiesForPut{}

	if properties, ok := oldUser.GetPropertiesOk(); ok && properties != nil {
		if c.Flags().Changed(cloudapiv6.ArgName) {
			n := c.Flags().String(cloudapiv6.ArgName)
			propertiesUpdated.SetName(n)

			c.Verbose("Property Name set: %v", n)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				propertiesUpdated.SetName(*name)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgK8sVersion) {
			v := c.Flags().String(cloudapiv6.ArgK8sVersion)
			propertiesUpdated.SetK8sVersion(v)

			c.Verbose("Property K8sVersion set: %v", v)
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgS3Bucket) {
			name := c.Flags().String(cloudapiv6.ArgS3Bucket)
			s3buckets := []ionoscloud.S3Bucket{{Name: &name}}

			propertiesUpdated.SetS3Buckets(s3buckets)
			c.Verbose("Property S3Buckets set: %v", s3buckets)
		} else {
			if bucketsOk, ok := properties.GetS3BucketsOk(); ok && bucketsOk != nil {
				propertiesUpdated.SetS3Buckets(*bucketsOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgApiSubnets) {
			propertiesUpdated.SetApiSubnetAllowList(c.Flags().StringSlice(cloudapiv6.ArgApiSubnets))

			c.Verbose("Property ApiSubnetAllowList set: %v", c.Flags().StringSlice(cloudapiv6.ArgApiSubnets))
		} else {
			if subnetAllowListOk, ok := properties.GetApiSubnetAllowListOk(); ok && subnetAllowListOk != nil {
				propertiesUpdated.SetApiSubnetAllowList(*subnetAllowListOk)
			}
		}

		if c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceDay) ||
			c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceTime) {
			if maintenance, ok := properties.GetMaintenanceWindowOk(); ok && maintenance != nil {
				newMaintenanceWindow := GetMaintenanceInfo(c, &resources.K8sMaintenanceWindow{
					KubernetesMaintenanceWindow: *maintenance,
				})
				propertiesUpdated.SetMaintenanceWindow(newMaintenanceWindow.KubernetesMaintenanceWindow)
			}
		}
	}

	return resources.K8sClusterForPut{
		KubernetesClusterForPut: ionoscloud.KubernetesClusterForPut{
			Properties: &propertiesUpdated.KubernetesClusterPropertiesForPut,
		},
	}
}

func DeleteAllK8sClusters(c *core.CommandConfig) error {
	c.Verbose("Getting K8sClusters...")

	k8Clusters, resp, err := c.CloudApiV6Services.K8s().ListClusters()
	if err != nil {
		return err
	}

	k8sClustersItems, ok := k8Clusters.GetItemsOk()
	if !ok || k8sClustersItems == nil {
		return fmt.Errorf("could not get items of K8sClusters")
	}

	if len(*k8sClustersItems) <= 0 {
		return fmt.Errorf("no K8sClusters found")
	}

	var multiErr error
	for _, k8sCluster := range *k8sClustersItems {
		id := k8sCluster.GetId()
		name := k8sCluster.GetProperties().GetName()
		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the K8sCluster with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.K8s().DeleteCluster(*id)
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

func GetK8sVersion(c *core.CommandConfig) (string, error) {
	k8sversion, resp, err := c.CloudApiV6Services.K8s().GetVersion()
	if err != nil {
		return "", err
	}

	k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
	k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	return k8sversion, nil
}

func GetK8sClusters(k8ss resources.K8sClusters) []resources.K8sCluster {
	u := make([]resources.K8sCluster, 0)

	if items, ok := k8ss.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.K8sCluster{KubernetesCluster: item})
		}
	}

	return u
}

func GetMaintenanceInfo(c *core.CommandConfig, maintenance *resources.K8sMaintenanceWindow) resources.K8sMaintenanceWindow {
	var day, time string
	if c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceDay) {
		day = c.Flags().String(cloudapiv6.ArgK8sMaintenanceDay)

		c.Verbose("Property DayOfTheWeek of MaintenanceWindow set: %v", day)
	} else {
		if d, ok := maintenance.GetDayOfTheWeekOk(); ok && d != nil {
			day = *d
		}
	}

	if c.Flags().Changed(cloudapiv6.ArgK8sMaintenanceTime) {
		time = c.Flags().String(cloudapiv6.ArgK8sMaintenanceTime)

		c.Verbose("Property Time of MaintenanceWindow set: %v", time)
	} else {
		if t, ok := maintenance.GetTimeOk(); ok && t != nil {
			time = *t
		}
	}

	return resources.K8sMaintenanceWindow{
		KubernetesMaintenanceWindow: ionoscloud.KubernetesMaintenanceWindow{
			DayOfTheWeek: &day,
			Time:         &time,
		},
	}
}
