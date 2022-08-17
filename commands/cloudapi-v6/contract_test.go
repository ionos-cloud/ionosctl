package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testContractResourceLimits = ionoscloud.ResourceLimits{
		CoresPerServer:         &testContractInt32,
		CoresPerContract:       &testContractInt32,
		CoresProvisioned:       &testContractInt32,
		RamPerServer:           &testContractInt32,
		RamPerContract:         &testContractInt32,
		RamProvisioned:         &testContractInt32,
		HddLimitPerVolume:      &testContractInt64,
		HddLimitPerContract:    &testContractInt64,
		HddVolumeProvisioned:   &testContractInt64,
		SsdLimitPerVolume:      &testContractInt64,
		SsdLimitPerContract:    &testContractInt64,
		SsdVolumeProvisioned:   &testContractInt64,
		DasVolumeProvisioned:   &testContractInt64,
		ReservableIps:          &testContractInt32,
		ReservedIpsOnContract:  &testContractInt32,
		ReservedIpsInUse:       &testContractInt32,
		K8sClusterLimitTotal:   &testContractInt32,
		K8sClustersProvisioned: &testContractInt32,
		NatGatewayProvisioned:  &testContractInt32,
		NatGatewayLimitTotal:   &testContractInt32,
		NlbProvisioned:         &testContractInt32,
		NlbLimitTotal:          &testContractInt32,
	}
	testContract = resources.Contract{
		Contract: ionoscloud.Contract{
			Properties: &ionoscloud.ContractProperties{
				ContractNumber: &testContractInt64,
				Owner:          &testContractVar,
				Status:         &testContractVar,
				RegDomain:      &testContractVar,
				ResourceLimits: &testContractResourceLimits,
			},
		},
	}
	testContracts = resources.Contracts{
		Contracts: ionoscloud.Contracts{
			Items: &[]ionoscloud.Contract{testContract.Contract},
		},
	}
	testContractInt64 = int64(2)
	testContractInt32 = int32(1)
	testContractVar   = "test-contract"
	testContractErr   = errors.New("contract test error occurred")
)

func TestContractCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ContractCmd())
	if ok := ContractCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestRunContractGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		rm.CloudApiV6Mocks.Contract.EXPECT().Get(gomock.AssignableToTypeOf(testQueryParamOther)).Return(testContracts, &testResponse, nil)
		err := RunContractGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunContractGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.CloudApiV6Mocks.Contract.EXPECT().Get(gomock.AssignableToTypeOf(testQueryParamOther)).Return(testContracts, nil, testContractErr)
		err := RunContractGet(cfg)
		assert.Error(t, err)
	})
}

func TestGetContractsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("contract", config.ArgCols), []string{"ContractNumber"})
	getContractCols(core.GetGlobalFlagName("contract", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetContractsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("contract", config.ArgCols), []string{"Unknown"})
	getContractCols(core.GetGlobalFlagName("contract", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
