package token

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func ServerTokenCmd() *core.Command {
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "token",
			Aliases:          []string{"t"},
			Short:            "Server Token Operations",
			Long:             "The sub-command of `ionosctl compute server token` allows you to get Token for specific Server.",
			TraverseChildren: true,
		},
	}

	tokenCmd.AddCommand(ServerTokenGetCmd())

	return core.WithConfigOverride(tokenCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
