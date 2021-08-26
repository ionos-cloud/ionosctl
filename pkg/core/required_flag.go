package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const RequiredFlagsAnnotation = "RequiredFlags"

type FlagOptionFunc func(c *Command, flagName string)

func RequiredFlagOption() FlagOptionFunc {
	return func(c *Command, flagName string) {
		c.Command.Flag(flagName).Usage = fmt.Sprintf("%s (required)", c.Command.Flag(flagName).Usage)
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

func CheckRequiredGlobalFlags(cmd *Command, cmdName string, globalFlagsName ...string) error {
	for _, flagName := range globalFlagsName {
		if !viper.IsSet(GetGlobalFlagName(cmdName, flagName)) {
			return RequiresMinOptionsErr(cmd, len(globalFlagsName))
		}
	}
	return nil
}

func CheckRequiredFlags(cmd *Command, ns string, localFlagsName ...string) error {
	for _, flagName := range localFlagsName {
		if !viper.IsSet(GetFlagName(ns, flagName)) {
			return RequiresMinOptionsErr(cmd, len(localFlagsName))
		}
	}
	return nil
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
