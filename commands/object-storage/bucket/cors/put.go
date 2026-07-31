package cors

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
		LongDesc: `Create or replace a bucket's CORS configuration. This is a full REPLACE, not a merge: the CORSRules array in the file becomes the bucket's entire rule set, so include every rule you want to keep.

Provide the configuration as a JSON file via --json-properties. The top-level object holds a "CORSRules" array; each rule supports:
  AllowedOrigins  Origins permitted to make cross-origin requests, e.g. "https://app.example.com". "*" allows any origin.
  AllowedMethods  HTTP methods allowed from those origins (GET, PUT, POST, DELETE, HEAD).
  AllowedHeaders  Request headers the browser may send in the actual request; "*" allows any. Matched against the Access-Control-Request-Headers preflight header.
  ExposeHeaders   Response headers that browser JavaScript is allowed to read (browsers hide all others by default).
  MaxAgeSeconds   How long the browser may cache the preflight (OPTIONS) result before asking again.

Run with --json-properties-example to print a ready-to-edit template.`,
		Example: `# Apply a CORS configuration from a file
ionosctl object-storage bucket cors put --name my-bucket --json-properties cors.json

# Print an example configuration to use as a starting point
ionosctl object-storage bucket cors put --json-properties-example`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
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

			_, err = client.MustObjectStorage().ObjectStorageClient.CORSApi.PutBucketCors(c.Context, name).
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

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the full CORS configuration ({\"CORSRules\":[...]}). Replaces any existing rules")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example CORS configuration JSON and exit without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
