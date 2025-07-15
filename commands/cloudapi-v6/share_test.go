package commands

import (
	"bufio"
	"bytes"
	"errors"
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
	shareTest = resources.GroupShare{
		GroupShare: compute.GroupShare{
			Properties: &compute.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
		},
	}
	shareTestGet = resources.GroupShare{
		GroupShare: compute.GroupShare{
			Id: &testShareVar,
			Properties: &compute.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
			Type: &testResourceType,
		},
	}
	shares = resources.GroupShares{
		GroupShares: compute.GroupShares{
			Id:    &testShareVar,
			Items: &[]compute.GroupShare{shareTest.GroupShare},
		},
	}
	shareTestId = resources.GroupShare{
		GroupShare: compute.GroupShare{
			Id: &testShareVar,
			Properties: &compute.GroupShareProperties{
				EditPrivilege:  &testShareBoolVar,
				SharePrivilege: &testShareBoolVar,
			},
		},
	}

	sharesList = resources.GroupShares{
		GroupShares: compute.GroupShares{
			Id: &testShareVar,
			Items: &[]compute.GroupShare{
				shareTestId.GroupShare,
				shareTestId.GroupShare,
			},
		},
	}
	shareProperties = resources.GroupShareProperties{
		GroupShareProperties: compute.GroupShareProperties{
			EditPrivilege:  &testShareBoolNewVar,
			SharePrivilege: &testShareBoolNewVar,
		},
	}
	shareNew = resources.GroupShare{
		GroupShare: compute.GroupShare{
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(groupsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testGroupVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(sharesList, &testResponse, nil).Times(len(getGroups(groupsList)))
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(shares, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(shares, nil, testShareErr)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTestGet, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTestGet, nil, testShareErr)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, &testResponse, testShareErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, testShareErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareNew, &testResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareNew, nil, testShareErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&shareTest, nil, testShareErr)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(sharesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(sharesList, nil, testShareErr)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.GroupShares{}, &testResponse, nil)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.GroupShares{GroupShares: compute.GroupShares{Items: &[]compute.GroupShare{}}}, &testResponse, nil)
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
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(sharesList, &testResponse, nil)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testShareErr)
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
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
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testShareErr)
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
		rm.CloudApiV6Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
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
