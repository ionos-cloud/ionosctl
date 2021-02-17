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

func generate() *builder.Command {
	const (
		completionBashLong = `Use this command to generate completion code for BASH terminal. IonosCTL supports completion for commands and flags.

Follow the next steps to enable it:

Linux:

- Generate completion code:
` + "`" + `ionosctl completion bash > $PWD/ionosctl_bash_completion.sh` + "`" + `
- Copy generated code to ` + "`" + `/etc/bash_completion.d/` + "`" + `:
` + "`" + `sudo cp $PWD/ionosctl_bash_completion.sh /etc/bash_completion.d/` + "`" + `
- Restart terminal to use auto-completion with TAB key.
- Clean-up:
` + "`" + `rm $PWD/ionosctl_bash_completion.sh` + "`" + `

Mac OS:

` + "`" + `ionosctl generate completion bash > /usr/local/etc/bash_completion.d/ionosctl_bash_completion.sh` + "`" + ``
		completionZshLong = `Use this command to generate completion code for ZSH terminal. IonosCTL supports completion for commands and flags.
Add the following line to your .profile or .bashrc.

` + "`" + `source  <(ionosctl completion zsh)` + "`" + `

Note:
- ZSH completions require zsh 5.2 or newer.`
	)

	completionCmd := &builder.Command{
		Command: &cobra.Command{
			Use:   "completion",
			Short: "Generate code to enable auto-completion with TAB key",
			Long:  "Use this command to generate completion code for specific shell for `ionosctl` commands and flags.",
		}}
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionBash, "bash", "Generate code to enable auto-completion with TAB key for BASH terminal", completionBashLong, "", false)
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionZsh, "zsh", "Generate code to enable auto-completion with TAB key for ZSH terminal", completionZshLong, "", false)

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

	fmt.Println(code)
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

	fmt.Println(code)
	return nil
}
