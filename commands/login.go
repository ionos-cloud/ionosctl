package commands

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

func login() *builder.Command {
	loginCmd := builder.NewCommand(context.TODO(), nil, noPreRun, RunLoginUser, "login", "Authentication command for SDK",
		`Use this command to authenticate. User data will be saved in `+"`"+`$XDG_CONFIG_HOME/ionosctl-config.json`+"`"+` file. 

You can use another configuration file for authentication with `+"`"+`--config`+"`"+` global option.

Note: The command can also be used without `+"`"+`--user`+"`"+` and `+"`"+`--password`+"`"+` flags (see Examples).`,
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
		c.Printer.Print("Enter your username:")
		in := bufio.NewReader(c.Stdin)
		user, err = in.ReadString('\n')
		if err != nil {
			return err
		}
		// Delete the delimiter from the string
		user = strings.TrimRight(user, "\r\n")
	}
	if pwd == "" {
		c.Printer.Print("Enter your password:")
		bytesPwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
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

	c.Printer.Print(utils.Result{
		Message: "Authentication successful!",
	})
	return nil
}
