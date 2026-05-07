package policy

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

// policyExample uses S3-compatible IAM policy syntax. The "s3:" action prefix
// and "arn:aws:s3:::" resource format are required by the S3-compatible API,
// not references to AWS services.
const policyExample = `{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "PublicReadGetObject",
      "Effect": "Allow",
      "Principal": {
        "AWS": ["*"]
      },
      "Action": ["s3:GetObject"],
      "Resource": ["arn:aws:s3:::BUCKET_NAME/*"]
    }
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "policy",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the bucket policy",
		LongDesc: "Create or replace the bucket policy. " +
			"The policy must be provided as a path to a JSON file via --json-properties. " +
			"Use --json-properties-example to see an example policy.",
		Example: "ionosctl object-storage bucket policy put --name my-bucket --json-properties policy.json\n" +
			"ionosctl object-storage bucket policy put --json-properties-example",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), policyExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := readInput(input)
			if err != nil {
				return fmt.Errorf("reading policy input: %w", err)
			}

			var bp objectstorage.BucketPolicy
			if err := json.Unmarshal(data, &bp); err != nil {
				return fmt.Errorf("parsing policy JSON: %w", err)
			}

			_, err = client.MustObjectStorage().ObjectStorageClient.PolicyApi.PutBucketPolicy(c.Context, name).
				BucketPolicy(bp).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Bucket policy for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file containing the bucket policy")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example bucket policy JSON and exit")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}

// readInput reads the policy JSON from a file path.
func readInput(path string) ([]byte, error) {
	return os.ReadFile(path)
}
