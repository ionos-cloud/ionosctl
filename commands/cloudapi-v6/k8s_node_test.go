package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
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

func TestRunK8sNodeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar).Return(nodesTest, &testResponse, nil)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar).Return(nodesTest, nil, testNodeErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, &testResponse, nil)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, testNodeErr)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, nil)
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, nil)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, testNodeErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(&testResponse, nil)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, testNodeErr)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodeId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sNodePoolId), testNodeVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgK8sClusterId), testNodeVar)
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(&testResponse, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, testNodeErr)
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
		rm.CloudApiV6Mocks.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
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
