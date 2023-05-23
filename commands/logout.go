package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"golang.org/x/exp/slices"
)

func LogoutCmd() *core.Command {
	loginCmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "logout",
		Resource:  "logout",
		Verb:      "logout",
		ShortDesc: "Convenience command for deletion of config file credentials",
		Example:   "ionosctl logout",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			data, err := config.ReadFile()
			if err != nil {
				return fmt.Errorf("logout intrerrupted: %w", err)
			}

			// Go through data struct and blank out all credentials, including old ones (userdata.username, userdata.password)
			msg := "De-authentication successful. Blanked out the following fields in your config file:\n"
			for k, _ := range data {
				if slices.Contains(config.FieldsWithSensitiveDataInConfigFile, k) {
					data[k] = ""
					msg += fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata."))
				}
			}

			err = config.WriteFile(data)
			if err != nil {
				return fmt.Errorf("failed updating config file: %w", err)
			}
			return c.Printer.Print(msg)
		},
		InitClient: false,
	})

	return loginCmd
}
