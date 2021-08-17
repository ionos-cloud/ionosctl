package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// RootCmd is the root level command that all other commands attach to
	rootCmd = &core.Command{
		Command: &cobra.Command{
			Use:              "ionosctl",
			Short:            "IONOS Cloud CLI",
			Long:             "IonosCTL is a command-line interface for the Ionos Cloud API.",
			TraverseChildren: true,
		},
	}
	noPreRun  = func(c *core.PreCommandConfig) error { return nil }
	ServerURL string
	Output    string
	Quiet     bool
	Force     bool
	Verbose   bool

	cfgFile string

	Major string
	Minor string
	Patch string
	// Label - If label is not set, the version will be: DEV
	// If label is set as `release`, it will show the version released
	Label string

	IonosctlVersion cliVersion
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *core.Command {
	return rootCmd
}

func init() {
	initConfig()
	rootCmd.Command.SetUsageTemplate(usageTemplate)

	// Init version
	initVersion()
	rootCmd.Command.Version = IonosctlVersion.GetVersion()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := rootCmd.GlobalFlags()
	rootPFlagSet.StringVarP(&cfgFile, config.ArgConfig, config.ArgConfigShort, config.GetConfigFile(), "Configuration file used for authentication")
	_ = viper.BindPFlag(config.ArgConfig, rootPFlagSet.Lookup(config.ArgConfig))

	rootPFlagSet.StringVarP(&ServerURL, config.ArgServerUrl, config.ArgServerUrlShort, config.DefaultApiURL, "Override default host url")
	_ = viper.BindPFlag(config.ArgServerUrl, rootPFlagSet.Lookup(config.ArgServerUrl))

	rootPFlagSet.StringVarP(&Output, config.ArgOutput, config.ArgOutputShort, config.DefaultOutputFormat, "Desired output format [text|json]")
	_ = viper.BindPFlag(config.ArgOutput, rootPFlagSet.Lookup(config.ArgOutput))
	_ = rootCmd.Command.RegisterFlagCompletionFunc(config.ArgOutput, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"json", "text"}, cobra.ShellCompDirectiveNoFileComp
	})

	rootPFlagSet.BoolVarP(&Quiet, config.ArgQuiet, config.ArgQuietShort, false, "Quiet output")
	_ = viper.BindPFlag(config.ArgQuiet, rootPFlagSet.Lookup(config.ArgQuiet))

	rootPFlagSet.BoolVarP(&Force, config.ArgForce, config.ArgForceShort, false, "Force command to execute without user input")
	_ = viper.BindPFlag(config.ArgForce, rootPFlagSet.Lookup(config.ArgForce))

	rootPFlagSet.BoolVarP(&Verbose, config.ArgVerbose, config.ArgVerboseShort, false, "see step by step process when running a command")
	_ = viper.BindPFlag(config.ArgVerbose, rootPFlagSet.Lookup(config.ArgVerbose))

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
		viper.SetConfigName("config")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match
}

func initVersion() {
	if Major != "" {
		i, _ := strconv.Atoi(Major)
		IonosctlVersion.major = i
	}
	if Minor != "" {
		i, _ := strconv.Atoi(Minor)
		IonosctlVersion.minor = i
	}
	if Patch != "" {
		i, _ := strconv.Atoi(Patch)
		IonosctlVersion.patch = i
	}
	if Label == "" {
		IonosctlVersion.label = "DEV"
	} else {
		IonosctlVersion.label = Label
	}
}

type cliVersion struct {
	major int
	minor int
	patch int
	label string
}

func (v cliVersion) GetVersion() string {
	if v.label != "release" {
		return fmt.Sprintf("%s", v.label)
	} else {
		return fmt.Sprintf("%d.%d.%d", v.major, v.minor, v.patch)
	}
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.AddCommand(version())
	// V6 Resources Commands
	rootCmd.AddCommand(login())
	rootCmd.AddCommand(location())
	rootCmd.AddCommand(datacenter())
	rootCmd.AddCommand(server())
	rootCmd.AddCommand(volume())
	rootCmd.AddCommand(lan())
	rootCmd.AddCommand(natgateway())
	rootCmd.AddCommand(networkloadbalancer())
	rootCmd.AddCommand(nic())
	rootCmd.AddCommand(loadBalancer())
	rootCmd.AddCommand(ipblock())
	rootCmd.AddCommand(ipconsumer())
	rootCmd.AddCommand(ipfailover())
	rootCmd.AddCommand(request())
	rootCmd.AddCommand(snapshot())
	rootCmd.AddCommand(image())
	rootCmd.AddCommand(firewallrule())
	rootCmd.AddCommand(flowlog())
	rootCmd.AddCommand(label())
	rootCmd.AddCommand(contract())
	rootCmd.AddCommand(user())
	rootCmd.AddCommand(group())
	rootCmd.AddCommand(resource())
	rootCmd.AddCommand(backupunit())
	rootCmd.AddCommand(pcc())
	rootCmd.AddCommand(share())
	rootCmd.AddCommand(k8s())
	rootCmd.AddCommand(template())
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
