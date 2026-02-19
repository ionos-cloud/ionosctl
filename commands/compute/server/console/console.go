package console

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func ServerConsoleCmd() *core.Command {
	consoleCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "console",
			Aliases:          []string{"url"},
			Short:            "Server Remote Console URL Operations",
			Long:             "The sub-command of `ionosctl compute server console` allows you to get the URL for Remote Console of a specific Server.",
			TraverseChildren: true,
		},
	}

	consoleCmd.AddCommand(ServerConsoleGetCmd())

	return core.WithConfigOverride(consoleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
