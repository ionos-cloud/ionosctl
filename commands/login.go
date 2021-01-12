package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"syscall"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
)

func login() *cobra.Command {
	var (
		user string
	)
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Authentication command for SDK",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			err := preRunLoginUser(user)
			return err
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runLoginUser()
			return err
		},
	}
	flags := loginCmd.Flags()
	flags.StringVar(&user, "user", "", "Username to login")

	return loginCmd
}

func preRunLoginUser(user string) error {
	if user == "" {
		fmt.Println("Enter your username:")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		user = scanner.Text()
	}

	fmt.Println("Enter your password:")
	bytesPwd, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	pwd := string(bytesPwd)

	viper.Set(config.Username, user)
	viper.Set(config.Password, pwd)
	return nil
}

func runLoginUser() error {
	cfg := config.GetAPIClientConfig()
	apiClient := ionoscloud.NewAPIClient(cfg)

	req := apiClient.DataCenterApi.DatacentersGet(context.Background())
	_, resp, _ := apiClient.DataCenterApi.DatacentersGetExecute(req)
	if resp != nil {
		if resp.StatusCode == 401 {
			return errors.New("user credentials are invalid")
		} else {
			fmt.Println("Authentication successful!")
		}
	}

	// Store credentials
	err := config.WriteConfigFile()
	return err
}
