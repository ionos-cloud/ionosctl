package cluster

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	k8ssConverted, err := resource2table.ConvertK8sClustersToTable(k8ss.KubernetesClusters)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(k8ss.KubernetesClusters, k8ssConverted,
		tabheaders.GetHeaders(allK8sClusterCols, defaultK8sClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sClusterGet(c *core.CommandConfig) error {
	if err := waitfor.WaitForState(c, waiter.K8sClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"K8s cluster with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	u, resp, err := c.CloudApiV6Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	uConverted, err := resource2table.ConvertK8sClusterToTable(u.KubernetesCluster)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(u.KubernetesCluster, uConverted,
		tabheaders.GetHeaders(allK8sClusterCols, defaultK8sClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sClusterCreate(c *core.CommandConfig) error {
	newCluster, err := getNewK8sCluster(c)
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Creating K8s Cluster..."))

	u, resp, err := c.CloudApiV6Services.K8s().CreateCluster(*newCluster)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		id, ok := u.GetIdOk()
		if !ok || id == nil {
			return fmt.Errorf("error getting new K8s Cluster id")
		}

		if err = waitfor.WaitForState(c, waiter.K8sClusterStateInterrogator, *id); err != nil {
			return err
		}

		if u, _, err = c.CloudApiV6Services.K8s().GetCluster(*id); err != nil {
			return err
		}

	}

	uConverted, err := resource2table.ConvertK8sClusterToTable(u.KubernetesCluster)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(u.KubernetesCluster, uConverted,
		tabheaders.GetHeaders(allK8sClusterCols, defaultK8sClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sClusterUpdate(c *core.CommandConfig) error {
	oldCluster, _, err := c.CloudApiV6Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))
	if err != nil {
		return err
	}

	newCluster := getK8sClusterInfo(oldCluster, c)

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Updating K8s cluster with ID: %v...", viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))))

	k8sUpd, resp, err := c.CloudApiV6Services.K8s().UpdateCluster(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)), newCluster)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForState)) {
		if err = waitfor.WaitForState(c, waiter.K8sClusterStateInterrogator, viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
		if k8sUpd, _, err = c.CloudApiV6Services.K8s().GetCluster(viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))); err != nil {
			return err
		}
	}

	k8sUpdConverted, err := resource2table.ConvertK8sClusterToTable(k8sUpd.KubernetesCluster)
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutputPreconverted(k8sUpd.KubernetesCluster, k8sUpdConverted,
		tabheaders.GetHeaders(allK8sClusterCols, defaultK8sClusterCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sClusterDelete(c *core.CommandConfig) error {
	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllK8sClusters(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete k8s cluster", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Starting deleting K8s cluster with id: %v...", k8sClusterId))

	resp, err := c.CloudApiV6Services.K8s().DeleteCluster(k8sClusterId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Kubernetes Cluster successfully deleted"))
	return nil
}

func getNewK8sCluster(c *core.CommandConfig) (*resources.K8sClusterForPost, error) {
	var (
		k8sversion string
		err        error
	)

	proper := resources.K8sClusterPropertiesForPost{}
	proper.SetName(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName)))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Property Name set: %v", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))))

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
		k8sversion = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", k8sversion))
	} else {
		if k8sversion, err = getK8sVersion(c); err != nil {
			return nil, err
		}
	}

	proper.SetK8sVersion(k8sversion)

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket)) {
		s3buckets := make([]ionoscloud.S3Bucket, 0)
		name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket))
		s3buckets = append(s3buckets, ionoscloud.S3Bucket{
			Name: &name,
		})
		proper.SetS3Buckets(s3buckets)

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property S3Buckets set: %v", s3buckets))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets)) {
		proper.SetApiSubnetAllowList(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets)))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"Property ApiSubnetAllowList set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgPublic)) {
		proper.SetPublic(viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgPublic)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)) {
		proper.SetLocation(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayIp)) {
		proper.SetNatGatewayIp(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgNatGatewayIp)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagNodeSubnet)) {
		proper.SetNodeSubnet(viper.GetString(core.GetFlagName(c.NS, constants.FlagNodeSubnet)))
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
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
			n := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
			propertiesUpdated.SetName(n)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Name set: %v", n))
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				propertiesUpdated.SetName(*name)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion)) {
			v := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sVersion))
			propertiesUpdated.SetK8sVersion(v)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property K8sVersion set: %v", v))
		} else {
			if vers, ok := properties.GetK8sVersionOk(); ok && vers != nil {
				propertiesUpdated.SetK8sVersion(*vers)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket)) {
			s3buckets := make([]ionoscloud.S3Bucket, 0)
			for _, name := range viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgS3Bucket)) {
				s3buckets = append(s3buckets, ionoscloud.S3Bucket{
					Name: &name,
				})
			}

			propertiesUpdated.SetS3Buckets(s3buckets)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property S3Buckets set: %v", s3buckets))
		} else {
			if bucketsOk, ok := properties.GetS3BucketsOk(); ok && bucketsOk != nil {
				propertiesUpdated.SetS3Buckets(*bucketsOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets)) {
			propertiesUpdated.SetApiSubnetAllowList(viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets)))

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Property ApiSubnetAllowList set: %v", viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgApiSubnets))))
		} else {
			if subnetAllowListOk, ok := properties.GetApiSubnetAllowListOk(); ok && subnetAllowListOk != nil {
				propertiesUpdated.SetApiSubnetAllowList(*subnetAllowListOk)
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceDay)) ||
			viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceTime)) {
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
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting K8sClusters..."))

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

func getK8sVersion(c *core.CommandConfig) (string, error) {
	k8sversion, resp, err := c.CloudApiV6Services.K8s().GetVersion()
	if err != nil {
		return "", err
	}

	k8sversion = strings.ReplaceAll(k8sversion, "\"", "")
	k8sversion = strings.ReplaceAll(k8sversion, "\n", "")
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceDay)) {
		day = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceDay))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property DayOfTheWeek of MaintenanceWindow set: %v", day))
	} else {
		if d, ok := maintenance.GetDayOfTheWeekOk(); ok && d != nil {
			day = *d
		}
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceTime)) {
		time = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sMaintenanceTime))

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Time of MaintenanceWindow set: %v", time))
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
