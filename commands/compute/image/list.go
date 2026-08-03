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
		Namespace: "image",
		Resource:  "image",
		Verb:      "list",
		Aliases:   []string{"l", "ls"},
		ShortDesc: "List public and private Images",
		LongDesc: `List the Images your contract can see: IONOS-provided PUBLIC images (OS distributions and ISOs) and your own PRIVATE images (uploaded or snapshotted). The Public column distinguishes the two; the ImageAliases column shows the stable names (e.g. "ubuntu:latest") you can pass instead of a UUID when creating a volume.

Because public images are replicated per location, the same OS appears once per region (see the Location column, e.g. de/fra, gb/lhr). Narrow the result with server-side filters.

You can filter the results using the ` + "`--filters`" + ` option. Use the following format to set filters: ` + "`--filters KEY1=VALUE1,KEY2=VALUE2`" + `. To list only your own images, filter on ` + "`public=false`" + `.
` + completer.ImagesFiltersUsage(),
		Example:    listImagesExample,
		PreCmdRun:  core.NoPreRun,
		CmdRun:     RunImageList,
		InitClient: true,
	})

	deprecatedMessage := "incompatible with --max-results. Use --filters --order-by --max-results options instead!"

	cmd.AddStringFlag(constants.FlagType, "", "", "Client-side filter by image type: HDD (a bootable disk image) or CDROM (an ISO you can attach as a virtual optical drive)", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CDROM", "HDD"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLicenceType, "", "", "Client-side filter by OS licence type (LINUX, RHEL, WINDOWS, WINDOWS2016/2019/2022/2025, UNKNOWN, OTHER). This is how the platform bills and configures the guest OS", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return constants.EnumLicenceType, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "", "Client-side filter by the physical location an image lives in, e.g. de/fra, de/txl, gb/lhr. Public images are replicated per location, so the same OS appears once per region", core.DeprecatedFlagOption(deprecatedMessage))
	_ = cmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(cloudapiv6.ArgImageAlias, "", "", "Client-side filter keeping only images whose image-alias contains this substring. An image-alias (e.g. \"ubuntu:20.04\", \"debian:latest\") is a stable human-friendly name you can use in place of a UUID when creating a volume", core.DeprecatedFlagOption(deprecatedMessage))
	cmd.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Client-side: keep only the N most recently created images (by createdDate, newest first). 0 (default) keeps all. Prefer --order-by createdDate --max-results N for server-side ordering", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))

	return cmd
}
