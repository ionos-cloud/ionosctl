package core

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const testConst = "test"

var (
	RootCmdTest = Command{
		Command: &cobra.Command{
			Use: testConst,
		},
	}
)

// SetFlag registers (if absent) and sets a local flag value on the test command,
// mirroring how production binds flags via AddXFlag. It is the test-side counterpart
// to the c.Flags() typed getters: it registers the flag with a type matching the
// value, so the typed getter the command-under-test calls won't panic. Use this in
// tests instead of the old viper.Set(core.GetFlagName(cfg.NS, name), value).
func (c *CommandConfig) SetFlag(name string, value any) {
	setTestFlag(c.Command.Command.Flags(), name, value)
}

// SetFlag registers (if absent) and sets a local flag value on the test command.
// See CommandConfig.SetFlag.
func (c *PreCommandConfig) SetFlag(name string, value any) {
	setTestFlag(c.Command.Command.Flags(), name, value)
}

func setTestFlag(fs *pflag.FlagSet, name string, value any) {
	switch v := value.(type) {
	case string:
		if fs.Lookup(name) == nil {
			fs.String(name, "", "")
		}
		_ = fs.Set(name, v)
	case bool:
		if fs.Lookup(name) == nil {
			fs.Bool(name, false, "")
		}
		_ = fs.Set(name, strconv.FormatBool(v))
	case int:
		if fs.Lookup(name) == nil {
			fs.Int(name, 0, "")
		}
		_ = fs.Set(name, strconv.Itoa(v))
	case int32:
		if fs.Lookup(name) == nil {
			fs.Int32(name, 0, "")
		}
		_ = fs.Set(name, strconv.FormatInt(int64(v), 10))
	case []string:
		if fs.Lookup(name) == nil {
			fs.StringSlice(name, nil, "")
		}
		_ = fs.Set(name, strings.Join(v, ","))
	case map[string]string:
		if fs.Lookup(name) == nil {
			fs.StringToString(name, nil, "")
		}
		for k, val := range v {
			_ = fs.Set(name, fmt.Sprintf("%s=%s", k, val))
		}
	default:
		panic(fmt.Sprintf("SetFlag: unsupported value type %T for flag %q", value, name))
	}
}

type PreCmdRunTest func(c *PreCommandConfig)

func PreCmdConfigTest(t *testing.T, writer io.Writer, preRunner PreCmdRunTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	if viper.GetString(constants.ArgOutput) == "" {
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
	}
	preCmdCfg := &PreCommandConfig{
		Command: &Command{
			Command: &cobra.Command{
				Use: testConst,
			},
		},
		NS:        testConst,
		Namespace: testConst,
		Resource:  testConst,
		Verb:      testConst,
	}

	preCmdCfg.Command.Command.SetOut(writer)
	preRunner(preCmdCfg)
}

type CmdRunnerTest func(c *CommandConfig, mocks *ResourcesMocksTest)

type ResourcesMocksTest struct {
	// Add New Services Resources Mocks
	CloudApiV6Mocks cloudapiv6.ResourcesMocks
}

type FlagValuePair struct {
	Flag  string
	Value interface{}
}

type TestCase struct {
	Name        string
	UserInput   io.Reader
	Args        []FlagValuePair
	Calls       func(...*gomock.Call)
	ExpectedErr bool // To be replaced by `error` type once it makes sense to do so (currently only one type of error is thrown)
}

func ExecuteTestCases(t *testing.T, funcToTest func(c *CommandConfig) error, testCases []TestCase, cfg *CommandConfig) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			viper.Reset()
			viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
			for _, argPair := range tc.Args {
				viper.Set(argPair.Flag, argPair.Value)
			}

			if tc.UserInput != nil {
				cfg.Command.Command.SetIn(tc.UserInput)
			}

			// Expected gomock calls, call order, call counts and returned values
			tc.Calls()

			err := funcToTest(cfg)

			if tc.ExpectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	if viper.GetString(constants.ArgOutput) == "" {
		viper.Set(constants.ArgOutput, constants.DefaultOutputFormat)
	}
	// Init Test Mock Resources and Services
	testMocks := initMockResources(ctrl)
	cmdConfig := &CommandConfig{
		Command: &Command{
			Command: &cobra.Command{
				Use: testConst,
			},
		},
		NS:        testConst,
		Namespace: testConst,
		Resource:  testConst,
		Verb:      testConst,
		Context:   context.TODO(),
		initCfg:   func(c *CommandConfig) error { return nil },
	}

	cmdConfig.Command.Command.SetOut(writer)
	cmdConfig = initMockServices(cmdConfig, testMocks)
	runner(cmdConfig, testMocks)
}

// Init Mock Resources for Test
func initMockResources(ctrl *gomock.Controller) *ResourcesMocksTest {
	return &ResourcesMocksTest{
		CloudApiV6Mocks: *cloudapiv6.InitMocksResources(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocksTest) *CommandConfig {
	c.CloudApiV6Services = *cloudapiv6.InitMockServices(&c.CloudApiV6Services, &tm.CloudApiV6Mocks)
	return c
}
