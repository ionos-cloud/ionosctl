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
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	userTest = resources.UserPost{
		UserPost: compute.UserPost{
			Properties: &compute.UserPropertiesPost{
				Firstname:     &testUserVar,
				Lastname:      &testUserVar,
				Email:         &testUserVar,
				Administrator: &testUserBoolVar,
				ForceSecAuth:  &testUserBoolVar,
				Password:      &testUserVar,
			},
		},
	}
	usersList = resources.Users{
		Users: compute.Users{
			Id: &testUserVar,
			Items: &[]compute.User{
				userTestGet.User,
				userTestGet.User,
			},
		},
	}
	userTestGet = resources.User{
		User: compute.User{
			Id: &testUserVar,
			Properties: &compute.UserProperties{
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
	userGroupAdd = compute.UserGroupPost{
		Id: &testUserVar,
	}
	users = resources.Users{
		Users: compute.Users{
			Id:    &testUserVar,
			Items: &[]compute.User{userTestGet.User},
		},
	}
	groupUsersTestList = resources.GroupMembers{
		GroupMembers: compute.GroupMembers{
			Items: &[]compute.User{
				userTestGet.User,
				userTestGet.User,
			},
		},
	}
	userProperties = resources.UserProperties{
		UserProperties: compute.UserProperties{
			Firstname:     &testUserNewVar,
			Lastname:      &testUserNewVar,
			Email:         &testUserNewVar,
			Administrator: &testUserBoolVar,
			ForceSecAuth:  &testUserBoolVar,
		},
	}
	userNew = resources.User{
		User: compute.User{
			Properties: &userProperties.UserProperties,
		},
	}
	userNewPut = resources.UserPut{
		UserPut: compute.UserPut{
			Properties: &compute.UserPropertiesPut{
				Firstname:     &testUserNewVar,
				Lastname:      &testUserNewVar,
				Email:         &testUserNewVar,
				Password:      &testUserNewVar,
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

func TestUserCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(UserCmd())
	if ok := UserCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("firstname=%s", testQueryParamVar))
		err := PreRunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		err := PreRunUserId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunUserId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunUserNameEmailPwd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testUserVar)
		err := PreRunUserNameEmailPwd(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserNameEmailPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunUserNameEmailPwd(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(users, &testResponse, nil)
		err := RunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserListQueryParams(t *testing.T) {
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
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Users{}, &testResponse, nil)
		err := RunUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(users, nil, testUserErr)
		err := RunUserList(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, &testResponse, nil)
		err := RunUserGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, testUserErr)
		err := RunUserGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Create(userTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, &testResponse, nil)
		err := RunUserCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Create(userTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, &testResponse, testUserErr)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Create(userTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, testUserErr)
		err := RunUserCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgForceSecAuth), testUserBoolVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAdmin), testUserBoolVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Update(testUserVar, userNewPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userNew, &testResponse, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testUserNewVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userNew, nil, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Update(testUserVar, userNewPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userNew, nil, nil)
		err := RunUserUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserNewVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Update(testUserVar, userNewPut, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userNew, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirstName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLastName), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testUserNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testUserNewVar)
		rm.CloudApiV6Mocks.User.EXPECT().Get(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, testUserErr)
		err := RunUserUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(usersList, &testResponse, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(usersList, nil, testUserErr)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Users{}, &testResponse, nil)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Users{Users: compute.Users{Items: &[]compute.User{}}}, &testResponse, nil)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.User.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(usersList, &testResponse, nil)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testUserErr)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testUserErr)
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.User.EXPECT().Delete(testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunUserDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunUserDelete(cfg)
		assert.Error(t, err)
	})
}

// Group Users Test

var (
	groupUsersTest = resources.GroupMembers{
		GroupMembers: compute.GroupMembers{
			Items: &[]compute.User{userTestGet.User},
		},
	}
	groupUserTest = resources.User{
		User: compute.User{
			Id: &testUserVar,
		},
	}
)

func TestRunGroupUserList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(groupUsersTest, &testResponse, nil)
		err := RunGroupUserList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(groupUsersTest, nil, testGroupErr)
		err := RunGroupUserList(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserAdd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddUser(testGroupVar, userGroupAdd, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, &testResponse, nil)
		err := RunGroupUserAdd(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserAddErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddUser(testGroupVar, userGroupAdd, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, nil, testUserErr)
		err := RunGroupUserAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserAddResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddUser(testGroupVar, userGroupAdd, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&userTestGet, &testResponse, testUserErr)
		err := RunGroupUserAdd(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemove(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupUserRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserRemoveAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(groupUsersTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupUserRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserRemoveAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(groupUsersTestList, nil, testUserErr)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.GroupMembers{}, &testResponse, nil)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.GroupMembers{GroupMembers: compute.GroupMembers{Items: &[]compute.User{}}}, &testResponse, nil)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListUsers(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(groupUsersTestList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testUserErr)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testUserErr)
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUserRemoveAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveUser(testGroupVar, testUserVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunGroupUserRemove(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUserRemoveAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testUserVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testGroupVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunGroupUserRemove(cfg)
		assert.Error(t, err)
	})
}
