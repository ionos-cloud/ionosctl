package publicaccessblock

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

const publicAccessBlockExample = `{
  "BlockPublicAcls": true,
  "IgnorePublicAcls": true,
  "BlockPublicPolicy": true,
  "RestrictPublicBuckets": true
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "public-access-block",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the public access block configuration for a bucket",
		LongDesc: `Create or replace the bucket's Public Access Block. These guardrails override ACLs and bucket policies, so turning them on is the dependable way to force a bucket private regardless of what its policy/ACLs say.

Provide the configuration as a JSON file via --json-properties - a flat object with four booleans:
  BlockPublicAcls        Reject requests that would apply a public ACL.
  IgnorePublicAcls       Ignore public ACLs already set.
  BlockPublicPolicy      Reject bucket policies that grant public access.
  RestrictPublicBuckets  Limit access via an existing public policy to authorized principals.

Set all four to true to lock a bucket down completely. This is a full REPLACE of the configuration. Run with --json-properties-example to print a template with every guardrail enabled.`,
		Example: `# Fully block public access (all four guardrails on)
ionosctl object-storage bucket public-access-block put --name my-bucket --json-properties config.json

# Print an example configuration
ionosctl object-storage bucket public-access-block put --json-properties-example`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), publicAccessBlockExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading public access block input: %w", err)
			}

			var req objectstorage.BlockPublicAccessPayload
			if err := json.Unmarshal(data, &req); err != nil {
				return fmt.Errorf("parsing public access block JSON: %w", err)
			}

			_, err = client.MustObjectStorage().ObjectStorageClient.PublicAccessBlockApi.PutPublicAccessBlock(c.Context, name).
				BlockPublicAccessPayload(req).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Public access block configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the four block flags (BlockPublicAcls, IgnorePublicAcls, BlockPublicPolicy, RestrictPublicBuckets). Replaces the existing configuration")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example public access block configuration JSON and exit without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
