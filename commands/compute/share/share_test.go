package share

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
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
	testGroupVar     = "test-group"
	testResourceType = ionoscloud.Type("datacenter")
	groupTestId      = resources.Group{
		Group: ionoscloud.Group{
			Id: &testGroupVar,
			Properties: &ionoscloud.GroupProperties{
				Name: &testGroupVar,
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
	shareTest = resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
		},
	}
	shareTestGet = resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Id: &testShareVar,
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
			Type: &testResourceType,
		},
	}
	shares = resources.GroupShares{
		GroupShares: ionoscloud.GroupShares{
			Id:    &testShareVar,
			Items: &[]ionoscloud.GroupShare{shareTest.GroupShare},
		},
	}
	shareTestId = resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Id: &testShareVar,
			Properties: &ionoscloud.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
		},
	}

	sharesList = resources.GroupShares{
		GroupShares: ionoscloud.GroupShares{
			Id: &testShareVar,
			Items: &[]ionoscloud.GroupShare{
				shareTestId.GroupShare,
				shareTestId.GroupShare,
			},
		},
	}
	shareProperties = resources.GroupShareProperties{
		GroupShareProperties: ionoscloud.GroupShareProperties{
			EditPrivilege:  &testShareBoolNewVar,
			SharePrivilege: &testShareBoolNewVar,
		},
	}
	shareNew = resources.GroupShare{
		GroupShare: ionoscloud.GroupShare{
			Properties: &shareProperties.GroupShareProperties,
		},
	}
	testShareBoolVar    = false
	testShareBoolNewVar = true
	testShareVar        = "test-share"
	testShareErr        = errors.New("share test error")
)

func TestShareCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ShareCmd())
	if ok := ShareCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunGroupResourceIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		err := PreRunGroupResourceIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunShareIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunGroupResourceIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List().Return(groupsList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testGroupVar).Return(sharesList, &testutil.TestResponse, nil).Times(len(helpers.GetGroups(groupsList)))
		err := RunShareListAll(cfg)
		assert.NoError(t, err)
	})
}
func TestRunShareList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(shares, &testutil.TestResponse, nil)
		err := RunShareList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(shares, nil, testShareErr)
		err := RunShareList(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, &testutil.TestResponse, nil)
		err := RunShareGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, nil, testShareErr)
		err := RunShareGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testutil.TestResponse, nil)
		err := RunShareCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testutil.TestResponse, testShareErr)
		err := RunShareCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, testShareErr)
		err := RunShareCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunShareCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, &testutil.TestResponse, nil)
		err := RunShareUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareUpdateOldShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
		err := RunShareUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, testShareErr)
		err := RunShareUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunShareUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, testShareErr)
		err := RunShareUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, nil)
		err := RunShareDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(sharesList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, nil)
		err := RunShareDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(sharesList, nil, testShareErr)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(resources.GroupShares{}, &testutil.TestResponse, nil)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(
			resources.GroupShares{GroupShares: ionoscloud.GroupShares{Items: &[]ionoscloud.GroupShare{}}}, &testutil.TestResponse, nil)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar).Return(sharesList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, testShareErr)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, nil)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testutil.TestRequestIdVar).Return(&testutil.TestRequestStatus, nil, testutil.TestRequestErr)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, testShareErr)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
		err := RunShareDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}
