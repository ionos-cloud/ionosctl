package commands

import (
	"context"
	"os"
	"strings"

	cobraprompt "github.com/avirtopeanu-ionos/cobra-prompt"
	"github.com/c-bata/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
)

var advancedPrompt = &cobraprompt.CobraPrompt{
	RootCmd:                  rootCmd.Command,
	PersistFlagValues:        true,
	ShowHelpCommandAndFlags:  true,
	DisableCompletionCommand: true,
	AddDefaultExitCommand:    true,
	GoPromptOptions: []prompt.Option{
		prompt.OptionTitle("ionosctl"),
		prompt.OptionPrefix("> "),
		prompt.OptionMaxSuggestion(10),
	},
	OnErrorFunc: func(err error) {
		if strings.Contains(err.Error(), "unknown command") {
			// Don't exit on unknown commands
			rootCmd.Command.PrintErrln(err)
			return
		}

		if strings.Contains(err.Error(), "401 Unauthorized") {
			// Don't exit on unknown commands
			rootCmd.Command.PrintErrln("Unauthorized. Please login. Use `ionosctl whoami` to debug your credentials")
			return
		}

		rootCmd.Command.PrintErr(err)
		os.Exit(1)
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
