package groups

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func GroupCreateCmd() *core.Command {

	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "user",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create MongoDB users.",
		Example: fmt.Sprintf("ionosctl autoscaling group create %s",
			core.FlagsUsage(constants.FlagDatacenterId, constants.FlagName)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired(constants.FlagClusterId)
			return err
		},
		CmdRun: func(c *core.CommandConfig) error {
			return nil
		},
	})

	return cmd
}
