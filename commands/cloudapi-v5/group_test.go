package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	groupTest = resources.Group{
		Group: ionoscloud.Group{
			Properties: &ionoscloud.GroupProperties{
				Name:                 &testGroupVar,
				CreateDataCenter:     &testGroupBoolVar,
				CreateSnapshot:       &testGroupBoolVar,
				ReserveIp:            &testGroupBoolVar,
				AccessActivityLog:    &testGroupBoolVar,
				CreatePcc:            &testGroupBoolVar,
				S3Privilege:          &testGroupBoolVar,
				CreateBackupUnit:     &testGroupBoolVar,
				CreateInternetAccess: &testGroupBoolVar,
				CreateK8sCluster:     &testGroupBoolVar,
			},
		},
	}
	groupTestId = resources.Group{
		Group: ionoscloud.Group{
			Id: &testGroupVar,
			Properties: &ionoscloud.GroupProperties{
				Name:                 &testGroupVar,
				CreateDataCenter:     &testGroupBoolVar,
				CreateSnapshot:       &testGroupBoolVar,
				ReserveIp:            &testGroupBoolVar,
				AccessActivityLog:    &testGroupBoolVar,
				CreatePcc:            &testGroupBoolVar,
				S3Privilege:          &testGroupBoolVar,
				CreateBackupUnit:     &testGroupBoolVar,
				CreateInternetAccess: &testGroupBoolVar,
				CreateK8sCluster:     &testGroupBoolVar,
			},
		},
	}
	groupTestNew = resources.Group{
		Group: ionoscloud.Group{
			Properties: &ionoscloud.GroupProperties{
				Name:                 &testGroupNewVar,
				CreateDataCenter:     &testGroupBoolNewVar,
				CreateSnapshot:       &testGroupBoolNewVar,
				ReserveIp:            &testGroupBoolNewVar,
				AccessActivityLog:    &testGroupBoolNewVar,
				CreatePcc:            &testGroupBoolNewVar,
				S3Privilege:          &testGroupBoolNewVar,
				CreateBackupUnit:     &testGroupBoolNewVar,
				CreateInternetAccess: &testGroupBoolNewVar,
				CreateK8sCluster:     &testGroupBoolNewVar,
			},
		},
	}
	groupNew = resources.Group{
		Group: ionoscloud.Group{
			Id:         &testGroupVar,
			Properties: groupTestNew.Properties,
		},
	}
	groupTestGet = resources.Group{
		Group: ionoscloud.Group{
			Id:         &testGroupVar,
			Properties: groupTest.Properties,
			Type:       &testGroupType,
		},
	}
	groups = resources.Groups{
		Groups: ionoscloud.Groups{
			Id:    &testGroupVar,
			Items: &[]ionoscloud.Group{groupTest.Group},
		},
	}
	groupsList = resources.Groups{
		Groups: ionoscloud.Groups{
			Id: &testGroupVar,
			Items: &[]ionoscloud.Group{
				groupTestId.Group,
				groupTestId.Group,
			},
		},
	}
	testGroupType       = ionoscloud.Type(testGroupVar)
	testGroupBoolVar    = false
	testGroupBoolNewVar = true
	testGroupVar        = "test-resource"
	testGroupNewVar     = "test-new-resource"
	testGroupErr        = errors.New("resource test error")
)

func TestGroupCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(GroupCmd())
	if ok := FirewallRuleCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("name=%s", testQueryParamVar)})
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		err := PreRunGroupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupUserIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgUserId), testGroupVar)
		err := PreRunGroupUserIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupUserIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunGroupUserIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.CloudApiV5Mocks.Group.EXPECT().List(resources.ListQueryParams{}).Return(groups, &testResponse, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV5Mocks.Group.EXPECT().List(testListQueryParam).Return(resources.Groups{}, &testResponse, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.CloudApiV5Mocks.Group.EXPECT().List(resources.ListQueryParams{}).Return(groups, nil, testGroupErr)
		err := RunGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgS3Privilege), testGroupBoolVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Create(groupTest).Return(&groupTest, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgS3Privilege), testGroupBoolVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Create(groupTest).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgS3Privilege), testGroupBoolVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgS3Privilege), testGroupBoolVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Create(groupTest).Return(&groupTestGet, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateNic), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateK8s), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateDc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreatePcc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateSnapshot), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgCreateBackUpUnit), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgReserveIp), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAccessLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgS3Privilege), testGroupBoolNewVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestNew, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Update(testGroupVar, groupTestNew).Return(&groupNew, &testResponse, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAll), true)
		rm.CloudApiV5Mocks.Group.EXPECT().List(resources.ListQueryParams{}).Return(groupsList, &testResponse, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(nil, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV5Mocks.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveGroupAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testGroupVar)
		cfg.Stdin = os.Stdin
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetGroupsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("resource", config.ArgCols), []string{"Name"})
	getGroupCols(core.GetGlobalFlagName("resource", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetGroupsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("resource", config.ArgCols), []string{"Unknown"})
	getGroupCols(core.GetGlobalFlagName("resource", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
