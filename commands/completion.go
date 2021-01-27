package commands

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/spf13/cobra"
)

func completion() *builder.Command {
	completionCmd := &builder.Command{
		Command: &cobra.Command{
			Use:   "completion",
			Short: "Modify your shell so `ionosctl` commands and flags autocomplete with <TAB>",
			Long: `
ionosctl completion helps you to enable autocompletion with <TAB> key for ionosctl commands and flags.
`,
		},
	}

	const (
		completionBashLong = `
# To load completions for each session, execute once:
Linux:
  $ ionosctl completion bash > $PWD/ionosctl_bash_completion.sh
  $ sudo cp $PWD/ionosctl_bash_completion.sh /etc/bash_completion.d/
  $ rm $PWD/ionosctl_bash_completion.sh

  Restart terminal to use auto-completion with <TAB>.

MacOS:
  $ ionosctl completion bash > /usr/local/etc/bash_completion.d/ionosctl_bash_completion.sh
`
		completionZshLong = `
# Add the following line to your .profile or .bashrc.

	source  <(ionosctl completion zsh)

Note:

- zsh completions requires zsh 5.2 or newer.
`
	)

	builder.NewCommand(context.TODO(), completionCmd, RunCompletionBash, "bash", "Generate completion code for bash", completionBashLong, false)
	builder.NewCommand(context.TODO(), completionCmd, RunCompletionZsh, "zsh", "Generate completion code for zsh", completionZshLong, false)

	return completionCmd
}

func RunCompletionBash(c *builder.CommandConfig) error {
	var buf bytes.Buffer

	err := rootCmd.Command.GenBashCompletion(&buf)
	if err != nil {
		return errors.New("error while generating bash completion:" + err.Error())
	}

	// remove the completion command from auto-completion
	code := buf.String()
	code = strings.Replace(code, `commands+=("completion")`, "", -1)

	fmt.Fprintf(c.Printer.Stdout, code)
	return nil
}

func RunCompletionZsh(c *builder.CommandConfig) error {
	var buf bytes.Buffer

	err := rootCmd.Command.GenZshCompletion(&buf)
	if err != nil {
		return errors.New("error while generating zsh completion:" + err.Error())
	}

	// remove the completion command from auto-completion
	code := buf.String()
	code = strings.Replace(code, `commands+=("completion")`, "", -1)

	fmt.Fprintf(c.Printer.Stdout, code)
	return nil
}
