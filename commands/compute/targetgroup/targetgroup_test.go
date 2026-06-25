package targetgroup

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
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
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
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
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		cfg.SetFlag(constants.FlagOrderBy, testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
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
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgLicenceType, testTargetGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar).Return(&targetGroupTestGet, nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Get(testTargetGroupVar).Return(&targetGroupTestGet, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgName, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgAlgorithm, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgProtocol, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgCheckTimeout, testTargetGroupTimeout)
		cfg.SetFlag(cloudapiv6.ArgCheckInterval, testTargetGroupTimeout)
		cfg.SetFlag(cloudapiv6.ArgRetries, testTargetGroupRetries)
		cfg.SetFlag(cloudapiv6.ArgPath, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgMethod, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgMatchType, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgResponse, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgRegex, testTargetGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgNegate, testTargetGroupBoolVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest).Return(&targetGroupTestGet, nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgName, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgAlgorithm, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgProtocol, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgCheckTimeout, testTargetGroupTimeout)
		cfg.SetFlag(cloudapiv6.ArgCheckInterval, testTargetGroupTimeout)
		cfg.SetFlag(cloudapiv6.ArgRetries, testTargetGroupRetries)
		cfg.SetFlag(cloudapiv6.ArgPath, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgMethod, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgMatchType, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgResponse, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgRegex, testTargetGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgNegate, testTargetGroupBoolVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Create(targetGroupTest).Return(&targetGroupTestGet, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgName, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgAlgorithm, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgProtocol, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgCheckTimeout, testTargetGroupNewTimeout)
		cfg.SetFlag(cloudapiv6.ArgCheckInterval, testTargetGroupNewTimeout)
		cfg.SetFlag(cloudapiv6.ArgRetries, testTargetGroupNewRetries)
		cfg.SetFlag(cloudapiv6.ArgPath, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgMethod, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgMatchType, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgResponse, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgRegex, testTargetGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgNegate, testTargetGroupBoolNewVar)
		viper.Set(constants.ArgWait, false)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgName, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgAlgorithm, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgProtocol, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgCheckTimeout, testTargetGroupNewTimeout)
		cfg.SetFlag(cloudapiv6.ArgCheckInterval, testTargetGroupNewTimeout)
		cfg.SetFlag(cloudapiv6.ArgRetries, testTargetGroupNewRetries)
		cfg.SetFlag(cloudapiv6.ArgPath, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgMethod, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgMatchType, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgResponse, testTargetGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgRegex, testTargetGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgNegate, testTargetGroupBoolNewVar)
		viper.Set(constants.ArgWait, false)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, nil)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(nil, testTargetGroupTargetErr)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().List().Return(targetGroups, nil, testTargetGroupTargetErr)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.TargetGroup.EXPECT().Delete(testTargetGroupVar).Return(&testutil.TestResponse, testTargetGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		viper.Set(constants.ArgWait, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		cfg.SetFlag(cloudapiv6.ArgTargetGroupId, testTargetGroupVar)
		viper.Set(constants.ArgWait, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunTargetGroupDelete(cfg)
		assert.Error(t, err)
	})
}
