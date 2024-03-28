package pipeline

import (
	"context"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	ionoscloud "github.com/ionos-cloud/sdk-go-logging"
	"github.com/spf13/viper"
)

var (
	pipelineToCreate      ionoscloud.PipelineCreate
	jsonPropertiesExample = "{\n  \"properties\": {\n    \"name\": \"examplePipe\",\n    \"logs\": [\n      {\n        \"source\": \"docker\",\n        \"tag\": \"tag1\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label1\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      },\n      {\n        \"source\": \"systemd\",\n        \"tag\": \"tag2\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label2\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 14\n          }\n        ]\n      },\n      {\n        \"source\": \"kubernetes\",\n        \"tag\": \"tag3\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label3\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 30\n          }\n        ]\n      },\n      {\n        \"source\": \"generic\",\n        \"tag\": \"tag4\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label4\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      }\n    ]\n  }\n}\n"
)

func PipelineCreateCmd() *core.Command {
	cmd := core.NewCommandWithJsonProperties(
		context.Background(), nil, jsonPropertiesExample, &pipelineToCreate, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "create",
			ShortDesc: "Create a logging pipeline",
			Example: `ionosctl logging-service pipeline create --json-properties PATH_TO_FILE
ionosctl logging-service pipeline create --json-properties-example
ionosctl logging-service pipeline create --name NAME --log-tag LOG_TAG --log-source LOG_SOURCE --log-protocol
LOG_PROTOCOL --log-retention-time LOG_RETENTION_TIMES`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)
	cmd.Command.Flags().StringSlice(constants.ArgCols, defaultCols, tabheaders.ColsMessage(defaultCols))
	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "Sets the name of the pipeline",
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Sets the tag for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogSource, "", "", constants.EnumLogSources,
		"Sets the source for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogProtocol, "", "", constants.EnumLogProtocols,
		"Sets the protocol for the pipeline log",
	)
	cmd.AddStringSliceFlag(constants.FlagLoggingPipelineLogLabels, "", nil, "Sets the labels for the pipeline log")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogType, "", "loki",
		"Sets the destination type for the pipeline log",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogRetentionTime, "", "30", constants.EnumLogRetentionTime,
		"Sets the retention time in days for the pipeline log",
	)

	return cmd
}

func runCreateCmd(c *core.CommandConfig) error {
	if f, err := c.Command.Command.Flags().GetString(constants.FlagJsonProperties); err == nil && f != "" {
		pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPost(context.Background()).Pipeline(
			pipelineToCreate,
		).Execute()
		if err != nil {
			return err
		}

		return handlePipelinePrint(pipeline, c)
	}

	return createFromFlags(c)
}

func preRunCreateCmd(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(
		c.Command, c.NS, []string{constants.FlagJsonProperties}, []string{constants.FlagJsonPropertiesExample},
		[]string{
			constants.FlagName, constants.FlagLoggingPipelineLogTag, constants.FlagLoggingPipelineLogSource,
			constants.FlagLoggingPipelineLogProtocol,
		},
	)
}

func createFromFlags(c *core.CommandConfig) error {
	name := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagName)))
	tag := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogTag))
	source := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogSource)))
	protocol := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogProtocol)))
	labels := viper.GetStringSlice(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogLabels))
	typ := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)))
	retentionTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime))

	retentionTimeInt, err := strconv.ParseInt(retentionTime, 10, 32)
	if err != nil {
		return err
	}

	retentionTimeInt32 := int32(retentionTimeInt)

	dest := ionoscloud.Destination{
		Type:            &typ,
		RetentionInDays: &retentionTimeInt32,
	}

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPost(context.Background()).Pipeline(
		ionoscloud.PipelineCreate{
			Properties: &ionoscloud.PipelineCreateProperties{
				Name: &name,
				Logs: &[]ionoscloud.PipelineCreatePropertiesLogs{
					{
						Tag:          &tag,
						Source:       &source,
						Protocol:     &protocol,
						Labels:       &labels,
						Destinations: &[]ionoscloud.Destination{dest},
					},
				},
			},
		},
	).Execute()
	if err != nil {
		return err
	}

	return handlePipelinePrint(pipeline, c)
}
