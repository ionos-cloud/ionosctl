package tagging

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
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
		LongDesc: `Create or replace the bucket's tag set. This REPLACES all existing tags with the ones in the file - it is not a merge - so include every tag you want the bucket to keep.

Provide the tags as a JSON file via --json-properties: a top-level object with a "TagSet" array of {"Key": ..., "Value": ...} pairs. Tags are commonly used for cost allocation and grouping resources by environment, team or project.

Run with --json-properties-example to print a ready-to-edit template.`,
		Example: `# Apply a set of tags from a file
ionosctl object-storage bucket tagging put --name my-bucket --json-properties tags.json

# Print an example tag set
ionosctl object-storage bucket tagging put --json-properties-example`,
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

			_, err = client.MustObjectStorage().ObjectStorageClient.TaggingApi.PutBucketTagging(c.Context, name).
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the tag set ({\"TagSet\":[{\"Key\":...,\"Value\":...}]}). Replaces all existing tags")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example tagging configuration JSON and exit without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
