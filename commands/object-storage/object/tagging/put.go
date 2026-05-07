package objecttagging

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

const objectTaggingExample = `{
  "TagSet": [
    {"Key": "Environment", "Value": "production"},
    {"Key": "Team", "Value": "platform"}
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "object-tagging",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the tagging configuration for an object",
		LongDesc: "Create or replace the tagging configuration for an object. " +
			"The configuration must be provided as a path to a JSON file via --json-properties. " +
			"Use --json-properties-example to see an example tagging configuration.",
		Example: "ionosctl object-storage object tagging put --name my-bucket --key my-object --json-properties tags.json\n" +
			"ionosctl object-storage object tagging put --json-properties-example",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, flagKey, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), objectTaggingExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			key := viper.GetString(core.GetFlagName(c.NS, flagKey))
			versionId := viper.GetString(core.GetFlagName(c.NS, flagVersionId))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading object tagging input: %w", err)
			}

			var tagReq objectstorage.PutObjectTaggingRequest
			if err := json.Unmarshal(data, &tagReq); err != nil {
				return fmt.Errorf("parsing object tagging JSON: %w", err)
			}

			req := client.MustObjectStorage().ObjectStorageClient.TaggingApi.PutObjectTagging(c.Context, name, key).
				PutObjectTaggingRequest(tagReq)
			if versionId != "" {
				req = req.VersionId(versionId)
			}

			_, _, err = req.Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Tagging configuration for object %q in bucket %q applied successfully\n", key, name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Object key", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Version ID of the object")
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file containing the tagging configuration")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example tagging configuration JSON and exit")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
