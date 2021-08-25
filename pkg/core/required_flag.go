package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/spf13/viper"
	multierror "go.uber.org/multierr"
)

func CheckRequiredGlobalFlags(cmdName string, globalFlagsName ...string) error {
	var multiErr error
	for _, flagName := range globalFlagsName {
		if !viper.IsSet(GetGlobalFlagName(cmdName, flagName)) {
			multiErr = multierror.Append(multiErr, clierror.NewRequiredFlagErr(flagName))
		}
	}
	if multiErr != nil {
		return multiErr
	}
	return nil
}

func CheckRequiredFlags(ns string, localFlagsName ...string) error {
	var multiErr error
	for _, flagName := range localFlagsName {
		if !viper.IsSet(GetFlagName(ns, flagName)) {
			multiErr = multierror.Append(multiErr, clierror.NewRequiredFlagErr(flagName))
		}
	}
	if multiErr != nil {
		return multiErr
	}
	return nil
}

func CheckRequiredFlagsWe(cmd *Command, ns string, localFlagsName ...string) error {
	for _, flagName := range localFlagsName {
		if !viper.IsSet(GetFlagName(ns, flagName)) {
			return RequiresMinOptionsErr(cmd, len(localFlagsName))
		}
	}
	return nil
}

const RequiredFlagsAnnotation = "RequiredFlags"

type FlagOptionFunc func(c *Command, flagName string)

func RequiredFlagOption() FlagOptionFunc {
	return func(c *Command, flagName string) {
		usg := c.Command.Flag(flagName).Usage
		c.Command.Flag(flagName).Usage = fmt.Sprintf("%s %s", usg, "(required)")
		c.Command.Annotations = map[string]string{RequiredFlagsAnnotation: fmt.Sprintf("%s --%s %s",
			c.Command.Annotations[RequiredFlagsAnnotation],
			flagName,
			strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))}
	}
}

func RequiresMinOptionsErr(cmd *Command, min int) error {
	return errors.New(
		fmt.Sprintf("%q requires at least %d %s.\n\nUsage: %s%s\n\nFor more details, see '%s --help'.",
			cmd.CommandPath(),
			min,
			pluralize("option", min),
			cmd.CommandPath(),
			cmd.Command.Annotations[RequiredFlagsAnnotation],
			cmd.CommandPath(),
		),
	)
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
