package lifecycle

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	objectstorage "github.com/ionos-cloud/sdk-go-bundle/products/objectstorage/v2"
	"github.com/spf13/viper"

	"github.com/ionos-cloud/ionosctl/v6/commands/object-storage/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

const lifecycleExample = `{
  "Rules": [
    {
      "ID": "expire-old-logs",
      "Prefix": "logs/",
      "Status": "Enabled",
      "Expiration": {
        "Days": 90
      }
    },
    {
      "ID": "cleanup-incomplete-uploads",
      "Prefix": "",
      "Status": "Enabled",
      "AbortIncompleteMultipartUpload": {
        "DaysAfterInitiation": 7
      }
    }
  ]
}`

func PutCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "object-storage",
		Resource:  "lifecycle",
		Verb:      "put",
		Aliases:   []string{"p"},
		ShortDesc: "Create or replace the lifecycle configuration for a bucket",
		LongDesc: `Create or replace a bucket's lifecycle configuration. This fully REPLACES any existing rules with the "Rules" array in the file.

Provide the configuration as a JSON file via --json-properties. Each rule needs an "ID", a "Prefix" (which object keys it applies to; "" = all objects), a "Status" ("Enabled" or "Disabled"), and at least one action:
  Expiration.Days / Expiration.Date            Delete current objects after N days, or on a specific date.
  NoncurrentVersionExpiration.NoncurrentDays   Delete old versions N days after they become noncurrent (versioned buckets).
  AbortIncompleteMultipartUpload.DaysAfterInitiation   Abort unfinished multipart uploads after N days.

Note: on a versioned bucket, Expiration alone only adds a delete marker and does not reclaim storage - pair it with NoncurrentVersionExpiration. ionosctl computes and sends the required Content-MD5 for the request automatically.

Run with --json-properties-example to print a ready-to-edit template (expire objects under logs/ after 90 days, plus abort stale multipart uploads after 7 days).`,
		Example: `# Apply lifecycle rules from a file
ionosctl object-storage bucket lifecycle put --name my-bucket --json-properties lifecycle.json

# Print an example configuration
ionosctl object-storage bucket lifecycle put --json-properties-example`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				return nil
			}
			return core.CheckRequiredFlags(c.Command, c.NS, constants.FlagName, constants.FlagJsonProperties)
		},
		CmdRun: func(c *core.CommandConfig) error {
			if viper.GetBool(core.GetFlagName(c.NS, constants.FlagJsonPropertiesExample)) {
				fmt.Fprintln(c.Command.Command.OutOrStdout(), lifecycleExample)
				return nil
			}

			name := viper.GetString(core.GetFlagName(c.NS, constants.FlagName))
			input := viper.GetString(core.GetFlagName(c.NS, constants.FlagJsonProperties))

			data, err := os.ReadFile(input)
			if err != nil {
				return fmt.Errorf("reading lifecycle input: %w", err)
			}

			var req objectstorage.PutBucketLifecycleRequest
			if err := json.Unmarshal(data, &req); err != nil {
				return fmt.Errorf("parsing lifecycle JSON: %w", err)
			}

			xmlBytes, err := xml.Marshal(req)
			if err != nil {
				return fmt.Errorf("serializing lifecycle configuration: %w", err)
			}
			hash := md5.Sum(xmlBytes)
			contentMD5 := base64.StdEncoding.EncodeToString(hash[:])

			_, err = client.MustObjectStorage().ObjectStorageClient.LifecycleApi.PutBucketLifecycle(c.Context, name).
				ContentMD5(contentMD5).
				PutBucketLifecycleRequest(req).
				Execute()
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "Lifecycle configuration for %q applied successfully\n", name)
			return nil
		},
		InitClient: false,
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "Name of the bucket", core.RequiredFlagOption(),
		core.WithCompletion(completer.BucketNames, constants.ObjectStorageApiRegionalURL, constants.ObjectStorageLocations))
	cmd.AddStringFlag(constants.FlagJsonProperties, "", "", "Path to a JSON file with the lifecycle rules ({\"Rules\":[...]}). Replaces all existing rules")
	cmd.AddBoolFlag(constants.FlagJsonPropertiesExample, "", false, "Print an example lifecycle configuration JSON and exit without contacting the API")

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false
	return cmd
}
