package user

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ionos-cloud/ionosctl/v6/commands/kafka/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/sdk-go-bundle/products/kafka/v2"
	"github.com/spf13/viper"
)

func GetAccess() *core.Command {
	cmd := core.NewCommand(
		context.Background(), nil, core.CommandBuilder{
			Verb:      "get-access",
			Namespace: "kafka",
			Resource:  "user",
			ShortDesc: "Get a user's credentials",
			LongDesc: `Get a Kafka user's credentials including certificate, private key, and CA certificate.
By default, the command writes three PEM files to the specified output directory (or current directory if not specified):
 - <username>-cert.pem
 - <username>-key.pem
 - <username>-ca.pem

You can also use '--output json' to print the full JSON response from the API to stdout instead of writing files.

IMPORTANT: Keep these credentials secure. The private key should never be shared or exposed publicly.`,
			Aliases: []string{"g", "get", "access"},
			Example: "ionosctl kafka user get-access " + core.FlagsUsage(constants.FlagLocation, constants.FlagClusterId, constants.FlagUserId),
			PreCmdRun: func(cmd *core.PreCommandConfig) error {

				return core.CheckRequiredFlags(
					cmd.Command, cmd.NS, constants.FlagLocation, constants.FlagClusterId, constants.FlagUserId,
				)
			},
			CmdRun: func(cmd *core.CommandConfig) error {
				clusterID, _ := cmd.Command.Command.Flags().GetString(constants.FlagClusterId)
				userID, _ := cmd.Command.Command.Flags().GetString(constants.FlagUserId)

				userAccess, _, err := client.Must().Kafka.UsersApi.ClustersUsersAccessGet(
					context.Background(), clusterID, userID,
				).Execute()
				if err != nil {
					return fmt.Errorf("unable to get user's credentials: %s", err)
				}

				cert, ok := userAccess.Metadata.GetCertificateOk()
				if !ok || cert == nil {
					return fmt.Errorf("certificate not found in the response")
				}
				priv, ok := userAccess.Metadata.GetPrivateKeyOk()
				if !ok || priv == nil {
					return fmt.Errorf("private key not found in the response")
				}
				ca, ok := userAccess.Metadata.GetCertificateAuthorityOk()
				if !ok || ca == nil {
					return fmt.Errorf("CA certificate not found in the response")
				}

				output := viper.GetString(constants.ArgOutput)
				if output == "json" || output == "api-json" {
					out, err := json.MarshalIndent(userAccess, "", "  ")
					if err != nil {
						return fmt.Errorf("failed to marshal credentials for stdout: %w", err)
					}
					cmd.Command.Command.Println(string(out))
					return nil
				}

				// Default path: write 3 PEM files with 0600 perms
				outputDir, _ := cmd.Command.Command.Flags().GetString("output-dir")
				if outputDir == "" {
					outputDir = "."
				}

				// Ensure output directory exists
				if err := os.MkdirAll(outputDir, 0o700); err != nil {
					return fmt.Errorf("unable to create output directory %s: %w", outputDir, err)
				}

				base := fmt.Sprintf("%s", userAccess.Properties.Name)
				certFile := filepath.Join(outputDir, fmt.Sprintf("%s-cert.pem", base))
				keyFile := filepath.Join(outputDir, fmt.Sprintf("%s-key.pem", base))
				caFile := filepath.Join(outputDir, fmt.Sprintf("%s-ca.pem", base))

				if err := writeFile(certFile, *cert); err != nil {
					return err
				}
				if err := writeFile(keyFile, *priv); err != nil {
					return err
				}
				if err := writeFile(caFile, *ca); err != nil {
					return err
				}

				// Print summary
				fmt.Fprintln(cmd.Command.Command.OutOrStdout(), "Wrote:")
				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), " - %s\n", caFile)
				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), " - %s\n", certFile)
				fmt.Fprintf(cmd.Command.Command.OutOrStdout(), " - %s", keyFile)

				return nil
			},
			InitClient: true,
		},
	)

	cmd.AddStringFlag(
		constants.FlagClusterId, "", "", "The ID of the cluster",
		core.RequiredFlagOption(), core.WithCompletion(
			func() []string {
				return completer.ClustersProperty(
					func(read kafka.ClusterRead) string {
						return read.Id
					},
				)
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)
	cmd.AddStringFlag(
		constants.FlagUserId, "", "", "The ID of the user", core.RequiredFlagOption(),
		core.WithCompletion(
			func() []string {
				return completer.Users(cmd.Command.Flag(constants.FlagClusterId).Value.String())
			}, constants.KafkaApiRegionalURL, constants.KafkaLocations,
		),
	)

	cmd.AddStringFlag("output-dir", "", ".", "Directory to save the user's credential PEM files")

	return cmd
}

func writeFile(path, content string) error {
	// ensure trailing newline for PEM readability
	if len(content) > 0 && content[len(content)-1] != '\n' {
		content = content + "\n"
	}
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
