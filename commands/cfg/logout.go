package cfg

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func LogoutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "logout",
		ShortDesc: "Convenience command for removing config file credentials",
		LongDesc: fmt.Sprintf(`This command is a 'Quality of Life' command which will parse your config file for fields that contain sensitive data.
If any such fields are found, their values will be replaced with an empty string.

%s`, constants.DescAuthenticationOrder),
		Example:   "ionosctl logout",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Println("todo")
			return nil
		},
		InitClient: false,
	})

	return cmd
}
