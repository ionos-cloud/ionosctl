package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go-bundle/products/compute"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	backupUnitTest = resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Email:    &testBackupUnitVar,
				Name:     &testBackupUnitVar,
				Password: &testBackupUnitVar,
			},
		},
	}
	backupUnitTestId = resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Id: &testBackUnitId,
			Properties: &ionoscloud.BackupUnitProperties{
				Email:    &testBackupUnitVar,
				Name:     &testBackupUnitVar,
				Password: &testBackupUnitVar,
			},
		},
	}
	backupUnitsList = resources.BackupUnits{
		BackupUnits: ionoscloud.BackupUnits{
			Id: &testBackUnitId,
			Items: &[]ionoscloud.BackupUnit{
				backupUnitTestId.BackupUnit,
				backupUnitTestId.BackupUnit,
			},
		},
	}
	backupUnitTestGet = resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Id:         &testBackupUnitVar,
			Properties: backupUnitTest.Properties,
			Metadata:   &ionoscloud.DatacenterElementMetadata{State: &testStateVar},
		},
	}
	backupUnitTestGetSSO = resources.BackupUnitSSO{
		BackupUnitSSO: ionoscloud.BackupUnitSSO{
			SsoUrl: &testBackupUnitVar,
		},
	}
	backupUnits = resources.BackupUnits{
		BackupUnits: ionoscloud.BackupUnits{
			Id:    &testBackupUnitVar,
			Items: &[]ionoscloud.BackupUnit{backupUnitTest.BackupUnit},
		},
	}
	backupUnitProperties = resources.BackupUnitProperties{
		BackupUnitProperties: ionoscloud.BackupUnitProperties{
			Email:    &testBackupUnitNewVar,
			Password: &testBackupUnitNewVar,
		},
	}
	backupUnitNew = resources.BackupUnit{
		BackupUnit: ionoscloud.BackupUnit{
			Properties: &ionoscloud.BackupUnitProperties{
				Name:     &testBackupUnitVar,
				Email:    &testBackupUnitNewVar,
				Password: &testBackupUnitNewVar,
			},
		},
	}
	testListQueryParamFilters = resources.ListQueryParams{
		Filters: &map[string]string{
			testQueryParamVar: testQueryParamVar,
		},
		OrderBy:    &testQueryParamVar,
		MaxResults: &testMaxResultsVar,
		QueryParams: resources.QueryParams{
			Depth: &testDepthListVar,
		},
	}
	testListQueryParam = resources.ListQueryParams{
		OrderBy:    &testOrderByVar,
		MaxResults: &testMaxResultsVar,
		QueryParams: resources.QueryParams{
			Depth: &testDepthListVar,
		},
	}
	testQueryParamOther = resources.QueryParams{
		Depth: &testDepthOtherVar,
	}
	testDepthListVar     = int32(1)
	testDepthOtherVar    = int32(0)
	testQueryParamVar    = "test-filter"
	testMaxResultsVar    = cloudapiv6.DefaultMaxResults
	testOrderByVar       = "" // default orderBy. Add to cloudapi constants?
	testBackupUnitVar    = "test-backup-unit"
	testBackUnitId       = "87aa25ec-5f74-4927-bd95-c8e42db06fe2"
	testBackupUnitNewVar = "test-new-backup-unit"
	testBackupUnitErr    = errors.New("backup-unit test error")
)

func TestBackupunitCmd(t *testing.T) {
	var err error
	core.RootCmdTest.AddCommand(BackupunitCmd())
	if ok := BackupunitCmd().IsAvailableCommand(); !ok {
		err = errors.New("non-available cmd")
	}
	assert.NoError(t, err)
}

func TestPreRunBackupUnitListNoFilter(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunBackupUnitList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("name=%s", testQueryParamVar)})
		err := PreRunBackupUnitList(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		err := PreRunBackupUnitList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunBackupUnitListFormatErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), testBackupUnitVar)
		err := PreRunBackupUnitList(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunBackupUnitId(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		err := PreRunBackupUnitId(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitIdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunBackupUnitId(cfg)
		assert.Error(t, err)
	})
}

func TestPreRunBackupUnitNameEmailPwd(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testBackupUnitVar)
		err := PreRunBackupUnitNameEmailPwd(cfg)
		assert.NoError(t, err)
	})
}

func TestPreRunBackupUnitNameEmailPwdErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := PreRunBackupUnitNameEmailPwd(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitList(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(backupUnits, nil, nil)
		err := RunBackupUnitList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgFilters), []string{fmt.Sprintf("%s=%s", testQueryParamVar, testQueryParamVar)})
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgOrderBy), testQueryParamVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.FlagMaxResults), testMaxResultsVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.BackupUnits{}, nil, nil)
		err := RunBackupUnitList(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(backupUnits, nil, testBackupUnitErr)
		err := RunBackupUnitList(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitGet(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Get(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTestGet, &testResponse, nil)
		err := RunBackupUnitGet(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitGetErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Get(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTestGet, nil, testBackupUnitErr)
		err := RunBackupUnitGet(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitGetSsoUrl(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().GetSsoUrl(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTestGetSSO, &testResponse, nil)
		err := RunBackupUnitGetSsoUrl(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitGetSsoUrlErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().GetSsoUrl(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTestGetSSO, nil, testBackupUnitErr)
		err := RunBackupUnitGetSsoUrl(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitCreate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Create(backupUnitTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTest, &testResponse, nil)
		err := RunBackupUnitCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitCreateResponseErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Create(backupUnitTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTest, &testResponseErr, testBackupUnitErr)
		err := RunBackupUnitCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitCreateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Create(backupUnitTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTest, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, nil)
		err := RunBackupUnitCreate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitCreateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgName), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Create(backupUnitTest, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitTest, nil, testBackupUnitErr)
		err := RunBackupUnitCreate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitUpdate(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgVerbose, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitNewVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitNew, &testResponse, nil)
		err := RunBackupUnitUpdate(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitUpdateWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitNewVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitNew, &testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, testRequestErr)
		err := RunBackupUnitUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitUpdateErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgPassword), testBackupUnitNewVar)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgEmail), testBackupUnitNewVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Update(testBackupUnitVar, backupUnitProperties, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&backupUnitNew, nil, testBackupUnitErr)
		err := RunBackupUnitUpdate(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDelete(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteAll(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(backupUnitsList, &testResponse, nil)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackUnitId, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackUnitId, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteAllListErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(backupUnitsList, nil, testBackupUnitErr)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteAllItemsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(resources.BackupUnits{}, &testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteAllLenErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(
			resources.BackupUnits{BackupUnits: ionoscloud.BackupUnits{Items: &[]ionoscloud.BackupUnit{}}}, &testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteAllErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgAll), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().List(gomock.AssignableToTypeOf(testListQueryParam)).Return(backupUnitsList, &testResponse, nil)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackUnitId, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, testBackupUnitErr)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackUnitId, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, testBackupUnitErr)
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}

func TestRunBackupUnitDeleteWaitErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, true)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		viper.Set(core.GetFlagName(cfg.NS, constants.ArgWaitForRequest), true)
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(&testResponse, nil)
		rm.CloudApiV6Mocks.Request.EXPECT().GetStatus(testRequestIdVar).Return(&testRequestStatus, nil, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteAskForConfirm(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		cfg.Stdin = bytes.NewReader([]byte("YES\n"))
		rm.CloudApiV6Mocks.BackupUnit.EXPECT().Delete(testBackupUnitVar, gomock.AssignableToTypeOf(testQueryParamOther)).Return(nil, nil)
		err := RunBackupUnitDelete(cfg)
		assert.NoError(t, err)
	})
}

func TestRunBackupUnitDeleteAskForConfirmErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgForce, false)
		viper.Set(core.GetFlagName(cfg.NS, cloudapiv6.ArgBackupUnitId), testBackupUnitVar)
		cfg.Stdin = os.Stdin
		err := RunBackupUnitDelete(cfg)
		assert.Error(t, err)
	})
}
