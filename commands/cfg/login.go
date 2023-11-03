package cfg

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	sdk "github.com/ionos-cloud/sdk-go/v6"
	"golang.org/x/term"

	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

func LoginCmd() *core.Command {

	const (
		loginExamples = `ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD

ionosctl login --token $IONOS_TOKEN

ionosctl login
Enter your username:
USERNAME
Enter your password:`
	)

	ctx := context.TODO()
	loginCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "login",
		Resource:  "login",
		Verb:      "login",
		Aliases:   []string{"log", "auth"},
		ShortDesc: "Use credentials to generate a config file in `ionosctl cfg location`",
		LongDesc: fmt.Sprintf(`The 'login' command allows you to authenticate with the IONOS Cloud APIs. There are three ways you can use it:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.
Note: If using '--token', you can skip verifying the used token with '--skip-verify'

If you specify a custom '--api-url', the custom URL will also be saved to the config file when you login successfully and used for future API calls.

If you use a username and password, this command will use these credentials to generate a token that will be saved in the config file. Please keep this token safe.

To find your config file location, use 'ionosctl cfg location'. If you want to use a different config file, use the '--config' global option. Changing the permissions of the config file will cause it to no longer work.

If a config file already exists, you will be asked to replace it. You can skip this verification with '--force'

%s`, constants.DescAuthenticationOrder),
		Example:    loginExamples,
		PreCmdRun:  PreRunLoginCmd,
		CmdRun:     RunLoginUser,
		InitClient: false,
	})
	loginCmd.AddStringFlag(constants.ArgUser, "", "", "Username to authenticate. Will be used to generate a token")
	loginCmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Password to authenticate. Will be used to generate a token")
	loginCmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", "Token to authenticate. If used, will be saved to the config file without generating a new token. Note: mutually exclusive with --user and --password")
	loginCmd.AddBoolFlag(constants.FlagSkipVerify, "", false, "Forcefully write the provided token to the config file without verifying if it is valid. Note: --token is required")

	return loginCmd
}

func PreRunLoginCmd(c *core.PreCommandConfig) error {
	if (viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) || viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword))) && viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return fmt.Errorf("use either --%s and --%s, either --%s", constants.ArgUser, constants.ArgPassword, constants.ArgToken)
	}
	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagSkipVerify)) && !viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return fmt.Errorf("--%s requires --%s to be set", constants.FlagSkipVerify, constants.ArgToken)
	}
	return nil
}

func RunLoginUser(c *core.CommandConfig) error {
	configPath := config.GetConfigFilePath()

	err := handleExistingConfigFile(c, configPath)
	if err != nil {
		return fmt.Errorf("failed handling existing config file: %w", err)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Note: The login command will save the credentials in a configuration file after the authentication is successful!"))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Note: As an alternative to this, ionosctl offers support for environment variables: $%s, $%s or $%s.", sdk.IonosUsernameEnvVar, sdk.IonosPasswordEnvVar, sdk.IonosTokenEnvVar))

	validCredentials := true
	data, err := buildConfigData(c)
	if err != nil {
		if !viper.GetBool(core.GetFlagName(c.NS, constants.FlagSkipVerify)) {
			return fmt.Errorf("failed building config data: %w", err)
		}
		validCredentials = false
	}

	// Store credentials
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Storing credentials to the configuration file: %s", viper.GetString(constants.ArgConfig)))
	err = config.Write(data)
	if err != nil {
		return fmt.Errorf("failed writing config data: %w", err)
	}

	var msg strings.Builder
	if validCredentials {
		ls, _, errTokens := client.NewClientFromCfgData(data).AuthClient.TokensApi.TokensGet(context.Background()).Execute()
		if errTokens != nil {
			return fmt.Errorf("failed retrieving current tokens: %w", err)
		}
		msgActiveTokens := fmt.Sprintf("Note: Your account has %d active tokens. ", len(*ls.Tokens))
		msg.WriteString(msgActiveTokens)
	}

	msg.WriteString(fmt.Sprintf("Config file updated successfully. Created the following fields in your config file:\n"))
	for k := range data {
		msg.WriteString(fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata.")))
	}

	_, err = fmt.Fprintln(c.Command.Command.OutOrStdout(), msg.String())
	return err
}

func handleExistingConfigFile(c *core.CommandConfig, configPath string) error {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	if !confirm.FAsk(
		c.Command.Command.InOrStdin(),
		fmt.Sprintf("Config file %s already exists. Do you want to replace it", configPath),
		viper.GetBool(constants.ArgForce),
	) {
		return fmt.Errorf(confirm.UserDenied)
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("error deleting config file %s: %w", configPath, err)
	}

	return nil
}

// buildConfigData returns map that will be written to config file, while also checking token is usable via testToken
func buildConfigData(c *core.CommandConfig) (map[string]string, error) {
	configData := map[string]string{} // This map is what we will write to the config file

	// API URL
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("API Url: %s", config.GetServerUrl()))
	if explicitUrl := config.GetServerUrl(); explicitUrl != "" {
		// Don't save the API url to the config if it's the default, since we don't want to revert to it if the user doesn't provide any value.
		// This was changed from old behaviour because some APIs (e.g. DNS API) [can] use a different server URL
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Saving API URL to config file: %s", explicitUrl))
		configData[constants.CfgServerUrl] = explicitUrl
	}

	// Explicit token
	token, err := c.Command.Command.Flags().GetString(constants.ArgToken)
	if err != nil {
		return nil, fmt.Errorf("failed reading token flag: %w", err)
	}
	if token != "" {
		configData[constants.CfgToken] = token
		return configData, client.TestCreds("", "", token) // Early return for preset token
	}

	// If here, user did not give a pre-set token. Generate one via username & password.
	// If the username & password are not given via flags, ask for them interactively

	var (
		// Set by flags. If not set, interactive mode will ask the user iff token not set.
		// They die here in this func block. Don't use them anywhere else!
		username string
		password string
	)

	// Note: Because of the prerun checks, we know for sure the following:
	// Either flags for token set, or username set, or password set, or username & password set.
	if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) && viper.GetString(fn) != "" {
		username = viper.GetString(fn)
	} else {
		// Interactively ask for username
		fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your username: ")
		reader := bufio.NewReader(c.Command.Command.InOrStdin())
		username, err = reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed reading username from set reader")
		}
		username = strings.TrimSpace(username) // remove trailing newline
	}

	if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) && viper.GetString(fn) != "" {
		password = viper.GetString(fn)
	} else {
		// Interactively ask for password
		fmt.Fprintln(c.Command.Command.OutOrStdout(), "Enter your password: ")

		if file, ok := c.Command.Command.InOrStdin().(*os.File); ok {
			bytePassword, err := term.ReadPassword(int(file.Fd()))
			if err != nil {
				return nil, fmt.Errorf("failed securely reading password from set file descriptor")
			}
			password = string(bytePassword)
		} else {
			return nil, fmt.Errorf("the set input does not have a file descriptor (is it set to a terminal?)")
		}
	}

	cl := client.NewClient(username, password, "", config.GetServerUrl())
	tok, _, err := cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed using username and password to generate a token: %w", err)
	}

	configData[constants.CfgToken] = *tok.Token
	return configData, nil
}
