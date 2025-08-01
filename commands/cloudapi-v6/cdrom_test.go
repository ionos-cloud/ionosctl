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
	testCdroms = resources.Cdroms{
		Cdroms: ionoscloud.Cdroms{
			Items: &[]ionoscloud.Image{testImage.Image},
		},
	}
	testCdromsList = resources.Cdroms{
		Cdroms: ionoscloud.Cdroms{
			Id: &testCdromVar,
			Items: &[]ionoscloud.Image{
				testImageCdRoms.Image,
				testImageCdRoms.Image,
			},
		},
	}
	testImageCdRoms = resources.Image{
		Image: ionoscloud.Image{
			Id: &testCdromVar,
			Properties: &ionoscloud.ImageProperties{
				Name:         &testImageVar,
				Location:     &testImageVar,
				Description:  &testImageVar,
				Size:         &testImageSize,
				LicenceType:  &testImageUpperVar,
				ImageType:    &testImageUpperVar,
				Public:       &testImagePublic,
				ImageAliases: &[]string{testImageVar},
				CloudInit:    &testImageVar,
			},
			Metadata: &ionoscloud.DatacenterElementMetadata{
				CreatedDate:     &testIonosTime,
				CreatedBy:       &testImageVar,
				CreatedByUserId: &testImageVar,
			},
		},
	}
	testCdromVar = "test-cdrom"
	testCdromErr = errors.New("cdrom test error")
)

func TestPreRunDcServerCdromIds(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		err := PreRunDcServerCdromIds(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDcServerCdromIdsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunDcServerCdromIds(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunServerCdromList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		err := PreRunServerCdromList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCdromListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		err := PreRunServerCdromList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunServerCdromListFilter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("createdBy=%s", testQueryParamVar))
		err := PreRunServerCdromList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromAttach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, &testResponse, nil)
		err := RunServerCdromAttach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromAttachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromAttachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().AttachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerCdromAttach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromsList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testCdroms, &testResponse, nil)
		err := RunServerCdromsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromsListListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Cdroms{}, &testResponse, nil)
		err := RunServerCdromsList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromsListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testCdroms, nil, testCdromErr)
		err := RunServerCdromsList(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, &testResponse, nil)
		err := RunServerCdromGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetach(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testCdromsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testCdromsList, nil, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Cdroms{}, &testResponse, nil)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.Cdroms{Cdroms: ionoscloud.Cdroms{Items: &[]ionoscloud.Image{}}}, &testResponse, nil)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.Server.EXPECT().ListCdroms(testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testListQueryParam)).Return(testCdromsList, &testResponse, nil)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponseErr, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunServerCdromDetachWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}

func TestRunCdromDetachAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("YES\n")))
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		rm.CloudApiV6Mocks.Server.EXPECT().DetachCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunServerCdromDetach(cfg)
		assert.NoError(t, err)
	})
}

func TestRunServerCdromDetachAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgDataCenterId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgServerId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgCdromId), testCdromVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), false)
		cfg.Command.Command.SetIn(bytes.NewReader([]byte("\n")))
		rm.CloudApiV6Mocks.Server.EXPECT().GetCdrom(testCdromVar, testCdromVar, testCdromVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testCdromErr)
		err := RunServerCdromDetach(cfg)
		assert.Error(t, err)
	})
}
