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
	compute "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	s3keyTest = resources.S3Key{
		S3Key: compute.S3Key{
			Properties: &compute.S3KeyProperties{
				Active: &testS3keyBoolVar,
			},
		},
	}
	s3keyTestId = resources.S3Key{
		S3Key: compute.S3Key{
			Id: &testS3keyVar,
			Properties: &compute.S3KeyProperties{
				Active: &testS3keyBoolVar,
			},
		},
	}
	s3keysList = resources.S3Keys{
		S3Keys: compute.S3Keys{
			Id: &testS3keyVar,
			Items: &[]compute.S3Key{
				s3keyTestId.S3Key,
				s3keyTestId.S3Key,
			},
		},
	}
	s3keyTestGet = resources.S3Key{
		S3Key: compute.S3Key{
			Id: &testS3keyVar,
			Properties: &compute.S3KeyProperties{
				SecretKey: &testS3keyVar,
				Active:    &testS3keyBoolVar,
			},
		},
	}
	s3keys = resources.S3Keys{
		S3Keys: compute.S3Keys{
			Id:    &testS3keyVar,
			Items: &[]compute.S3Key{s3keyTest.S3Key},
		},
	}
	s3keyProperties = compute.S3KeyProperties{
		Active: &testS3keyBoolNewVar,
	}
	s3keyNew = resources.S3Key{
		S3Key: compute.S3Key{
			Properties: &s3keyProperties,
		},
	}
	testS3keyBoolVar    = false
	testS3keyBoolNewVar = true
	testS3keyVar        = "test-s3key"
	testS3keyErr        = errors.New("s3key test error")
)

func TestPreRunUserKeyIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		err := PreRunUserKeyIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserKeyIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunUserKeyIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(s3keys, &testResponse, nil)
		err := RunUserS3KeyList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(s3keys, nil, testS3keyErr)
		err := RunUserS3KeyList(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTestGet, &testResponse, nil)
		err := RunUserS3KeyGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTestGet, nil, testS3keyErr)
		err := RunUserS3KeyGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Create(testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTest, &testResponse, nil)
		err := RunUserS3KeyCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Create(testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTest, &testResponse, testS3keyErr)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Create(testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Create(testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyTest, nil, testS3keyErr)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyNew, &testResponse, nil)
		err := RunUserS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyNew, nil, nil)
		err := RunUserS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunUserS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&s3keyNew, nil, testS3keyErr)
		err := RunUserS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(s3keysList, &testResponse, nil)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(s3keysList, nil, testS3keyErr)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.S3Keys{}, &testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.S3Keys{S3Keys: compute.S3Keys{Items: &[]compute.S3Key{}}}, &testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().List(testS3keyVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(s3keysList, &testResponse, nil)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testS3keyErr)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testS3keyErr)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgS3KeyId), testS3keyVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgUserId), testS3keyVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}
