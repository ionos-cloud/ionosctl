package commands

import (
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
	networkloadbalancerTest = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Properties: &ionoscloud.NetworkLoadBalancerProperties{
				Name:         &testNetworkLoadBalancerVar,
				Ips:          &[]string{testNetworkLoadBalancerVar},
				TargetLan:    &testNetworkLoadBalancerIntVar,
				ListenerLan:  &testNetworkLoadBalancerIntVar,
				LbPrivateIps: &[]string{testNetworkLoadBalancerVar},
			},
		},
	}
	networkloadbalancerTestId = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Id: &testNetworkLoadBalancerVar,
			Properties: &ionoscloud.NetworkLoadBalancerProperties{
				Name:         &testNetworkLoadBalancerVar,
				Ips:          &[]string{testNetworkLoadBalancerVar},
				TargetLan:    &testNetworkLoadBalancerIntVar,
				ListenerLan:  &testNetworkLoadBalancerIntVar,
				LbPrivateIps: &[]string{testNetworkLoadBalancerVar},
			},
		},
	}
	networkloadbalancerTestGet = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Id:         &testNetworkLoadBalancerVar,
			Properties: networkloadbalancerTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	networkloadbalancers = resources.NetworkLoadBalancers{
		NetworkLoadBalancers: ionoscloud.NetworkLoadBalancers{
			Id:    &testNetworkLoadBalancerVar,
			Items: &[]ionoscloud.NetworkLoadBalancer{networkloadbalancerTest.NetworkLoadBalancer},
		},
	}
	networkloadbalancersList = resources.NetworkLoadBalancers{
		NetworkLoadBalancers: ionoscloud.NetworkLoadBalancers{
			Id: &testNetworkLoadBalancerVar,
			Items: &[]ionoscloud.NetworkLoadBalancer{
				networkloadbalancerTestId.NetworkLoadBalancer,
				networkloadbalancerTestId.NetworkLoadBalancer,
			},
		},
	}
	networkloadbalancerProperties = resources.NetworkLoadBalancerProperties{
		NetworkLoadBalancerProperties: ionoscloud.NetworkLoadBalancerProperties{
			Name:         &testNetworkLoadBalancerNewVar,
			Ips:          &[]string{testNetworkLoadBalancerNewVar},
			TargetLan:    &testNetworkLoadBalancerNewIntVar,
			ListenerLan:  &testNetworkLoadBalancerNewIntVar,
			LbPrivateIps: &[]string{testNetworkLoadBalancerNewVar},
		},
	}
	networkloadbalancerNew = resources.NetworkLoadBalancer{
		NetworkLoadBalancer: ionoscloud.NetworkLoadBalancer{
			Properties: &networkloadbalancerProperties.NetworkLoadBalancerProperties,
		},
	}
	testNetworkLoadBalancerIntVar    = int32(1)
	testNetworkLoadBalancerNewIntVar = int32(2)
	testNetworkLoadBalancerVar       = "test-networkloadbalancer"
	testNetworkLoadBalancerNewVar    = "test-new-networkloadbalancer"
	testNetworkLoadBalancerErr       = errors.New("networkloadbalancer test error")
)

func TestNetworkloadbalancerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(NetworkloadbalancerCmd())
	if ok := NetworkloadbalancerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunNetworkLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		err := PreRunNetworkLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerListFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunNetworkLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunNetworkLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunNetworkLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		err := PreRunDcNetworkLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcNetworkLoadBalancerIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		err := PreRunDcNetworkLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(resources.ListQueryParams{}).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testDatacenterVar, resources.ListQueryParams{}).Return(networkloadbalancersList, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunNetworkLoadBalancerListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(networkloadbalancers, &testResponse, nil)
		err := RunNetworkLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, testListQueryParam).Return(resources.NetworkLoadBalancers{}, &testResponse, nil)
		err := RunNetworkLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(networkloadbalancers, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&networkloadbalancerTestGet, &testResponse, nil)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&networkloadbalancerTestGet, nil, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&networkloadbalancerTestGet, nil, nil)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForState), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&networkloadbalancerTestGet, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Get(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&networkloadbalancerTestGet, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest, resources.QueryParams{}).Return(&networkloadbalancerTest, &testResponse, nil)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest, resources.QueryParams{}).Return(&networkloadbalancerTest, &testResponse, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest, resources.QueryParams{}).Return(&networkloadbalancerTest, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Create(testNetworkLoadBalancerVar, networkloadbalancerTest, resources.QueryParams{}).Return(&networkloadbalancerTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties, resources.QueryParams{}).
			Return(&networkloadbalancerNew, &testResponse, nil)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties, resources.QueryParams{}).Return(&networkloadbalancerNew, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPrivateIps), testNetworkLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgTargetLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgListenerLan), testNetworkLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Update(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, networkloadbalancerProperties, resources.QueryParams{}).Return(&networkloadbalancerNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(networkloadbalancersList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(networkloadbalancersList, nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(resources.NetworkLoadBalancers{}, &testResponse, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(
			resources.NetworkLoadBalancers{NetworkLoadBalancers: ionoscloud.NetworkLoadBalancers{Items: &[]ionoscloud.NetworkLoadBalancer{}}}, &testResponse, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(config.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().List(testNetworkLoadBalancerVar, resources.ListQueryParams{}).Return(networkloadbalancersList, &testResponse, nil)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, testNetworkLoadBalancerErr)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(nil, testNetworkLoadBalancerErr)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar, resources.QueryParams{}).Return(&testRequestStatus, nil, testRequestErr)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.NetworkLoadBalancer.EXPECT().Delete(testNetworkLoadBalancerVar, testNetworkLoadBalancerVar, resources.QueryParams{}).Return(nil, nil)
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunNetworkLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testNetworkLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgNetworkLoadBalancerId), testNetworkLoadBalancerVar)
		cfg.Stdin = os.Stdin
		err := RunNetworkLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetNetworkLoadBalancersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), []string{"Name"})
	getNetworkLoadBalancersCols(
		core.GetGlobalFlagName("networkloadbalancer", config.ArgCols),
		core.GetFlagName("networkloadbalancer", cloudapiv6.ArgAll),
		w,
	)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetNetworkLoadBalancersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("networkloadbalancer", config.ArgCols), []string{"Unknown"})
	getNetworkLoadBalancersCols(
		core.GetGlobalFlagName("networkloadbalancer", config.ArgCols),
		core.GetFlagName("networkloadbalancer", cloudapiv6.ArgAll),
		w,
	)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}
