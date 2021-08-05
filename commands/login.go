package commands

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources/v5"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func login() *core.Command {
	ctx := context.TODO()
	loginCmd := core.NewCommand(ctx, nil, core.CommandBuilder{
		Namespace: "login",
		Resource:  "login",
		Verb:      "login",
		Aliases:   []string{"log", "auth"},
		ShortDesc: "Authentication command for SDK",
		LongDesc: `Use this command to authenticate. By default, user data will be saved in:

* macOS: ` + "`" + `${HOME}/Library/Application Support/ionosctl/config.json` + "`" + `
* Linux: ` + "`" + `${XDG_CONFIG_HOME}/ionosctl/config.json` + "`" + `
* Windows: ` + "`" + `%APPDATA%\ionosctl\config.json` + "`" + `.

You can use another configuration file for authentication with ` + "`" + `--config` + "`" + ` global option.

Note: The command can also be used without ` + "`" + `--user` + "`" + ` and ` + "`" + `--password` + "`" + ` flags. For more details, see Examples.`,
		Example:    loginExamples,
		PreCmdRun:  noPreRun,
		CmdRun:     RunLoginUser,
		InitClient: false,
	})
	loginCmd.AddStringFlag(config.ArgUser, "", "", "Username to authenticate")
	loginCmd.AddStringFlag(config.ArgPassword, config.ArgPasswordShort, "", "Password to authenticate")
	loginCmd.AddStringFlag(config.ArgToken, "", "", "Token to authenticate")

	return loginCmd
}

func RunLoginUser(c *core.CommandConfig) error {
	var err error
	username := viper.GetString(core.GetFlagName(c.NS, config.ArgUser))
	pwd := viper.GetString(core.GetFlagName(c.NS, config.ArgPassword))
	token := viper.GetString(core.GetFlagName(c.NS, config.ArgToken))

	if username == "" {
		err := c.Printer.Print("Enter your username:")
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
		err := c.Printer.Print("Enter your password:")
		if err != nil {
			return err
		}
		bytesPwd, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return err
		}
		pwd = string(bytesPwd)
	}
	viper.Set(config.Username, username)
	viper.Set(config.Password, pwd)
	viper.Set(config.Token, token)
	viper.Set(config.ServerUrl, viper.GetString(config.ArgServerUrl))

	clientSvc, err := v5.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		config.GetServerUrl(),
	)
	if err != nil {
		return err
	}

	dcsSvc := v5.NewDataCenterService(clientSvc.Get(), context.TODO())
	_, _, err = dcsSvc.List()
	if err != nil {
		return err
	}

	// Store credentials
	err = config.WriteFile()
	if err != nil {
		return err
	}

	return c.Printer.Print(printer.Result{
		Message: "Authentication successful!",
	})
}
