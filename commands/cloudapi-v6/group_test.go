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
	groupTest = resources.Group{
		Group: compute.Group{
			Properties: compute.GroupProperties{
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
		Group: compute.Group{
			Id: &testGroupVar,
			Properties: compute.GroupProperties{
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
		Groups: compute.Groups{
			Id: &testGroupVar,
			Items: &[]compute.Group{
				groupTestId.Group,
				groupTestId.Group,
			},
		},
	}
	groupTestNew = resources.Group{
		Group: compute.Group{
			Properties: compute.GroupProperties{
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
		Group: compute.Group{
			Id:         &testGroupVar,
			Properties: groupTestNew.Properties,
		},
	}
	groupTestGet = resources.Group{
		Group: compute.Group{
			Id:         &testGroupVar,
			Properties: groupTest.Properties,
			Type:       &testGroupType,
		},
	}
	groups = resources.Groups{
		Groups: compute.Groups{
			Id:    &testGroupVar,
			Items: &[]compute.Group{groupTest.Group},
		},
	}
	testGroupType       = compute.Type(testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("name=%s", testQueryParamVar))
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		err := PreRunGroupId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunGroupId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupUserIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testGroupVar)
		err := PreRunGroupUserIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupUserIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunGroupUserIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessCerts), testGroupBoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessCerts), testGroupBoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessCerts), testGroupBoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateFlowLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessMonitoring), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessCerts), testGroupBoolVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateNic), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateK8s), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateDc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreatePcc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateSnapshot), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateBackUpUnit), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgReserveIp), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3Privilege), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCreateFlowLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessMonitoring), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAccessCerts), testGroupBoolNewVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Groups{Groups: compute.Groups{Items: &[]compute.Group{}}}, &testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}
