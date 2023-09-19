package core

import (
	"context"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	authservice "github.com/ionos-cloud/ionosctl/v6/services/auth-v1"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	cloudapidbaaspgsql "github.com/ionos-cloud/ionosctl/v6/services/dbaas-postgres"
	"github.com/spf13/cobra"
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
	CloudApiV6Mocks         cloudapiv6.ResourcesMocks
	CloudApiDbaasPgsqlMocks cloudapidbaaspgsql.ResourcesMocks
	AuthV1Mocks             authservice.ResourcesMocks
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
			for _, argPair := range tc.Args {
				viper.Set(argPair.Flag, argPair.Value)
			}

			if tc.UserInput != nil {
				cfg.Stdin = tc.UserInput
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
		CloudApiV6Mocks:         *cloudapiv6.InitMocksResources(ctrl),
		AuthV1Mocks:             *authservice.InitMocksResources(ctrl),
		CloudApiDbaasPgsqlMocks: *cloudapidbaaspgsql.InitMocksResources(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocksTest) *CommandConfig {
	c.CloudApiV6Services = *cloudapiv6.InitMockServices(&c.CloudApiV6Services, &tm.CloudApiV6Mocks)
	c.AuthV1Services = *authservice.InitMockServices(&c.AuthV1Services, &tm.AuthV1Mocks)
	c.CloudApiDbaasPgsqlServices = *cloudapidbaaspgsql.InitMockServices(&c.CloudApiDbaasPgsqlServices, &tm.CloudApiDbaasPgsqlMocks)
	return c
}
