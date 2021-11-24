package core

import (
	"context"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	authv1 "github.com/ionos-cloud/ionosctl/services/auth-v1"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
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
	if viper.GetString(config.ArgOutput) == "" {
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	}
	p, _ := printer.NewPrinterRegistry(writer, writer)
	prt := p[viper.GetString(config.ArgOutput)]
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
		Printer:   prt,
	}
	preRunner(preCmdCfg)
}

type CmdRunnerTest func(c *CommandConfig, mocks *ResourcesMocksTest)

type ResourcesMocksTest struct {
	// Add New Services Resources Mocks
	CloudApiV6Mocks cloudapi_v6.ResourcesMocks
	AuthV1Mocks     authv1.ResourcesMocks
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	if viper.GetString(config.ArgOutput) == "" {
		viper.Set(config.ArgOutput, config.DefaultOutputFormat)
	}
	printReg, _ := printer.NewPrinterRegistry(writer, writer)
	prt := printReg[viper.GetString(config.ArgOutput)]
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
		Printer:   prt,
		Context:   context.TODO(),
		initCfg:   func(c *CommandConfig) error { return nil },
	}
	cmdConfig = initMockServices(cmdConfig, testMocks)
	runner(cmdConfig, testMocks)
}

// Init Mock Resources for Test
func initMockResources(ctrl *gomock.Controller) *ResourcesMocksTest {
	return &ResourcesMocksTest{
		CloudApiV6Mocks: *cloudapi_v6.InitMocksResources(ctrl),
		AuthV1Mocks:     *authv1.InitMocksResources(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocksTest) *CommandConfig {
	c.CloudApiV6Services = *cloudapi_v6.InitMockServices(&c.CloudApiV6Services, &tm.CloudApiV6Mocks)
	c.AuthV1Services = *authv1.InitMockServices(&c.AuthV1Services, &tm.AuthV1Mocks)
	return c
}
