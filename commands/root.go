package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/monitoring"

	"github.com/ionos-cloud/ionosctl/v6/commands/cdn"
	"github.com/ionos-cloud/ionosctl/v6/commands/kafka"

	certificates "github.com/ionos-cloud/ionosctl/v6/commands/cert"
	"github.com/ionos-cloud/ionosctl/v6/commands/cfg"
	"github.com/ionos-cloud/ionosctl/v6/commands/compute"
	container_registry "github.com/ionos-cloud/ionosctl/v6/commands/container-registry"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas"
	"github.com/ionos-cloud/ionosctl/v6/commands/dns"
	logging_service "github.com/ionos-cloud/ionosctl/v6/commands/logging-service"
	"github.com/ionos-cloud/ionosctl/v6/commands/token"
	vm_autoscaling "github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/globalwait"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/version"
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
	Output       string
	VerboseLevel int
	Quiet        bool
	Force        bool

	cfgFile string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Command.Execute(); err != nil {
		os.Exit(1)
	}
	if viper.GetBool(constants.ArgWait) {
		if err := globalwait.WaitForAvailable(os.Stderr); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Re-render output with fresh data showing the final state.
		// Try new table.Table path first, then fall back to legacy jsontabwriter path.
		freshData, err := globalwait.FetchResource()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not fetch updated resource: %v\n", err)
			return
		}

		globalwait.SetRerendering(true)
		defer globalwait.SetRerendering(false)

		if r, cols := globalwait.GetRerenderable(); r != nil {
			// New path: table.Table implements Rerenderable
			if err := r.Extract(freshData); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not extract fresh data: %v\n", err)
				return
			}
			out, err := r.Render(cols)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not re-render output: %v\n", err)
				return
			}
			fmt.Fprint(os.Stdout, out)
		} else if ri := globalwait.GetRenderInfo(); ri != nil {
			// Legacy path: jsontabwriter with RenderInfo
			out, err := jsontabwriter.GenerateOutput(ri.Prefix, ri.Mapping, freshData, ri.Cols)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: could not re-render output: %v\n", err)
				return
			}
			fmt.Fprint(os.Stdout, out)
		}
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

	rootCmd.Command.Version = version.Get() // Send the current version to Cobra
	viper.Set(constants.CLIHttpUserAgent, fmt.Sprintf("ionosctl/%v", rootCmd.Command.Version))

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := rootCmd.Command.PersistentFlags()

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
	_ = rootCmd.Command.RegisterFlagCompletionFunc(
		constants.ArgOutput,
		func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return []string{
				jsontabwriter.JSONFormat, jsontabwriter.TextFormat, jsontabwriter.APIFormat,
			}, cobra.ShellCompDirectiveNoFileComp
		},
	)
	rootCmd.GlobalFlags().StringP(
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

	rootPFlagSet.Bool(constants.ArgWait, false,
		"Wait for the resource to reach AVAILABLE state after the command completes. "+
			"Works for create/update commands that return API resources with an href field")
	_ = viper.BindPFlag(constants.ArgWait, rootPFlagSet.Lookup(constants.ArgWait))

	rootPFlagSet.Int(constants.ArgWaitTimeout, constants.DefaultWaitTimeoutSeconds,
		"Timeout in seconds for the global --wait flag")
	_ = viper.BindPFlag(constants.ArgWaitTimeout, rootPFlagSet.Lookup(constants.ArgWaitTimeout))

	rootPFlagSet.SortFlags = false

	// Add SubCommands to RootCmd
	addCommands()

	// Deprecate old per-command wait flags in favour of global --wait
	deprecateWaitFlags(rootCmd.Command)

	// because of Viper Shenanigans, we have to bind it last, after any commands, to avoid overwriting the default...
	_ = viper.BindPFlag(constants.ArgServerUrl, rootCmd.GlobalFlags().Lookup(constants.ArgServerUrl))

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

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.AddCommand(Shell())
	rootCmd.AddCommand(VersionCmd())
	rootCmd.AddCommand(Man())
	rootCmd.AddCommand(cfg.Login())

	// cfg
	rootCmd.AddCommand(cfg.ConfigCmd())
	// Config namespace commands are also available via the root command, but are hidden
	for _, cmd := range cfg.ConfigCmd().SubCommands() {
		if cmd.Name() == "location" {
			// This one is confusing without `cfg` namespace;
			// It also would override CPU Architecture locations command, so skip it.
			continue
		}
		cmd.Command.Hidden = true
		rootCmd.AddCommand(cmd)
	}

	// V6 Resources Commands under "compute" namespace
	rootCmd.AddCommand(compute.Root())
	// Hidden backward-compat aliases at root level (e.g. "ionosctl server" still works)
	for _, cmd := range compute.HiddenAliases() {
		rootCmd.AddCommand(cmd)
	}
	// Auth Command
	rootCmd.AddCommand(token.TokenCmd())
	// Add DBaaS Commands
	rootCmd.AddCommand(dbaas.DataBaseServiceCmd())
	// Add Certificate Manager Commands
	rootCmd.AddCommand(certificates.Root())
	// Add Container Registry Commands
	rootCmd.AddCommand(container_registry.ContainerRegistryCmd())
	// VM-Autoscaling commands
	rootCmd.AddCommand(vm_autoscaling.Root())

	rootCmd.AddCommand(dns.Root())
	rootCmd.AddCommand(logging_service.Root())

	rootCmd.AddCommand(monitoring.Root())

	rootCmd.AddCommand(cdn.Command())

	rootCmd.AddCommand(vpn.Root())

	rootCmd.AddCommand(kafka.Command())
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

// deprecateWaitFlags walks the command tree and marks old per-command wait flags
// as deprecated in favour of the global --wait flag.
func deprecateWaitFlags(cmd *cobra.Command) {
	for _, flagName := range []string{
		constants.ArgWaitForRequest,
		constants.ArgWaitForState,
		constants.ArgWaitForDelete,
	} {
		if f := cmd.Flags().Lookup(flagName); f != nil {
			cmd.Flags().MarkDeprecated(flagName, "use global --wait flag instead")
		}
	}
	for _, child := range cmd.Commands() {
		deprecateWaitFlags(child)
	}
}
