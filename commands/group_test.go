package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
	testGroupType       = ionoscloud.Type(testGroupVar)
	testGroupBoolVar    = false
	testGroupBoolNewVar = true
	testGroupVar        = "test-resource"
	testGroupNewVar     = "test-new-resource"
	testGroupErr        = errors.New("resource test error")
)

func TestPreRunGroupIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		err := PreRunGroupIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), "")
		err := PreRunGroupIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupUserIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testGroupVar)
		err := PreRunGroupUserIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupUserIdsValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), "")
		err := PreRunGroupUserIdsValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGroupNameValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupVar)
		err := PreRunGroupNameValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGroupNameValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), "")
		err := PreRunGroupNameValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Group.EXPECT().List().Return(groups, nil, nil)
		err := RunGroupList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Group.EXPECT().List().Return(groups, nil, testGroupErr)
		err := RunGroupList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateNic), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateK8s), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateDc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreatePcc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateSnapshot), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateBackUpUnit), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupReserveIp), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupAccessLog), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTest, nil, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateNic), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateK8s), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateDc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreatePcc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateSnapshot), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateBackUpUnit), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupReserveIp), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupAccessLog), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, &testResponse, nil)
		err := RunGroupCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateNic), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateK8s), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateDc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreatePcc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateSnapshot), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateBackUpUnit), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupReserveIp), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupAccessLog), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, testGroupErr)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateNic), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateK8s), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateDc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreatePcc), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateSnapshot), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateBackUpUnit), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupReserveIp), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupAccessLog), testGroupBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupS3Privilege), testGroupBoolVar)
		rm.Group.EXPECT().Create(groupTest).Return(&groupTestGet, nil, nil)
		err := RunGroupCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupName), testGroupNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateNic), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateK8s), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateDc), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreatePcc), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateSnapshot), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupCreateBackUpUnit), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupReserveIp), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupAccessLog), testGroupBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupS3Privilege), testGroupBoolNewVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTestNew, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTestNew).Return(&groupNew, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, nil)
		rm.Group.EXPECT().Update(testGroupVar, groupTest).Return(&groupTest, nil, nil)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Get(testGroupVar).Return(&groupTest, nil, testGroupErr)
		err := RunGroupUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupDeleteResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(&testResponse, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, testGroupErr)
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Group.EXPECT().Delete(testGroupVar).Return(nil, nil)
		err := RunGroupDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveGroupAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		cfg.Stdin = os.Stdin
		err := RunGroupDelete(cfg)
		assert.Error(t, err)
	})
}

// Group Users Test

var (
	groupUsersTest = resources.GroupMembers{
		GroupMembers: ionoscloud.GroupMembers{
			Items: &[]ionoscloud.User{userTest.User},
		},
	}
	groupUserTest = resources.User{
		User: ionoscloud.User{
			Id: &testUserVar,
		},
	}
)

func TestRunGroupListUsers(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().ListUsers(testGroupVar).Return(groupUsersTest, nil, nil)
		err := RunGroupListUsers(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListUsersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().ListUsers(testGroupVar).Return(groupUsersTest, nil, testGroupErr)
		err := RunGroupListUsers(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupAddUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, nil, nil)
		err := RunGroupAddUser(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupAddUserErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, nil, testUserErr)
		err := RunGroupAddUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupAddUserResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, &testResponse, nil)
		err := RunGroupAddUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupRemoveUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, nil)
		err := RunGroupRemoveUser(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveUserErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, testUserErr)
		err := RunGroupRemoveUser(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupRemoveUserAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, nil)
		err := RunGroupRemoveUser(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveUserAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testGroupVar)
		cfg.Stdin = os.Stdin
		err := RunGroupRemoveUser(cfg)
		assert.Error(t, err)
	})
}

func TestGetGroupsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("resource", config.ArgCols), []string{"Name"})
	getGroupCols(builder.GetGlobalFlagName("resource", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetGroupsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("resource", config.ArgCols), []string{"Unknown"})
	getGroupCols(builder.GetGlobalFlagName("resource", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getGroupsIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
