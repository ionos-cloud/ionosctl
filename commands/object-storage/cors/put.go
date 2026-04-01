package cors

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

const corsExample = `{
  "CORSRules": [
    {
      "AllowedOrigins": ["http://www.example.com"],
      "AllowedMethods": ["GET", "PUT", "POST"],
      "AllowedHeaders": ["*"],
      "ExposeHeaders": ["x-amz-request-id"],
      "MaxAgeSeconds": 3600
    }
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "cors",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the CORS configuration for a bucket",
		LongDesc: "Create or replace the CORS configuration for a bucket. " +
			"The configuration must be provided as a path to a JSON file via --json-properties. " +
			"Use --json-properties-example to see an example CORS configuration.",
		Example: "ionosctl object-storage cors put --name my-bucket --json-properties cors.json\n" +
			"ionosctl object-storage cors put --json-properties-example",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			exampleFlag := viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample))
			if exampleFlag {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), corsExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading CORS input: %w", err)
			}

			var corsReq objectstorage.PutBucketCorsRequest
			if err := json.Unmarshal(data, &corsReq); err != nil {
				return fmt.Errorf("parsing CORS JSON: %w", err)
			}

			s3, _, err := client.GetRegionalObjectStorageClient(context.Background(), name)
			if err != nil {
				return err
			}

			_, err = s3.CORSApi.PutBucketCors(context.Background(), name).
				PutBucketCorsRequest(corsReq).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "CORS configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file containing the CORS configuration")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example CORS configuration JSON and exit")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
