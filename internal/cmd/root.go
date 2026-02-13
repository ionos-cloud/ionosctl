package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/version"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// this is a temporary package that will be removed once all commands have been migrated to the new structure.
// It exists to avoid circular dependencies between the old and new command structures.
//
// It is responsible for defining the root command and global flags, as well as the help template.

var (
	// RootCmd is the root level command that all other commands attach to
	RootCmd = &core.Command{
		Command: &cobra.Command{
			Use:              "ionosctl",
			Short:            "IONOS Cloud CLI",
			Long:             "IonosCTL is a command-line interface for the Ionos Cloud API.",
			TraverseChildren: true,
		},
	}

	// Exported variables for command flags
	Output       string
	VerboseLevel int
	Quiet        bool
	Force        bool

	cfgFile string
)

// GetRootCmd returns the root command
func GetRootCmd() *core.Command {
	return RootCmd
}

// Execute runs the root command
func Execute() {
	if err := RootCmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
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
	RootCmd.Command.SetUsageTemplate(helpTemplate)
	RootCmd.Command.SetHelpCommand(helpCommand)

	RootCmd.Command.Version = version.Get() // Send the current version to Cobra
	viper.Set(constants.CLIHttpUserAgent, fmt.Sprintf("ionosctl/%v", RootCmd.Command.Version))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := RootCmd.Command.PersistentFlags()

	// Customize Help Flag
	rootPFlagSet.BoolP("help", "h", false, "Print usage")
	// Add Custom Flags
	rootPFlagSet.StringVarP(
		&cfgFile, constants.ArgConfig, constants.ArgConfigShort, config.GetConfigFilePath(),
		"Configuration file used for authentication",
	)
	_ = viper.BindPFlag(constants.ArgConfig, rootPFlagSet.Lookup(constants.ArgConfig))
	rootPFlagSet.StringVarP(
		&Output, constants.ArgOutput, constants.ArgOutputShort, constants.DefaultOutputFormat,
		"Desired output format [text|json|api-json]",
	)
	_ = viper.BindPFlag(constants.ArgOutput, rootPFlagSet.Lookup(constants.ArgOutput))
	_ = RootCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgOutput,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				jsontabwriter.JSONFormat, jsontabwriter.TextFormat, jsontabwriter.APIFormat,
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	RootCmd.GlobalFlags().StringP(
		constants.ArgServerUrl, constants.ArgServerUrlShort, constants.DefaultApiURL,
		"Override default host url",
	)
	rootPFlagSet.BoolVarP(&Quiet, constants.ArgQuiet, constants.ArgQuietShort, false, "Quiet output")
	_ = viper.BindPFlag(constants.ArgQuiet, rootPFlagSet.Lookup(constants.ArgQuiet))
	rootPFlagSet.BoolVarP(
		&Force, constants.ArgForce, constants.ArgForceShort, false, "Force command to execute without user input",
	)
	_ = viper.BindPFlag(constants.ArgForce, rootPFlagSet.Lookup(constants.ArgForce))
	rootPFlagSet.CountVarP(
		&VerboseLevel, constants.ArgVerbose, constants.ArgVerboseShort,
		"Increase verbosity level [-v, -vv, -vvv]",
	)
	_ = viper.BindPFlag(constants.ArgVerbose, rootPFlagSet.Lookup(constants.ArgVerbose))

	rootPFlagSet.Bool(constants.ArgNoHeaders, false, "Don't print table headers when table output is used")
	_ = viper.BindPFlag(constants.ArgNoHeaders, rootPFlagSet.Lookup(constants.ArgNoHeaders))

	rootPFlagSet.IntP(constants.DeprecatedFlagMaxResults, "M", 50, "DEPRECATED: Setting '--max-results' just sets '--limit' with the same value")
	_ = viper.BindPFlag(constants.FlagLimit, rootPFlagSet.Lookup(constants.DeprecatedFlagMaxResults)) // bind to limit
	rootPFlagSet.MarkHidden(constants.DeprecatedFlagMaxResults)

	rootPFlagSet.IntP(constants.FlagLimit, "", 50, "Maximum number of items to return per request")
	_ = viper.BindPFlag(constants.FlagLimit, rootPFlagSet.Lookup(constants.FlagLimit))

	rootPFlagSet.IntP(constants.FlagOffset, "", 0, "Number of items to skip before starting to collect the results")
	_ = viper.BindPFlag(constants.FlagOffset, rootPFlagSet.Lookup(constants.FlagOffset))

	rootPFlagSet.String(constants.FlagQuery, "", "JMESPath query string to filter the output")
	_ = viper.BindPFlag(constants.FlagQuery, rootPFlagSet.Lookup(constants.FlagQuery))

	rootPFlagSet.IntP(constants.FlagDepth, constants.FlagDepthShort, 1, "Level of detail for response objects")
	_ = viper.BindPFlag(constants.FlagDepth, rootPFlagSet.Lookup(constants.FlagDepth))

	rootPFlagSet.String(constants.FlagOrderBy, "", "Property to order the results by")
	_ = viper.BindPFlag(constants.FlagOrderBy, rootPFlagSet.Lookup(constants.FlagOrderBy))

	rootPFlagSet.StringSliceP(constants.FlagFilters, constants.FlagFiltersShort, []string{}, "Limit results to results containing the specified filter:"+
		"KEY1=VALUE1,KEY2=VALUE2")
	_ = viper.BindPFlag(constants.FlagFilters, rootPFlagSet.Lookup(constants.FlagFilters))

	rootPFlagSet.SortFlags = false

	// because of Viper Shenanigans, we have to bind it last, after any commands, to avoid overwriting the default...
	_ = viper.BindPFlag(constants.ArgServerUrl, RootCmd.GlobalFlags().Lookup(constants.ArgServerUrl))

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

const (
	availableCommands = `{{- if .HasAvailableSubCommands}}
AVAILABLE COMMANDS:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short | trim}}{{end}}{{end}}{{print "\n"}}{{end}}`

	seeAlso = `{{- if .HasHelpSubCommands}}
SEE ALSO:
{{.Annotations.SeeAlsos | trim}}{{print "\n"}}{{end}}`

	additionalHelpTopics = `{{- if .HasHelpSubCommands}}
Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short | trim}}{{end}}{{end}}{{print "\n"}}{{end}}`

	localFlags = `{{- if .HasAvailableLocalFlags}}
FLAGS:
{{.LocalFlags.FlagUsages}}{{end}}`

	globalFlags = `{{- if .HasAvailableInheritedFlags}}
GLOBAL FLAGS:
{{.InheritedFlags.FlagUsages}}{{end}}`

	usage = `{{- if .Runnable}}
USAGE:
  {{.UseLine | trim}}{{print "\n"}}{{end}}`

	aliases = `{{- if gt (len .Aliases) 0}}
ALIASES:
  {{.NameAndAliases | trim}}{{print "\n"}}{{end}}`

	examples = `{{- if .HasExample}}
EXAMPLES:
{{.Example | trim}}{{print "\n"}}{{end}}`

	moreInfo = `{{- if .HasAvailableSubCommands}}
Use "{{.CommandPath}} [command] --help" for more information about a command.{{print "\n"}}{{end}}`
)

var helpTemplate = strings.Join(
	[]string{
		seeAlso,
		additionalHelpTopics,
		globalFlags,
		localFlags,
		aliases,
		examples,
		usage,
		availableCommands,
		moreInfo,
	}, "",
)
