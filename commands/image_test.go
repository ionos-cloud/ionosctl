package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	img = resources.Image{
		Image: ionoscloud.Image{
			Id: &testImageVar,
			Properties: &ionoscloud.ImageProperties{
				Name:        &testImageVar,
				Location:    &testImageVar,
				Description: &testImageVar,
				Size:        &testImageSize,
				LicenceType: &testImageVar,
				ImageType:   &testImageVar,
				Public:      &testImagePublic,
			},
		},
	}
	images = resources.Images{
		Images: ionoscloud.Images{
			Id:    &testImageVar,
			Items: &[]ionoscloud.Image{img.Image},
		},
	}
	testImageSize   = float32(2)
	testImagePublic = true
	testImageVar    = "test-image"
	testImageErr    = errors.New("image test error")
)

func TestPreImageIdValidate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		err := PreRunImageIdValidate(cfg)
		assert.NoError(t, err)
	})
}

func TestPreImageIdValidateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.PreCmdConfigTest(t, w, func(cfg *builder.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), "")
		err := PreRunImageIdValidate(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Image.EXPECT().List().Return(images, nil, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		rm.Image.EXPECT().List().Return(images, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageListSort(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageLocation), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageLicenceType), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageType), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageSize), testImageSize)
		rm.Image.EXPECT().List().Return(images, nil, nil)
		err := RunImageList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageListSortErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageLocation), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageLicenceType), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageType), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageSize), testImageSize)
		rm.Image.EXPECT().List().Return(images, nil, testImageErr)
		err := RunImageList(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		rm.Image.EXPECT().Get(testImageVar).Return(&img, nil, nil)
		err := RunImageGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		rm.Image.EXPECT().Get(testImageVar).Return(&img, nil, testImageErr)
		err := RunImageGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Image.EXPECT().Delete(testImageVar).Return(nil, nil)
		err := RunImageDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		rm.Image.EXPECT().Delete(testImageVar).Return(nil, testImageErr)
		err := RunImageDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunImageDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Image.EXPECT().Delete(testImageVar).Return(nil, nil)
		err := RunImageDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunImageDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgImageId), testImageVar)
		viper.Set(builder.GetFlagName(cfg.ParentName, cfg.Name, config.ArgWait), false)
		cfg.Stdin = os.Stdin
		err := RunImageDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetImagesCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("image", config.ArgCols), []string{"Name"})
	getImageCols(builder.GetGlobalFlagName("image", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetImagesColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() {}
	w := bufio.NewWriter(&b)
	viper.Set(builder.GetGlobalFlagName("image", config.ArgCols), []string{"Unknown"})
	getImageCols(builder.GetGlobalFlagName("image", config.ArgCols), w)
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
