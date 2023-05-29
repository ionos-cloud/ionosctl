package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"golang.org/x/exp/slices"
)

func LogoutCmd() *core.Command {
	loginCmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "logout",
		Resource:  "logout",
		Verb:      "logout",
		ShortDesc: "Convenience command for deletion of config file credentials. To also remove your account's active tokens, use `ionosctl token delete --all`",
		Example:   "ionosctl logout",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			data, err := config.Read()
			if err != nil {
				return fmt.Errorf("logout intrerrupted: %w", err)
			}

			ls, _, err := client.Must().AuthClient.TokensApi.TokensGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			// Go through data struct and blank out all credentials, including old ones (userdata.username, userdata.password)
			msg := fmt.Sprintf("De-authentication successful. Note: Your account has %d active tokens. Affected fields:\n", len(*ls.Tokens))
			for k, _ := range data {
				if slices.Contains(config.FieldsWithSensitiveDataInConfigFile, k) {
					data[k] = ""
					msg += fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata."))
				}
			}

			err = config.Write(data)
			if err != nil {
				return fmt.Errorf("failed updating config file: %w", err)
			}
			return c.Printer.Print(msg)
		},
		InitClient: false,
	})

	return loginCmd
}
