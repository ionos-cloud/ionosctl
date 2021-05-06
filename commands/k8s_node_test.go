package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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
			},
		},
	}
	nodeTestGet = resources.K8sNode{
		KubernetesNode: ionoscloud.KubernetesNode{
			Id:         &testNodeVar,
			Properties: nodeTest.Properties,
			Metadata: &ionoscloud.KubernetesNodeMetadata{
				State: &testNodeVar,
			},
		},
	}
	nodesTest = resources.K8sNodes{
		KubernetesNodes: ionoscloud.KubernetesNodes{
			Id:    &testNodeVar,
			Items: &[]ionoscloud.KubernetesNode{nodeTest.KubernetesNode},
		},
	}
	testNodeVar = "test-node"
	testNodeErr = errors.New("node test error")
)

func TestPreRunK8sClusterNodesIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		err := PreRunK8sClusterNodesIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunK8sClusterNodesIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), "")
		err := PreRunK8sClusterNodesIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		rm.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar).Return(nodesTest, nil, nil)
		err := RunK8sNodeList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		rm.K8s.EXPECT().ListNodes(testNodeVar, testNodeVar).Return(nodesTest, nil, testNodeErr)
		err := RunK8sNodeList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, nil)
		err := RunK8sNodeGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForState), false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().GetNode(testNodeVar, testNodeVar, testNodeVar).Return(&nodeTestGet, nil, testNodeErr)
		err := RunK8sNodeGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, testNodeErr)
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.K8s.EXPECT().RecreateNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
		err := RunK8sNodeRecreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeRecreateAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = os.Stdin
		err := RunK8sNodeRecreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		rm.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, testNodeErr)
		err := RunK8sNodeDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.K8s.EXPECT().DeleteNode(testNodeVar, testNodeVar, testNodeVar).Return(nil, nil)
		err := RunK8sNodeDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sNodeDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodeId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sNodePoolId), testNodeVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgK8sClusterId), testNodeVar)
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
	viper.Set(builder.GetGlobalFlagName("node", config.ArgCols), []string{"Name"})
	getK8sNodeCols(builder.GetGlobalFlagName("node", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetK8sNodeColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("node", config.ArgCols), []string{"Unknown"})
	getK8sNodeCols(builder.GetGlobalFlagName("node", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetK8sNodesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getK8sNodesIds(w, testNodeVar, testNodeVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
