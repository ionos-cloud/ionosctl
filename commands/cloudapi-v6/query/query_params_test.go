package query

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	testFiltersVar = []string{"test", "testing", "filter", "validate"}
	testFilterVar  = "name"
)

func TestValidateFilters(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.PreCmdConfigTest(t, w, func(cfg *core.PreCommandConfig) {
		viper.Reset()
		viper.Set(constants.ArgQuiet, false)
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		cfg.Command.Command.Flags().StringSlice(constants.FlagFilters, []string{}, "")
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "test=test")
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
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "test")
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
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "name=test,location=test")
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
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "name=test")
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
		cfg.Command.Command.Flags().StringSlice(constants.FlagFilters, []string{}, "")
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "name=test,location=test")
		result, err := GetListQueryParams(cfg)
		assert.NoError(t, err)
		assert.True(t, result.Filters != nil)
		filtersKV := *result.Filters
		assert.True(t, filtersKV["name"][0] == "test")
		assert.True(t, filtersKV["location"][0] == "test")
	})
}

func TestGetListQueryParamsErr(t *testing.T) {
	var b bytes.Buffer
	w := bufio.NewWriter(&b)
	core.CmdConfigTest(t, w, func(cfg *core.CommandConfig, rm *core.ResourcesMocksTest) {
		viper.Reset()
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
		viper.Set(constants.ArgQuiet, false)
		cfg.Command.Command.Flags().Set(constants.FlagFilters, "name")
		_, err := GetListQueryParams(cfg)
		assert.NoError(t, err)
	})
}
