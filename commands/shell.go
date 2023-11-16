package commands

import (
	"context"

	comptplus "github.com/avirtopeanu-ionos/comptplus"
	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
		PreCmdRun: core.NoPreRun,
		CmdRun: func(c *core.CommandConfig) error {
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	return versionCmd
}
