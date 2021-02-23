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
	const (
		completionBashLong = `Use this command to generate completion code for BASH terminal. IonosCTL supports completion for commands and flags.

To load completions for the current session, execute: 

    source <(ionosctl completion bash)

To make these changes permanent, append the above line to your ` + "`" + `.bashrc` + "`" + ` file and use:

    source ~/.bashrc

You will need to start a new shell for this setup to take effect.`
		completionZshLong = `Use this command to generate completion code for ZSH terminal. IonosCTL supports completion for commands and flags.

If shell completions are not already enabled for your environment, you need to enable them. 
Add the following line to your ` + "`" + `~/.zshrc` + "`" + ` file:

    autoload -Uz compinit; compinit

To load completions for each session execute the following commands:

    mkdir -p ~/.config/ionosctl/completion/zsh
    ionosctl completion zsh > ~/.config/ionosctl/completion/zsh/_ionosctl

Finally add the following line to your ` + "`" + `~/.zshrc` + "`" + `file, *before* you
call the ` + "`" + `compinit` + "`" + ` function:

    fpath+=(~/.config/ionosctl/completion/zsh)

In the end your ` + "`" + `~/.zshrc` + "`" + ` file should contain the following two lines in the order given here:

    fpath+=(~/.config/ionosctl/completion/zsh)
    #  ... anything else that needs to be done before compinit
    autoload -Uz compinit; compinit
    # ...

You will need to start a new shell for this setup to take effect.
Note: ZSH completions require zsh 5.2 or newer.`
		completionPowerShellLong = `Use this command to generate completion code for PowerShell terminal. IonosCTL supports completion for commands and flags.

PowerShell supports three different completion modes:

- TabCompleteNext (default Windows style - on each key press the next option is displayed)
- Complete (works like Bash)
- MenuComplete (works like Zsh)

You set the mode with ` + "`" + `Set-PSReadLineKeyHandler -Key Tab -Function <mode>` + "`" + `

Descriptions will only be supported for Complete and MenuComplete.

Follow the next steps to enable it:

To load completions for the current session, execute: 

    PS> ionosctl completion powershell | Out-String | Invoke-Expression

To load completions for every new session, run:

    PS> ionosctl completion powershell > ionosctl.ps1

and source this file from your PowerShell profile or you can append the above line to your PowerShell profile file. 

You will need to start a new PowerShell for this setup to take effect.

Note: PowerShell completions require version 5.0 or above, which comes with Windows 10 and can be downloaded separately for Windows 7 or 8.1.`
		completionFishLong = `Use this command to generate completion code for Fish terminal. IonosCTL supports completion for commands and flags.

To load completions into the current shell execute:

    ionosctl completion fish | source

In order to make the completions permanent execute once:

    ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish`
	)

	completionCmd := &builder.Command{
		Command: &cobra.Command{
			Use:   "completion",
			Short: "Generate code to enable auto-completion with `TAB` key",
			Long:  "Use this command to generate completion code for specific shell for `ionosctl` commands and flags.",
		}}
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionBash, "bash", "Generate code to enable auto-completion with `TAB` key for BASH terminal", completionBashLong, "", false)
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionZsh, "zsh", "Generate code to enable auto-completion with `TAB` key for ZSH terminal", completionZshLong, "", false)
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionPowerShell, "powershell", "Generate code to enable auto-completion with `TAB` key for PowerShell terminal", completionPowerShellLong, "", false)
	builder.NewCommand(context.TODO(), completionCmd, noPreRun, RunCompletionFish, "fish", "Generate code to enable auto-completion with `TAB` key for Fish terminal", completionFishLong, "", false)

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

func RunCompletionPowerShell(c *builder.CommandConfig) error {
	var buf bytes.Buffer

	err := rootCmd.Command.GenPowerShellCompletionWithDesc(&buf)
	if err != nil {
		return errors.New("error while generating PowerShell completion:" + err.Error())
	}
	fmt.Println(buf.String())
	return nil
}

func RunCompletionFish(c *builder.CommandConfig) error {
	var buf bytes.Buffer

	err := rootCmd.Command.GenFishCompletion(&buf, false)
	if err != nil {
		return errors.New("error while generating Fish completion:" + err.Error())
	}

	// remove the completion command from auto-completion
	code := buf.String()
	code = strings.Replace(code, `commands+=("completion")`, "", -1)

	fmt.Println(code)
	return nil
}
