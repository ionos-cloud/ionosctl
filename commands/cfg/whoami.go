package cfg

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func WhoamiCmd() *core.Command {

	const (
		FlagProvenance      = "provenance"
		FlagProvenanceShort = "p"
	)

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "whoami",
		ShortDesc: "Tells you who you are logged in as. Use `--provenance` to debug where your credentials are being used from",
		LongDesc: fmt.Sprintf(`This command will tell you the email of the user you are logged in as.
You can use '--provenance' flag to see which of these sources are being used.
If using a token, it will use the JWT's claims payload to find out your user UUID, then use the Users API on that UUID to find out your e-mail address.
If no token is present, the command will fall back to using the username and password for authentication.

%s`, constants.DescAuthenticationOrder),
		Example: `ionosctl cfg whoami
ionosctl cfg whoami --provenance`,
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			cl, err := client.Get()

			// Does user want to see provenance of his configuration? i.e. where does each key get its value from.
			if fn := core.GetFlagName(c.NS, FlagProvenance); viper.GetBool(fn) {
				out := "Authentication layers, in order of priority:\n"
				for i, layer := range client.ConfigurationPriorityRules {
					if layer == cl.UsedLayer {
						out += fmt.Sprintf("* [%d] %s (USED)\n", i+1, layer.Help)
						if cl.IsTokenAuth() {
							out += "    - Using token for authentication.\n"
						} else {
							out += "    - Using username and password for authentication.\n"
						}
					} else {
						out += fmt.Sprintf("  [%d] %s\n", i+1, layer.Help)
					}
				}
				_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), out)
				return err
			}

			// Get Client error. Intentionally handled after provenance rule prints
			if err != nil {
				return err
			}

			// -- Below this point, we are 100% certain the client is using valid credentials. --

			token := cl.CloudClient.GetConfig().Token
			if jwt.Valid(token) {
				// Valid token
				usernameViaToken, err := jwt.Username(token)
				if err != nil {
					return fmt.Errorf("failed getting username via token: %w", err)
				}
				_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), usernameViaToken)
				return err
			}

			// -- Below this point, we are 100% certain the client is using valid username & password. --

			_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), cl.CloudClient.GetConfig().Username)
			return err

		},
		InitClient: false,
	})

	cmd.AddBoolFlag(FlagProvenance, FlagProvenanceShort, false, "If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.")

	return cmd
}
