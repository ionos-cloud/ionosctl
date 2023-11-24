package commands

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionoscloudsdk/comptplus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var advancedPrompt = &comptplus.CobraPrompt{
	RootCmd:                  rootCmd.Command,
	ShowHelpCommandAndFlags:  true,
	DisableCompletionCommand: true,
	AddDefaultExitCommand:    true,
	ShowHiddenCommands:       true,
	ShowHiddenFlags:          true,
	GoPromptOptions: []prompt.Option{
		prompt.OptionTitle("ionosctl"),
		prompt.OptionPrefix("> "),
		prompt.OptionShowCompletionAtStart(),

		prompt.OptionDescriptionTextColor(prompt.Black),
		prompt.OptionSuggestionTextColor(prompt.White),
		prompt.OptionDescriptionBGColor(prompt.LightGray),
		prompt.OptionSuggestionBGColor(prompt.DarkGray),

		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),

		prompt.OptionPreviewSuggestionTextColor(prompt.DarkBlue),
		prompt.OptionPrefixTextColor(prompt.DefaultColor),
		prompt.OptionScrollbarThumbColor(prompt.DarkGray),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
	},

	OnErrorFunc: func(err error) {
		// rootCmd.Command.PrintErr(err)
		return
	},

	HookBefore: func(cmd *cobra.Command, input string) {
		forceExists := false
		// TODO: Why forceExists is always true, even for commands that don't have --force?
		cmd.Flags().VisitAll(func(flag *pflag.Flag) {
			if flag.Name == constants.ArgForce {
				forceExists = true
			}
		})

		if forceSet, err := cmd.Flags().GetBool(constants.ArgForce); forceExists && (!forceSet || err != nil) {
			fmt.Println("Warning: '--force' needs to be set for this command! Interactive shell does not support user inputs. Repeat this command with --force to continue.")
			cmd = rootCmd.Command
			cmd.SetArgs([]string{"helpgi"})
		}
	},

	CustomFlagResetBehaviour: func(flag *pflag.Flag) {
		sliceValue, ok := flag.Value.(pflag.SliceValue)
		if !ok {
			// For non-slice flags, just set to the default value
			flag.Value.Set(flag.DefValue)
			return
		}

		err := sliceValue.Replace([]string{})
		if err != nil {
			flag.Value.Set(flag.DefValue)
		}
	},
}

func Shell() *core.Command {
	flagPersistFlagValues := "persist-flag-values"

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "shell",
		Resource:  "shell",
		Verb:      "shell",
		ShortDesc: "Interactive shell - BETA",
		LongDesc: `The ionosctl shell command launches an interactive shell environment, enabling a more dynamic and intuitive way to interact with the ionosctl CLI.
This shell is designed to enhance your command-line experience with advanced features and customizations, powered by the comptplus library.

DEFAULT CONTROLS:
Ctrl + A\tGo to the beginning of the line (Home)
Ctrl + E\tGo to the end of the line (End)
Ctrl + P\tPrevious command (Up arrow)
Ctrl + N\tNext command (Down arrow)
Ctrl + F\tForward one character
Ctrl + B\tBackward one character
Ctrl + D\tDelete character under the cursor
Ctrl + H\tDelete character before the cursor (Backspace)
Ctrl + W\tCut the word before the cursor to the clipboard
Ctrl + K\tCut the line after the cursor to the clipboard
Ctrl + U\tCut the line before the cursor to the clipboard
Ctrl + L\tClear the screen`,
		Example: "ionosctl shell",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			_, err := client.Get()
			if err != nil {
				return fmt.Errorf("usage of the interactive shell requires valid credentials. "+
					"You can use `ionosctl whoami` to debug your configuration: %w", err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			fmt.Printf("ionosctl v%s\n", Version)
			fmt.Println("Warning: We recommend keeping usage of this interactive shell to non-production critical applications.")
			fmt.Println("   - DANGER: '--force' will always be set! Deleting resources will not ask for confirmation!")
			fmt.Println("   - DANGER: Certain commands that require user input may freeze the shell!")
			fmt.Println("   - This is a BETA feature. Please report any bugs to github.com/ionos-cloud/ionosctl/issues/new/choose")
			advancedPrompt.PersistFlagValues = viper.GetBool(flagPersistFlagValues)
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	cmd.AddBoolFlag(flagPersistFlagValues, "p", false, "Persist flag values between commands")

	return cmd
}
