package commands

import (
	"bufio"
	"bytes"
	"errors"
	"regexp"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
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

func TestPreImageId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageId), testImageVar)
		err := PreRunImageId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreImageIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunImageId(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgCols), allImageCols)
		rm.Image.EXPECT().List().Return(testImages, nil, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Image.EXPECT().List().Return(testImages, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageAlias), testImageVar)
		rm.Image.EXPECT().List().Return(testImages, nil, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListSortOptionErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageAlias), "no alias")
		rm.Image.EXPECT().List().Return(testImages, nil, nil)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageListSortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLicenceType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgType), testImageVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLatest), 1)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageAlias), testImageVar)
		rm.Image.EXPECT().List().Return(testImages, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageId), testImageVar)
		rm.Image.EXPECT().Get(testImageVar).Return(&testImage, nil, nil)
		err := RunImageGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgImageId), testImageVar)
		rm.Image.EXPECT().Get(testImageVar).Return(&testImage, nil, testImageErr)
		err := RunImageGet(cfg)
		assert.Error(t, err)
	})
}

func TestGetImagesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("image", config.ArgCols), []string{"Name"})
	getImageCols(core.GetGlobalFlagName("image", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetImagesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("image", config.ArgCols), []string{"Unknown"})
	getImageCols(core.GetGlobalFlagName("image", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetImagesIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/config.json")
	getImageIds(w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
