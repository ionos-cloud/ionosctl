package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	cloudapiv5 "github.com/ionos-cloud/ionosctl/commands/cloudapi-v5"
	"github.com/ionos-cloud/ionosctl/internal/config"
	"github.com/ionos-cloud/ionosctl/internal/core"
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
	ServerURL string
	Output    string
	Quiet     bool
	Force     bool
	Verbose   bool
	All       bool

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

// Customize Help Command
var helpCommand = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about the command",
	RunE: func(c *cobra.Command, args []string) error {
		cmd, args, e := c.Root().Find(args)
		if cmd == nil || e != nil || len(args) > 0 {
			return fmt.Errorf("unknown help topic: %v", strings.Join(args, " "))
		}
		helpFunc := cmd.HelpFunc()
		helpFunc(cmd, args)
		return nil
	},
}

func init() {
	initConfig()

	rootCmd.Command.SetUsageTemplate(helpTemplate)
	rootCmd.Command.SetHelpCommand(helpCommand)

	// Init version
	initVersion()
	rootCmd.Command.Version = IonosctlVersion.GetVersion()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := rootCmd.GlobalFlags()
	// Customize Help Flag
	rootPFlagSet.BoolP("help", "h", false, "Print usage")
	// Add Custom Flags
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
	rootPFlagSet.BoolVarP(&Verbose, config.ArgVerbose, config.ArgVerboseShort, false, "Print step-by-step process when running command")
	_ = viper.BindPFlag(config.ArgVerbose, rootPFlagSet.Lookup(config.ArgVerbose))

	// Add SubCommands to RootCmd
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
	// Read Environment Variables.
	// For authentication, there are used the following ENV:
	// IONOS_USERNAME, IONOS_PASSWORD or IONOS_TOKEN
	// The user can also overwrite the endpoint: IONOS_API_URL
	viper.AutomaticEnv()
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
		IonosctlVersion.patch = Patch
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
	patch string
	label string
}

func (v cliVersion) GetVersion() string {
	if v.label != "release" {
		return fmt.Sprintf("%s", v.label)
	} else {
		return fmt.Sprintf("%d.%d.%s", v.major, v.minor, v.patch)
	}
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(LoginCmd())
	// Add CloudApiV5 Resources Commands
	rootCmd.AddCommand(cloudapiv5.DatacenterCmd())
	rootCmd.AddCommand(cloudapiv5.BackupunitCmd())
	rootCmd.AddCommand(cloudapiv5.ContractCmd())
	rootCmd.AddCommand(cloudapiv5.FirewallRuleCmd())
	rootCmd.AddCommand(cloudapiv5.GroupCmd())
	rootCmd.AddCommand(cloudapiv5.ImageCmd())
	rootCmd.AddCommand(cloudapiv5.IpblockCmd())
	rootCmd.AddCommand(cloudapiv5.IpconsumerCmd())
	rootCmd.AddCommand(cloudapiv5.IpfailoverCmd())
	rootCmd.AddCommand(cloudapiv5.K8sCmd())
	rootCmd.AddCommand(cloudapiv5.LabelCmd())
	rootCmd.AddCommand(cloudapiv5.LanCmd())
	rootCmd.AddCommand(cloudapiv5.LoadBalancerCmd())
	rootCmd.AddCommand(cloudapiv5.LocationCmd())
	rootCmd.AddCommand(cloudapiv5.NicCmd())
	rootCmd.AddCommand(cloudapiv5.PccCmd())
	rootCmd.AddCommand(cloudapiv5.ResourceCmd())
	rootCmd.AddCommand(cloudapiv5.RequestCmd())
	rootCmd.AddCommand(cloudapiv5.ServerCmd())
	rootCmd.AddCommand(cloudapiv5.ShareCmd())
	rootCmd.AddCommand(cloudapiv5.UserCmd())
	rootCmd.AddCommand(cloudapiv5.VolumeCmd())
}

const helpTemplate = `USAGE: {{if .Runnable}}
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
