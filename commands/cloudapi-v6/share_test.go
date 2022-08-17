package commands

import (
	"github.com/golang/mock/gomock"
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunGroupResourceIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgSharePrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Group.EXPECT().ListShares(testShareVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.GroupShares{GroupShares: ionoscloud.GroupShares{Items: &[]ionoscloud.GroupShare{}}}, &testResponse, nil)
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunShareDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
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
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgResourceId), testShareVar)
		cfg.Stdin = os.Stdin
		err := RunShareDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetSharesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("share", config.ArgCols), []string{"Type"})
	getGroupShareCols(core.GetGlobalFlagName("share", config.ArgCols), core.GetFlagName("share", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetSharesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("share", config.ArgCols), []string{"Unknown"})
	getGroupShareCols(core.GetGlobalFlagName("share", config.ArgCols), core.GetFlagName("share", cloudapiv6.ArgAll), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
