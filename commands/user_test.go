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
	userTest = resources.User{
		User: ionoscloud.User{
			Properties: &ionoscloud.UserProperties{
				Firstname:     &testUserVar,
				Lastname:      &testUserVar,
				Email:         &testUserVar,
				Administrator: &testUserBoolVar,
				ForceSecAuth:  &testUserBoolVar,
				Password:      &testUserVar,
			},
		},
	}
	userTestGet = resources.User{
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
	users = resources.Users{
		Users: ionoscloud.Users{
			Id:    &testUserVar,
			Items: &[]ionoscloud.User{userTest.User},
		},
	}
	userProperties = resources.UserProperties{
		UserProperties: ionoscloud.UserProperties{
			Firstname:     &testUserNewVar,
			Lastname:      &testUserNewVar,
			Email:         &testUserNewVar,
			Administrator: &testUserBoolVar,
			ForceSecAuth:  &testUserBoolVar,
		},
	}
	userNew = resources.User{
		User: ionoscloud.User{
			Properties: &userProperties.UserProperties,
		},
	}
	testUserBoolVar = false
	testUserVar     = "test-user"
	testUserNewVar  = "test-new-user"
	testUserErr     = errors.New("user test error")
)

func TestPreRunUserIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		err := PreRunUserIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), "")
		err := PreRunUserIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunUserNameEmailPwdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserPassword), testUserVar)
		err := PreRunUserNameEmailPwdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserNameEmailPwdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserPassword), "")
		err := PreRunUserNameEmailPwdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
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
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, nil)
		err := RunUserGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTestGet, nil, testUserErr)
		err := RunUserGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTest, nil, nil)
		err := RunUserCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTest, &testResponse, nil)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserPassword), testUserVar)
		rm.User.EXPECT().Create(userTest).Return(&userTest, nil, testUserErr)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserForceSecAuth), testUserBoolVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserAdministrator), testUserBoolVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTest, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNew).Return(&userNew, nil, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userNew, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNew).Return(&userNew, nil, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserNewVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTest, nil, nil)
		rm.User.EXPECT().Update(testUserVar, userNew).Return(&userNew, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserFirstName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserLastName), testUserNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserEmail), testUserNewVar)
		rm.User.EXPECT().Get(testUserVar).Return(&userTest, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Delete(testUserVar).Return(nil, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		rm.User.EXPECT().Delete(testUserVar).Return(nil, testUserErr)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.User.EXPECT().Delete(testUserVar).Return(nil, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testUserVar)
		cfg.Stdin = os.Stdin
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetUsersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("user", config.ArgCols), []string{"Firstname"})
	getUserCols(builder.GetGlobalFlagName("user", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetUsersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("user", config.ArgCols), []string{"Unknown"})
	getUserCols(builder.GetGlobalFlagName("user", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getUsersIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetGroupUsersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getGroupUsersIds(w, testResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
