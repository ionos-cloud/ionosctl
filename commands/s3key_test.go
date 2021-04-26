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
	s3keyTest = resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &ionoscloud.S3KeyProperties{
				Active: &testS3keyBoolVar,
			},
		},
	}
	s3keyTestGet = resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Id: &testS3keyVar,
			Properties: &ionoscloud.S3KeyProperties{
				SecretKey: &testS3keyVar,
				Active:    &testS3keyBoolVar,
			},
		},
	}
	s3keys = resources.S3Keys{
		S3Keys: ionoscloud.S3Keys{
			Id:    &testS3keyVar,
			Items: &[]ionoscloud.S3Key{s3keyTest.S3Key},
		},
	}
	s3keyProperties = ionoscloud.S3KeyProperties{
		Active: &testS3keyBoolNewVar,
	}
	s3keyNew = resources.S3Key{
		S3Key: ionoscloud.S3Key{
			Properties: &s3keyProperties,
		},
	}
	testS3keyBoolVar    = false
	testS3keyBoolNewVar = true
	testS3keyVar        = "test-s3key"
	testS3keyErr        = errors.New("s3key test error")
)

func TestPreRunGlobalUserIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		err := PreRunGlobalUserIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalUserIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), "")
		err := PreRunGlobalUserIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalUserIdKeyIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		err := PreRunGlobalUserIdKeyIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalUserIdKeyIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), "")
		err := PreRunGlobalUserIdKeyIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunGlobalUserIdKeyIdActiveValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyActive), testS3keyBoolVar)
		err := PreRunGlobalUserIdKeyIdActiveValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunGlobalUserIdKeyIdActiveValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), "")
		err := PreRunGlobalUserIdKeyIdActiveValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.S3Key.EXPECT().List(testS3keyVar).Return(s3keys, nil, nil)
		err := RunS3KeyList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().List(testS3keyVar).Return(s3keys, nil, testS3keyErr)
		err := RunS3KeyList(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar).Return(&s3keyTestGet, nil, nil)
		err := RunS3KeyGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Get(testS3keyVar, testS3keyVar).Return(&s3keyTestGet, nil, testS3keyErr)
		err := RunS3KeyGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, nil, nil)
		err := RunS3KeyCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, &testResponse, nil)
		err := RunS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, nil, nil)
		err := RunS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		rm.S3Key.EXPECT().Create(testS3keyVar).Return(&s3keyTest, nil, testS3keyErr)
		err := RunS3KeyCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, nil)
		err := RunS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyUpdateOld(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, nil)
		err := RunS3KeyUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, nil)
		err := RunS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyActive), testS3keyBoolNewVar)
		rm.S3Key.EXPECT().Update(testS3keyVar, testS3keyVar, s3keyNew).Return(&s3keyNew, nil, testS3keyErr)
		err := RunS3KeyUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, nil)
		err := RunS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, testS3keyErr)
		err := RunS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, nil)
		err := RunS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunS3KeyDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetGlobalFlagName(cfg.ParentName, config.ArgUserId), testS3keyVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgS3KeyId), testS3keyVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.S3Key.EXPECT().Delete(testS3keyVar, testS3keyVar).Return(nil, nil)
		err := RunS3KeyDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunS3KeyDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgUserId), testS3keyVar)
		cfg.Stdin = os.Stdin
		err := RunS3KeyDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetS3KeyCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("s3key", config.ArgCols), []string{"Active"})
	getS3KeyCols(builder.GetGlobalFlagName("s3key", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetS3KeyColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("s3key", config.ArgCols), []string{"Unknown"})
	getS3KeyCols(builder.GetGlobalFlagName("s3key", config.ArgCols), w)
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
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getS3KeyIds(w, testS3keyVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
