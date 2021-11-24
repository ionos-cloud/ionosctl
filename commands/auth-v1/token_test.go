package authv1

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	authv1 "github.com/ionos-cloud/ionosctl/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/services/auth-v1/resources"
	sdkgoauth "github.com/ionos-cloud/sdk-go-auth"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testToken = sdkgoauth.Token{
		Id:             &testTokenVar,
		Href:           &testTokenVar,
		CreatedDate:    &testTokenVar,
		ExpirationDate: &testTokenVar,
	}
	testJwt = resources.Jwt{
		Jwt: sdkgoauth.Jwt{
			Token: &testTokenNewVar,
		},
	}
	testTokens = resources.Tokens{
		Tokens: sdkgoauth.Tokens{
			Tokens: &[]sdkgoauth.Token{testToken},
		},
	}
	testDeleteResponse = resources.DeleteResponse{
		DeleteResponse: sdkgoauth.DeleteResponse{
			Success: &testTokenBoolVar,
		},
	}
	testTokenVar     = "test-token"
	testTokenBoolVar = true
	testTokenNewVar  = "test-new-token"
	testTokenErr     = errors.New("token test error occurred")
)

func TestTokenCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(TokenCmd())
	if ok := TokenCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunTokenId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		err := PreRunTokenId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTokenIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunTokenId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTokenDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		err := PreRunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTokenDeleteExpired(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		err := PreRunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTokenDeleteCurrent(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		err := PreRunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunTokenDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		err := PreRunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.AuthV1Mocks.Token.EXPECT().List().Return(testTokens, nil, nil)
		err := RunTokenList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.AuthV1Mocks.Token.EXPECT().List().Return(testTokens, nil, testTokenErr)
		err := RunTokenList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		res := resources.Token{Token: testToken}
		rm.AuthV1Mocks.Token.EXPECT().Get(testTokenVar).Return(&res, nil, nil)
		err := RunTokenGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		res := resources.Token{Token: testToken}
		rm.AuthV1Mocks.Token.EXPECT().Get(testTokenVar).Return(&res, nil, testTokenErr)
		err := RunTokenGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.AuthV1Mocks.Token.EXPECT().Create().Return(&testJwt, nil, nil)
		err := RunTokenCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenCreateNilErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.AuthV1Mocks.Token.EXPECT().Create().Return(nil, nil, nil)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreateTokenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		rm.AuthV1Mocks.Token.EXPECT().Create().Return(&resources.Jwt{}, nil, nil)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.AuthV1Mocks.Token.EXPECT().Create().Return(&testJwt, nil, testTokenErr)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteByIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar).Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL").Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL").Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteExpired(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("EXPIRED").Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteExpiredErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("EXPIRED").Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteCurrent(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("CURRENT").Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteCurrentErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("CURRENT").Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteByIdAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		cfg.Stdin = os.Stdin
		err := RunTokenDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetTokensCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("token", config.ArgCols), []string{"ExpirationDate"})
	getTokenCols(core.GetGlobalFlagName("token", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetTokensColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("token", config.ArgCols), []string{"Unknown"})
	getTokenCols(core.GetGlobalFlagName("token", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
