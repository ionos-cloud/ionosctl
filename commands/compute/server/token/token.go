package token

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

func ServerTokenCmd() *core.Command {
	tokenCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "token",
			Aliases: []string{"t"},
			Short:   "Get a Server's JWT access token",
			Long:    "The sub-command of `ionosctl compute server token` returns the JSON Web Token (JWT) associated with a Server. This token authorizes direct access to that server's endpoints — for example the VNC remote console — so treat it as a secret and do not share or log it.",

			TraverseChildren: true,
		},
	}

	tokenCmd.AddCommand(ServerTokenGetCmd())

	return core.WithConfigOverride(tokenCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
