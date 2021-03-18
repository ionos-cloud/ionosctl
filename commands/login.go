package commands

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

func login() *builder.Command {
	loginCmd := builder.NewCommand(context.TODO(), nil, noPreRun, RunLoginUser, "login", "Authentication command for SDK",
		`Use this command to authenticate. By default, user data will be saved in:

* macOS: `+"`"+`${HOME}/Library/Application Support/ionosctl/config.json`+"`"+`
* Linux: `+"`"+`${XDG_CONFIG_HOME}/ionosctl/config.json`+"`"+`
* Windows: `+"`"+`%APPDATA%\ionosctl\config.json`+"`"+`.

You can use another configuration file for authentication with `+"`"+`--config`+"`"+` global option.

Note: The command can also be used without `+"`"+`--user`+"`"+` and `+"`"+`--password`+"`"+` flags. For more details, see Examples.`,
		loginExamples, false)
	loginCmd.AddStringFlag("user", "", "", "Username to authenticate")
	loginCmd.AddStringFlag("password", "", "", "Password to authenticate")

	return loginCmd
}

func RunLoginUser(c *builder.CommandConfig) error {
	var err error
	user := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, "user"))
	pwd := viper.GetString(builder.GetFlagName(c.ParentName, c.Name, "password"))

	if user == "" {
		err := c.Printer.Print("Enter your username:")
		if err != nil {
			return err
		}
		in := bufio.NewReader(c.Stdin)
		user, err = in.ReadString('\n')
		if err != nil {
			return err
		}
		user = strings.TrimRight(user, "\r\n")
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
	viper.Set(config.Username, user)
	viper.Set(config.Password, pwd)

	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.ArgServerUrl),
	)
	if err != nil {
		return err
	}

	dcsSvc := resources.NewDataCenterService(clientSvc.Get(), context.TODO())
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
