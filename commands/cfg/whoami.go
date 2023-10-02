package cfg

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func WhoamiCmd() *core.Command {

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Verb:      "whoami",
		ShortDesc: "Tells you who you are logged in as. Use `--provenance` to debug where your credentials are being used from",
		LongDesc: fmt.Sprintf(`This command will tell you the email of the user you are logged in as.
You can use '--provenance' flag to see which of these sources are being used. Note that If authentication fails, this flag is set by default.
If using a token, it will use the JWT's claims payload to find out your user UUID, then use the Users API on that UUID to find out your e-mail address.
If no token is present, the command will fall back to using the username and password for authentication.

%s`, constants.DescAuthenticationOrder),
		Example: `ionosctl cfg whoami
ionosctl cfg whoami --provenance`,
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			cl, err := client.Get()

			if fn := core.GetFlagName(c.NS, constants.FlagProvenance); err != nil || viper.GetBool(fn) {
				return handleProvenance(c, cl, err)
			}

			// - - - Below this point, we are sure that client.Get() has returned a valid client object.

			if !cl.IsTokenAuth() {
				// Handle Username & Password Authentication
				_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), cl.CloudClient.GetConfig().Username)
				return err
			}

			// Handle token authentication
			usernameViaToken, jwtParseErr := jwt.Username(cl.CloudClient.GetConfig().Token)
			if jwtParseErr != nil {
				return fmt.Errorf("failed getting username via token: %w", err)
			}
			_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), usernameViaToken)
			return err

		},
		InitClient: false,
	})

	cmd.AddBoolFlag(constants.FlagProvenance, constants.FlagProvenanceShort, false, "If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.")

	return cmd
}

func handleProvenance(c *core.CommandConfig, cl *client.Client, authErr error) error {
	var builder strings.Builder

	if authErr != nil {
		builder.WriteString("Note: Authentication failed!")
		if cl.UsedLayer() == nil {
			builder.WriteString(" None of the authentication layers had a token, or both username & password set.")
		}
		builder.WriteString("\n")
	}

	builder.WriteString("Authentication layers, in order of priority:\n")
	for i, layer := range client.ConfigurationPriorityRules {
		if cl.UsedLayer() != nil && *cl.UsedLayer() == layer {
			builder.WriteString(fmt.Sprintf("* [%d] %s (USED)\n", i+1, layer.Description))
			authType := "token"
			if !cl.IsTokenAuth() {
				authType = "username and password"
			}
			builder.WriteString(fmt.Sprintf("    - Using %s for authentication.\n    - Using %s as the API URL.\n", authType, config.GetServerUrlOrApiIonos()))
		} else {
			builder.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, layer.Description))
		}
	}
	_, err := fmt.Fprintln(c.Command.Command.OutOrStdout(), builder.String())
	return err
}
