package commands

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionoscloudsdk/comptplus"
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

		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ShiftLeft, Fn: prompt.GoLeftWord}),
		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ShiftRight, Fn: prompt.GoRightWord}),
		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ShiftDown, Fn: prompt.GoLineBeginning}),
		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ShiftUp, Fn: prompt.GoLineEnd}),
		prompt.OptionAddKeyBind(prompt.KeyBind{Key: prompt.ShiftDelete, Fn: prompt.DeleteWord}),

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
	CustomFlagResetBehaviour: func(flag *pflag.Flag) {
		if flag.Name == "cols" {
			flag.Value.Set("")
			return
		}
		flag.Value.Set(flag.DefValue)
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

CUSTOM CONTROLS: (your usual shell controls might not work)
- SHIFT + LEFT/RIGHT: Quickly navigate words left/right
- SHIFT + UP/DOWN: Quickly navigate to the beginning/end of the line
- SHIFT + DELETE: Delete previous word`,
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
			fmt.Println("Warning: This interactive shell is a BETA feature. We recommend keeping usage and testing to non-production critical applications.")
			advancedPrompt.PersistFlagValues = viper.GetBool(flagPersistFlagValues)
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	cmd.AddBoolFlag(flagPersistFlagValues, "p", false, "Persist flag values between commands")

	return cmd
}
