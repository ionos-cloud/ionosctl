package console

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func ServerConsoleCmd() *core.Command {
	consoleCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "console",
			Aliases: []string{"url"},
			Short:   "Get a Server's VNC remote console URL",
			Long:    "The sub-command of `ionosctl compute server console` returns the browser-based VNC remote console URL for a Server. The remote console gives you keyboard/screen access to the machine as if you were at its physical monitor — useful before networking/SSH is configured, or to recover a server you cannot reach over the network.",

			TraverseChildren: true,
		},
	}

	consoleCmd.AddCommand(ServerConsoleGetCmd())

	return core.WithConfigOverride(consoleCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
