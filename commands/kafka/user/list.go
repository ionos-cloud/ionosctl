package user

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
)

func List() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "list",
			Namespace: "kafka",
			Resource:  "user",
			ShortDesc: "List a cluster's users",
			Aliases:   []string{"l", "ls"},
			Example:   "ionosctl kafka user list",
			PreCmdRun: func(cmd *core.PreCommandConfig) error {
				return core.CheckRequiredFlags(
					cmd.Command, cmd.NS, constants.FlagLocation, constants.FlagClusterId,
				)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)

				usersList, _, err := client.Must().Kafka.UsersApi.ClustersUsersGet(
					context.Background(), clusterID,
				).Execute()
				if err != nil {
					return fmt.Errorf("unable to list users: %s", err)
				}

				cols, _ := cmd.Command.Command.Flags().GetStringSlice(constants.ArgCols)
				out, err := jsontabwriter.GenerateOutput("items", jsonpaths.KafkaUser,
					usersList, tabheaders.GetHeadersAllDefault(allCols, cols))
				if err != nil {
					return err
				}

				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), "%s", out)

				return nil
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, "", "", "The ID of the cluster",
		core.RequiredFlagOption(), core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(read kafka.ClusterRead) string {
						return read.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)
	cmd.AddStringFlag(
		constants.FlagUserId, "", "", "The ID of the user", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.Users(cmd.Command.Flag(constants.FlagClusterId).Value.String())
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.AddBoolFlag("stdout", "", false, "Output the credentials to stdout in a JSON format")

	return cmd
}
