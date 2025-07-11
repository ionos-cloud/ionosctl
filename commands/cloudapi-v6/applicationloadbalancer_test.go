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
	applicationloadbalancerTest = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Properties: &ionoscloud.ApplicationLoadBalancerProperties{
				Name:         &testApplicationLoadBalancerVar,
				Ips:          &[]string{testApplicationLoadBalancerVar},
				TargetLan:    &testApplicationLoadBalancerIntVar,
				ListenerLan:  &testApplicationLoadBalancerIntVar,
				LbPrivateIps: &[]string{testApplicationLoadBalancerVar},
			},
		},
	}
	applicationloadbalancerTestGet = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Id:         &testApplicationLoadBalancerVar,
			Properties: applicationloadbalancerTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	applicationloadbalancers = resources.ApplicationLoadBalancers{
		ApplicationLoadBalancers: ionoscloud.ApplicationLoadBalancers{
			Id:    &testApplicationLoadBalancerVar,
			Items: &[]ionoscloud.ApplicationLoadBalancer{applicationloadbalancerTestGet.ApplicationLoadBalancer},
		},
	}
	applicationloadbalancerProperties = resources.ApplicationLoadBalancerProperties{
		ApplicationLoadBalancerProperties: ionoscloud.ApplicationLoadBalancerProperties{
			Name:         &testApplicationLoadBalancerNewVar,
			Ips:          &[]string{testApplicationLoadBalancerNewVar},
			TargetLan:    &testApplicationLoadBalancerNewIntVar,
			ListenerLan:  &testApplicationLoadBalancerNewIntVar,
			LbPrivateIps: &[]string{testApplicationLoadBalancerNewVar},
		},
	}
	applicationloadbalancerNew = resources.ApplicationLoadBalancer{
		ApplicationLoadBalancer: ionoscloud.ApplicationLoadBalancer{
			Properties: &applicationloadbalancerProperties.ApplicationLoadBalancerProperties,
		},
	}
	testApplicationLoadBalancerIntVar    = int32(1)
	testApplicationLoadBalancerNewIntVar = int32(2)
	testApplicationLoadBalancerVar       = "test-applicationloadbalancer"
	testApplicationLoadBalancerNewVar    = "test-new-applicationloadbalancer"
	testApplicationLoadBalancerErr       = errors.New("applicationloadbalancer test error")
)

func TestApplicationLoadBalancerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ApplicationLoadBalancerCmd())
	if ok := ApplicationLoadBalancerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunDcApplicationLoadBalancerIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		err := PreRunDcApplicationLoadBalancerIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcApplicationLoadBalancerIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		err := PreRunDcApplicationLoadBalancerIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunApplicationLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		err := PreRunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunApplicationLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		err := PreRunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerListAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.Datacenter.EXPECT().List(cloudapiv6.ParentResourceListQueryParams).Return(dcs, &testResponse, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testDatacenterVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil).Times(len(getDataCenters(dcs)))
		err := RunApplicationLoadBalancerListAll(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil)
		err := RunApplicationLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.FlagFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil)
		err := RunApplicationLoadBalancerList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerList(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, nil, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, &testResponse, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetWait(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForState), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, nil, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, nil, nil)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerGetWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForState), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Get(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTestGet, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTest, nil, nil)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTest, &testResponse, nil)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTest, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Create(testApplicationLoadBalancerVar, applicationloadbalancerTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerNew, nil, nil)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerNewIntVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerNew, nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagName), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagPrivateIps), testApplicationLoadBalancerNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagTargetLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagListenerLan), testApplicationLoadBalancerNewIntVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Update(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, applicationloadbalancerProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&applicationloadbalancerNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllResponse(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(constants.FlagVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAllAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagAll), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().List(testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(applicationloadbalancers, &testResponse, nil)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testApplicationLoadBalancerErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagWaitForRequest), true)
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		rm.CloudApiV6Mocks.ApplicationLoadBalancer.EXPECT().Delete(testApplicationLoadBalancerVar, testApplicationLoadBalancerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunApplicationLoadBalancerDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.FlagQuiet, false)
		viper.Set(constants.FlagOutput, constants.DefaultOutputFormat)
		viper.Set(constants.FlagForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagDataCenterId), testApplicationLoadBalancerVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.FlagApplicationLoadBalancerId), testApplicationLoadBalancerVar)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		err := RunApplicationLoadBalancerDelete(cfg)
		assert.Error(t, err)
	})
}
