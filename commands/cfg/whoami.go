package cfg

import (
	"context"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/jwt"
)

func WhoamiCmd() *core.Command {

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "whoami",
		Resource:  "whoami",
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
			cl, authErr := client.Get()

			if showProv, _ := c.Command.Command.Flags().GetBool(constants.FlagProvenance); authErr != nil || showProv {
				return handleProvenance(c, cl, authErr)
			}

			// - - - Below this point, we are sure that client.Get() has returned a valid client object.

			if cl.AuthSource == client.AuthSourceCfgBasic || cl.AuthSource == client.AuthSourceEnvBasic {
				// Handle Username & Password Authentication
				_, err := fmt.Fprintln(c.Command.Command.OutOrStdout(), cl.CloudClient.GetConfig().Username)
				return err
			}

			// Handle token authentication
			usernameViaToken, jwtParseErr := jwt.Username(cl.CloudClient.GetConfig().Token)
			if jwtParseErr != nil {
				return fmt.Errorf("failed getting username via token: %w", jwtParseErr)
			}
			fmt.Fprintln(c.Command.Command.OutOrStdout(), usernameViaToken)
			return nil
		},
		InitClient: false,
	})

	cmd.AddBoolFlag(constants.FlagProvenance, constants.FlagProvenanceShort, false, "If set, the command prints the layers of authentication sources, their order of priority, and which one was used. It also tells you if a token or username and password are being used for authentication.")

	return cmd
}

// handleProvenance prints out all authentication layers in priority order,
// marks which one was actually used, and shows whether it’s token vs. user/pass
// plus the effective API URL.
func handleProvenance(c *core.CommandConfig, cl *client.Client, authErr error) error {
	var b strings.Builder

	// If auth itself failed, note it
	if authErr != nil {
		b.WriteString("Note: Authentication failed!\n")
	}

	// List all possible sources in priority order
	order := []client.AuthSource{
		client.AuthSourceEnvBearer,
		client.AuthSourceEnvBasic,
		client.AuthSourceCfgBearer,
		client.AuthSourceCfgBasic,
	}

	b.WriteString("Authentication layers, in order of priority:\n")
	for i, src := range order {
		// highlight the one actually used
		if cl.AuthSource == src {
			b.WriteString(fmt.Sprintf("* [%d] %s (USED)\n", i+1, src))
		} else {
			b.WriteString(fmt.Sprintf("  [%d] %s\n", i+1, src))
		}
	}

	// Finally, print it all out
	_, err := fmt.Fprintln(c.Command.Command.OutOrStdout(), b.String())
	return err
}
