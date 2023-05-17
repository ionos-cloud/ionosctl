package commands

import (
	"context"
	"fmt"

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
		Example:   loginExamples,
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			data, err := config.ReadFile()
			if err != nil {
				return fmt.Errorf("logout intrerrupted: %w", err)
			}

			// Go through data struct and blank out all credentials, including old ones (userdata.username, userdata.password)
			modificationsExist := false
			for k, _ := range data {
				if slices.Contains(config.FieldsWithSensitiveDataInConfigFile, k) {
					data[k] = ""
					modificationsExist = true
				}
			}

			wereModifsOut := "no credentials were found in the config file"
			if modificationsExist {
				wereModifsOut = "some credentials were found in the config file"
			}

			err = config.WriteFile(data)
			if err != nil {
				return fmt.Errorf("failed updating config file, %s: %w", wereModifsOut, err)
			}
			return nil
		},
		InitClient: false,
	})

	return loginCmd
}
