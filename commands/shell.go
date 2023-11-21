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
	GoPromptOptions: []prompt.Option{
		prompt.OptionTitle("ionosctl"),
		prompt.OptionPrefix("> "),
		prompt.OptionShowCompletionAtStart(),
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
