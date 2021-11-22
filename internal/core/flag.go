package core

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const (
	RequiredFlagsAnnotation   = "RequiredFlags"
	DeprecatedFlagsAnnotation = "DeprecatedFlags"
)

var (
	flagNamePrintF  = "%s --%s %s"
	requiredFlagErr = errors.New("error checking required flags on command")
)

type FlagOptionFunc func(cmd *Command, flagName string)

func DeprecatedFlagOption() FlagOptionFunc {
	return func(cmd *Command, flagName string) {
		cmd.Command.Flag(flagName).Usage = fmt.Sprintf("%s (deprecated)", cmd.Command.Flag(flagName).Usage)
		// For documentation purposes, add flag to command Annotation
		cmd.Command.Annotations = map[string]string{DeprecatedFlagsAnnotation: fmt.Sprintf(flagNamePrintF,
			cmd.Command.Annotations[DeprecatedFlagsAnnotation],
			flagName,
			strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))}
	}
}

func RequiredFlagOption() FlagOptionFunc {
	return func(cmd *Command, flagName string) {
		cmd.Command.Flag(flagName).Usage = fmt.Sprintf("%s (required)", cmd.Command.Flag(flagName).Usage)
		// For documentation purposes, add flag to command Annotation
		cmd.Command.Annotations = map[string]string{RequiredFlagsAnnotation: fmt.Sprintf(flagNamePrintF,
			cmd.Command.Annotations[RequiredFlagsAnnotation],
			flagName,
			strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))}
	}
}

func RequiresMinOptionsErr(cmd *Command, flagNames ...string) error {
	if cmd == nil {
		return requiredFlagErr
	}
	var usage string
	usage = cmd.CommandPath()
	for _, flagName := range flagNames {
		usage = fmt.Sprintf(flagNamePrintF, usage, flagName,
			strings.ReplaceAll(strings.ToUpper(flagName), "-", "_"),
		)
	}
	return errors.New(
		fmt.Sprintf("%q requires at least %d %s.\n\nUsage: %s\n\nFor more details, see '%s --help'.",
			cmd.CommandPath(),
			len(flagNames),
			pluralize("option", len(flagNames)),
			usage,
			cmd.CommandPath(),
		),
	)
}

func RequiresMultipleOptionsErr(cmd *Command, flagNamesSets ...[]string) error {
	if cmd == nil {
		return requiredFlagErr
	}
	var usage string
	for _, flagNamesSet := range flagNamesSets {
		usage = fmt.Sprintf("%s%s", usage, cmd.CommandPath())
		for _, flagName := range flagNamesSet {
			usage = fmt.Sprintf(flagNamePrintF, usage, flagName,
				strings.ReplaceAll(strings.ToUpper(flagName), "-", "_"),
			)
		}
		usage = fmt.Sprintf("%s\n", usage)
	}
	return errors.New(
		fmt.Sprintf("%q requires at least %d %s.\n\nUsage:\n%s\nFor more details, see '%s --help'.",
			cmd.CommandPath(),
			minLen(flagNamesSets...),
			pluralize("option", minLen(flagNamesSets...)),
			usage,
			cmd.CommandPath(),
		),
	)
}

func CheckRequiredFlags(cmd *Command, ns string, localFlagsName ...string) error {
	for _, flagName := range localFlagsName {
		if !viper.IsSet(GetFlagName(ns, flagName)) {
			return RequiresMinOptionsErr(cmd, localFlagsName...)
		}
	}
	return nil
}

// CheckRequiredFlagsSets focuses on commands that support multiple ways to run,
// and it is required at least one of the ways to be set.
// For example: for `datacenter delete`, it is required to be set
// either `--datacenter-id` option, either `--all` option.
func CheckRequiredFlagsSets(cmd *Command, ns string, localFlagsNameSets ...[]string) error {
	checkSet := false
	for _, flagNameSet := range localFlagsNameSets {
		err := CheckRequiredFlags(cmd, ns, flagNameSet...)
		if err == nil {
			checkSet = true
		}
	}
	// If one of the flags set provided is set, return nil.
	// If none of the flags sets are not set, return error message.
	if checkSet == true {
		return nil
	} else {
		return RequiresMultipleOptionsErr(cmd, localFlagsNameSets...)
	}
}

// minLen gets the minimum length of the arrays provided as input
func minLen(sets ...[]string) int {
	var min int
	if len(sets) > 0 {
		min = len(sets[0])
	}
	for _, set := range sets {
		if len(set) < min {
			min = len(set)
		}
	}
	return min
}

func pluralize(word string, number int) string {
	if number == 1 {
		return word
	}
	return word + "s"
}
