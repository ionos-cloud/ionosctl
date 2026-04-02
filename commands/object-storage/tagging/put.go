package tagging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

const taggingExample = `{
  "TagSet": [
    {"Key": "Environment", "Value": "production"},
    {"Key": "Team", "Value": "platform"}
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "tagging",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the tagging configuration for a bucket",
		LongDesc: "Create or replace the tagging configuration for a bucket. " +
			"The configuration must be provided as a path to a JSON file via --json-properties. " +
			"Use --json-properties-example to see an example tagging configuration.",
		Example: "ionosctl object-storage tagging put --name my-bucket --json-properties tags.json\n" +
			"ionosctl object-storage tagging put --json-properties-example",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), taggingExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading tagging input: %w", err)
			}

			var tagReq objectstorage.PutBucketTaggingRequest
			if err := json.Unmarshal(data, &tagReq); err != nil {
				return fmt.Errorf("parsing tagging JSON: %w", err)
			}

			s3, _, err := client.GetRegionalObjectStorageClient(c.Context, name)
			if err != nil {
				return err
			}

			_, err = s3.TaggingApi.PutBucketTagging(c.Context, name).
				PutBucketTaggingRequest(tagReq).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Tagging configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagName, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.BucketNames(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file containing the tagging configuration")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example tagging configuration JSON and exit")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
