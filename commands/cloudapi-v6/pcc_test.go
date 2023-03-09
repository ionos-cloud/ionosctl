package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
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
	pccsList = resources.PrivateCrossConnects{
		PrivateCrossConnects: ionoscloud.PrivateCrossConnects{
			Id: &testPccVar,
			Items: &[]ionoscloud.PrivateCrossConnect{
				pccTestId.PrivateCrossConnect,
				pccTestId.PrivateCrossConnect,
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
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunPccList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunPccId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		err := PreRunPccId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunPccIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunPccId(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(pccs, &testResponse, nil)
		err := RunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.PrivateCrossConnects{}, &testResponse, nil)
		err := RunPccList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(pccs, nil, testPccErr)
		err := RunPccList(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTestGet, &testResponse, nil)
		err := RunPccGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTestGet, nil, testPccErr)
		err := RunPccGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccPeersList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, &testResponse, nil)
		err := RunPccPeersList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccPeersListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().GetPeers(testPccVar).Return(&[]resources.Peer{pccPeerTest}, nil, testPccErr)
		err := RunPccPeersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Create(pccTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, &testResponse, nil)
		err := RunPccCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Create(pccTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, &testResponse, testPccErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Create(pccTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, nil, testPccErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Create(pccTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccNewVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, nil, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccNew, &testResponse, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateOldUser(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccNew, nil, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccNew, nil, nil)
		err := RunPccUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccNewVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, nil, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccNew, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, nil, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Update(testPccVar, pccProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccUpdateGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDescription), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Get(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&pccTest, nil, testPccErr)
		err := RunPccUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(pccsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(pccsList, nil, testPccErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.PrivateCrossConnects{}, &testResponse, nil)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.PrivateCrossConnects{PrivateCrossConnects: ionoscloud.PrivateCrossConnects{Items: &[]ionoscloud.PrivateCrossConnect{}}}, &testResponse, nil)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgServerUrl, constants.DefaultApiURL)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(pccsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testPccErr)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testPccErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunPccDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.Pcc.EXPECT().Delete(testPccVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunPccDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunPccDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPccId), testPccVar)
		cfg.Stdin = os.Stdin
		err := RunPccDelete(cfg)
		assert.Error(t, err)
	})
}
