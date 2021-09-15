package commands

import (
	"bufio"
	"bytes"
	"errors"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	targetGroupTestGet = resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Id: &testTargetGroupVar,
			Properties: &ionoscloud.TargetGroupProperties{
				Name:      &testTargetGroupVar,
				Algorithm: &testTargetGroupVar,
				Protocol:  &testTargetGroupVar,
				HealthCheck: &ionoscloud.TargetGroupHealthCheck{
					CheckTimeout:   &testTargetGroupTimeout,
					ConnectTimeout: &testTargetGroupTimeout,
					TargetTimeout:  &testTargetGroupTimeout,
					Retries:        &testTargetGroupRetries,
				},
				HttpHealthCheck: &ionoscloud.TargetGroupHttpHealthCheck{
					Path:      &testTargetGroupVar,
					Method:    &testTargetGroupVar,
					MatchType: &testTargetGroupVar,
					Response:  &testTargetGroupVar,
					Regex:     &testTargetGroupBoolVar,
					Negate:    &testTargetGroupBoolVar,
				},
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				State: &targetGroupState,
			},
		},
	}
	targetGroupTest = resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Properties: &ionoscloud.TargetGroupProperties{
				Name:      &testTargetGroupVar,
				Algorithm: &testTargetGroupVar,
				Protocol:  &testTargetGroupVar,
				HealthCheck: &ionoscloud.TargetGroupHealthCheck{
					CheckTimeout:   &testTargetGroupTimeout,
					ConnectTimeout: &testTargetGroupTimeout,
					TargetTimeout:  &testTargetGroupTimeout,
					Retries:        &testTargetGroupRetries,
				},
				HttpHealthCheck: &ionoscloud.TargetGroupHttpHealthCheck{
					Path:      &testTargetGroupVar,
					Method:    &testTargetGroupVar,
					MatchType: &testTargetGroupVar,
					Response:  &testTargetGroupVar,
					Regex:     &testTargetGroupBoolVar,
					Negate:    &testTargetGroupBoolVar,
				},
			},
		},
	}
	targetGroups = resources.TargetGroups{
		TargetGroups: ionoscloud.TargetGroups{
			Id:    &testTargetGroupVar,
			Items: &[]ionoscloud.TargetGroup{targetGroupTestGet.TargetGroup},
		},
	}
	targetGroupNewProperties = resources.TargetGroupProperties{
		TargetGroupProperties: ionoscloud.TargetGroupProperties{
			Name:      &testTargetGroupNewVar,
			Algorithm: &testTargetGroupNewVar,
			Protocol:  &testTargetGroupNewVar,
			HealthCheck: &ionoscloud.TargetGroupHealthCheck{
				CheckTimeout:   &testTargetGroupNewTimeout,
				ConnectTimeout: &testTargetGroupNewTimeout,
				TargetTimeout:  &testTargetGroupNewTimeout,
				Retries:        &testTargetGroupNewRetries,
			},
			HttpHealthCheck: &ionoscloud.TargetGroupHttpHealthCheck{
				Path:      &testTargetGroupNewVar,
				Method:    &testTargetGroupNewVar,
				MatchType: &testTargetGroupNewVar,
				Response:  &testTargetGroupNewVar,
				Regex:     &testTargetGroupBoolNewVar,
				Negate:    &testTargetGroupBoolNewVar,
			},
		},
	}
	targetGroupNew = resources.TargetGroup{
		TargetGroup: ionoscloud.TargetGroup{
			Id:         &testTargetGroupVar,
			Properties: &targetGroupNewProperties.TargetGroupProperties,
		},
	}
	targetGroupState          = "BUSY"
	testTargetGroupBoolVar    = false
	testTargetGroupBoolNewVar = true
	testTargetGroupRetries    = int32(3)
	testTargetGroupNewRetries = int32(5)
	testTargetGroupTimeout    = int32(5000)
	testTargetGroupNewTimeout = int32(5500)
	testTargetGroupVar        = "test-targetgroup"
	testTargetGroupNewVar     = "test-new-targetgroup"
	testTargetGroupErr        = errors.New("targetgroup test error")
)

func TestPreRunTargetGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunTargetGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, testTargetGroupErr)
		err := RunTargetGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar).Return(&targetGroupTestGet, nil, nil)
		err := RunTargetGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar).Return(&targetGroupTestGet, nil, testTargetGroupErr)
		err := RunTargetGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConnectionTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest).Return(&targetGroupTestGet, nil, nil)
		err := RunTargetGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConnectionTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest).Return(&targetGroupTestGet, &testResponse, nil)
		err := RunTargetGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConnectionTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupVar, &targetGroupNewProperties).Return(&targetGroupNew, nil, nil)
		err := RunTargetGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgConnectionTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupVar, &targetGroupNewProperties).Return(&targetGroupNew, nil, testTargetGroupErr)
		err := RunTargetGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, nil)
		err := RunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, testTargetGroupErr)
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, nil)
		err := RunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetTargetGroupsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("targetgroup", config.ArgCols), []string{"Name"})
	getTargetGroupCols(core.GetGlobalFlagName("targetgroup", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetTargetGroupsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("targetgroup", config.ArgCols), []string{"Unknown"})
	getTargetGroupCols(core.GetGlobalFlagName("targetgroup", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetTargetGroupsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getTargetGroupIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
