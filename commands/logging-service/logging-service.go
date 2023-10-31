package logging_service

import (
	"github.com/ionos-cloud/ionosctl/v6/commands/logging-service/pipeline"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/die"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoggingServiceCmd() *core.Command {
	cmd := &core.Command{
		Command: &cobra.Command{
			Use:              "logging-service",
			Aliases:          []string{"log-svc"},
			Short:            "LaaS Operations. Manage and centralize your application/infrastructure's logs",
			TraverseChildren: true,
			PersistentPreRun: func(cmd *cobra.Command, args []string) {
				checkForToken()
			},
		},
	}

	cmd.AddCommand(pipeline.PipelineCmd())
	return cmd
}

func checkForToken() {
	for _, layer := range client.ConfigurationPriorityRules {
		token := viper.GetString(layer.TokenKey)

		if token != "" {
			return
		}
	}

	die.Die("authorization failed, use a token for authentication or use the `--token` flag\n")
}
