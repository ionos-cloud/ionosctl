package commands

import (
	"bufio"
	"bytes"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	"github.com/ionos-cloud/sdk-go-bundle/products/compute/v2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testIpConsumer = compute.IpBlock{
		Id: &testIpConsumerVar,
		Properties: &compute.IpBlockProperties{
			IpConsumers: &[]compute.IpConsumer{
				{
					Ip:              &testIpConsumerVar,
					Mac:             &testIpConsumerVar,
					NicId:           &testIpConsumerVar,
					ServerId:        &testIpConsumerVar,
					ServerName:      &testIpConsumerVar,
					DatacenterId:    &testIpConsumerVar,
					DatacenterName:  &testIpConsumerVar,
					K8sNodePoolUuid: &testIpConsumerVar,
					K8sClusterUuid:  &testIpConsumerVar,
				},
			},
		},
	}
	testIpConsumerProperties = compute.IpBlock{
		Id: &testIpConsumerVar,
	}
	testIpConsumerGet = compute.IpBlock{
		Id: &testIpConsumerVar,
		Properties: &compute.IpBlockProperties{
			Name: &testIpConsumerVar,
		},
	}
	testIpConsumerVar = "test-ip-consumer"
	testIpConsumerErr = errors.New("ip consumer test: error occurred")
)

func TestIpconsumerCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(IpconsumerCmd())
	if ok := IpconsumerCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestRunIpConsumersList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpConsumerVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpConsumerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.IpBlock{IpBlock: testIpConsumer}, &testResponse, nil)
		err := RunIpConsumersList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunIpConsumersListPropertiesErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpConsumerVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpConsumerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.IpBlock{IpBlock: testIpConsumerProperties}, nil, nil)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpConsumersListGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpConsumerVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpConsumerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.IpBlock{IpBlock: testIpConsumerGet}, nil, nil)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}

func TestRunIpConsumersListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgIpBlockId), testIpConsumerVar)
		rm.CloudApiV6Mocks.IpBlocks.EXPECT().Get(testIpConsumerVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&resources.IpBlock{IpBlock: testIpConsumer}, nil, testIpConsumerErr)
		err := RunIpConsumersList(cfg)
		assert.Error(t, err)
	})
}
