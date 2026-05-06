package node

import (
	"errors"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
)

func PreRunK8sNodesList(c *core.PreCommandConfig) error {
	if err := core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId); err != nil {
		return err
	}
	return nil
}

func PreRunK8sClusterNodesIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgK8sNodeId)
}

func PreRunK8sClusterNodesIdsAll(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgAll},
		[]string{constants.FlagClusterId, constants.FlagNodepoolId, cloudapiv6.ArgK8sNodeId},
	)
}

func RunK8sNodeList(c *core.CommandConfig) error {
	c.Verbose("Listing Nodes from K8s NodePool ID: %v from K8s Cluster ID: %v",
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))

	k8ss, resp, err := c.CloudApiV6Services.K8s().ListNodes(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}

	if err != nil {
		return err
	}

	return c.Printer(allK8sNodeCols).Prefix("items").Print(k8ss.KubernetesNodes)
}

func RunK8sNodeGet(c *core.CommandConfig) error {
	if err := waitfor.WaitForState(c, waiter.K8sNodeStateInterrogator, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeId))); err != nil {
		return err
	}

	c.Verbose("Getting K8s Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v......",
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)))

	u, resp, err := c.CloudApiV6Services.K8s().GetNode(
		viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId)),
		viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId)),
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeId)),
	)
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allK8sNodeCols).Print(u.KubernetesNode)
}

func RunK8sNodeRecreate(c *core.CommandConfig) error {
	if !confirm.FAsk(c.Command.Command.InOrStdin(), "recreate k8s node", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	k8sClusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	k8sNodePoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
	k8sNodeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeId))

	c.Verbose("K8sClusterId: %v, K8sNodePoolId: %v, K8sNodeId: %v",
		k8sClusterId, k8sNodePoolId, k8sNodeId)
	c.Verbose("Recreating Node...")

	resp, err := c.CloudApiV6Services.K8s().RecreateNode(k8sClusterId, k8sNodePoolId, k8sNodeId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Status: Command node recreate has been successfully executed")

	return nil
}

func RunK8sNodeDelete(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodepoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))
	nodeId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgK8sNodeId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllK8sNodes(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete k8s node", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v...", nodeId, nodepoolId, clusterId)

	resp, err := c.CloudApiV6Services.K8s().DeleteNode(clusterId, nodepoolId, nodeId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Kubernetes Node successfully deleted")
	return nil
}

func DeleteAllK8sNodes(c *core.CommandConfig) error {
	clusterId := viper.GetString(core.GetFlagName(c.NS, constants.FlagClusterId))
	nodepoolId := viper.GetString(core.GetFlagName(c.NS, constants.FlagNodepoolId))

	c.Verbose("K8sCluster ID: %v", clusterId)
	c.Verbose("K8sNodePool ID: %v", nodepoolId)
	c.Verbose("Getting K8sNodes...")

	k8sNodes, resp, err := c.CloudApiV6Services.K8s().ListNodes(clusterId, nodepoolId)
	if err != nil {
		return err
	}

	k8sNodesItems, ok := k8sNodes.GetItemsOk()
	if !ok || k8sNodesItems == nil {
		return fmt.Errorf("could not get items of Kubernetes Nodes")
	}

	if len(*k8sNodesItems) <= 0 {
		return fmt.Errorf("no Kubernetes Nodes found")
	}

	var multiErr error
	for _, dc := range *k8sNodesItems {
		id := dc.GetId()

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Deleting Node with ID: %v from K8s NodePool ID: %v from K8s Cluster ID: %v...", *id, nodepoolId, clusterId), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.K8s().DeleteNode(clusterId, nodepoolId, *id)
		if resp != nil && request.GetId(resp) != "" {
			c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
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
