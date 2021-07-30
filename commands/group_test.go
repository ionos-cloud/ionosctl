package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	groupTest = v5.Group{
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
	groupTestNew = v5.Group{
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
	groupNew = v5.Group{
		Group: ionoscloud.Group{
			Id:         &testGroupVar,
			Properties: groupTestNew.Properties,
		},
	}
	groupTestGet = v5.Group{
		Group: ionoscloud.Group{
			Id:         &testGroupVar,
			Properties: groupTest.Properties,
			Type:       &testGroupType,
		},
	}
	groups = v5.Groups{
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

func TestPreRunGroupId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
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
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testGroupVar)
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

func TestPreRunGroupName(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupVar)
		err := PreRunGroupName(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupNameErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunGroupName(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Group.EXPECT().List().Return(groups, nil, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Group.EXPECT().List().Return(groups, nil, testGroupErr)
		err := RunGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTest, nil, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateNic), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateK8s), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateDc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreatePcc), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateSnapshot), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateBackUpUnit), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgReserveIp), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAccessLog), testGroupBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, nil)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testGroupNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateNic), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateK8s), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateDc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreatePcc), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateSnapshot), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCreateBackUpUnit), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgReserveIp), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAccessLog), testGroupBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3Privilege), testGroupBoolNewVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestNew, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTestNew).Return(&groupNew, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveGroupAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
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

func TestGetGroupsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getGroupsIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
