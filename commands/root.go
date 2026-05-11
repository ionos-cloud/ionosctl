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
	objectstorage "github.com/ionos-cloud/ionosctl/v6/commands/object-storage"
	"github.com/ionos-cloud/ionosctl/v6/commands/token"
	vm_autoscaling "github.com/ionos-cloud/ionosctl/v6/commands/vm-autoscaling"
	"github.com/ionos-cloud/ionosctl/v6/commands/vpn"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/config"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/globalwait"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
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
			Short:            "IONOS CLOUD CLI",
			Long:             "IonosCTL is a command-line interface for the IONOS CLOUD API.",
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
	err := rootCmd.Command.Execute()

	if err == nil && viper.GetBool(constants.ArgWait) {
		token, username, password := getAuthCreds()

		if waitErr := globalwait.WaitForAvailable(os.Stderr, token, username, password); waitErr != nil {
			fmt.Fprintf(os.Stderr, "Error waiting: %v\n", waitErr)
			os.Exit(1)
		}

		// Re-render output with fresh data showing final state
		if !viper.GetBool(constants.ArgQuiet) {
			if r, cols := globalwait.GetRerenderable(); r != nil {
				freshData, fetchErr := globalwait.FetchResource(token, username, password)
				if fetchErr != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not fetch updated resource: %v\n", fetchErr)
				} else {
					globalwait.SetRerendering(true)
					defer globalwait.SetRerendering(false)
					if extractErr := r.Extract(freshData); extractErr != nil {
						fmt.Fprintf(os.Stderr, "Warning: could not extract fresh data: %v\n", extractErr)
					} else if out, renderErr := r.Render(cols); renderErr != nil {
						fmt.Fprintf(os.Stderr, "Warning: could not re-render output: %v\n", renderErr)
					} else {
						fmt.Fprint(os.Stdout, out)
					}
				}
			}
		}
	}

	if err != nil {
		os.Exit(1)
	}
}

func GetRootCmd() *core.Command {
	return rootCmd
}

// Customize Help Command
var helpCommand = &cobra.Command{
	Use:   "help [command]",
	Short: "Help about any command",
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
			return []string{"json", "text", "api-json"}, cobra.ShellCompDirectiveNoFileComp
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
	_ = viper.BindPFlag(constants.DeprecatedFlagMaxResults, rootPFlagSet.Lookup(constants.DeprecatedFlagMaxResults))
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

	rootPFlagSet.BoolP(constants.ArgWait, "w", false,
		"Wait for the resource to reach AVAILABLE state after the command completes")
	_ = viper.BindPFlag(constants.ArgWait, rootPFlagSet.Lookup(constants.ArgWait))

	rootPFlagSet.IntP(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds,
		"Timeout in seconds for --wait and other wait operations")
	_ = viper.BindPFlag(constants.ArgTimeout, rootPFlagSet.Lookup(constants.ArgTimeout))

	// Deprecated aliases: old per-command wait flags now just set --wait.
	// Kept for backward compatibility so existing scripts don't break.
	for _, old := range []string{
		constants.ArgWaitForRequest,
		constants.ArgWaitForState,
		constants.ArgWaitForDelete,
	} {
		rootPFlagSet.Bool(old, false, "DEPRECATED: use --wait instead")
		_ = viper.BindPFlag(old, rootPFlagSet.Lookup(old))
		rootPFlagSet.MarkHidden(old)
		rootPFlagSet.MarkDeprecated(old, "use --wait instead")
	}

	// Deprecated -W shorthand removed: conflicts with --weight (-W) in NLB/targetgroup commands.
	// The long form --wait-for-state (above) still works. -W was rarely used directly.
	rootPFlagSet.Bool("wait-for-state-deprecated", false, "DEPRECATED: use --wait instead")
	_ = viper.BindPFlag("wait-for-state-deprecated", rootPFlagSet.Lookup("wait-for-state-deprecated"))
	rootPFlagSet.MarkHidden("wait-for-state-deprecated")
	rootPFlagSet.MarkDeprecated("wait-for-state-deprecated", "use --wait (-w) instead")

	// If any old flag is set, activate --wait
	cobra.OnInitialize(func() {
		for _, old := range []string{
			constants.ArgWaitForRequest,
			constants.ArgWaitForState,
			constants.ArgWaitForDelete,
			"wait-for-state-deprecated",
		} {
			if viper.GetBool(old) {
				viper.Set(constants.ArgWait, true)
			}
		}
	})

	// Wire the BeforeRender hook: when --wait is set, capture href and suppress
	// initial output so we can re-render with the final AVAILABLE state.
	table.BeforeRender = func(t *table.Table, visibleCols []string) bool {
		if !viper.GetBool(constants.ArgWait) || globalwait.IsRerendering() {
			return true // render normally
		}
		// Only suppress output for known valid formats. Invalid formats
		// (e.g. typo "-o jso") should render normally so the error surfaces
		// immediately instead of being lost after wait + re-render failure.
		switch viper.GetString(constants.ArgOutput) {
		case "text", "json", "api-json":
		default:
			return true
		}
		href := globalwait.ExtractHref(t.Raw())
		if href == "" {
			// No href in response (e.g. postgres-v1, mongo, DNS).
			id := globalwait.ExtractID(t.Raw())
			if id == "" {
				return true // list or unrecognized format - render normally
			}
			// For GET, the transport-captured URL is already the resource URL.
			// For POST/PUT/PATCH, it's the collection URL - append the id.
			if base := globalwait.GetHref(); base != "" && !globalwait.IsGetOperation() {
				globalwait.CaptureHref(strings.TrimRight(base, "/") + "/" + id)
			}
			if globalwait.GetHref() == "" {
				return true // no href and no fallback, render normally
			}
		} else {
			// Response has href, use it directly. More specific than the
			// transport-captured URL. buildFullURL resolves relative hrefs.
			globalwait.CaptureHref(href)
		}
		globalwait.CaptureRerenderable(t, visibleCols)
		return false // suppress initial output
	}

	rootPFlagSet.SortFlags = false

	// Add SubCommands to RootCmd
	addCommands()

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

const (
	groupServices = "services"
	groupAuth     = "auth"
	groupOther    = "other"
)

// addCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.Command.AddGroup(
		&cobra.Group{ID: groupServices, Title: "Cloud Services:"},
		&cobra.Group{ID: groupAuth, Title: "Authentication & Configuration:"},
		&cobra.Group{ID: groupOther, Title: "Other:"},
	)
	rootCmd.Command.SetHelpCommandGroupID(groupOther)
	rootCmd.Command.SetCompletionCommandGroupID(groupOther)

	// Cloud Services
	addServiceCmd(cdn.Command())
	addServiceCmd(certificates.Root())
	addServiceCmd(compute.Root())
	addServiceCmd(container_registry.ContainerRegistryCmd())
	addServiceCmd(dbaas.DataBaseServiceCmd())
	addServiceCmd(dns.Root())
	addServiceCmd(kafka.Command())
	addServiceCmd(logging_service.Root())
	addServiceCmd(monitoring.Root())
	addServiceCmd(vm_autoscaling.Root())
	addServiceCmd(vpn.Root())
	addServiceCmd(objectstorage.Root())

	// Hidden backward-compat aliases at root level (e.g. "ionosctl server" still works)
	for _, cmd := range compute.HiddenAliases() {
		rootCmd.AddCommand(cmd)
	}

	// Authentication & Configuration
	addAuthCmd(cfg.Login())
	addAuthCmd(token.TokenCmd())
	addAuthCmd(cfg.ConfigCmd())
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

	// Other
	addOtherCmd(Shell())
	addOtherCmd(VersionCmd())
	addOtherCmd(Man())
}

func addServiceCmd(cmd *core.Command) {
	cmd.Command.GroupID = groupServices
	rootCmd.AddCommand(cmd)
}

func addAuthCmd(cmd *core.Command) {
	cmd.Command.GroupID = groupAuth
	rootCmd.AddCommand(cmd)
}

func addOtherCmd(cmd *core.Command) {
	cmd.Command.GroupID = groupOther
	rootCmd.AddCommand(cmd)
}

const (
	availableCommands = `{{- if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}
AVAILABLE COMMANDS:{{range $cmds}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short | trim}}{{end}}{{end}}{{else}}{{range $group := .Groups}}

{{$group.Title}}{{range $cmds}}{{if (and (eq .GroupID $group.ID) (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short | trim}}{{end}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}

Additional Commands:{{range $cmds}}{{if (and (eq .GroupID "") (or .IsAvailableCommand (eq .Name "help")))}}
  {{rpad .Name .NamePadding }} {{.Short | trim}}{{end}}{{end}}{{end}}{{end}}{{print "\n"}}{{end}}`

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

// getAuthCreds extracts auth credentials from the already-initialized client.
func getAuthCreds() (token, username, password string) {
	cl, err := client.Get()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: could not obtain credentials for --wait: %v\n", err)
		return "", "", ""
	}
	cfg := cl.CloudClient.GetConfig()
	return cfg.Token, cfg.Username, cfg.Password
}

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
