package core

import (
	"context"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/printer"
	"github.com/ionos-cloud/ionosctl/pkg/cloudapi-v5"
	"io"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const testConst = "test"

type PreCmdRunTest func(c *PreCommandConfig)

func PreCmdConfigTest(t *testing.T, writer io.Writer, preRunner PreCmdRunTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
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
	CloudApiV5Mocks cloudapi_v5.ResourcesMocks
}

func CmdConfigTest(t *testing.T, writer io.Writer, runner CmdRunnerTest) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	printReg, _ := printer.NewPrinterRegistry(writer, writer)
	prt := printReg[viper.GetString(config.ArgOutput)]
	// Init Test Mock Resources and Services
	testMocks := initMockResources(ctrl)
	cmdConfig := &CommandConfig{
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
		CloudApiV5Mocks: *cloudapi_v5.InitMocksResources(ctrl),
	}
}

// Init Mock Services for Command Test
func initMockServices(c *CommandConfig, tm *ResourcesMocksTest) *CommandConfig {
	c.CloudApiV5Services = *cloudapi_v5.InitMockServices(&c.CloudApiV5Services, &tm.CloudApiV5Mocks)
	return c
}
