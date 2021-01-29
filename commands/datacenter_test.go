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
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testVar = "test"
	dc      = ionoscloud.Datacenter{
		Id: &testVar,
		Properties: &ionoscloud.DatacenterProperties{
			Name:        &testVar,
			Description: &testVar,
			Location:    &testVar,
		},
	}
	dcs = resources.Datacenters{
		Datacenters: ionoscloud.Datacenters{
			Id:    &testVar,
			Items: &[]ionoscloud.Datacenter{dc, dc},
		},
	}

	testErr = errors.New("error occurred")
)

func TestRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Datacenter.EXPECT().List().Return(dcs, nil, nil)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), "")
		err := RunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterList_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		rm.Datacenter.EXPECT().List().Return(dcs, nil, testErr)
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), "")
		err := RunDataCenterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		res := resources.Datacenter{Datacenter: dc}
		rm.Datacenter.EXPECT().Get("test").Return(&res, nil, nil)
		err := RunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterGet_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		res := resources.Datacenter{Datacenter: dc}
		rm.Datacenter.EXPECT().Get(testVar).Return(&res, nil, testErr)
		err := RunDataCenterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterName), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterDescription), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterRegion), testVar)
		rm.Datacenter.EXPECT().Create(testVar, testVar, testVar).Return(&resources.Datacenter{dc}, nil, nil)
		err := RunDataCenterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterCreate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterName), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterDescription), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterRegion), testVar)
		rm.Datacenter.EXPECT().Create(testVar, testVar, testVar).Return(&resources.Datacenter{dc}, nil, testErr)
		err := RunDataCenterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterDescription), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterRegion), testVar)
		rm.Datacenter.EXPECT().Update(testVar, testVar, testVar).Return(&resources.Datacenter{dc}, nil, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterUpdate_RequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), "")
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == utils.NewRequiredFlagErr(config.ArgDataCenterId).Error())
	})
}

func TestRunDataCenterUpdate_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterDescription), testVar)
		rm.Datacenter.EXPECT().Update(testVar, testVar, testVar).Return(&resources.Datacenter{dc}, nil, testErr)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		rm.Datacenter.EXPECT().Delete(testVar).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDelete_RequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), "")
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
		assert.True(t, err.Error() == utils.NewRequiredFlagErr(config.ArgDataCenterId).Error())
	})
}

func TestRunDataCenterDelete_Err(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, true)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		rm.Datacenter.EXPECT().Delete(testVar).Return(nil, testErr)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDelete_AskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Datacenter.EXPECT().Delete(testVar).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDelete_AskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	builder.CmdConfigTest(t, w, func(cfg *builder.CommandConfig, rm *builder.ResourcesMocks) {
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgIgnoreStdin, false)
		viper.Set(builder.GetFlagName(cfg.Name, config.ArgDataCenterId), testVar)
		cfg.Stdin = os.Stdin
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetDataCentersIds(t *testing.T) {
	defer func(a func()) { utils.ErrAction = a }(utils.ErrAction)
	var b bytes.Buffer
	utils.ErrAction = func() {}

	w := bufio.NewWriter(&b)
	viper.Set(config.ArgConfig, "../pkg/testdata/ionosctl-config.json")
	getDataCentersIds(w)
	err := w.Flush()
	assert.NoError(t, err)

	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
