package commands

import (
	"fmt"
	"os"

	"github.com/common-nighthawk/go-figure"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// RootCmd is the root level command that all other commands attach to
	rootCmd = &cobra.Command{
		Use:              "ionosctl",
		Short:            "ionosctl is a command line interface (CLI) for the Ionos Cloud SDK",
		Long:             asciiLogo.String(),
		TraverseChildren: true,
	}

	ServerURL string
	Output    string
	Verbose   bool

	cfgFile   string
	asciiLogo = figure.NewFigure("ionosctl", "slant", true)
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	initConfig()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootPFlagSet := rootCmd.PersistentFlags()
	rootPFlagSet.StringVarP(&cfgFile, config.ArgConfig, "c", config.GetConfigFilePath(), "Specify a custom config file")
	viper.BindPFlag(config.ArgConfig, rootPFlagSet.Lookup(config.ArgConfig))

	rootPFlagSet.StringVarP(&ServerURL, config.ArgServerUrl, "u", config.DefaultApiUrl, "Override default API endpoint")
	viper.BindPFlag(config.ArgServerUrl, rootPFlagSet.Lookup(config.ArgServerUrl))

	rootPFlagSet.StringVarP(&Output, config.ArgOutput, "o", config.DefaultOutputFormat, "Desired output format [text|json]")
	viper.BindPFlag(config.ArgOutput, rootPFlagSet.Lookup(config.ArgOutput))

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
		viper.SetConfigType("yaml")
	}

	viper.AutomaticEnv() // read in environment variables that match
	_ = viper.ReadInConfig()
}

// AddCommands adds sub commands to the base command.
func addCommands() {
	rootCmd.AddCommand(login())
	rootCmd.AddCommand(completion())
	rootCmd.AddCommand(list())
	rootCmd.AddCommand(create())
	rootCmd.AddCommand(update())
	rootCmd.AddCommand(delete())
}
