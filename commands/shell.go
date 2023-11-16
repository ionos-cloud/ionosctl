package commands

import (
	"context"
	"fmt"

	"github.com/avirtopeanu-ionos/comptplus"
	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/spf13/viper"
)

var advancedPrompt = &comptplus.CobraPrompt{
	RootCmd:                  rootCmd.Command,
	PersistFlagValues:        true,
	ShowHelpCommandAndFlags:  true,
	DisableCompletionCommand: true,
	AddDefaultExitCommand:    true,
	GoPromptOptions: []prompt.Option{
		prompt.OptionTitle("ionosctl"),
		prompt.OptionPrefix("> "),
		prompt.OptionShowCompletionAtStart(),
		prompt.OptionAddKeyBind(
			prompt.KeyBind{Key: prompt.Enter, Fn: func(buf *prompt.Buffer) {
				viper.Reset()
			}},
		),
	},
	OnErrorFunc: func(err error) {
		rootCmd.Command.PrintErr(err)
		return
	},
}

func Shell() *core.Command {
	versionCmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
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
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	return versionCmd
}
