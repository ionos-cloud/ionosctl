package commands

import (
	"github.com/golang/mock/gomock"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	lanNicId    = int32(1)
	lanNewNicId = int32(2)
	dhcpNic     = false
	dhcpNewNic  = true
	ipsNic      = []string{"x.x.x.x"}
	n           = ionoscloud.Nic{
		Id: &testNicVar,
		Properties: &ionoscloud.NicProperties{
			Name:           &testNicVar,
			Lan:            &lanNicId,
			Dhcp:           &dhcpNic,
			Ips:            &ipsNic,
			FirewallActive: &dhcpNic,
			FirewallType:   &testNicVar,
			DeviceNumber:   &lanNicId,
			PciSlot:        &lanNicId,
			Mac:            &testNicVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
	}
	nicsList = resources.Nics{
		Nics: ionoscloud.Nics{
			Id: &testNicVar,
			Items: &[]ionoscloud.Nic{
				n,
				n,
			},
		},
	}
	nicLoadBalancer = ionoscloud.Nic{
		Id: &testLoadbalancerVar,
		Properties: &ionoscloud.NicProperties{
			Name:           &testNicVar,
			Lan:            &lanNicId,
			Dhcp:           &dhcpNic,
			Ips:            &ipsNic,
			FirewallActive: &dhcpNic,
			Mac:            &testNicVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
	}
	balancedNicsList = resources.BalancedNics{
		BalancedNics: ionoscloud.BalancedNics{
			Id:    &testNicVar,
			Items: &[]ionoscloud.Nic{nicLoadBalancer, nicLoadBalancer},
		},
	}
	nicProperties = resources.NicProperties{
		NicProperties: ionoscloud.NicProperties{
			Name: &testNicNewVar,
			Dhcp: &dhcpNewNic,
			Lan:  &lanNewNicId,
		},
	}
	nicNew = resources.Nic{
		Nic: ionoscloud.Nic{
			Id: &testNicVar,
			Properties: &ionoscloud.NicProperties{
				Name:           nicProperties.NicProperties.Name,
				Lan:            nicProperties.NicProperties.Lan,
				Dhcp:           nicProperties.NicProperties.Dhcp,
				Ips:            &ipsNic,
				FirewallActive: &dhcpNic,
			},
		},
	}
	testNicCreate = resources.Nic{
		Nic: ionoscloud.Nic{
			Properties: &ionoscloud.NicProperties{
				Name:           &testNicVar,
				Lan:            &lanNicId,
				Dhcp:           &dhcpNic,
				Ips:            &ipsNic,
				FirewallActive: &dhcpNic,
				FirewallType:   &testNicVar,
			},
		},
	}
	ns = resources.Nics{
		Nics: ionoscloud.Nics{
			Id:    &testNicVar,
			Items: &[]ionoscloud.Nic{n},
		},
	}
	balancedns = resources.BalancedNics{
		BalancedNics: ionoscloud.BalancedNics{
			Id:    &testNicVar,
			Items: &[]ionoscloud.Nic{n},
		},
	}
	testNicVar    = "test-nic"
	testNicNewVar = "test-new-nic"
	testNicErr    = errors.New("nic test: error occurred")
)

func TestNicCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(NicCmd())
	if ok := NicCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		err := PreRunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNicListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ns, &testResponse, nil)
		err := RunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Nics{}, &testResponse, nil)
		err := RunNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(ns, nil, testNicErr)
		err := RunNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, &testResponse, nil)
		err := RunNicGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().Get(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNicId)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallType), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().Create(testNicVar, testNicVar, testNicCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, &testResponse, nil)
		err := RunNicCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNicId)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallType), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().Create(testNicVar, testNicVar, testNicCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, testNicErr)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), ipsNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNicId)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFirewallType), testNicVar)
		rm.CloudApiV6Mocks.Nic.EXPECT().Create(testNicVar, testNicVar, testNicCreate, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNicCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNewNicId)
		rm.CloudApiV6Mocks.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nicNew, &testResponse, nil)
		err := RunNicUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNewNicId)
		rm.CloudApiV6Mocks.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nicNew, nil, testNicErr)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNicNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDhcp), dhcpNewNic)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLanId), lanNewNicId)
		rm.CloudApiV6Mocks.Nic.EXPECT().Update(testNicVar, testNicVar, testNicVar, nicProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&nicNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNicUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nicsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nicsList, nil, testNicErr)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Nics{}, &testResponse, nil)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Nics{Nics: ionoscloud.Nics{Items: &[]ionoscloud.Nic{}}}, &testResponse, nil)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().List(testNicVar, testNicVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(nicsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testNicErr)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testNicErr)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNicDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Nic.EXPECT().Delete(testNicVar, testNicVar, testNicVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunNicDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNicDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testNicVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunNicDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNicsCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgCols), []string{"Name"})
	getNicsCols(core.GetGlobalFlagName("nic", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNicsColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("nic", config.ArgCols), []string{"Unknown"})
	getNicsCols(core.GetGlobalFlagName("nic", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

// LoadBalancer Nic

func TestPreRunDcNicLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		err := PreRunDcNicLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLoadBalancerNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		err := PreRunLoadBalancerNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLoadBalancerNicListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunLoadBalancerNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunLoadBalancerNicListFiltersErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunLoadBalancerNicList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNicLoadBalancerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDcNicLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerNicAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().AttachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLoadBalancerNicAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(balancedns, nil, nil)
		err := RunLoadBalancerNicList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(balancedns, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, nil)
		err := RunLoadBalancerNicGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().GetNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Nic{Nic: n}, nil, testLoadbalancerErr)
		err := RunLoadBalancerNicGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicDetachAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(balancedNicsList, nil, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicDetachAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(balancedNicsList, nil, testNicErr)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.BalancedNics{}, nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.BalancedNics{BalancedNics: ionoscloud.BalancedNics{Items: &[]ionoscloud.Nic{}}}, nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().ListNics(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(balancedNicsList, nil, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testNicErr)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testLoadbalancerErr)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerNicDetachAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().DetachNic(testLoadbalancerVar, testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunLoadBalancerNicDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerNicDetachAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		cfg.Stdin = os.Stdin
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNicId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		err := RunLoadBalancerNicDetach(cfg)
		assert.Error(t, err)
	})
}
