package commands

import (
	"context"
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/stromland/cobra-prompt"
)

var simplePrompt = &cobraprompt.CobraPrompt{
	RootCmd:                  GetRootCmd().Command,
	AddDefaultExitCommand:    true,
	DisableCompletionCommand: true,
	OnErrorFunc: func(err error) {
		fmt.Printf(err.Error())
	},
}

func InteractiveShellCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "shell",
		Resource:  "shell",
		Verb:      "shell",
		ShortDesc: "Use ionosctl interactively",
		Example:   "ionosctl shell",
		PreCmdRun: core.NoPreRun,
		CmdRun: func(commandConfig *core.CommandConfig) error {
			fmt.Printf("Use `exit` to leave the shell.\n")
			simplePrompt.Run()
			return nil
		},
		InitClient: false,
	})

	return cmd
}
