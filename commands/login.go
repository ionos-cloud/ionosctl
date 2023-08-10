package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	client2 "github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/jsontabwriter"

	"github.com/ionos-cloud/ionosctl/v6/pkg/config"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	sdk "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

const (
	loginExamples = `ionosctl login --user $IONOS_USERNAME --password $IONOS_PASSWORD
Status: Authentication successful!

ionosctl login --token $IONOS_TOKEN
Status: Authentication successful!

ionosctl login
Enter your username:
USERNAME
Enter your password:

Status: Authentication successful!`
)

func LoginCmd() *core.Command {
	ctx := context.TODO()
	loginCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "login",
		Resource:  "login",
		Verb:      "login",
		Aliases:   []string{"log", "auth"},
		ShortDesc: "Authentication command for SDK",
		LongDesc: `Use this command to authenticate. You can use  ` + "`" + `--user` + "`" + ` and ` + "`" + `--password` + "`" + ` flags or you can use  ` + "`" + `--token` + "`" + ` flag to set the credentials.

By default, the user data after running this command will be saved in:

* macOS: ` + "`" + `${HOME}/Library/Application Support/ionosctl/config.json` + "`" + `
* Linux: ` + "`" + `${XDG_CONFIG_HOME}/ionosctl/config.json` + "`" + `
* Windows: ` + "`" + `%APPDATA%\ionosctl\config.json` + "`" + `.

You can use another configuration file for authentication with the ` + "`" + `--config` + "`" + ` global option.

Note: The IONOS Cloud CLI supports also authentication with environment variables: $IONOS_USERNAME, $IONOS_PASSWORD or $IONOS_TOKEN.`,
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
	if viper.IsSet(core.GetFlagName(c.NS, constants.ArgUser)) && viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) && viper.IsSet(core.GetFlagName(c.NS, constants.ArgToken)) {
		return errors.New("it is recommended to use either username + password, either token")
	}
	return nil
}

func RunLoginUser(c *core.CommandConfig) error {
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput(
		"Note: The login command will save the credentials in a configuration file after the authentication is successful!"))
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput(
		"Note: As an alternative to this, ionosctl offers support for environment variables: $%s, $%s or $%s.",
		sdk.IonosUsernameEnvVar, sdk.IonosPasswordEnvVar, sdk.IonosTokenEnvVar))

	username := viper.GetString(core.GetFlagName(c.NS, constants.ArgUser))
	pwd := viper.GetString(core.GetFlagName(c.NS, constants.ArgPassword))
	token := viper.GetString(core.GetFlagName(c.NS, constants.ArgToken))

	if token != "" {
		// If token is set, use only token
		viper.Set(constants.Token, token)
		fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Token is set."))
	} else {
		// If token and username are not set, display messages
		if username == "" {
			_, err := fmt.Fprintf(c.Stdout, jsontabwriter.GenerateLogOutput("Enter your username:"))
			if err != nil {
				return err
			}

			in := bufio.NewReader(c.Stdin)
			username, err = in.ReadString('\n')
			if err != nil {
				return err
			}

			username = strings.TrimRight(username, "\r\n")
		}

		if pwd == "" {
			_, err := fmt.Fprintf(c.Stdout, jsontabwriter.GenerateLogOutput("Enter your password:"))
			if err != nil {
				return err
			}

			bytesPwd, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}

			pwd = string(bytesPwd)
		}

		viper.Set(constants.Username, username)
		fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Username is set %s", viper.GetString(constants.Username)))

		viper.Set(constants.Password, pwd)
		fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Password is set"))
	}

	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("ServerUrl: %s", config.GetServerUrl()))
	viper.Set(constants.ServerUrl, viper.GetString(constants.ArgServerUrl))

	client, err := client2.Get()
	if err != nil {
		return err
	}

	// Check the auth is correct
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput("Checking authentication..."))
	dcsSvc := resources.NewDataCenterService(client, context.TODO())
	_, _, err = dcsSvc.List(resources.ListQueryParams{})
	if err != nil {
		return err
	}

	// Store credentials
	fmt.Fprintf(c.Stderr, jsontabwriter.GenerateVerboseOutput(
		"Storing credentials to the configuration file: %v", viper.GetString(constants.ArgConfig)))
	err = config.WriteFile()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Stdout, jsontabwriter.GenerateLogOutput("Authentication successful!"))

	return nil
}
