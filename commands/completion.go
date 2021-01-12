package commands

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func completion() *cobra.Command {
	completionCmd := &cobra.Command{
		Use:   "completion",
		Short: "Modify your shell so ionosctl commands autocomplete with TAB",
	}

	// TODO: to be updated with completion subcommand for zsh, powershell?
	completionCmd.AddCommand(
		bashCompletion(),
	)

	return completionCmd
}

func bashCompletion() *cobra.Command {
	bashCompletionCmd := &cobra.Command{
		Use:   "bash",
		Short: "Generate completion code for bash",
		Long: `
# To load completions for each session, execute once:
Linux:
  $ ionosctl completion bash > $PWD/ionosctl_bash_completion.sh
  $ sudo cp $PWD/ionosctl_bash_completion.sh /etc/bash_completion.d/
  $ rm $PWD/ionosctl_bash_completion.sh

  Restart terminal to use auto-completion with TAB.

MacOS:
  $ ionosctl completion bash > /usr/local/etc/bash_completion.d/ionosctl_bash_completion.sh
`,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			var buf bytes.Buffer

			err := cmd.Root().GenBashCompletion(&buf)
			if err != nil {
				return fmt.Errorf("error while generating bash completion: %v", err)
			}

			// remove the completion command from auto-completion
			code := buf.String()
			code = strings.Replace(code, `commands+=("completion")`, "", -1)

			fmt.Print(code)
			return nil
		},
	}
	return bashCompletionCmd
}
