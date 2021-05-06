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
	pccTest = resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Properties: &ionoscloud.PrivateCrossConnectProperties{
				Name:        &testPccVar,
				Description: &testPccVar,
			},
		},
	}
	pccTestGet = resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Id:         &testPccVar,
			Properties: pccTest.Properties,
		},
	}
	pccPeerTest = resources.Peer{
		Peer: ionoscloud.Peer{
			Id:             &testPccVar,
			Name:           &testPccVar,
			DatacenterId:   &testPccVar,
			DatacenterName: &testPccVar,
			Location:       &testPccVar,
		},
	}
	pccs = resources.PrivateCrossConnects{
		PrivateCrossConnects: ionoscloud.PrivateCrossConnects{
			Id:    &testPccVar,
			Items: &[]ionoscloud.PrivateCrossConnect{pccTest.PrivateCrossConnect},
		},
	}
	pccProperties = resources.PrivateCrossConnectProperties{
		PrivateCrossConnectProperties: ionoscloud.PrivateCrossConnectProperties{
			Name:        &testPccNewVar,
			Description: &testPccNewVar,
		},
	}
	pccNew = resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Properties: &pccProperties.PrivateCrossConnectProperties,
		},
	}
	testPccVar    = "test-pcc"
	testPccNewVar = "test-new-pcc"
	testPccErr    = errors.New("pcc test error")
)

func TestPreRunPccId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		err := PreRunPccId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), "")
		err := PreRunPccId(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Pcc.EXPECT().List().Return(pccs, nil, nil)
		err := RunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.Pcc.EXPECT().List().Return(pccs, nil, testPccErr)
		err := RunPccList(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTestGet, nil, nil)
		err := RunPccGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTestGet, nil, testPccErr)
		err := RunPccGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccGetPeers(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, nil, nil)
		err := RunPccGetPeers(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccGetPeersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, nil, testPccErr)
		err := RunPccGetPeers(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccVar)
		rm.Pcc.EXPECT().Create(pccTest).Return(&pccTest, nil, nil)
		err := RunPccCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccVar)
		rm.Pcc.EXPECT().Create(pccTest).Return(&pccTest, &testResponse, nil)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccVar)
		rm.Pcc.EXPECT().Create(pccTest).Return(&pccTest, nil, testPccErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
		rm.Pcc.EXPECT().Create(pccTest).Return(&pccTest, nil, nil)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccNewVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccNew, nil, nil)
		rm.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccNewVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccNewVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, nil)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccName), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccDescription), testPccVar)
		rm.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().Delete(testPccVar).Return(nil, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		rm.Pcc.EXPECT().Delete(testPccVar).Return(nil, testPccErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWaitForRequest), true)
		rm.Pcc.EXPECT().Delete(testPccVar).Return(nil, nil)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Pcc.EXPECT().Delete(testPccVar).Return(nil, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgPccId), testPccVar)
		cfg.Stdin = os.Stdin
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetPccsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("pcc", config.ArgCols), []string{"Name"})
	getPccCols(builder.GetGlobalFlagName("pcc", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetPccsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("pcc", config.ArgCols), []string{"Unknown"})
	getPccCols(builder.GetGlobalFlagName("pcc", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetPccsIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getPccsIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
