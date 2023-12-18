package version

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

func List() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dataplatform",
		Resource:  "version",
		Verb:      "list",
		Aliases:   []string{"ls"},
		ShortDesc: "List all Dataplatform Cluster versions, including deprecated ones. To view the latest version, use the 'version active' command",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			ls, err := VersionsE()
			if err != nil {
				return fmt.Errorf("failed returning versions: %w", err)
			}
			fmt.Fprintf(c.Command.Command.OutOrStdout(), strings.Join(ls, ","))
			return nil
		},
		InitClient: true,
	})

	return cmd
}
