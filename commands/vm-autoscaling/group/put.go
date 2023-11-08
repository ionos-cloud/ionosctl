package group

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	vmasc "github.com/ionos-cloud/sdk-go-vm-autoscaling"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Put() *core.Command {
	var jsonStruct vmasc.GroupPut
	cmd := core.NewCommandWithJsonProperties(context.Background(), nil, exampleJson, &jsonStruct, core.CommandBuilder{
		Namespace: "vm-autoscaling",
		Resource:  "groups",
		Verb:      "put",
		Aliases:   []string{"p", "update"},
		ShortDesc: "Perform a PUT operation to modify an existing group",
		Example: fmt.Sprintf("ionosctl vm-autoscaling group put %s",
			core.FlagsUsage(constants.FlagJsonProperties)),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return core.CheckRequiredFlagsSets(c.Command, c.NS,
				[]string{constants.FlagGroupId, constants.FlagJsonProperties},
				[]string{constants.FlagGroupId, constants.FlagJsonPropertiesExample},
			)
		},
		CmdRun: func(c *core.CommandConfig) error {
			group, _, err := client.Must().VMAscClient.GroupsPut(context.Background(),
				viper.GetString(core.GetFlagName(c.NS, constants.FlagGroupId))).
				GroupPut(jsonStruct).Execute()
			if err != nil {
				return err
			}

			colsDesired := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))
			out, err := jsontabwriter.GenerateOutput("", allJSONPaths, group,
				tabheaders.GetHeaders(allCols, defaultCols, colsDesired))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

			return nil
		},
	})

	cmd.AddStringFlag(constants.FlagGroupId, constants.FlagIdShort, "", "ID of the autoscaling group to modify", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		// get ID of all groups
		return GroupsProperty(func(r vmasc.GroupResource) string {
			return fmt.Sprintf(*r.Id) // + "\t" + *r.Properties.Name) // Commented because this SDK functionality currently broken
		}), cobra.ShellCompDirectiveNoFileComp
	})

	return cmd
}
