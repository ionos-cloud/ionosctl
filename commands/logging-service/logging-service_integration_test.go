//go:build integration
// +build integration

package logging_service_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/logs"
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	temporaryToken = ""
	pipelineId     = ""
	logTag         = "testlog"
)

func TestLoggingServiceCmd(t *testing.T) {
	if err := setup(); err != nil {
		assert.FailNow(t, err.Error())
	}
	defer teardown()

	t.Run("test pipeline commands", testPipeline)
	t.Run("test logs commands", testLogs)
}

func testPipeline(t *testing.T) {
	outBuff := bytes.NewBuffer([]byte{})
	viper.Set(constants.ArgOutput, jsontabwriter.TextFormat)

	t.Run(
		"test pipeline create", func(t *testing.T) {
			cmd := pipeline.PipelineCreateCmd()
			cmd.Command.SetOut(outBuff)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagName), "testpipe")
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogTag), logTag)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogSource), "docker")
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogProtocol), "http")

			err := cmd.Command.Execute()
			assert.NoError(t, err)

			uuidRegex, err := regexp.Compile("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}")
			assert.NoError(t, err)

			pipelineId = uuidRegex.FindStringSubmatch(outBuff.String())[0]
		},
	)

	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.TextFormat)
	viper.Set(constants.ArgQuiet, true)

	t.Run(
		"test pipeline list", func(t *testing.T) {
			cmd := pipeline.PipelineListCmd()
			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"test pipeline get", func(t *testing.T) {
			cmd := pipeline.PipelineGetCmd()
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	viper.Reset()
}

func testLogs(t *testing.T) {
	viper.Set(constants.ArgOutput, jsontabwriter.TextFormat)
	viper.Set(constants.ArgQuiet, true)

	// this wastes a lot of time, but pipelines take around 5 minutes to provision (at least at the moment)
	time.Sleep(300 * time.Second)

	t.Run(
		"test logs add", func(t *testing.T) {
			cmd := logs.LogsAddCmd()
			fmt.Println(pipelineId)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogTag), "new"+logTag)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogSource), "kubernetes")
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogProtocol), "http")

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	viper.Reset()
	viper.Set(constants.ArgOutput, jsontabwriter.TextFormat)
	viper.Set(constants.ArgQuiet, true)

	t.Run(
		"test logs list", func(t *testing.T) {
			cmd := logs.LogsListCmd()
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"test logs get", func(t *testing.T) {
			cmd := logs.LogsGetCmd()
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogTag), logTag)

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"test logs remove", func(t *testing.T) {
			cmd := logs.LogsRemoveCmd()
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogTag), logTag)
			viper.Set(constants.ArgForce, true)

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	t.Run(
		"test logs update", func(t *testing.T) {
			cmd := logs.LogsUpdateCmd()
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogTag), "new"+logTag)
			viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineLogSource), "systemd")

			err := cmd.Command.Execute()
			assert.NoError(t, err)
		},
	)

	viper.Reset()
}

func setup() error {
	if os.Getenv("IONOS_TOKEN") != "" {
		return nil
	}

	username := os.Getenv("IONOS_USERNAME")
	password := os.Getenv("IONOS_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("empty user/password and no token set")
	}

	tok, _, err := client.Must().AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return err
	}

	temporaryToken = *tok.GetToken()

	err = os.Setenv("IONOS_TOKEN", temporaryToken)
	if err != nil {
		return nil
	}

	return nil
}

func teardown() {
	viper.Set(constants.ArgOutput, jsontabwriter.TextFormat)
	viper.Set(constants.ArgQuiet, true)
	viper.Set(constants.ArgForce, true)

	cmd := pipeline.PipelineDeleteCmd()
	viper.Set(core.GetFlagName(cmd.NS, constants.FlagLoggingPipelineId), pipelineId)

	err := cmd.Command.Execute()
	if err != nil {
		fmt.Printf("failed to clean up logging pipeline (%v): %v", pipelineId, err)
	}

	if temporaryToken == "" {
		return
	}

	headers, err := jwt.Headers(temporaryToken)
	if err != nil {
		log.Fatal(err.Error())
	}

	tokenId, err := jwt.Kid(headers)
	if err != nil {
		log.Fatal(err.Error())
	}

	_, _, err = client.Must().AuthClient.TokensApi.TokensDeleteById(context.Background(), tokenId).Execute()
	if err != nil {
		fmt.Printf("failed to delete temporary token: %v", err)
	}
}
