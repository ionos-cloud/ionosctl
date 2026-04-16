package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/elk-language/go-prompt"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/version"
	"github.com/ionoscloudsdk/comptplus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var advancedPrompt = &comptplus.CobraPrompt{
	RootCmd:                   rootCmd.Command,
	ShowHelpCommandAndFlags:   true,
	DisableCompletionCommand:  true,
	AddDefaultExitCommand:     true,
	ShowHiddenCommands:        true,
	ShowHiddenFlags:           true,
	AsyncFlagValueSuggestions: true,
	GoPromptOptions: []prompt.Option{
		prompt.WithTitle("ionosctl"),
		prompt.WithPrefix("> "),
		prompt.WithShowCompletionAtStart(),

		prompt.WithDescriptionTextColor(prompt.Black),
		prompt.WithSuggestionTextColor(prompt.White),
		prompt.WithDescriptionBGColor(prompt.LightGray),
		prompt.WithSuggestionBGColor(prompt.DarkGray),

		prompt.WithSelectedDescriptionTextColor(prompt.White),
		prompt.WithSelectedSuggestionTextColor(prompt.Black),
		prompt.WithSelectedDescriptionBGColor(prompt.DarkGray),
		prompt.WithSelectedSuggestionBGColor(prompt.LightGray),

		prompt.WithPrefixTextColor(prompt.DefaultColor),
		prompt.WithScrollbarThumbColor(prompt.DarkGray),
		prompt.WithScrollbarBGColor(prompt.DefaultColor),
	},

	OnErrorFunc: func(err error) {
		// Printing this would lead to duplicated errors
		// TODO: Fix me
		// rootCmd.Command.PrintErr(err)
	},

	CustomFlagResetBehaviour: func(flag *pflag.Flag) {
		sliceValue, ok := flag.Value.(pflag.SliceValue)
		if !ok {
			// For non-slice flags, just set to the default value
			flag.Value.Set(flag.DefValue)
			return
		}

		err := sliceValue.Replace([]string{})
		if err != nil {
			flag.Value.Set(flag.DefValue)
		}
	},
}

func Shell() *core.Command {
	flagPersistFlagValues := "persist-flag-values"

	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "shell",
		Resource:  "shell",
		Verb:      "shell",
		ShortDesc: "Interactive shell",
		LongDesc: `The ionosctl shell command launches an interactive shell environment, enabling a more dynamic and intuitive way to interact with the ionosctl CLI.
This shell is designed to enhance your command-line experience with advanced features and customizations, powered by the comptplus library.

DEFAULT CONTROLS:
Ctrl + A\tGo to the beginning of the line (Home)
Ctrl + E\tGo to the end of the line (End)
Ctrl + F\tForward one character
Ctrl + B\tBackward one character
Ctrl + D\tDelete character under the cursor
Ctrl + H\tDelete character before the cursor (Backspace)
Ctrl + W\tCut the word before the cursor to the clipboard
Ctrl + K\tCut the line after the cursor to the clipboard
Ctrl + U\tCut the line before the cursor to the clipboard
Ctrl + L\tClear the screen`,
		Example: "ionosctl shell",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if os.Getenv("__IONOSCTL_SHELL_ACTIVE") == "1" {
				return fmt.Errorf("already inside an ionosctl shell session")
			}
			_, err := client.Get()
			if err != nil {
				return fmt.Errorf("usage of the interactive shell requires valid credentials. "+
					"You can use `ionosctl whoami` to debug your configuration: %w", err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			os.Setenv("__IONOSCTL_SHELL_ACTIVE", "1")

			fmt.Printf("ionosctl %s\n", version.Get())
			fmt.Println("Controls:")
			fmt.Println("   Ctrl+A  Go to beginning of line   Ctrl+K  Cut line after cursor")
			fmt.Println("   Ctrl+E  Go to end of line         Ctrl+U  Cut line before cursor")
			fmt.Println("   Ctrl+F  Forward one char          Ctrl+W  Cut word before cursor")
			fmt.Println("   Ctrl+B  Backward one char         Ctrl+H  Backspace")
			fmt.Println("   Ctrl+D  Delete char under cursor  Ctrl+L  Clear screen")
			advancedPrompt.PersistFlagValues = viper.GetBool(flagPersistFlagValues)
			advancedPrompt.Run()
			return nil
		},
		InitClient: false,
	})

	cmd.AddBoolFlag(flagPersistFlagValues, "p", false, "Persist flag values between commands")

	return cmd
}
