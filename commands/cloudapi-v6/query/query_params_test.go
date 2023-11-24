package query

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testFiltersVar    = []string{"test", "testing", "filter", "validate"}
	testFilterVar     = "name"
	testMaxResultsVar = int32(2)
)

func TestValidateFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().StringSlice(cloudapiv6.ArgFilters, []string{}, "")
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "test=test")
		err := ValidateFilters(cfg, testFiltersVar, "")
		assert.NoError(t, err)
	})
}

func TestValidateFiltersLengthErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		err := ValidateFilters(cfg, testFiltersVar, "")
		assert.Error(t, err)
	})
}

func TestValidateFiltersFormatErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "test")
		err := ValidateFilters(cfg, testFiltersVar, "")
		assert.Error(t, err)
	})
}

func TestValidateFiltersInvalidsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "name=test,location=test")
		err := ValidateFilters(cfg, testFiltersVar, "")
		assert.Error(t, err)
	})
}

func TestValidateFiltersInvalidErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "name=test")
		err := ValidateFilters(cfg, testFiltersVar, "")
		assert.Error(t, err)
	})
}

func TestGetListQueryParams(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().StringSlice(cloudapiv6.ArgFilters, []string{}, "")
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "name=test,location=test")
		// cfg.Command.Command.Flags().Set(cloudapiv6.ArgOrderBy, testFilterVar)
		// cfg.Command.Command.Flags().(constants.FlagMaxResults, testMaxResultsVar)
		result, err := GetListQueryParams(cfg)
		assert.NoError(t, err)
		assert.True(t, result.Filters != nil)
		filtersKV := *result.Filters
		assert.True(t, filtersKV["name"][0] == "test")
		assert.True(t, filtersKV["location"][0] == "test")
		// assert.True(t, *result.OrderBy == testFilterVar) Muted temporarily due to viper pflag mapping removal
		// assert.True(t, *result.MaxResults == testMaxResultsVar)
	})
}

func TestGetListQueryParamsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(cloudapiv6.ArgFilters, "name")
		_, err := GetListQueryParams(cfg)
		assert.NoError(t, err)
	})
}
