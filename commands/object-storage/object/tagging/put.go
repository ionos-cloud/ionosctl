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
		ShortDesc: "Create or replace the entire tag set on an object",
		LongDesc: "Create or replace the tag set on an object.\n\n" +
			"This REPLACES the object's whole tag set (it is not a merge): any existing tags not present in the supplied file are removed. To keep existing tags, include them in the file. An object may hold at most 10 tags.\n\n" +
			"The tag set is supplied as a JSON file via --json-properties, containing a \"TagSet\" array of {\"Key\", \"Value\"} objects. Run the command with --json-properties-example to print a ready-to-edit template.",
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket holding the object", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagKey, flagKeyShort, "", "Key of the object to tag", core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.ObjectKeys(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagName)))
		}, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(flagVersionId, "", "", "Tag this specific object version instead of the current one (versioned buckets only)")
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the full tag set to apply (a \"TagSet\" array of {\"Key\",\"Value\"} pairs). See --json-properties-example")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example tag-set JSON to stdout and exit, without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
