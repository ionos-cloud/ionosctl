package lan

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/helpers"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute/testutil"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	publicLan    = true
	publicNewLan = false
	lanPostTest  = ionoscloud.Lan{
		Properties: &ionoscloud.LanProperties{
			Name:       &testLanVar,
			IpFailover: nil,
			Pcc:        &testLanVar,
			Public:     &publicLan,
		},
	}
	lp = ionoscloud.Lan{
		Id:         &testLanVar,
		Properties: lanPostTest.Properties,
		Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testutil.TestStateVar},
	}
	l = ionoscloud.Lan{
		Id: &testLanVar,
		Properties: &ionoscloud.LanProperties{
			Name: &testLanVar,
			Pcc:  &testLanVar,
		},
	}
	lanProperties = resources.LanProperties{
		LanProperties: ionoscloud.LanProperties{
			Name:   &testLanNewVar,
			Pcc:    &testLanNewVar,
			Public: &publicNewLan,
		},
	}
	lanNew = resources.Lan{
		Lan: ionoscloud.Lan{
			Id: &testLanVar,
			Properties: &ionoscloud.LanProperties{
				Name:       lanProperties.LanProperties.Name,
				Public:     lanProperties.LanProperties.Public,
				IpFailover: nil,
				Pcc:        &testLanNewVar,
			},
		},
	}
	ls = resources.Lans{
		Lans: ionoscloud.Lans{
			Id:    &testLanVar,
			Items: &[]ionoscloud.Lan{l},
		},
	}
	lansList = resources.Lans{
		Lans: ionoscloud.Lans{
			Id: &testLanVar,
			Items: &[]ionoscloud.Lan{
				l,
				l,
			},
		},
	}
	testLanVar    = "test-lan"
	testLanNewVar = "test-new-lan"
	testLanErr    = errors.New("lan test: error occurred")
)

func TestLanCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LanCmd())
	if ok := LanCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunLansList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		err := PreRunLansList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLansListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("createdBy=%s", testutil.TestQueryParamVar))
		err := PreRunLansList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLansListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		err := PreRunLansList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List().Return(testutil.TestDcs, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testutil.TestDatacenterVar).Return(lansList, &testutil.TestResponse, nil).Times(len(helpers.GetDataCenters(testutil.TestDcs)))
		err := RunLanListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(ls, &testutil.TestResponse, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, fmt.Sprintf("%s=%s", testutil.TestQueryParamVar, testutil.TestQueryParamVar))
		cfg.SetFlag(constants.FlagOrderBy, testutil.TestQueryParamVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(resources.Lans{}, &testutil.TestResponse, nil)
		err := RunLanList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(ls, nil, testLanErr)
		err := RunLanList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, &testutil.TestResponse, nil)
		err := RunLanGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Get(testLanVar, testLanVar).Return(&resources.Lan{Lan: l}, nil, testLanErr)
		err := RunLanGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgName, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgPccId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgPublic, publicLan)
		rm.CloudApiV6Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{Lan: lanPostTest}).Return(&resources.LanPost{Lan: lp}, &testutil.TestResponse, nil)
		err := RunLanCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgName, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgPublic, publicLan)
		cfg.SetFlag(cloudapiv6.ArgPccId, testLanVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Create(testLanVar, resources.LanPost{Lan: lanPostTest}).Return(&resources.LanPost{Lan: lp}, nil, testLanErr)
		err := RunLanCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgName, testLanNewVar)
		cfg.SetFlag(cloudapiv6.ArgPublic, publicNewLan)
		cfg.SetFlag(cloudapiv6.ArgPccId, testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testutil.TestResponse, nil)
		err := RunLanUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgName, testLanNewVar)
		cfg.SetFlag(cloudapiv6.ArgPublic, publicNewLan)
		cfg.SetFlag(cloudapiv6.ArgPccId, testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, nil, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgServerId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgName, testLanNewVar)
		cfg.SetFlag(cloudapiv6.ArgPublic, publicNewLan)
		cfg.SetFlag(cloudapiv6.ArgPccId, testLanNewVar)
		rm.CloudApiV6Mocks.Lan.EXPECT().Update(testLanVar, testLanVar, lanProperties).Return(&lanNew, &testutil.TestResponse, testLanErr)
		err := RunLanUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testutil.TestResponse, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(lansList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testutil.TestResponse, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(lansList, nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(resources.Lans{}, &testutil.TestResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(
			resources.Lans{Lans: ionoscloud.Lans{Items: &[]ionoscloud.Lan{}}}, &testutil.TestResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.SetFlag(cloudapiv6.ArgAll, true)
		rm.CloudApiV6Mocks.Lan.EXPECT().List(testLanVar).Return(lansList, &testutil.TestResponse, nil)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testutil.TestResponse, testLanErr)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(&testutil.TestResponse, nil)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, testLanErr)
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLanDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		viper.Set(constants.ArgWait, false)
		rm.CloudApiV6Mocks.Lan.EXPECT().Delete(testLanVar, testLanVar).Return(nil, nil)
		err := RunLanDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLanDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.SetFlag(cloudapiv6.ArgDataCenterId, testLanVar)
		cfg.SetFlag(cloudapiv6.ArgLanId, testLanVar)
		viper.Set(constants.ArgWait, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunLanDelete(cfg)
		assert.Error(t, err)
	})
}
