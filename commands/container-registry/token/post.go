package token

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/commands/container-registry/registry"
	"github.com/ionos-cloud/ionosctl/pkg/constants"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	sdkgo "github.com/ionos-cloud/sdk-go-container-registry"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var tokenPostProperties = sdkgo.NewPostTokenPropertiesWithDefaults()

func TokenPostCmd() *core.Command {
	cmd := core.NewCommand(
		context.TODO(), nil, core.CommandBuilder{
			Namespace:  "container-registry",
			Resource:   "token",
			Verb:       "create",
			Aliases:    []string{"c"},
			ShortDesc:  "Create a new token",
			LongDesc:   "Create a new token used to access a container registry",
			Example:    "ionosctl container-registry token create --registry-id [REGISTRY-ID] --name [TOKEN-NAME]",
			PreCmdRun:  PreCmdPostToken,
			CmdRun:     CmdPostToken,
			InitClient: true,
		},
	)

	cmd.AddStringFlag(FlagName, "", "", "Name of the Token", core.RequiredFlagOption())
	cmd.AddStringFlag(FlagExpiryDate, "", "", "Expiry date of the Token")
	cmd.AddStringFlag(FlagStatus, "", "", "Status of the Token")
	_ = cmd.Command.RegisterFlagCompletionFunc(
		FlagStatus, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				"enabled", "disabled",
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.AddStringFlag(FlagTimeUntilExpiry, "", "", "Set the amount of time until the token expires (e.g. 1h, 1d, 1w, 1m, 1y)")

	cmd.AddStringFlag(FlagRegId, "r", "", "Registry ID", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(
		"registry-id", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return registry.RegsIds(), cobra.ShellCompDirectiveNoFileComp
		},
	)

	cmd.Command.Flags().StringSlice(constants.ArgCols, nil, printer.ColsMessage(allCols))
	_ = cmd.Command.RegisterFlagCompletionFunc(
		constants.ArgCols,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return allCols, cobra.ShellCompDirectiveNoFileComp
		},
	)
	cmd.Command.MarkFlagsMutuallyExclusive(FlagExpiryDate, FlagTimeUntilExpiry)
	return cmd
}

func PreCmdPostToken(c *core.PreCommandConfig) error {
	err := core.CheckRequiredFlags(c.Command, c.NS, FlagName, FlagRegId)
	if err != nil {
		return err
	}

	return nil
}

func CmdPostToken(c *core.CommandConfig) error {
	var err error
	id, err := c.Command.Command.Flags().GetString("registry-id")
	if err != nil {
		return err
	}
	name, err := c.Command.Command.Flags().GetString(FlagName)
	if err != nil {
		return err
	}
	tokenPostProperties.SetName(name)

	if viper.IsSet(core.GetFlagName(c.NS, FlagExpiryDate)) {
		var expiryDate time.Time
		expiryDateString, err := c.Command.Command.Flags().GetString(FlagExpiryDate)
		if err != nil {
			return err
		}
		expiryDate, err = time.Parse(time.RFC3339, expiryDateString)
		if err != nil {
			return err
		}
		tokenPostProperties.SetExpiryDate(expiryDate)

	} else if viper.IsSet(core.GetFlagName(c.NS, FlagTimeUntilExpiry)) {
		var timeUntilExpiry string
		timeUntilExpiry, err = c.Command.Command.Flags().GetString(FlagTimeUntilExpiry)
		if err != nil {
			return err
		}
		timeNow := time.Now()
		duration, err := ParseExpiryTime(timeUntilExpiry)
		if err != nil {
			return err
		}
		timeNow.Add(duration)
	}

	if viper.IsSet(core.GetFlagName(c.NS, FlagStatus)) {
		var status string
		status, err = c.Command.Command.Flags().GetString(FlagStatus)
		if err != nil {
			return err
		}
		tokenPostProperties.SetStatus(status)
	}

	tokenInput := sdkgo.NewPostTokenInputWithDefaults()
	tokenInput.SetProperties(*tokenPostProperties)

	token, _, err := c.ContainerRegistryServices.Token().Post(*tokenInput, id)
	if err != nil {
		return err
	}

	tokenPrint := sdkgo.NewTokenResponseWithDefaults()
	tokenPrint.SetProperties(*token.GetProperties())

	return c.Printer.Print(getTokenPrint(nil, c, &[]sdkgo.TokenResponse{*tokenPrint}, true))
}

func ParseExpiryTime(expiryTime string) (time.Duration, error) {
	years := 0
	months := 0
	days := 0
	hours := 0

	if !strings.ContainsAny(expiryTime, "0123456789") {
		return 0, fmt.Errorf("invalid expiry time format")
	}

	number := ""

	for i := 0; i < len(expiryTime); i++ {
		if string(expiryTime[i]) != "y" && string(expiryTime[i]) != "m" && string(expiryTime[i]) != "d" && string(expiryTime[i]) != "h" {
			number += string(expiryTime[i])
		} else if expiryTime[i] == 'y' {
			years, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'm' {
			months, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'd' {
			days, _ = strconv.Atoi(number)
			number = ""
		} else if expiryTime[i] == 'h' {
			hours, _ = strconv.Atoi(number)
			number = ""
		}
	}

	return time.Duration(years)*time.Hour*24*365 + time.Duration(months)*time.Hour*24*30 + time.Duration(days)*time.Hour*24 + time.Duration(hours)*time.Hour, nil
}
