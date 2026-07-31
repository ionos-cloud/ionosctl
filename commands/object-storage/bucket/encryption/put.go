package encryption

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

const encryptionExample = `{
  "Rules": [
    {
      "ApplyServerSideEncryptionByDefault": {
        "SSEAlgorithm": "AES256"
      }
    }
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "encryption",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the default encryption configuration for a bucket",
		LongDesc: `Create or replace a bucket's default server-side encryption rule. From then on, new objects are encrypted at rest by the server unless the individual upload specifies its own encryption. Existing objects are not re-encrypted.

Provide the configuration as a JSON file via --json-properties. The top-level object holds a "Rules" array; each rule carries "ApplyServerSideEncryptionByDefault" with an "SSEAlgorithm". For SSE with server-managed keys use "AES256". (This is the algorithm returned by 'get'.)

Run with --json-properties-example to print a ready-to-edit template.`,
		Example: `# Set AES256 default encryption from a file
ionosctl object-storage bucket encryption put --name my-bucket --json-properties encryption.json

# Print an example configuration
ionosctl object-storage bucket encryption put --json-properties-example`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), encryptionExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading encryption input: %w", err)
			}

			var encReq objectstorage.PutBucketEncryptionRequest
			if err := json.Unmarshal(data, &encReq); err != nil {
				return fmt.Errorf("parsing encryption JSON: %w", err)
			}

			_, err = client.MustObjectStorage().ObjectStorageClient.EncryptionApi.PutBucketEncryption(c.Context, name).
				PutBucketEncryptionRequest(encReq).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Encryption configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the encryption rules ({\"Rules\":[{\"ApplyServerSideEncryptionByDefault\":{\"SSEAlgorithm\":\"AES256\"}}]}). Replaces any existing rule")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example encryption configuration JSON and exit without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
