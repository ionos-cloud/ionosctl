package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testK8sVersionsVar = []string{"test-k8s-version"}
	testK8sVersionVar  = "test-k8s-version"
	testK8sVersionErr  = errors.New("k8s-version test error")
)

func TestRunK8sVersionList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListVersions().Return(testK8sVersionsVar, &testResponse, nil)
		err := RunK8sVersionList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sVersionListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.K8s.EXPECT().ListVersions().Return(testK8sVersionsVar, nil, testK8sVersionErr)
		err := RunK8sVersionList(cfg)
		assert.Error(t, err)
	})
}

func TestRunK8sVersionGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetVersion().Return(testK8sVersionVar, &testResponse, nil)
		err := RunK8sVersionGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunK8sVersionGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.K8s.EXPECT().GetVersion().Return(testK8sVersionVar, nil, testK8sVersionErr)
		err := RunK8sVersionGet(cfg)
		assert.Error(t, err)
	})
}
