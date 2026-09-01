package pipeline

import (
	"context"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/logging/v2"
	"github.com/spf13/viper"
)

var (
	pipelineToCreate      logging.PipelineCreate
	jsonPropertiesExample = "{\n  \"properties\": {\n    \"name\": \"examplePipe\",\n    \"logs\": [\n      {\n        \"source\": \"docker\",\n        \"tag\": \"tag1\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label1\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      },\n      {\n        \"source\": \"systemd\",\n        \"tag\": \"tag2\",\n        \"protocol\": \"http\",\n        \"labels\": [\n          \"label2\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 14\n          }\n        ]\n      },\n      {\n        \"source\": \"kubernetes\",\n        \"tag\": \"tag3\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label3\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 30\n          }\n        ]\n      },\n      {\n        \"source\": \"generic\",\n        \"tag\": \"tag4\",\n        \"protocol\": \"tcp\",\n        \"labels\": [\n          \"label4\"\n        ],\n        \"destinations\": [\n          {\n            \"type\": \"loki\",\n            \"retentionInDays\": 7\n          }\n        ]\n      }\n    ]\n  }\n}\n"
)

func PipelineCreateCmd() *core.Command {
	cmd := core.NewCommandWithJsonProperties(
		context.Background(), nil, jsonPropertiesExample, &pipelineToCreate, core.CommandBuilder{
			Namespace: "logging-service",
			Resource:  "pipeline",
			Verb:      "create",
			ShortDesc: "Create a logging pipeline",
			LongDesc: `Create a logging pipeline together with its first log stream. On success the pipeline's ingestion addresses (TCP/HTTP) and Grafana address are returned; generate a key with 'pipeline key' before shipping logs.

There are two ways to define the pipeline:
  1. Flags: --name plus a single log described by --log-tag, --log-source, --log-protocol (and optionally --log-type, --log-retention-time, --log-labels). This creates a pipeline with exactly one log.
  2. --json-properties: a path to a JSON file describing the full pipeline, which lets you define MULTIPLE logs at once. Use --json-properties-example to print a ready-to-edit template.

A pipeline must contain at least one log, which is why source, tag and protocol are required in flag mode.

Valid values:
  --log-source:   docker, systemd, kubernetes, generic (the kind of workload the log comes from)
  --log-protocol: http or tcp (how your shipper connects to the ingestion endpoint)
  --log-type:     the destination backend; currently loki
  --log-retention-time: 7, 14 or 30 days

Note: --name and --log-source are normalised to lower-case before being sent.`,
			Example: `# Create a pipeline with a single Kubernetes log stream shipped over TCP
ionosctl logging-service pipeline create --location de/txl --name my-pipeline --log-tag k8s --log-source kubernetes --log-protocol tcp

# Advanced: pin retention and attach labels, shipping Docker logs over HTTP
ionosctl logging-service pipeline create --location de/txl --name app-logs --log-tag docker-app --log-source docker --log-protocol http --log-retention-time 30 --log-labels env=prod,team=core

# Create from a JSON file (allows multiple log streams in one pipeline)
ionosctl logging-service pipeline create --location de/txl --json-properties ./pipeline.json

# Print a JSON template you can edit and feed back via --json-properties
ionosctl logging-service pipeline create --json-properties-example`,
			PreCmdRun: preRunCreateCmd,
			CmdRun:    runCreateCmd,
		},
	)
	cmd.AddStringFlag(
		constants.FlagName, constants.FlagNameShort, "", "Human-readable name of the pipeline, shown in listings and Grafana. Normalised to lower-case",
	)
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogTag, "", "", "Tag that identifies this log stream within the pipeline; used later to reference the log (e.g. in 'logs get/update/remove') and to route entries",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogSource, "", "", constants.EnumLogSources,
		"The kind of workload producing the logs, which selects how they are parsed. One of: docker, systemd, kubernetes, generic (use 'generic' for anything that does not fit the others)",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogProtocol, "", "", constants.EnumLogProtocols,
		"Transport your shipper uses to push logs to the ingestion endpoint: 'http' targets the pipeline's HTTP address, 'tcp' its TCP address",
	)
	cmd.AddStringSliceFlag(constants.FlagLoggingPipelineLogLabels, "", nil, "Optional labels attached to every entry of this log stream, useful for filtering in Grafana. Comma-separated, e.g. --log-labels env=prod,team=core")
	cmd.AddStringFlag(
		constants.FlagLoggingPipelineLogType, "", "loki",
		"Destination backend the logs are stored in and queried from. Currently 'loki'",
	)
	cmd.AddSetFlag(
		constants.FlagLoggingPipelineLogRetentionTime, "", "30", constants.EnumLogRetentionTime,
		"How many days logs are kept before being deleted. One of: 7, 14, 30",
	)

	return cmd
}

func runCreateCmd(c *core.CommandConfig) error {
	if f, err := c.Command.Command.Flags().GetString(constants.FlagJsonProperties); err == nil && f != "" {
		pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPost(context.Background()).
			PipelineCreate(
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
	return c.CheckRequiredFlagsSetsAndLocation(
		[]string{constants.FlagJsonProperties}, []string{constants.FlagJsonPropertiesExample},
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
	typ := strings.ToLower(viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogType)))
	retentionTime := viper.GetString(core.GetFlagName(c.NS, constants.FlagLoggingPipelineLogRetentionTime))

	retentionTimeInt, err := strconv.ParseInt(retentionTime, 10, 32)
	if err != nil {
		return err
	}

	retentionTimeInt32 := int32(retentionTimeInt)

	dest := logging.PipelineNoAddrLogsDestinations{
		Type:            typ,
		RetentionInDays: retentionTimeInt32,
	}

	pipeline, _, err := client.Must().LoggingServiceClient.PipelinesApi.PipelinesPost(context.Background()).
		PipelineCreate(
			logging.PipelineCreate{
				Properties: logging.PipelineNoAddr{
					Name: name,
					Logs: []logging.PipelineNoAddrLogs{
						{
							Tag:          tag,
							Source:       source,
							Protocol:     protocol,
							Destinations: []logging.PipelineNoAddrLogsDestinations{dest},
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
