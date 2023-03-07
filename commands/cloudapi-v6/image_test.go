package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strings"
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
	testImage = resources.Image{
		Image: ionoscloud.Image{
			Id: &testImageVar,
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
	testImages = resources.Images{
		Images: ionoscloud.Images{
			Id:    &testImageVar,
			Items: &[]ionoscloud.Image{testImage.Image, testImage.Image},
		},
	}
	testImageSize     = float32(2)
	testImagePublic   = true
	testImageVar      = "test-image"
	testImageUpperVar = strings.ToUpper(testImageVar)
	testImageErr      = errors.New("image test error")
)

func TestImageCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(ImageCmd())
	if ok := ImageCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreImageId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testImageVar)
		err := PreRunImageId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreImageIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunImageId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunImageList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		err := PreRunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunImageListFilter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("createdBy=%s", testQueryParamVar)})
		err := PreRunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunImageListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), allImageCols)
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testImages, &testResponse, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgCols), allImageCols)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.Images{}, &testResponse, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testImages, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testImageVar)
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testImages, nil, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListSortOptionErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), "no alias")
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testImages, nil, nil)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageListSortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageAlias), testImageVar)
		rm.CloudApiV6Mocks.Image.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(testImages, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testImageVar)
		rm.CloudApiV6Mocks.Image.EXPECT().Get(testImageVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, &testResponse, nil)
		err := RunImageGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgImageId), testImageVar)
		rm.CloudApiV6Mocks.Image.EXPECT().Get(testImageVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testImage, nil, testImageErr)
		err := RunImageGet(cfg)
		assert.Error(t, err)
	})
}
