package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	"os"
	"regexp"
	"testing"

	cloudapiv5 "github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
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

func TestPreRunGroupResourceIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
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

func TestRunShareList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().ListShares(testShareVar).Return(shares, &testResponse, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().ListShares(testShareVar).Return(shares, nil, testShareErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, &testResponse, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, nil, testShareErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testResponse, nil)
		err := RunShareCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testResponseErr, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, testShareErr)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, &testResponse, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, testShareErr)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSharePrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.CloudApiV5Mocks.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgSharePrivilege), testShareBoolNewVar)
		rm.CloudApiV5Mocks.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, testShareErr)
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
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(&testResponse, nil)
		err := RunShareDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		rm.CloudApiV5Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, testShareErr)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV5Mocks.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
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
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgGroupId), testShareVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgResourceId), testShareVar)
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
	getGroupShareCols(core.GetGlobalFlagName("share", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetSharesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("share", config.ArgCols), []string{"Unknown"})
	getGroupShareCols(core.GetGlobalFlagName("share", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
