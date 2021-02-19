package commands

import (
	"fmt"
	"os"

	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// RootCmd is the root level command that all other commands attach to
	rootCmd = &builder.Command{
		Command: &cobra.Command{
			Use:              "ionosctl",
			Short:            "Ionos Cloud CLI",
			Long:             "IonosCTL is a command-line interface (CLI) for the Ionos Cloud API.",
			TraverseChildren: true,
		},
	}
	noPreRun  = func(c *builder.PreCommandConfig) error { return nil }
	ServerURL string
	Output    string
	Quiet     bool
	Verbose   bool

	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *builder.Command {
	return rootCmd
}

func init() {
	initConfig()

	rootCmd.Command.SetUsageTemplate(usageTemplate)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := rootCmd.Command.PersistentFlags()
	rootPFlagSet.StringVarP(&cfgFile, config.ArgConfig, "c", config.GetConfigFilePath(), "Configuration file used for authentication")
	viper.BindPFlag(config.ArgConfig, rootPFlagSet.Lookup(config.ArgConfig))

	rootPFlagSet.StringVarP(&ServerURL, config.ArgServerUrl, "u", config.DefaultApiURL, "Override default API endpoint")
	viper.BindPFlag(config.ArgServerUrl, rootPFlagSet.Lookup(config.ArgServerUrl))

	rootPFlagSet.StringVarP(&Output, config.ArgOutput, "o", config.DefaultOutputFormat, "Desired output format [text|json]")
	viper.BindPFlag(config.ArgOutput, rootPFlagSet.Lookup(config.ArgOutput))
	rootCmd.Command.RegisterFlagCompletionFunc(config.ArgOutput, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "text"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootPFlagSet.BoolVarP(&Quiet, config.ArgQuiet, "q", false, "Quiet output")
	viper.BindPFlag(config.ArgQuiet, rootPFlagSet.Lookup(config.ArgQuiet))

	rootPFlagSet.Bool(config.ArgIgnoreStdin, false, "Force command to execute without user input")
	viper.BindPFlag(config.ArgIgnoreStdin, rootPFlagSet.Lookup(config.ArgIgnoreStdin))

	rootPFlagSet.BoolVarP(&Verbose, config.ArgVerbose, "v", false, "Enable verbose output")

	addCommands()

	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName("ionosctl-config")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match
	_ = viper.ReadInConfig()
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.AddCommand(login())
	rootCmd.AddCommand(version())
	rootCmd.AddCommand(completion())
	rootCmd.AddCommand(location())
	rootCmd.AddCommand(datacenter())
	rootCmd.AddCommand(server())
	rootCmd.AddCommand(volume())
	rootCmd.AddCommand(lan())
	rootCmd.AddCommand(nic())
	rootCmd.AddCommand(loadbalancer())
	rootCmd.AddCommand(request())
}

const usageTemplate = `USAGE: {{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

ALIASES:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

EXAMPLES:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

AVAILABLE COMMANDS:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

FLAGS:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

GLOBAL FLAGS:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

SEE ALSO:
{{.Annotations.SeeAlsos}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`
