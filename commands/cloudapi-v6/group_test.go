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
	groupTest = resources.Group{
		Group: ionoscloud.Group{
			Properties: &ionoscloud.GroupProperties{
				Name:                        &testGroupVar,
				CreateDataCenter:            &testGroupBoolVar,
				CreateSnapshot:              &testGroupBoolVar,
				ReserveIp:                   &testGroupBoolVar,
				AccessActivityLog:           &testGroupBoolVar,
				CreatePcc:                   &testGroupBoolVar,
				S3Privilege:                 &testGroupBoolVar,
				CreateBackupUnit:            &testGroupBoolVar,
				CreateInternetAccess:        &testGroupBoolVar,
				CreateK8sCluster:            &testGroupBoolVar,
				CreateFlowLog:               &testGroupBoolVar,
				AccessAndManageMonitoring:   &testGroupBoolVar,
				AccessAndManageCertificates: &testGroupBoolVar,
				AccessAndManageDns:          &testGroupBoolVar,
				ManageRegistry:              &testGroupBoolVar,
				ManageDBaaS:                 &testGroupBoolVar,
				ManageDataplatform:          &testGroupBoolVar,
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
	groupsList = resources.Groups{
		Groups: ionoscloud.Groups{
			Id: &testGroupVar,
			Items: &[]ionoscloud.Group{
				groupTestId.Group,
				groupTestId.Group,
			},
		},
	}
	groupTestNew = resources.Group{
		Group: ionoscloud.Group{
			Properties: &ionoscloud.GroupProperties{
				Name:                        &testGroupNewVar,
				CreateDataCenter:            &testGroupBoolNewVar,
				CreateSnapshot:              &testGroupBoolNewVar,
				ReserveIp:                   &testGroupBoolNewVar,
				AccessActivityLog:           &testGroupBoolNewVar,
				CreatePcc:                   &testGroupBoolNewVar,
				S3Privilege:                 &testGroupBoolNewVar,
				CreateBackupUnit:            &testGroupBoolNewVar,
				CreateInternetAccess:        &testGroupBoolNewVar,
				CreateK8sCluster:            &testGroupBoolNewVar,
				CreateFlowLog:               &testGroupBoolNewVar,
				AccessAndManageMonitoring:   &testGroupBoolNewVar,
				AccessAndManageCertificates: &testGroupBoolNewVar,
				AccessAndManageDns:          &testGroupBoolNewVar,
				ManageRegistry:              &testGroupBoolNewVar,
				ManageDBaaS:                 &testGroupBoolNewVar,
				ManageDataplatform:          &testGroupBoolNewVar,
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
	if ok := GroupCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("name=%s", testQueryParamVar))
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		err := PreRunGroupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupUserIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagUserId), testGroupVar)
		err := PreRunGroupUserIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupUserIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunGroupUserIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groups, &testResponse, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Groups{}, &testResponse, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groups, nil, testGroupErr)
		err := RunGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessCerts), testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessCerts), testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessCerts), testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessCerts), testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestGet, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateNic), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateDc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreatePcc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateSnapshot), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateBackUpUnit), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagReserveIp), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagCreateFlowLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessMonitoring), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAccessCerts), testGroupBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTestNew, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTestNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupNew, &testResponse, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groupsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groupsList, nil, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Groups{}, &testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Groups{Groups: ionoscloud.Groups{Items: &[]ionoscloud.Group{}}}, &testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groupsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testGroupErr)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveGroupAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagGroupId), testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}
