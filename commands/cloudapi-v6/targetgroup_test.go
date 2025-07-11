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
					CheckTimeout:  &testTargetGroupTimeout,
					CheckInterval: &testTargetGroupTimeout,
					Retries:       &testTargetGroupRetries,
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
					CheckTimeout:  &testTargetGroupTimeout,
					CheckInterval: &testTargetGroupTimeout,
					Retries:       &testTargetGroupRetries,
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
				CheckTimeout:  &testTargetGroupNewTimeout,
				CheckInterval: &testTargetGroupNewTimeout,
				Retries:       &testTargetGroupNewRetries,
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunTargetGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTargetGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTargetGroupDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		err := PreRunTargetGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTargetGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLicenceType), testTargetGroupVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckInterval), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckTimeout), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckInterval), testTargetGroupTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testTargetGroupRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPath), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMatchType), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResponse), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRegex), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNegate), testTargetGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckInterval), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAlgorithm), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagProtocol), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckTimeout), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCheckInterval), testTargetGroupNewTimeout)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRetries), testTargetGroupNewRetries)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPath), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMethod), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagMatchType), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagResponse), testTargetGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagRegex), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagNegate), testTargetGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
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
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetGroupId), testTargetGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}
