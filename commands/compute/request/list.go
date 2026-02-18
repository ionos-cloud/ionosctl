package request

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func RequestListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "request",
		Resource:   "request",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Requests",
		LongDesc:   "Use this command to list all Requests on your account.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.RequestsFiltersUsage(),
		Example:    `ionosctl request list --latest N`,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunRequestList,
		InitClient: true,
	})
	cmd.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show latest N Requests. If it is not set, all Requests will be printed", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	cmd.AddStringFlag(cloudapiv6.ArgMethod, "", "", "Show only the Requests with this method. E.g CREATE, UPDATE, DELETE", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgMethod, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"POST", "PUT", "DELETE", "PATCH", "CREATE", "UPDATE"}, cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
