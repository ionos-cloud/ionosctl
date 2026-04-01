package encryption

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

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
		LongDesc: "Create or replace the default encryption configuration for a bucket. " +
			"The configuration must be provided as a path to a JSON file via --json-properties. " +
			"Use --json-properties-example to see an example encryption configuration.",
		Example: "ionosctl object-storage encryption put --name my-bucket --json-properties encryption.json\n" +
			"ionosctl object-storage encryption put --json-properties-example",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			exampleFlag := viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample))
			if exampleFlag {
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

			s3, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			_, err = s3.EncryptionApi.PutBucketEncryption(context.Background(), name).
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file containing the encryption configuration")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example encryption configuration JSON and exit")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
