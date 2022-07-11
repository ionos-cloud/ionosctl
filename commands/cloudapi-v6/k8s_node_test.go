package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	nodeTest = resources.K8sNode{
		KubernetesNode: ionoscloud.KubernetesNode{
			Properties: &ionoscloud.KubernetesNodeProperties{
				Name:       &testNodeVar,
				K8sVersion: &testNodeVar,
				PublicIP:   &testNodeVar,
				PrivateIP:  &testNodeVar,
			},
		},
	}
	nodeTestId = resources.K8sNode{
		KubernetesNode: ionoscloud.KubernetesNode{
			Id: &testNodeVar,
			Properties: &ionoscloud.KubernetesNodeProperties{
				Name:       &testNodeVar,
				K8sVersion: &testNodeVar,
				PublicIP:   &testNodeVar,
				PrivateIP:  &testNodeVar,
			},
		},
	}
	nodesTestList = resources.K8sNodes{
		KubernetesNodes: ionoscloud.KubernetesNodes{
			Id: &testNodeVar,
			Items: &[]ionoscloud.KubernetesNode{
				nodeTestId.KubernetesNode,
				nodeTestId.KubernetesNode,
			},
		},
	}
	nodeTestGet = resources.K8sNode{
		KubernetesNode: ionoscloud.KubernetesNode{
			Id:         &testNodeVar,
			Properties: nodeTest.Properties,
			Metadata: &ionoscloud.KubernetesNodeMetadata{
				State: &testStateVar,
			},
		},
	}
	nodesTest = resources.K8sNodes{
		KubernetesNodes: ionoscloud.KubernetesNodes{
			Id:    &testNodeVar,
			Items: &[]ionoscloud.KubernetesNode{nodeTest.KubernetesNode},
		},
	}
	testNodeVar  = "test-node"
	testStateVar = "ACTIVE"
	testNodeErr  = errors.New("node test error")
)

func TestPreRunK8sClusterNodesIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		err := PreRunK8sClusterNodesIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodesIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunK8sClusterNodesIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sNodesList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		err := PreRunK8sNodesList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodesListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunK8sNodesList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodesListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunK8sNodesList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(nodesTest, &testResponse, nil)
		err := RunK8sNodeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, testListQueryParam).Return(resources.K8sNodes{}, &testResponse, nil)
		err := RunK8sNodeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(nodesTest, nil, testNodeErr)
		err := RunK8sNodeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&nodeTestGet, &testResponse, nil)
		err := RunK8sNodeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&nodeTestGet, nil, testNodeErr)
		err := RunK8sNodeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&nodeTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&nodeTestGet, nil, nil)
		err := RunK8sNodeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&nodeTestGet, nil, testNodeErr)
		err := RunK8sNodeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(nil, testNodeErr)
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(nil, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = os.Stdin
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(nodesTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(nodesTestList, nil, testNodeErr)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(resources.K8sNodes{}, &testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(
			resources.K8sNodes{KubernetesNodes: ionoscloud.KubernetesNodes{Items: &[]ionoscloud.KubernetesNode{}}}, &testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, resources.ListQueryParams{}).Return(nodesTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, testNodeErr)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(nil, testNodeErr)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, resources.QueryParams{}).Return(nil, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = os.Stdin
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetK8sNodeCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("node", config.ArgCols), []string{"Name"})
	getK8sNodeCols(core.GetGlobalFlagName("node", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodeColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("node", config.ArgCols), []string{"Unknown"})
	getK8sNodeCols(core.GetGlobalFlagName("node", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
