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
	targetGroupTestGet = resources.TargetGroup{
		TargetGroup: compute.TargetGroup{
			Id: &testTargetGroupVar,
			Properties: compute.TargetGroupProperties{
				Name:      &testTargetGroupVar,
				Algorithm: &testTargetGroupVar,
				Protocol:  &testTargetGroupVar,
				HealthCheck: &compute.TargetGroupHealthCheck{
					CheckTimeout:  &testTargetGroupTimeout,
					CheckInterval: &testTargetGroupTimeout,
					Retries:       &testTargetGroupRetries,
				},
				HttpHealthCheck: &compute.TargetGroupHttpHealthCheck{
					Path:      &testTargetGroupVar,
					Method:    &testTargetGroupVar,
					MatchType: &testTargetGroupVar,
					Response:  &testTargetGroupVar,
					Regex:     &testTargetGroupBoolVar,
					Negate:    &testTargetGroupBoolVar,
				},
			},
			Metadata: &compute.DatacenterElementMetadata{
				State: &targetGroupState,
			},
		},
	}
	targetGroupTest = resources.TargetGroup{
		TargetGroup: compute.TargetGroup{
			Properties: compute.TargetGroupProperties{
				Name:      &testTargetGroupVar,
				Algorithm: &testTargetGroupVar,
				Protocol:  &testTargetGroupVar,
				HealthCheck: &compute.TargetGroupHealthCheck{
					CheckTimeout:  &testTargetGroupTimeout,
					CheckInterval: &testTargetGroupTimeout,
					Retries:       &testTargetGroupRetries,
				},
				HttpHealthCheck: &compute.TargetGroupHttpHealthCheck{
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
		TargetGroups: compute.TargetGroups{
			Id:    &testTargetGroupVar,
			Items: []compute.TargetGroup{targetGroupTestGet.TargetGroup},
		},
	}
	targetGroupNewProperties = resources.TargetGroupProperties{
		TargetGroupProperties: compute.TargetGroupProperties{
			Name:      &testTargetGroupNewVar,
			Algorithm: &testTargetGroupNewVar,
			Protocol:  &testTargetGroupNewVar,
			HealthCheck: &compute.TargetGroupHealthCheck{
				CheckTimeout:  &testTargetGroupNewTimeout,
				CheckInterval: &testTargetGroupNewTimeout,
				Retries:       &testTargetGroupNewRetries,
			},
			HttpHealthCheck: &compute.TargetGroupHttpHealthCheck{
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
		TargetGroup: compute.TargetGroup{
			Id:         &testTargetGroupVar,
			Properties: targetGroupNewProperties.TargetGroupProperties,
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

func TestTargetGroupCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(TargetGroupCmd())
	if ok := TargetGroupCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunTargetGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunTargetGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTargetGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupListResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, &testResponse, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, testTargetGroupErr)
		err := RunTargetGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, nil)
		err := RunTargetGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&targetGroupTestGet, nil, nil)
		err := RunTargetGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&targetGroupTestGet, &testResponse, nil)
		err := RunTargetGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&targetGroupTestGet, nil, testTargetGroupErr)
		err := RunTargetGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&targetGroupTestGet, nil, nil)
		err := RunTargetGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupCreateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&targetGroupTestGet, &testResponse, nil)
		err := RunTargetGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupVar, &targetGroupNewProperties, testQueryParamOther).Return(&targetGroupNew, nil, nil)
		err := RunTargetGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCheckInterval), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Update(testTargetGroupVar, &targetGroupNewProperties, testQueryParamOther).Return(&targetGroupNew, nil, testTargetGroupErr)
		err := RunTargetGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testTargetGroupTargetErr)
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(targetGroups, nil, testTargetGroupTargetErr)
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testTargetGroupErr)
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTargetGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}
