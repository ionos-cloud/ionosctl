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
	userTest = v5.UserPost{
		UserPost: ionoscloud.UserPost{
			Properties: &ionoscloud.UserPropertiesPost{
				Firstname:     &testUserVar,
				Lastname:      &testUserVar,
				Email:         &testUserVar,
				Administrator: &testUserBoolVar,
				ForceSecAuth:  &testUserBoolVar,
				Password:      &testUserVar,
			},
		},
	}
	userTestGet = v5.User{
		User: ionoscloud.User{
			Id: &testUserVar,
			Properties: &ionoscloud.UserProperties{
				Firstname:         &testUserVar,
				Lastname:          &testUserVar,
				Email:             &testUserVar,
				Administrator:     &testUserBoolVar,
				ForceSecAuth:      &testUserBoolVar,
				SecAuthActive:     &testUserBoolVar,
				S3CanonicalUserId: &testUserVar,
				Active:            &testUserBoolVar,
			},
		},
	}
	users = v5.Users{
		Users: ionoscloud.Users{
			Id:    &testUserVar,
			Items: &[]ionoscloud.User{userTestGet.User},
		},
	}
	userProperties = v5.UserProperties{
		UserProperties: ionoscloud.UserProperties{
			Firstname:     &testUserNewVar,
			Lastname:      &testUserNewVar,
			Email:         &testUserNewVar,
			Administrator: &testUserBoolVar,
			ForceSecAuth:  &testUserBoolVar,
		},
	}
	userNew = v5.User{
		User: ionoscloud.User{
			Properties: &userProperties.UserProperties,
		},
	}
	userNewPut = v5.UserPut{
		UserPut: ionoscloud.UserPut{
			Properties: &ionoscloud.UserPropertiesPut{
				Firstname:     &testUserNewVar,
				Lastname:      &testUserNewVar,
				Email:         &testUserNewVar,
				Administrator: &testUserBoolVar,
				ForceSecAuth:  &testUserBoolVar,
			},
		},
	}
	testUserBoolVar = false
	testUserVar     = "test-user"
	testUserNewVar  = "test-new-user"
	testUserErr     = errors.New("user test error")
)

func TestPreRunUserId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		err := PreRunUserId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunUserId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunUserNameEmailPwd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testUserVar)
		err := PreRunUserNameEmailPwd(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserNameEmailPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunUserNameEmailPwd(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.User.EXPECT().List().Return(users, nil, nil)
		err := RunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.User.EXPECT().List().Return(users, nil, testUserErr)
		err := RunUserList(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, nil)
		err := RunUserGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, testUserErr)
		err := RunUserGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTestGet, nil, nil)
		err := RunUserCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTestGet, &testResponseErr, nil)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTestGet, nil, testUserErr)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgForceSecAuth), testUserBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgAdmin), testUserBoolVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNewPut).Return(&userNew, nil, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userNew, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNewPut).Return(&userNew, nil, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserNewVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNewPut).Return(&userNew, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgEmail), testUserNewVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Delete(testUserVar).Return(nil, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Delete(testUserVar).Return(nil, testUserErr)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.User.EXPECT().Delete(testUserVar).Return(nil, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		cfg.Stdin = os.Stdin
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

// Group Users Test

var (
	groupUsersTest = v5.GroupMembers{
		GroupMembers: ionoscloud.GroupMembers{
			Items: &[]ionoscloud.User{userTestGet.User},
		},
	}
	groupUserTest = v5.User{
		User: ionoscloud.User{
			Id: &testUserVar,
		},
	}
)

func TestRunGroupUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().ListUsers(testGroupVar).Return(groupUsersTest, nil, nil)
		err := RunGroupUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		rm.Group.EXPECT().ListUsers(testGroupVar).Return(groupUsersTest, nil, testGroupErr)
		err := RunGroupUserList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, nil, nil)
		err := RunGroupUserAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, nil, testUserErr)
		err := RunGroupUserAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserAddResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().AddUser(testGroupVar, groupUserTest).Return(&userTestGet, &testResponseErr, nil)
		err := RunGroupUserAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, nil)
		err := RunGroupUserRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, testUserErr)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Group.EXPECT().RemoveUser(testGroupVar, testUserVar).Return(nil, nil)
		err := RunGroupUserRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgGroupId), testGroupVar)
		cfg.Stdin = os.Stdin
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestGetUsersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("user", config.ArgCols), []string{"Firstname"})
	getUserCols(core.GetGlobalFlagName("user", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetUsersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("user", config.ArgCols), []string{"Unknown"})
	getUserCols(core.GetGlobalFlagName("user", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetUsersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getUsersIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetGroupUsersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getGroupUsersIds(w, testResourceVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
