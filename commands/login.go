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
		LongDesc: fmt.Sprintf(`Use this command to authenticate.
You can either use the interactive mode, or you can use "--user" and "--password" flags or "--token" flag to set the credentials.
If using username & password, this command will generate a JWT token which will be saved in the config file. Please safeguard your token.
The config file, by default, will be created at %s. You can use another configuration file for authentication with the "--config" global option.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN, these override the config file token.`, config.GetConfigFile()),
		Example:    loginExamples,
		PreCmdRun:  PreRunLoginCmd,
		CmdRun:     RunLoginUser,
		InitClient: false,
	})
	loginCmd.AddStringFlag(constants.ArgUser, "", "", "Username to authenticate")
	loginCmd.AddStringFlag(constants.ArgPassword, constants.ArgPasswordShort, "", "Password to authenticate")
	loginCmd.AddStringFlag(constants.ArgToken, constants.ArgTokenShort, "", "Token to authenticate")

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
	c.Printer.Verbose("Storing credentials to the configuration file: %v", viper.GetString(constants.ArgConfig))
	err = config.WriteFile(data)
	if err != nil {
		return fmt.Errorf("failed writing config data: %w", err)
	}

	msg := "Authentication successful. Created the following fields in your config file:\n"
	for k, _ := range data {
		msg += fmt.Sprintf(" • %s\n", strings.TrimPrefix(k, "userdata."))
	}

	return c.Printer.Print(msg)
}

func testToken(token string) error {
	cl, err := client.NewClient("", "", token, config.GetServerUrl())
	if err != nil {
		return fmt.Errorf("failed getting client via token: %w", err)
	}

	_, _, err = cl.CloudClient.DataCentersApi.DatacentersGet(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("failed running a test SDK func (DatacentersGet): %w", err)
	}

	return nil
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
		configData[constants.ServerUrl] = explicitUrl
	}

	// Explicit token
	if fn := core.GetFlagName(c.NS, constants.ArgToken); viper.IsSet(fn) {
		tok := viper.GetString(fn)
		configData[constants.Token] = tok
		return configData, testToken(tok) // Early return for preset token
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
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username) // remove trailing newline
	}

	if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
		password = viper.GetString(fn)
	} else {
		// Interactively ask for password
		c.Printer.Print("Enter your password: ")

		if file, ok := c.Stdin.(*os.File); ok {
			bytePassword, _ := term.ReadPassword(int(file.Fd()))
			password = string(bytePassword)
			fmt.Println() // print a newline after password input
		} else {
			return nil, fmt.Errorf("the set input does not have a file descriptor (is it set to a terminal?)")
		}
	}

	cl, err := client.NewClient(username, password, "", config.GetServerUrl())
	tok, _, err := cl.AuthClient.TokensApi.TokensGenerate(context.Background()).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed using username and password to generate a token: %w", err)
	}

	configData[constants.Token] = *tok.Token
	return configData, testToken(*tok.Token)
}
