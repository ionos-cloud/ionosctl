package commands

import (
	"context"
	"fmt"

	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionoscloudsdk/comptplus"
)

var advancedPrompt = &comptplus.CobraPrompt{
	RootCmd:                  rootCmd.Command,
	PersistFlagValues:        true, // Adds flag which allows persisting flag values between commands
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
		prompt.OptionSuggestionBGColor(prompt.DarkBlue),
		prompt.OptionDescriptionBGColor(prompt.Blue),

		prompt.OptionSelectedDescriptionTextColor(prompt.White),
		prompt.OptionSelectedSuggestionTextColor(prompt.Black),
		prompt.OptionSelectedSuggestionBGColor(prompt.Blue),
		prompt.OptionSelectedDescriptionBGColor(prompt.DarkBlue),

		prompt.OptionScrollbarThumbColor(prompt.LightGray),
		prompt.OptionScrollbarBGColor(prompt.DefaultColor),
	},
	HookBefore: func(_ string) {
		// initConfig()
	},
	HookAfter: func(_ string) {
	},

	OnErrorFunc: func(err error) {
		rootCmd.Command.PrintErr(err)
		return
	},
}

func Shell() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "shell",
		Resource:  "shell",
		Verb:      "shell",
		ShortDesc: "Interactive shell - BETA",
		Example:   "ionosctl shell",
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
			fmt.Println("Warning: This interactive shell is a BETA feature and may not work as expected.")
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	return cmd
}
