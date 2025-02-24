package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	nodeTest = resources.K8sNode{
		KubernetesNode: compute.KubernetesNode{
			Properties: compute.KubernetesNodeProperties{
				Name:       &testNodeVar,
				K8sVersion: &testNodeVar,
				PublicIP:   &testNodeVar,
				PrivateIP:  &testNodeVar,
			},
		},
	}
	nodeTestId = resources.K8sNode{
		KubernetesNode: compute.KubernetesNode{
			Id: &testNodeVar,
			Properties: compute.KubernetesNodeProperties{
				Name:       &testNodeVar,
				K8sVersion: &testNodeVar,
				PublicIP:   &testNodeVar,
				PrivateIP:  &testNodeVar,
			},
		},
	}
	nodesTestList = resources.K8sNodes{
		KubernetesNodes: compute.KubernetesNodes{
			Id: &testNodeVar,
			Items: []compute.KubernetesNode{
				nodeTestId.KubernetesNode,
				nodeTestId.KubernetesNode,
			},
		},
	}
	nodeTestGet = resources.K8sNode{
		KubernetesNode: compute.KubernetesNode{
			Id:         &testNodeVar,
			Properties: nodeTest.Properties,
			Metadata: &compute.KubernetesNodeMetadata{
				State: &testStateVar,
			},
		},
	}
	nodesTest = resources.K8sNodes{
		KubernetesNodes: compute.KubernetesNodes{
			Id:    &testNodeVar,
			Items: []compute.KubernetesNode{nodeTest.KubernetesNode},
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunK8sClusterNodesIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunK8sNodesList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		err := PreRunK8sNodesList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodesListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunK8sNodesList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sNodesListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunK8sNodesList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodesTest, &testResponse, nil)
		err := RunK8sNodeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.K8sNodes{}, &testResponse, nil)
		err := RunK8sNodeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodesTest, nil, testNodeErr)
		err := RunK8sNodeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodeTestGet, &testResponse, nil)
		err := RunK8sNodeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodeTestGet, nil, testNodeErr)
		err := RunK8sNodeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodeTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodeTestGet, nil, nil)
		err := RunK8sNodeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nodeTestGet, nil, testNodeErr)
		err := RunK8sNodeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testNodeErr)
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodesTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodesTestList, nil, testNodeErr)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.K8sNodes{}, &testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.K8sNodes{KubernetesNodes: compute.KubernetesNodes{Items: []compute.KubernetesNode{}}}, &testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nodesTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testNodeErr)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testNodeErr)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagNodepoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagClusterId), testNodeVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}
