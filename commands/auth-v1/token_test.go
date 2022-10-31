package authv1

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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
	testTokenVar        = "test-token"
	testTokenBoolVar    = true
	testTokenContractNo = int32(1)
	testTokenNewVar     = "test-new-token"
	testTokenErr        = errors.New("token test error occurred")
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunTokenId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunTokenDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
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
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().List(testTokenContractNo).Return(testTokens, nil, nil)
		err := RunTokenList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		rm.AuthV1Mocks.Token.EXPECT().List(int32(0)).Return(testTokens, nil, testTokenErr)
		err := RunTokenList(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		res := resources.Token{Token: testToken}
		rm.AuthV1Mocks.Token.EXPECT().Get(testTokenVar, testTokenContractNo).Return(&res, nil, nil)
		err := RunTokenGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		res := resources.Token{Token: testToken}
		rm.AuthV1Mocks.Token.EXPECT().Get(testTokenVar, int32(0)).Return(&res, nil, testTokenErr)
		err := RunTokenGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().Create(testTokenContractNo).Return(&testJwt, nil, nil)
		err := RunTokenCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenCreateNilErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().Create(testTokenContractNo).Return(nil, nil, nil)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreateTokenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().Create(testTokenContractNo).Return(&resources.Jwt{}, nil, nil)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().Create(testTokenContractNo).Return(&testJwt, nil, testTokenErr)
		err := RunTokenCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteTokenId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar, testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteCriteriaCurrent(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(config.Token, testTokenVar)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("CURRENT", int32(0)).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteCriteriaExpired(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("EXPIRED", testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteCriteriaAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL", testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteById(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar, testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDeleteById(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteByIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar, testTokenContractNo).Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDeleteById(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL", testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDeleteAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL", testTokenContractNo).Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDeleteAll(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteAllResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("ALL", testTokenContractNo).Return(nil, nil, nil)
		err := RunTokenDeleteAll(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteExpired(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("EXPIRED", testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDeleteExpired(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteExpiredErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("EXPIRED", testTokenContractNo).Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDeleteExpired(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteCurrent(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		viper.Set(config.Token, testTokenVar)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("CURRENT", testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDeleteCurrent(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteCurrentNoErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		err := RunTokenDeleteCurrent(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteCurrentErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(config.Token, testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		rm.AuthV1Mocks.Token.EXPECT().DeleteByCriteria("CURRENT", testTokenContractNo).Return(&testDeleteResponse, nil, testTokenErr)
		err := RunTokenDeleteCurrent(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteByIdAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgContractNo), testTokenContractNo)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.AuthV1Mocks.Token.EXPECT().DeleteByID(testTokenVar, testTokenContractNo).Return(&testDeleteResponse, nil, nil)
		err := RunTokenDeleteById(cfg)
		assert.NoError(t, err)
	})
}

func TestRunTokenDeleteByIdAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgTokenId), testTokenVar)
		cfg.Stdin = os.Stdin
		err := RunTokenDeleteById(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteCurrentAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgCurrent), true)
		viper.Set(config.Token, testTokenVar)
		cfg.Stdin = os.Stdin
		err := RunTokenDeleteCurrent(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteExpiredAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgExpired), true)
		cfg.Stdin = os.Stdin
		err := RunTokenDeleteExpired(cfg)
		assert.Error(t, err)
	})
}

func TestRunTokenDeleteAllAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, authv1.ArgAll), true)
		cfg.Stdin = os.Stdin
		err := RunTokenDeleteAll(cfg)
		assert.Error(t, err)
	})
}

func TestGetTokensCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("token", constants.ArgCols), []string{"ExpirationDate"})
	getTokenCols(core.GetGlobalFlagName("token", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetTokensColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("token", constants.ArgCols), []string{"Unknown"})
	getTokenCols(core.GetGlobalFlagName("token", constants.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
