package lan

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/die"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunK8sClusterNodePoolLanIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgLanId)
}

func PreRunK8sClusterNodePoolLanRemove(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgLanId},
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgAll},
	)
}

func RunK8sNodePoolLanList(c *core.CommandConfig) error {
	k8ss, resp, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
	)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	properties, ok := k8ss.GetPropertiesOk()
	if !ok || properties == nil {
		return fmt.Errorf("error getting node pool properties")
	}

	lans, ok := properties.GetLansOk()
	if !ok || lans == nil {
		return fmt.Errorf("error getting node pool lans")
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.K8sNodePoolLan, *lans,
		tabheaders.GetHeadersAllDefault(defaultK8sNodePoolLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolLanAdd(c *core.CommandConfig) error {
	ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
	)
	if err != nil {
		return err
	}

	input := getNewK8sNodePoolLanInfo(c, ng)
	ngNew, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.K8sNodePoolLan, getK8sNodePoolLansForPut(ngNew),
		tabheaders.GetHeadersAllDefault(defaultK8sNodePoolLanCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunK8sNodePoolLanRemove(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := RemoveAllK8sNodePoolsLans(c); err != nil {
			return err
		}

		return nil
	}

	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove node pool lan", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	ng, _, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId)
	if err != nil {
		return err
	}

	input := removeK8sNodePoolLanInfo(c, ng)
	_, resp, err := c.CloudApiV6Services.K8s().UpdateNodePool(clusterId, nodePoolId, input)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Kubernetes Node Pool Lan successfully deleted"))
	return nil
}

func RemoveAllK8sNodePoolsLans(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("K8sCluster ID: %v", clusterId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("K8sNodePool ID: %v", nodePoolId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting K8sNodePool Lans..."))

	k8sNodepool, resp, err := c.CloudApiV6Services.K8s().GetNodePool(clusterId, nodePoolId)
	if err != nil {
		return err
	}

	nodePoolProperties, ok := k8sNodepool.GetPropertiesOk()
	if !ok || nodePoolProperties == nil {
		return fmt.Errorf("could not get Node Pool properties")
	}

	lans, ok := nodePoolProperties.GetLansOk()
	if !ok || lans == nil {
		return fmt.Errorf("could not get Lans items")
	}

	if len(*lans) <= 0 {
		return fmt.Errorf("no Lans found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("K8s NodePool Lans to be removed:"))
	for _, lan := range *lans {
		if id, ok := lan.GetIdOk(); ok && id != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("K8s NodePool Lan Id: %s", string(*id)))
		}
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "remove all the K8sNodePool Lans", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing all the K8sNodePool Lans..."))

	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}
	if n, ok := nodePoolProperties.GetNodeCountOk(); ok && n != nil {
		propertiesUpdated.SetNodeCount(*n)
	}

	if n, ok := nodePoolProperties.GetAutoScalingOk(); ok && n != nil {
		propertiesUpdated.SetAutoScaling(*n)
	}

	if n, ok := nodePoolProperties.GetMaintenanceWindowOk(); ok && n != nil {
		propertiesUpdated.SetMaintenanceWindow(*n)
	}

	if n, ok := nodePoolProperties.GetK8sVersionOk(); ok && n != nil {
		propertiesUpdated.SetK8sVersion(*n)
	}

	newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
	propertiesUpdated.SetLans(newLans)
	k8sNodePoolUpdated := resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}

	_, resp, err = c.CloudApiV6Services.K8s().UpdateNodePool(clusterId, nodePoolId, k8sNodePoolUpdated)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Kubernetes Node Pool Lans successfully deleted"))
	return nil
}

func getNewK8sNodePoolLanInfo(c *core.CommandConfig, oldNg *resources.K8sNodePool) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}

	if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
		if n, ok := properties.GetNodeCountOk(); ok && n != nil {
			propertiesUpdated.SetNodeCount(*n)
		}

		if n, ok := properties.GetAutoScalingOk(); ok && n != nil {
			propertiesUpdated.SetAutoScaling(*n)
		}

		if n, ok := properties.GetMaintenanceWindowOk(); ok && n != nil {
			propertiesUpdated.SetMaintenanceWindow(*n)
		}

		if n, ok := properties.GetK8sVersionOk(); ok && n != nil {
			propertiesUpdated.SetK8sVersion(*n)
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)
			// Append existing LANs
			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					newLans = append(newLans, existingLan)
				}
			}

			// Add new LANs
			lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
			dhcp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgDhcp))
			newLan := ionoscloud.KubernetesNodePoolLan{
				Id:   &lanId,
				Dhcp: &dhcp,
			}

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Adding a Kubernetes NodePool LAN with id: %v and dhcp: %v", lanId, dhcp))

			if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork)) {
				network := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgNetwork))
				gatewayIp := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgGatewayIp))

				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property Network set: %v", network))
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property GatewayIp set: %v", gatewayIp))

				if len(network) != len(gatewayIp) {
					die.Die(fmt.Sprintf("Flags %s, %s have different number of arguments, must be the same", cloudapiv6.ArgNetwork, cloudapiv6.ArgGatewayIp))
				}

				routes := make([]ionoscloud.KubernetesNodePoolLanRoutes, 0)
				for i, net := range network {
					routes = append(routes,
						ionoscloud.KubernetesNodePoolLanRoutes{
							Network:   pointer.From(net), // Copy the loop variable and take its address. See #289 - always same address would be used
							GatewayIp: &gatewayIp[i],
						},
					)
				}

				newLan.SetRoutes(routes)
			}

			newLans = append(newLans, newLan)
			propertiesUpdated.SetLans(newLans)
		}
	}

	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func removeK8sNodePoolLanInfo(c *core.CommandConfig, oldNg *resources.K8sNodePool) resources.K8sNodePoolForPut {
	propertiesUpdated := resources.K8sNodePoolPropertiesForPut{}

	if properties, ok := oldNg.GetPropertiesOk(); ok && properties != nil {
		if n, ok := properties.GetNodeCountOk(); ok && n != nil {
			propertiesUpdated.SetNodeCount(*n)
		}

		if n, ok := properties.GetAutoScalingOk(); ok && n != nil {
			propertiesUpdated.SetAutoScaling(*n)
		}

		if n, ok := properties.GetMaintenanceWindowOk(); ok && n != nil {
			propertiesUpdated.SetMaintenanceWindow(*n)
		}

		if n, ok := properties.GetK8sVersionOk(); ok && n != nil {
			propertiesUpdated.SetK8sVersion(*n)
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLanId)) {
			lanId := viper.GetInt32(core.GetFlagName(c.NS, cloudapiv6.ArgLanId))
			newLans := make([]ionoscloud.KubernetesNodePoolLan, 0)

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Removing a Kubernetes NodePool LAN with id: %v", lanId))

			if existingLans, ok := properties.GetLansOk(); ok && existingLans != nil {
				for _, existingLan := range *existingLans {
					if id, ok := existingLan.GetIdOk(); ok && id != nil {
						if *id != lanId {
							newLans = append(newLans, existingLan)
						}
					}
				}
			}

			propertiesUpdated.SetLans(newLans)
		}
	}

	return resources.K8sNodePoolForPut{
		KubernetesNodePoolForPut: ionoscloud.KubernetesNodePoolForPut{
			Properties: &propertiesUpdated.KubernetesNodePoolPropertiesForPut,
		},
	}
}

func getK8sNodePoolLansForPut(ng *resources.K8sNodePool) []ionoscloud.KubernetesNodePoolLan {
	ss := make([]ionoscloud.KubernetesNodePoolLan, 0)

	if ng != nil {
		if properties, ok := ng.GetPropertiesOk(); ok && properties != nil {
			if lans, ok := properties.GetLansOk(); ok && lans != nil {
				for _, lanItem := range *lans {
					ss = append(ss, lanItem)
				}
			}
		}
	}

	return ss
}

func PreRunK8sClusterNodePoolIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId)
}
