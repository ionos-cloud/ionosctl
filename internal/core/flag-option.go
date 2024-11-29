package core

import (
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	RequiredFlagsAnnotation   = "RequiredFlags"
	DeprecatedFlagsAnnotation = "DeprecatedFlags"
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

// WithCompletionComplex is a FlagOptionFunc that allows for more complex completion logic.
// It is recommended to use one of the simpler helper functions WithCompletion or WithCompletionE instead.
// Only use this complex function if you need to handle more complex logic, like args-based completion, or custom filtering based on already typed keys (from toComplete)
func WithCompletionComplex(
	completionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective),
	baseURL string,
	allowedLocations []string,
) FlagOptionFunc {
	// Pre-generate the map only if allowedLocations is provided
	var locationToURL map[string]string
	if allowedLocations != nil {
		locationToURL = make(map[string]string, len(allowedLocations))
		for _, loc := range allowedLocations {
			normalizedLoc := strings.ReplaceAll(loc, "/", "-") // Replace `/` with `-`
			locationToURL[normalizedLoc] = fmt.Sprintf(baseURL, normalizedLoc)
		}
	}

	return func(cmdToRegister *Command, flagName string) {
		cmdToRegister.Command.RegisterFlagCompletionFunc(flagName,
			func(passedCmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
				viper.AutomaticEnv()

				// Check if ArgServerURL is already set manually
				if viper.IsSet(constants.ArgServerUrl) || viper.IsSet(constants.EnvServerUrl) {
					// If manually set, do nothing and directly call completionFunc
					return completionFunc(passedCmd, args, toComplete)
				}

				// Handle location-based logic if allowedLocations is provided
				if locationToURL != nil {
					if location, _ := passedCmd.Flags().GetString(constants.FlagLocation); location != "" {
						if url, ok := locationToURL[location]; ok {
							viper.Set(constants.ArgServerUrl, url)
						} else {
							// Return an error directive if location is invalid
							return nil, cobra.ShellCompDirectiveError
						}
					}
				} else {
					// Use the baseURL directly if no locations are provided
					viper.Set(constants.ArgServerUrl, baseURL)
				}

				return completionFunc(passedCmd, args, toComplete)
			},
		)
	}
}

// WithCompletionE is a FlagOptionFunc that allows for a completion function that can return an error.
//
// Usage:
//
// - WithCompletionE(completionFuncE, "api.%s.ionos.com", allowedLocations) for a regional API
//
// - WithCompletionE(completionFuncE, "api.ionos.com", nil) for an API with a single endpoint
func WithCompletionE(completionFunc func() ([]string, error), baseURL string, allowedLocations []string) FlagOptionFunc {
	return WithCompletionComplex(func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		results, err := completionFunc()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		return results, cobra.ShellCompDirectiveNoFileComp
	}, baseURL, allowedLocations)
}

// WithCompletion is a FlagOptionFunc that allows for a completion function that returns a list of strings.
//
// Usage:
//
// - WithCompletion(completionFunc, "api.%s.ionos.com", allowedLocations) for a regional API
//
// - WithCompletion(completionFunc, "api.ionos.com", nil) for an API with a single endpoint
func WithCompletion(completionFunc func() []string, baseURL string, allowedLocations []string) FlagOptionFunc {
	return WithCompletionComplex(func(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
		return completionFunc(), cobra.ShellCompDirectiveNoFileComp
	}, baseURL, allowedLocations)
}
