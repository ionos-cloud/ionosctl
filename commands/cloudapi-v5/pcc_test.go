package cloudapi_v5

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
	"github.com/ionos-cloud/ionosctl/internal/utils/clierror"
	cloudapiv5 "github.com/ionos-cloud/ionosctl/services/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v5/resources"
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
	pccTestId = resources.PrivateCrossConnect{
		PrivateCrossConnect: ionoscloud.PrivateCrossConnect{
			Id: &testPccVar,
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
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
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
	pccsList = resources.PrivateCrossConnects{
		PrivateCrossConnects: ionoscloud.PrivateCrossConnects{
			Id: &testPccVar,
			Items: &[]ionoscloud.PrivateCrossConnect{
				pccTestId.PrivateCrossConnect,
				pccTestId.PrivateCrossConnect,
			},
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

func TestPccCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(PccCmd())
	if ok := PccCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunPccList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunPccList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunPccId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		err := PreRunPccId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunPccId(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.Pcc.EXPECT().List(resources.ListQueryParams{}).Return(pccs, &testResponse, nil)
		err := RunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().List(testListQueryParam).Return(resources.PrivateCrossConnects{}, &testResponse, nil)
		err := RunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		rm.CloudApiV5Mocks.Pcc.EXPECT().List(resources.ListQueryParams{}).Return(pccs, nil, testPccErr)
		err := RunPccList(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTestGet, &testResponse, nil)
		err := RunPccGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTestGet, nil, testPccErr)
		err := RunPccGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccPeersList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, nil, nil)
		err := RunPccPeersList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccPeersListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, nil, testPccErr)
		err := RunPccPeersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Create(pccTest).Return(&pccTest, &testResponse, nil)
		err := RunPccCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Create(pccTest).Return(&pccTest, &testResponseErr, nil)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Create(pccTest).Return(&pccTest, nil, testPccErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Create(pccTest).Return(&pccTest, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccNewVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, &testResponse, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccNew, nil, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccNewVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties).Return(&pccNew, &testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgDescription), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Get(testPccVar).Return(&pccTest, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(&testResponse, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgAll), true)
		rm.CloudApiV5Mocks.Pcc.EXPECT().List(resources.ListQueryParams{}).Return(pccsList, &testResponse, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(&testResponse, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(nil, testPccErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(&testResponse, nil)
		rm.CloudApiV5Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV5Mocks.Pcc.EXPECT().Delete(testPccVar).Return(nil, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv5.ArgPccId), testPccVar)
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
	viper.Set(core.GetGlobalFlagName("pcc", config.ArgCols), []string{"Name"})
	getPccCols(core.GetGlobalFlagName("pcc", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetPccsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("pcc", config.ArgCols), []string{"Unknown"})
	getPccCols(core.GetGlobalFlagName("pcc", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
