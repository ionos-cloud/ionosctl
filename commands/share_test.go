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

func TestPreRunGroupResourceIdsValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		err := PreRunGroupResourceIdsValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunShareIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), "")
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), "")
		err := PreRunGroupResourceIdsValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupListShares(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		rm.Group.EXPECT().ListShares(testShareVar).Return(shares, nil, nil)
		err := RunGroupListShares(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupListSharesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		rm.Group.EXPECT().ListShares(testShareVar).Return(shares, nil, testShareErr)
		err := RunGroupListShares(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupGetShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, nil, nil)
		err := RunGroupGetShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunShareGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTestGet, nil, testShareErr)
		err := RunGroupGetShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupAddShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
		err := RunGroupAddShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupAddShareResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, &testResponse, nil)
		err := RunGroupAddShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupAddShareErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, testShareErr)
		err := RunGroupAddShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupAddShareWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Group.EXPECT().AddShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
		err := RunGroupAddShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSharePrivilege), testShareBoolNewVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, nil)
		err := RunGroupUpdateShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateShareOldShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareTest).Return(&shareTest, nil, nil)
		err := RunGroupUpdateShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupUpdateShareErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSharePrivilege), testShareBoolNewVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, testShareErr)
		err := RunGroupUpdateShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateShareWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSharePrivilege), testShareBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, nil)
		rm.Group.EXPECT().UpdateShare(testShareVar, testShareVar, shareNew).Return(&shareNew, nil, nil)
		err := RunGroupUpdateShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupUpdateShareGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgEditPrivilege), testShareBoolNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgSharePrivilege), testShareBoolNewVar)
		rm.Group.EXPECT().GetShare(testShareVar, testShareVar).Return(&shareTest, nil, testShareErr)
		err := RunGroupUpdateShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupRemoveShare(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
		err := RunGroupRemoveShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveShareWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), true)
		rm.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
		err := RunGroupRemoveShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupRemoveShareErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		rm.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, testShareErr)
		err := RunGroupRemoveShare(cfg)
		assert.Error(t, err)
	})
}

func TestRunGroupRemoveShareAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Group.EXPECT().RemoveShare(testShareVar, testShareVar).Return(nil, nil)
		err := RunGroupRemoveShare(cfg)
		assert.NoError(t, err)
	})
}

func TestRunGroupRemoveShareAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgGroupId), testShareVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgResourceId), testShareVar)
		cfg.Stdin = os.Stdin
		err := RunGroupRemoveShare(cfg)
		assert.Error(t, err)
	})
}

func TestGetSharesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("share", config.ArgCols), []string{"Type"})
	getGroupShareCols(builder.GetGlobalFlagName("share", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetSharesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("share", config.ArgCols), []string{"Unknown"})
	getGroupShareCols(builder.GetGlobalFlagName("share", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetGroupResourcesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getGroupResourcesIds(w, testResourceVar)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
