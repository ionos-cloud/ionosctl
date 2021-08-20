package commands

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"regexp"
	"testing"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	dcVersion = int32(1)
	dc        = ionoscloud.Datacenter{
		Id: &testDatacenterVar,
		Properties: &ionoscloud.DatacenterProperties{
			Name:              &testDatacenterVar,
			Description:       &testDatacenterVar,
			Location:          &testDatacenterVar,
			Version:           &dcVersion,
			Features:          &[]string{testDatacenterVar},
			SecAuthProtection: &testDatacenterBoolVar,
		},
		Metadata: &ionoscloud.DatacenterElementMetadata{
			State: &testStateVar,
		},
	}
	dcProperties = v5.DatacenterProperties{
		DatacenterProperties: ionoscloud.DatacenterProperties{
			Name:        &testDatacenterNewVar,
			Description: &testDatacenterNewVar,
		},
	}
	dcNew = v5.Datacenter{
		Datacenter: ionoscloud.Datacenter{
			Id: &testDatacenterVar,
			Properties: &ionoscloud.DatacenterProperties{
				Name:        dcProperties.DatacenterProperties.Name,
				Description: dcProperties.DatacenterProperties.Description,
				Location:    &testDatacenterVar,
			},
		},
	}
	dcs = v5.Datacenters{
		Datacenters: ionoscloud.Datacenters{
			Id:    &testDatacenterVar,
			Items: &[]ionoscloud.Datacenter{dc, dc},
		},
	}
	testDatacenterVar     = "test-datacenter"
	testDatacenterBoolVar = false
	testDatacenterNewVar  = "test-new-datacenter"
	testDatacenterErr     = errors.New("datacenter test error occurred")
)

func TestPreRunDataCenterId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		err := PreRunDataCenterId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunDataCenterIdRequiredFlagErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		err := PreRunDataCenterId(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().List().Return(dcs, nil, nil)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		err := RunDataCenterList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), "")
		rm.Datacenter.EXPECT().List().Return(dcs, nil, testDatacenterErr)
		err := RunDataCenterList(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		res := v5.Datacenter{Datacenter: dc}
		rm.Datacenter.EXPECT().Get(testDatacenterVar).Return(&res, nil, nil)
		err := RunDataCenterGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		res := v5.Datacenter{Datacenter: dc}
		rm.Datacenter.EXPECT().Get(testDatacenterVar).Return(&res, nil, testDatacenterErr)
		err := RunDataCenterGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		rm.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar).Return(&v5.Datacenter{dc}, nil, nil)
		err := RunDataCenterCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar).Return(&v5.Datacenter{dc}, nil, nil)
		err := RunDataCenterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgLocation), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Create(testDatacenterVar, testDatacenterVar, testDatacenterVar).Return(&v5.Datacenter{dc}, nil, testDatacenterErr)
		err := RunDataCenterCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties).Return(&dcNew, nil, nil)
		err := RunDataCenterUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties).Return(&dcNew, nil, nil)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties).Return(&dcNew, nil, testDatacenterErr)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterUpdateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgName), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDescription), testDatacenterNewVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Update(testDatacenterVar, dcProperties).Return(&dcNew, &testResponseErr, nil)
		err := RunDataCenterUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Delete(testDatacenterVar).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Delete(testDatacenterVar).Return(nil, testDatacenterErr)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), true)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		rm.Datacenter.EXPECT().Delete(testDatacenterVar).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunDataCenterDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.Datacenter.EXPECT().Delete(testDatacenterVar).Return(nil, nil)
		err := RunDataCenterDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunDataCenterDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocks) {
		viper.Reset()
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
		viper.Set(config.ArgQuiet, false)
		viper.Set(config.ArgForce, false)
		viper.Set(config.ArgServerUrl, config.DefaultApiURL)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgDataCenterId), testDatacenterVar)
		viper.Set(core.GetFlagName(cfg.NS, config.ArgWaitForRequest), false)
		cfg.Stdin = os.Stdin
		err := RunDataCenterDelete(cfg)
		assert.Error(t, err)
	})
}

func TestGetDatacentersCols(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("datacenter", config.ArgCols), []string{"Name"})
	getDataCenterCols(core.GetGlobalFlagName("datacenter", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
}

func TestGetDatacentersColsErr(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	viper.Set(core.GetGlobalFlagName("datacenter", config.ArgCols), []string{"Unknown"})
	getDataCenterCols(core.GetGlobalFlagName("datacenter", config.ArgCols), w)
	err := w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`unknown column Unknown`)
	assert.True(t, re.Match(b.Bytes()))
}

func TestGetDataCentersIds(t *testing.T) {
	defer func(a func()) { clierror.ErrAction = a }(clierror.ErrAction)
	var b bytes.Buffer
	clierror.ErrAction = func() { return }
	w := bufio.NewWriter(&b)
	err := os.Setenv(ionoscloud.IonosUsernameEnvVar, "user")
	assert.NoError(t, err)
	err = os.Setenv(ionoscloud.IonosPasswordEnvVar, "pass")
	assert.NoError(t, err)
	viper.Set(config.ArgServerUrl, config.DefaultApiURL)
	getDataCentersIds(w)
	err = w.Flush()
	assert.NoError(t, err)
	re := regexp.MustCompile(`401 Unauthorized`)
	assert.True(t, re.Match(b.Bytes()))
}
