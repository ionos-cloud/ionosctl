package group

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

func TestPreRunGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		cfg.SetFlag(cloudapiv6.ArgUserId, testGroupVar)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groups, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		cfg.SetFlag(constants.FlagOrderBy, testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(resources.Groups{}, &testutil.TestResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groups, nil, testGroupErr)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testutil.TestResponse, nil)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgName, testGroupVar)
		cfg.SetFlag(cloudapiv6.ArgCreateNic, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateK8s, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateDc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreatePcc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateSnapshot, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateBackUpUnit, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgReserveIp, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgS3Privilege, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateFlowLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessMonitoring, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessCerts, testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest).Return(&groupTest, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgName, testGroupVar)
		cfg.SetFlag(cloudapiv6.ArgCreateNic, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateK8s, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateDc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreatePcc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateSnapshot, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateBackUpUnit, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgReserveIp, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgS3Privilege, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateFlowLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessMonitoring, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessCerts, testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest).Return(&groupTestGet, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgName, testGroupVar)
		cfg.SetFlag(cloudapiv6.ArgCreateNic, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateK8s, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateDc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreatePcc, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateSnapshot, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateBackUpUnit, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgReserveIp, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgS3Privilege, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgCreateFlowLog, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessMonitoring, testGroupBoolVar)
		cfg.SetFlag(cloudapiv6.ArgAccessCerts, testGroupBoolVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, testGroupErr)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		cfg.SetFlag(cloudapiv6.ArgName, testGroupNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateNic, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateK8s, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateDc, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreatePcc, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateSnapshot, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateBackUpUnit, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgReserveIp, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgAccessLog, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgS3Privilege, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgCreateFlowLog, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgAccessMonitoring, testGroupBoolNewVar)
		cfg.SetFlag(cloudapiv6.ArgAccessCerts, testGroupBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTestNew, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTestNew).Return(&groupNew, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groupsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groupsList, nil, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(resources.Groups{}, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(
			resources.Groups{Groups: ionoscloud.Groups{Items: &[]ionoscloud.Group{}}}, &testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groupsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, testGroupErr)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, nil)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(&testutil.TestResponse, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(nil, testGroupErr)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
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
		cfg.SetFlag(cloudapiv6.ArgGroupId, testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}
