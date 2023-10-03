package cfg

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"golang.org/x/exp/slices"
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
			data, err := config.Read()
			if err != nil {
				return fmt.Errorf("logout intrerrupted: %w", err)
			}

			printNumberOfTokens := true
			ls, _, err := client.Must(func(_ error) {
				// If some error in creating the client, don't fail the command, but disable client-related functionality
				printNumberOfTokens = false
			}).AuthClient.TokensApi.TokensGet(context.Background()).Execute()
			if err != nil {
				return err
			}

			msg := "De-authentication successful."
			if printNumberOfTokens {
				msg += fmt.Sprintf(" Note: Your account has %d active tokens.", len(*ls.Tokens))
			}
			msg += " Affected fields:\n"

			for k, _ := range data {
				// Go through data struct and blank out all credentials, including old ones (userdata.username, userdata.password)
				if slices.Contains(config.FieldsWithSensitiveDataInConfigFile, k) {
					data[k] = ""
					msg += fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata."))
				}
			}

			err = config.Write(data)
			if err != nil {
				return fmt.Errorf("failed updating config file: %w", err)
			}
			_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), msg)
			return err
		},
		InitClient: false,
	})

	return cmd
}
