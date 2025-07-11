package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	dhcpLoadbalancer    = true
	dhcpLoadbalancerNew = false
	loadb               = ionoscloud.Loadbalancer{
		Id: &testLoadbalancerVar,
		Properties: &ionoscloud.LoadbalancerProperties{
			Name: &testLoadbalancerVar,
			Dhcp: &dhcpLoadbalancer,
			Ip:   &testLoadbalancerVar,
		},
	}
	loadbs = resources.Loadbalancers{
		Loadbalancers: ionoscloud.Loadbalancers{
			Id:    &testLoadbalancerVar,
			Items: &[]ionoscloud.Loadbalancer{loadb},
		},
	}
	lbList = resources.Loadbalancers{
		Loadbalancers: ionoscloud.Loadbalancers{
			Id: &testLoadbalancerVar,
			Items: &[]ionoscloud.Loadbalancer{
				loadb,
				loadb,
			},
		},
	}
	loadbalancerProperties = resources.LoadbalancerProperties{
		LoadbalancerProperties: ionoscloud.LoadbalancerProperties{
			Name: &testLoadbalancerNewVar,
			Dhcp: &dhcpLoadbalancerNew,
			Ip:   &testLoadbalancerNewVar,
		},
	}
	loadbalancerNew = resources.Loadbalancer{
		Loadbalancer: ionoscloud.Loadbalancer{
			Id:         &testLoadbalancerVar,
			Properties: &loadbalancerProperties.LoadbalancerProperties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	testLoadbalancerVar    = "test-loadbalancer"
	testLoadbalancerNewVar = "test-new-loadbalancer"
	testLoadbalancerErr    = errors.New("loadbalancer test: error occurred")
)

func TestLoadBalancerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(LoadBalancerCmd())
	if ok := LoadBalancerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDcLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		err := PreRunDcLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLoadBalancerIdsRequiredFlagsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		err := PreRunDcLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		err := PreRunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLoadBalancerListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testDatacenterVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(lbList, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunLoadBalancerListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(loadbs, &testResponse, nil)
		err := RunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Loadbalancers{}, &testResponse, nil)
		err := RunLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(loadbs, nil, testLoadbalancerErr)
		err := RunLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Loadbalancer{Loadbalancer: loadb}, &testResponse, nil)
		err := RunLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Get(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Loadbalancer{Loadbalancer: loadb}, &testResponse, nil)
		err := RunLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Loadbalancer{Loadbalancer: loadb}, nil, testLoadbalancerErr)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancer)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Create(testLoadbalancerVar, testLoadbalancerVar, dhcpLoadbalancer, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.Loadbalancer{Loadbalancer: loadb}, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loadbalancerNew, &testResponse, nil)
		err := RunLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loadbalancerNew, nil, testLoadbalancerErr)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loadbalancerNew, &testResponse, testLoadbalancerErr)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDhcp), dhcpLoadbalancerNew)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIp), testLoadbalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Update(testLoadbalancerVar, testLoadbalancerVar, loadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&loadbalancerNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(lbList, &testResponse, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(lbList, nil, testLanErr)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Loadbalancers{}, &testResponse, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Loadbalancers{Loadbalancers: ionoscloud.Loadbalancers{Items: &[]ionoscloud.Loadbalancer{}}}, &testResponse, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagServerUrl, constants.DefaultApiURL)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().List(testLoadbalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(lbList, &testResponse, nil)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testLoadbalancerErr)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testLoadbalancerErr)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		rm.CloudApiV6Mocks.Loadbalancer.EXPECT().Delete(testLoadbalancerVar, testLoadbalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagLoadBalancerId), testLoadbalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}
