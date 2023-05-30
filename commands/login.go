package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"golang.org/x/term"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	sdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

const (
	loginFlagUseApiUrl = "use-default-api-url"

	loginExamples = `ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD

ionosctl login --token $IONOS_TOKEN

ionosctl login
Enter your username:
USERNAME
Enter your password:`
)

func LoginCmd() *core.Command {
	ctx := context.TODO()
	loginCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "login",
		Resource:  "login",
		Verb:      "login",
		Aliases:   []string{"log", "auth"},
		ShortDesc: "Authentication command for SDK",
		LongDesc: fmt.Sprintf(`The 'login' command allows you to authenticate with the IONOS Cloud APIs. There are three ways you can use it:
  1. Interactive mode: Just type 'ionosctl login' and you'll be prompted to enter your username and password.
  2. Use the '--user' and '--password' flags: Enter your credentials in the command.
  3. Use the '--token' flag: Provide an authentication token.

If you use a username and password, this command generates a token that's saved in the config file. Please keep this token safe. If you specify a custom '--api-url', it'll be saved to the config file when you login successfully and used for future API calls.

By default, the config file is located at %s. If you want to use a different config file, use the '--config' global option. Changing the permissions of the config file will cause it to no longer work.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN, these override the config file token.`, config.GetConfigFile()),
		Example:    loginExamples,
		PreCmdRun:  PreRunLoginCmd,
		CmdRun:     RunLoginUser,
		InitClient: false,
	})
	loginCmd.AddStringFlag(constants.ArgUser, "", "", "Username to authenticate")
	loginCmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Password to authenticate")
	loginCmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", "Token to authenticate")
	loginCmd.AddBoolFlag(loginFlagUseApiUrl, "", false, fmt.Sprintf(
		"Use the default authentication URL (%s) for auth checking, even if you specify a different '--%s'", constants.DefaultApiURL, constants.ArgServerUrl))

	return loginCmd
}

func PreRunLoginCmd(c *core.PreCommandConfig) error {
	if (viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) || viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword))) && viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return fmt.Errorf("use either --%s and/or --%s, either --%s", constants.ArgUser, constants.ArgPassword, constants.ArgToken)
	}
	return nil
}

func RunLoginUser(c *core.CommandConfig) error {
	c.Printer.Verbose("Note: The login command will save the credentials in a configuration file after the authentication is successful!")
	c.Printer.Verbose("Note: As an alternative to this, ionosctl offers support for environment variables: $%s, $%s or $%s.",
		sdk.IonosUsernameEnvVar, sdk.IonosPasswordEnvVar, sdk.IonosTokenEnvVar)

	data, err := buildConfigData(c)
	if err != nil {
		return fmt.Errorf("failed building config data: %w", err)
	}

	// Store credentials
	c.Printer.Verbose("Storing credentials to the configuration file: %s", viper.GetString(constants.ArgConfig))
	err = config.Write(data)
	if err != nil {
		return fmt.Errorf("failed writing config data: %w", err)
	}

	ls, _, err := client.Must().AuthClient.TokensApi.TokensGet(context.Background()).Execute()
	if err != nil {
		return err
	}

	// Go through data struct and blank out all credentials, including old ones (userdata.username, userdata.password)
	msg := fmt.Sprintf("Authentication successful. Note: Your account has %d active tokens. Created the following fields in your config file:\n", len(*ls.Tokens))
	for k, _ := range data {
		msg += fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata."))
	}

	return c.Printer.Print(msg)
}

// buildConfigData returns map that will be written to config file, while also checking token is usable via testToken
func buildConfigData(c *core.CommandConfig) (map[string]string, error) {
	configData := map[string]string{} // This map is what we will write to the config file

	// API URL
	c.Printer.Verbose("API Url: %s", config.GetServerUrl())
	if explicitUrl := config.GetServerUrl(); strings.TrimLeft(explicitUrl, "https://") != strings.TrimLeft(constants.DefaultApiURL, "https://") {
		// Don't save the API url to the config if it's the default, since we don't want to revert to that value if the user doesn't provide any.
		// This was changed from old behaviour because some APIs (e.g. DNS API) [can] use a different server URL
		c.Printer.Verbose("Saving API URL to config file")
		configData[constants.CfgServerUrl] = explicitUrl
	}

	// Explicit token
	if fn := core.GetFlagName(c.NS, constants.ArgToken); viper.IsSet(fn) {
		tok := viper.GetString(fn)
		configData[constants.CfgToken] = tok
		return configData, client.TestCreds("", "", tok) // Early return for preset token
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
	if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) {
		username = viper.GetString(fn)
	} else {
		// Interactively ask for username
		c.Printer.Print("Enter your username: ")
		reader := bufio.NewReader(c.Stdin)
		var err error
		username, err = reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed reading username from set reader")
		}
		username = strings.TrimSpace(username) // remove trailing newline
	}

	if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
		password = viper.GetString(fn)
	} else {
		// Interactively ask for password
		c.Printer.Print("Enter your password: ")

		if file, ok := c.Stdin.(*os.File); ok {
			bytePassword, err := term.ReadPassword(int(file.Fd()))
			if err != nil {
				return nil, fmt.Errorf("failed securely reading password from set file descriptor")
			}
			password = string(bytePassword)
		} else {
			return nil, fmt.Errorf("the set input does not have a file descriptor (is it set to a terminal?)")
		}
	}

	cl, err := client.NewClient(username, password, "", config.GetServerUrl())
	tok, _, err := cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed using username and password to generate a token: %w", err)
	}

	configData[constants.CfgToken] = *tok.Token
	return configData, client.TestCreds("", "", *tok.Token)
}
