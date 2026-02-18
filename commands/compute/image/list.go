package image

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/spf13/cobra"
)

func ImageListCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Images",
		LongDesc:   "Use this command to get a full list of available public Images.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ImagesFiltersUsage(),
		Example:    listImagesExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunImageList,
		InitClient: true,
	})

	deprecatedMessage := "incompatible with --max-results. Use --filters --order-by --max-results options instead!"

	cmd.AddStringFlag(constants.FlagType, "", "", "The type of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CDROM", "HDD"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLicenceType, "", "", "The licence type of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return constants.EnumLicenceType, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "", "The location of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgImageAlias, "", "", "Image Alias or part of Image Alias to sort Images by", core.DeprecatedFlagOption(deprecatedMessage))
	cmd.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))

	return cmd
}
