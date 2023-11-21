package comptplus

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// DynamicSuggestionsAnnotation for dynamic suggestions.
const DynamicSuggestionsAnnotation = "cobra-prompt-dynamic-suggestions"

// PersistFlagValuesFlag the flag that will be available when PersistFlagValues is true
const PersistFlagValuesFlag = "persist-flag-values"

const CacheIntervalFlag = "cache-interval"

// CobraPrompt given a Cobra command it will make every flag and sub commands available as suggestions.
// Command.Short will be used as description for the suggestion.
type CobraPrompt struct {
	// RootCmd is the start point, all its sub commands and flags will be available as suggestions
	RootCmd *cobra.Command

	// GoPromptOptions is for customize go-prompt
	// see https://github.com/c-bata/go-prompt/blob/master/option.go
	GoPromptOptions []prompt.Option

	// DynamicSuggestionsFunc will be executed if a command has CallbackAnnotation as an annotation. If it's included
	// the value will be provided to the DynamicSuggestionsFunc function.
	DynamicSuggestionsFunc func(annotationValue string, document *prompt.Document) []prompt.Suggest

	// PersistFlagValues will persist flags. For example have verbose turned on every command.
	PersistFlagValues bool

	// ShowHelpCommandAndFlags will make help command and flag for every command available.
	ShowHelpCommandAndFlags bool

	// DisableCompletionCommand will disable the default completion command for cobra
	DisableCompletionCommand bool

	// ShowHiddenCommands makes hidden commands available
	ShowHiddenCommands bool

	// ShowHiddenFlags makes hidden flags available
	ShowHiddenFlags bool

	// AddDefaultExitCommand adds a command for exiting prompt loop
	AddDefaultExitCommand bool

	// OnErrorFunc handle error for command.Execute, if not set print error and exit
	OnErrorFunc func(err error)

	// HookAfter is a hook that will be executed every time after a command has been executed
	HookAfter func(input string)

	// HookBefore is a hook that will be executed every time  a command has been executed
	HookBefore func(input string)

	// InArgsParser adds a custom parser for the command line arguments (default: strings.Fields)
	InArgsParser func(args string) []string

	// SuggestionFilter will be uses when filtering suggestions as typing
	SuggestionFilter func(suggestions []prompt.Suggest, document *prompt.Document) []prompt.Suggest

	lastFlagValueSuggestionsTime time.Time

	lastFlagValueSuggestions []prompt.Suggest
}

// Run will automatically generate suggestions for all cobra commands and flags defined by RootCmd
// and execute the selected commands. Run will also reset all given flags by default, see PersistFlagValues
func (co *CobraPrompt) Run() {
	co.RunContext(context.Background())
}

func (co *CobraPrompt) RunContext(ctx context.Context) {
	if co.RootCmd == nil {
		panic("RootCmd is not set. Please set RootCmd")
	}
	co.prepareCommands()

	p := prompt.New(
		co.executeCommand(ctx),
		co.findSuggestions,
		co.GoPromptOptions...,
	)
	p.Run()
}

func (co *CobraPrompt) executeCommand(ctx context.Context) func(string) {
	return func(input string) {
		co.HookBefore(input)
		args := co.parseInput(input)
		os.Args = append([]string{os.Args[0]}, args...)
		if err := co.RootCmd.ExecuteContext(ctx); err != nil {
			if co.OnErrorFunc != nil {
				co.OnErrorFunc(err)
			} else {
				co.RootCmd.PrintErrln(err)
				os.Exit(1)
			}
		}
		co.HookAfter(input)
	}
}

func (co *CobraPrompt) parseInput(input string) []string {
	if co.InArgsParser != nil {
		return co.InArgsParser(input)
	}
	return strings.Fields(input)
}

func (co *CobraPrompt) prepareCommands() {
	if co.ShowHelpCommandAndFlags {
		co.RootCmd.InitDefaultHelpCmd()
	}
	if co.DisableCompletionCommand {
		co.RootCmd.CompletionOptions.DisableDefaultCmd = true
	}
	if co.AddDefaultExitCommand {
		co.RootCmd.AddCommand(&cobra.Command{
			Use:   "exit",
			Short: "Exit prompt",
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(0)
			},
		})
	}
	if co.PersistFlagValues {
		co.RootCmd.PersistentFlags().BoolP(PersistFlagValuesFlag, "", false, "Persist flag values")
	}
}

// findSuggestions generates command and flag suggestions for the prompt.
func (co *CobraPrompt) findSuggestions(d prompt.Document) []prompt.Suggest {
	command := co.RootCmd
	args := strings.Fields(d.CurrentLine())

	if found, _, err := command.Find(args); err == nil {
		command = found
	}

	interval, err := command.Flags().GetDuration(CacheIntervalFlag)
	if err != nil || interval == 0 {
		interval = 500 * time.Millisecond
	}

	var suggestions []prompt.Suggest
	currentFlag, isFlagValueContext := getCurrentFlagAndValueContext(d, command)

	if !isFlagValueContext {
		suggestions = append(suggestions, getFlagSuggestions(command, co, d)...)
		suggestions = append(suggestions, getCommandSuggestions(command, co)...)
		suggestions = append(suggestions, getDynamicSuggestions(command, co, d)...)
	} else {
		suggestions = co.lastFlagValueSuggestions
		if time.Since(co.lastFlagValueSuggestionsTime) > interval {
			suggestions = getFlagValueSuggestions(command, d, currentFlag)
			co.lastFlagValueSuggestions = suggestions
			co.lastFlagValueSuggestionsTime = time.Now()
		}
	}

	if co.SuggestionFilter != nil {
		return co.SuggestionFilter(suggestions, &d)
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

// getFlagSuggestions returns a slice of flag suggestions.
func getFlagSuggestions(cmd *cobra.Command, co *CobraPrompt, d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest
	persistFlagValues, _ := cmd.Flags().GetBool(PersistFlagValuesFlag)

	addFlags := func(flag *pflag.Flag) {
		if len(d.Text) < 1 && flag.Changed && !persistFlagValues {
			flag.Value.Set(flag.DefValue)
		}
		if flag.Hidden && !co.ShowHiddenFlags {
			return
		}

		if strings.HasPrefix(d.GetWordBeforeCursor(), "--") {
			suggestions = append(suggestions, prompt.Suggest{Text: "--" + flag.Name, Description: flag.Usage})
		} else if strings.HasPrefix(d.GetWordBeforeCursor(), "-") && flag.Shorthand != "" {
			suggestions = append(suggestions, prompt.Suggest{Text: "-" + flag.Shorthand, Description: flag.Usage})
		}
	}

	cmd.LocalFlags().VisitAll(addFlags)
	cmd.InheritedFlags().VisitAll(addFlags)
	return suggestions
}

// getCommandSuggestions returns a slice of command suggestions.
func getCommandSuggestions(cmd *cobra.Command, co *CobraPrompt) []prompt.Suggest {
	var suggestions []prompt.Suggest
	if cmd.HasAvailableSubCommands() {
		for _, c := range cmd.Commands() {
			if !c.Hidden || co.ShowHiddenCommands {
				suggestions = append(suggestions, prompt.Suggest{Text: c.Name(), Description: c.Short})
			}
		}
	}
	return suggestions
}

// getDynamicSuggestions returns a slice of dynamic arg completions.
func getDynamicSuggestions(cmd *cobra.Command, co *CobraPrompt, d prompt.Document) []prompt.Suggest {
	var suggestions []prompt.Suggest
	if dynamicSuggestionKey, ok := cmd.Annotations[DynamicSuggestionsAnnotation]; ok {
		if co.DynamicSuggestionsFunc != nil {
			dynamicSuggestions := co.DynamicSuggestionsFunc(dynamicSuggestionKey, &d)
			suggestions = append(suggestions, dynamicSuggestions...)
		}
	}
	return suggestions
}

// getFlagValueSuggestions returns a slice of flag value suggestions.
func getFlagValueSuggestions(cmd *cobra.Command, d prompt.Document, currentFlag string) []prompt.Suggest {
	var suggestions []prompt.Suggest

	// Check if the current flag is boolean. If so, do not suggest values.
	if flag := cmd.Flags().Lookup(currentFlag); flag != nil && flag.Value.Type() == "bool" {
		return suggestions
	}

	if compFunc, exists := cmd.GetFlagCompletionFunc(currentFlag); exists {
		completions, _ := compFunc(cmd, strings.Fields(d.CurrentLine()), currentFlag)
		for _, completion := range completions {
			text, description, _ := strings.Cut(completion, "\t")
			suggestions = append(suggestions, prompt.Suggest{Text: text, Description: description})
		}
	}
	return suggestions
}

// getCurrentFlagAndValueContext parses the document to find:
//   - current flag
//   - partial value
//   - and whether the context is suitable for flag value suggestions.
func getCurrentFlagAndValueContext(d prompt.Document, cmd *cobra.Command) (string, bool) {
	prevWords := strings.Fields(d.TextBeforeCursor())
	textBeforeCursor := d.TextBeforeCursor()
	hasSpaceSuffix := strings.HasSuffix(textBeforeCursor, " ")

	lastWord := ""
	secondLastWord := ""
	if len(prevWords) > 0 {
		lastWord = prevWords[len(prevWords)-1]
		if len(prevWords) > 1 {
			secondLastWord = prevWords[len(prevWords)-2]
		}
	}

	// Case where the last word is a partial value -- second last word is a flag (non-bool)
	if !hasSpaceSuffix && strings.HasPrefix(secondLastWord, "-") {
		flagName := getFlagNameFromArg(secondLastWord, cmd)
		if flag := cmd.Flags().Lookup(flagName); flag != nil && flag.Value.Type() != "bool" {
			return flagName, true
		}
	}

	// Done with writing a flag (`--arg `) -> appropriate context
	if hasSpaceSuffix && len(lastWord) > 0 && strings.HasPrefix(lastWord, "-") {
		return getFlagNameFromArg(lastWord, cmd), true
	}

	// Not done typing a flag -> not appropriate context
	if !hasSpaceSuffix && len(lastWord) > 0 && !strings.HasPrefix(lastWord, "-") {
		return "",  false
	}

	// Done with writing a flag value (`--arg MyArg `) -> not appropriate context
	if hasSpaceSuffix && len(secondLastWord) > 0 && strings.HasPrefix(secondLastWord, "-") {
		return "", false
	}

	return "", false
}

// getFlagNameFromArg extracts the flag name from a given argument, handling both shorthand and full flag names.
func getFlagNameFromArg(arg string, cmd *cobra.Command) string {
	trimmedArg := strings.TrimLeft(arg, "-")
	if len(trimmedArg) == 1 { // Shorthand flag
		if shorthandFlag := cmd.Flags().ShorthandLookup(trimmedArg); shorthandFlag != nil {
			return shorthandFlag.Name
		}
	} else { // Full flag name
		if fullFlag := cmd.Flags().Lookup(trimmedArg); fullFlag != nil {
			return fullFlag.Name
		}
	}
	return ""
}
