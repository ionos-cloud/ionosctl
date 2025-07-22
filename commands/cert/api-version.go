package cert

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
)

func CertGetApiVersionCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "certmanager",
		Resource:  "certmanager",
		Verb:      "api-version",
		Aliases:   []string{"api", "info"},
		ShortDesc: "[DEPRECATED] Get Certificate Manager API Version",
		LongDesc: "[DEPRECATED] This command is deprecated and will be removed in the future" +
			"It will simply output v2.0 and exit, as this API route has been removed.",
		Example:    "ionosctl certificate-manager api-version",
		PreCmdRun:  core.NoPreRun,
		CmdRun:     CmdGetApiVersion,
		InitClient: true,
	})
	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, "[DEPRECATED]")

	cmd.Command.Hidden = true
	cmd.Command.Deprecated = "This command is deprecated and will be removed in the future. It will simply output v2.0 and exit, as this API route has been removed."

	return cmd
}
func CmdGetApiVersion(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Getting Api Version..."))

	// v2.0 API version removed this route. This is a dummy response. to avoid errs and breaking changes.
	fmt.Fprintf(c.Command.Command.OutOrStdout(), "v2.0")
	return nil
}
