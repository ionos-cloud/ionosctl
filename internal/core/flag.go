package core

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/spf13/viper"
)

const (
	RequiredFlagsAnnotation   = "RequiredFlags"
	DeprecatedFlagsAnnotation = "DeprecatedFlags"
)

var (
	flagNamePrintF     = "%s --%s %s"
	flagNameBoolPrintF = "%s --%s"
	requiredFlagErr    = errors.New("error checking required flags on command")
)

type FlagOptionFunc func(cmd *Command, flagName string)

func DeprecatedFlagOption(help string) FlagOptionFunc {
	return func(cmd *Command, flagName string) {
		cmd.Command.Flag(flagName).Deprecated = help
		// For documentation purposes, add flag to command Annotation
		if len(cmd.Command.Annotations) > 0 {
			cmd.Command.Annotations[DeprecatedFlagsAnnotation] = fmt.Sprintf(flagNamePrintF, cmd.Command.Annotations[DeprecatedFlagsAnnotation], flagName, strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))
		} else {
			cmd.Command.Annotations = map[string]string{DeprecatedFlagsAnnotation: fmt.Sprintf(flagNamePrintF, cmd.Command.Annotations[DeprecatedFlagsAnnotation], flagName, strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))}
		}
	}
}

func RequiredFlagOption() FlagOptionFunc {
	return func(cmd *Command, flagName string) {
		cmd.Command.Flag(flagName).Usage = fmt.Sprintf("%s (required)", cmd.Command.Flag(flagName).Usage)
		// For documentation purposes, add flag to command Annotation
		if len(cmd.Command.Annotations) > 0 {
			cmd.Command.Annotations[RequiredFlagsAnnotation] = fmt.Sprintf(flagNamePrintF, cmd.Command.Annotations[RequiredFlagsAnnotation], flagName, strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))
		} else {
			cmd.Command.Annotations = map[string]string{RequiredFlagsAnnotation: fmt.Sprintf(flagNamePrintF, cmd.Command.Annotations[RequiredFlagsAnnotation], flagName, strings.ToUpper(strings.ReplaceAll(flagName, "-", "_")))}
		}
	}
}

// FlagAsVariable takes a flag name and returns it as a screaming camel case
//
// e.g. `FlagAsVariable("datacenter-id") -> "DATACENTER_ID"
func FlagAsVariable(flag string) string {
	return strings.ReplaceAll(strings.ToUpper(flag), "-", "_")
}

// FlagUsage ("datacenter-id") -> "--datacenter-id DATACENTER_ID"
//
// Used as a convenience func
func FlagUsage(flag string) string {
	return fmt.Sprintf("--%s %s", flag, FlagAsVariable(flag))
}

// FlagsUsage calls FlagUsage for every flag in the slice
func FlagsUsage(flags ...string) string {
	usage := ""
	for _, flagName := range flags {
		usage += FlagUsage(flagName) + " "
	}
	return usage
}

func RequiresMinOptionsErr(cmd *Command, flagNames ...string) error {
	if cmd == nil || cmd.Command == nil {
		return requiredFlagErr
	}
	usage := fmt.Sprintf("%s %s", cmd.CommandPath(), FlagsUsage(flagNames...))
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
	if cmd == nil || cmd.Command == nil {
		return requiredFlagErr
	}
	var usage string
	for _, flagNamesSet := range flagNamesSets {
		usage = fmt.Sprintf("%s%s", usage, cmd.CommandPath())
		for _, flagName := range flagNamesSet {
			if cmd.Command.Flag(flagName) != nil && cmd.Command.Flag(flagName).Value.Type() == "bool" {
				usage = fmt.Sprintf(flagNameBoolPrintF, usage, flagName)
			} else {
				usage = fmt.Sprintf(flagNamePrintF, usage, flagName,
					strings.ReplaceAll(strings.ToUpper(flagName), "-", "_"),
				)
			}
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
		if !viper.IsSet(GetFlagName(ns, flagName)) && !viper.IsSet(flagName) {
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

type FlagNameSetWithPredicate struct {
	FlagNameSet    []string
	Predicate      func(interface{}) bool
	PredicateParam interface{}
}

// If a flag being set to a certain value creates some extra flag dependencies, then use this function!
func CheckRequiredFlagsSetsIfPredicate(cmd *Command, ns string, localFlagsNameSets ...FlagNameSetWithPredicate) error {
	anyFlagSet := false
	flagSetsValidPredicate := [][]string{}
	for _, flagNameSet := range localFlagsNameSets {
		if !flagNameSet.Predicate(flagNameSet.PredicateParam) {
			continue
		}
		err := CheckRequiredFlags(cmd, ns, flagNameSet.FlagNameSet...)
		flagSetsValidPredicate = append(flagSetsValidPredicate, flagNameSet.FlagNameSet)
		if err == nil {
			anyFlagSet = true
		}
	}
	// If none of the flags sets are set, return error message.
	if !anyFlagSet {
		return RequiresMultipleOptionsErr(cmd, flagSetsValidPredicate...)
	}

	return nil
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

// //
// --- CUSTOM FLAG TYPES ---
// For custom validation and error handling within pflag's Set function
// Use pflag's Var and VarP respectively in conjunction with the custom flag's constructor to add these custom types to a command.
//

// / -- START UUID FLAG TYPE --
type uuidFlag struct {
	Value string
}

// Instantiate an empty uuidFlag type. Use this in pflag's Var/VarP funcs first argument
func newUuidFlag(defaultValue string) *uuidFlag {
	return &uuidFlag{Value: defaultValue}
}

// PFlag calls this function when it finds an argument provided by the user of uuidFlag type.
func (u *uuidFlag) Set(p string) error {
	IsValidUUID := func(u string) bool {
		_, err := uuid.FromString(u)
		return err == nil
	}

	if !IsValidUUID(p) {
		// return fmt.Errorf("%s does not match UUID-4 format", p)
		fmt.Fprintf(os.Stderr, jsontabwriter.GenerateLogOutput(fmt.Sprintf("WARNING: %s does not match UUID-4 format", p)))
	}

	// Valid UUID if passed above check
	u.Value = p
	return nil
}

func (u *uuidFlag) Type() string {
	return "string"
}

func (u uuidFlag) String() string {
	return u.Value
}

// / -- END UUID FLAG TYPE --

// SetFlag /
// Values set for this flag must be part of allowed values
// NOTE: Track progress of https://github.com/spf13/pflag/issues/236 : Might be implemented in pflag
type SetFlag struct {
	Value   string
	Allowed []string
}

func newSetFlag(defaultValue string, Allowed []string) *SetFlag {
	return &SetFlag{
		Value:   defaultValue,
		Allowed: Allowed,
	}
}

func (a *SetFlag) Set(p string) error {
	isIncluded := func(opts []string, val string) bool {
		for _, opt := range opts {
			if val == opt {
				return true
			}
		}
		return false
	}
	if !isIncluded(a.Allowed, p) {
		// Error: invalid argument "las" for "-l, --location" flag: please pick one of these values: fra, fkb, txl, lhr, las, ewr, vit
		return fmt.Errorf(
			"please pick one of these values: %s",
			strings.Join(a.Allowed, ","),
		)
	}
	a.Value = p
	return nil
}

func (a *SetFlag) Type() string {
	return "string"
}

func (a SetFlag) String() string {
	return a.Value
}
