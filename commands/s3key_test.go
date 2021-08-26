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
	s3keyTest = v5.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &testS3keyBoolVar,
			},
		},
	}
	s3keyTestGet = v5.S3Key{
		S3Key: ionoscloud.S3Key{
			Id: &testS3keyVar,
			Properties: &ionoscloud.S3KeyProperties{
				SecretKey: &testS3keyVar,
				Active:    &testS3keyBoolVar,
			},
		},
	}
	s3keys = v5.S3Keys{
		S3Keys: ionoscloud.S3Keys{
			Id:    &testS3keyVar,
			Items: &[]ionoscloud.S3Key{s3keyTest.S3Key},
		},
	}
	s3keyProperties = ionoscloud.S3KeyProperties{
		Active: &testS3keyBoolNewVar,
	}
	s3keyNew = v5.S3Key{
		S3Key: ionoscloud.S3Key{
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
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		err := PreRunUserKeyIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunUserKeyIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunUserKeyIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.S3Key.EXPECT().List(testS3keyVar).Return(s3keys, &testResponse, nil)
		err := RunUserS3KeyList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().List(testS3keyVar).Return(s3keys, nil, testS3keyErr)
		err := RunUserS3KeyList(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar).Return(&s3keyTestGet, &testResponse, nil)
		err := RunUserS3KeyGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar).Return(&s3keyTestGet, nil, testS3keyErr)
		err := RunUserS3KeyGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, &testResponse, nil)
		err := RunUserS3KeyCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, &testResponseErr, nil)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, nil, nil)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, nil, testS3keyErr)
		err := RunUserS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, &testResponse, nil)
		err := RunUserS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, nil)
		err := RunUserS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, nil)
		err := RunUserS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, testS3keyErr)
		err := RunUserS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(&testResponse, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, testS3keyErr)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunUserS3KeyDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgS3KeyId), testS3keyVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, nil)
		err := RunUserS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunUserS3KeyDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgUserId), testS3keyVar)
		cfg.Stdin = os.Stdin
		err := RunUserS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetS3KeyCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("s3key", config.ArgCols), []string{"Active"})
	getS3KeyCols(core.GetGlobalFlagName("s3key", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetS3KeyColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("s3key", config.ArgCols), []string{"Unknown"})
	getS3KeyCols(core.GetGlobalFlagName("s3key", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetS3KeyIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosTokenEnvVar, "tok")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getS3KeyIds(w, testS3keyVar)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
