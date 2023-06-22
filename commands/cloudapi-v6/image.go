package commands

import (
	"bufio"
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/die"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"golang.org/x/exp/slices"

	"github.com/ionos-cloud/ionosctl/v6/internal/functional"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/printer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/multierr"
	"golang.org/x/sync/errgroup"
)

func ImageCmd() *core.Command {
	ctx := context.TODO()
	imageCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "image",
			Aliases:          []string{"img"},
			Short:            "Image Operations",
			Long:             "The sub-commands of `ionosctl image` allow you to see information about the Images available.",
			TraverseChildren: true,
		},
	}
	globalFlags := imageCmd.GlobalFlags()
	globalFlags.StringSliceP(constants.ArgCols, "", defaultImageCols, printer.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetFlagName(imageCmd.Name(), constants.ArgCols), globalFlags.Lookup(constants.ArgCols))
	_ = imageCmd.Command.RegisterFlagCompletionFunc(constants.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return allImageCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	list := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "list",
		Aliases:    []string{"l", "ls"},
		ShortDesc:  "List Images",
		LongDesc:   "Use this command to get a full list of available public Images.\n\nYou can filter the results using `--filters` option. Use the following format to set filters: `--filters KEY1=VALUE1,KEY2=VALUE2`.\n" + completer.ImagesFiltersUsage(),
		Example:    listImagesExample,
		PreCmdRun:  PreRunImageList,
		CmdRun:     RunImageList,
		InitClient: true,
	})
	list.AddStringFlag(cloudapiv6.ArgType, "", "", "The type of the Image", core.DeprecatedFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CDROM", "HDD"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgLicenceType, "", "", "The licence type of the Image", core.DeprecatedFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"WINDOWS", "WINDOWS2016", "LINUX", "OTHER", "UNKNOWN"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "", "The location of the Image", core.DeprecatedFlagOption())
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgImageAlias, "", "", "Image Alias or part of Image Alias to sort Images by", core.DeprecatedFlagOption())
	list.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed", core.DeprecatedFlagOption())
	list.AddInt32Flag(constants.FlagMaxResults, constants.FlagMaxResultsShort, cloudapiv6.DefaultMaxResults, constants.DescMaxResults)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "get",
		Aliases:    []string{"g"},
		ShortDesc:  "Get a specified Image",
		LongDesc:   "Use this command to get information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    getImageExample,
		PreCmdRun:  PreRunImageId,
		CmdRun:     RunImageGet,
		InitClient: true,
	})
	get.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = get.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	get.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	update := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "update",
		Aliases:    []string{"u", "up"},
		ShortDesc:  "Update a specified Image",
		LongDesc:   "Use this command to update information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    "ionosctl image update --image-id IMAGE_ID",
		PreCmdRun:  PreRunImageId,
		CmdRun:     RunImageUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	update.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end

	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")
	update.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	addPropertiesFlags := func(command *core.Command) {
		command.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Image")
		command.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Image")
		command.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", []string{"UNKNOWN", "WINDOWS", "WINDOWS2016", "WINDOWS2022", "LINUX", "OTHER"}, "The OS type of this image")
		command.AddSetFlag("cloud-init", "", "V1", []string{"V1", "NONE"}, "Cloud init compatibility")
		command.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", true, "'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
		command.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", true, "'Hot-Plug' RAM")
		command.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", true, "'Hot-Plug' NIC")
		command.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", true, "'Hot-Plug' Virt-IO drive")
		command.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", true, "'Hot-Plug' SCSI drive")
		command.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
		command.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "'Hot-Unplug' RAM")
		command.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "'Hot-Unplug' NIC")
		command.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "'Hot-Unplug' Virt-IO drive")
		command.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "'Hot-Unplug' SCSI drive")
	}

	addPropertiesFlags(update)

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "delete",
		Aliases:    []string{"d"},
		ShortDesc:  "Delete an image",
		LongDesc:   "Use this command to delete a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    "ionosctl image delete --image-id IMAGE_ID",
		PreCmdRun:  PreRunImageDelete,
		CmdRun:     RunImageDelete,
		InitClient: true,
	})
	deleteCmd.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all non-public images")

	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")
	deleteCmd.AddBoolFlag(constants.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		upload Command
		Perform upload to given FTP server. All locations from `location list` + karlsruhe (fkb)
		- ftp://ftp-fkb.ionos.com/hdd-images
		- ftp://ftp-fkb.ionos.com/iso-images
		https://docs.ionos.com/cloud/compute-engine/block-storage/block-storage-faq#how-do-i-upload-my-own-images-with-ftp
	*/
	upload := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "upload",
		Aliases:    []string{"ftp-upload", "ftp", "upl"},
		ShortDesc:  "Upload an image to FTP server using FTP over TLS (FTPS)",
		LongDesc:   "Use this command to upload an HDD or ISO image.\n\nRequired values to run command:\n\n* Location\n",
		Example:    "ionosctl img u -i kolibri.iso -l fkb,fra,vit",
		PreCmdRun:  PreRunImageUpload,
		CmdRun:     RunImageUpload,
		InitClient: true,
	})

	upload.AddStringFlag("ftp-user", "", "", "Override username for FTP server")
	upload.AddStringFlag("ftp-pass", "", "", "Override password for FTP server")

	upload.AddStringSliceFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, nil, "Location to upload to. Must be an array containing only fra, fkb, txl, lhr, las, ewr, vit", core.RequiredFlagOption())
	upload.AddStringSliceFlag("image", "i", nil, "Slice of paths to images, can be absolute path or relative to current working directory", core.RequiredFlagOption())
	upload.AddStringFlag("ftp-url", "", "ftp-%s.ionos.com", "URL of FTP server, with %s flag if location is embedded into url")
	upload.AddBoolFlag("skip-verify", "", false, "Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!")
	upload.AddBoolFlag("skip-update", "", false, "After the image is uploaded to the FTP server, send a PATCH to the API with the contents of the image properties flags and emulate a \"create\" command.")
	upload.AddStringFlag("crt-path", "", "", "(Unneeded for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.")
	upload.AddStringSliceFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, nil, "Rename the uploaded images. These names should not contain any extension. By default, this is the base of the image path")
	upload.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, 300, "(seconds) Context Deadline. FTP connection will time out after this many seconds")

	addPropertiesFlags(upload)

	upload.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end
	upload.Command.SilenceUsage = true       // Don't print help if setting only 1 out of 2 required flags - too many flags. Help must be invoked manually via --help

	return imageCmd
}

func PreRunImageDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{cloudapiv6.ArgImageId}, []string{cloudapiv6.ArgAll})
}

func RunImageDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNonPublicImages(c); err != nil {
			return err
		}
		return c.Printer.Print(printer.Result{Resource: c.Resource, Verb: c.Verb})
	} else {
		listQueryParams, err := query.GetListQueryParams(c)
		if err != nil {
			return err
		}
		queryParams := listQueryParams.QueryParams
		if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete image"); err != nil {
			return err
		}
		imgId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
		c.Printer.Verbose("Starting deletion on image with ID: %v...", imgId)
		resp, err := c.CloudApiV6Services.Images().Delete(imgId, queryParams)
		if resp != nil && printer.GetId(resp) != "" {
			c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
		}
		if err != nil {
			return err
		}
		if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
			return err
		}
		return c.Printer.Print(getImagePrint(resp, c, nil))
	}
}

// Util func - Given a slice of public & non-public images, return only those images that are non-public.
// If any image in the slice has null properties, or "Properties.Public" field is nil, the image is skipped (and a verbose message is shown)
func getNonPublicImages(imgs []ionoscloud.Image, printService printer.PrintService) ([]ionoscloud.Image, error) {
	var nonPublicImgs []ionoscloud.Image
	for _, i := range imgs {
		properties, ok := i.GetPropertiesOk()
		if !ok {
			printService.Verbose("skipping %s: properties are nil\n", *i.GetId())
			continue
		}
		isPublic, ok := properties.GetPublicOk()
		if !ok {
			printService.Verbose("skipping %s: field `public` is nil\n", *i.GetId())
			continue
		}
		if !*isPublic {
			nonPublicImgs = append(nonPublicImgs, i)
		}
	}
	return nonPublicImgs, nil
}

// DeleteAllNonPublicImages deletes non-public images, as deleting public images is forbidden by the API.
func DeleteAllNonPublicImages(c *core.CommandConfig) error {
	depth := int32(1)
	images, resp, err := c.CloudApiV6Services.Images().List(
		resources.ListQueryParams{QueryParams: resources.QueryParams{Depth: &depth}},
	)
	if err != nil {
		return err
	}
	allItems, ok := images.GetItemsOk()
	if !(ok && len(*allItems) > 0 && allItems != nil) {
		return errors.New("could not retrieve images")
	}

	items, err := getNonPublicImages(*allItems, c.Printer)
	if err != nil {
		return err
	}
	if len(items) < 1 {
		return errors.New("no non-public images found")
	}

	_ = c.Printer.Warn("Images to be deleted:")
	// TODO: this is duplicated across all resources - refactor this (across all resources)
	for _, img := range items {
		delIdAndName := ""
		if id, ok := img.GetIdOk(); ok && id != nil {
			delIdAndName += "ID: `" + *id
		}
		if properties, ok := img.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				delIdAndName += "`, Name: " + *name
			}
		}
		_ = c.Printer.Warn(delIdAndName)
	}

	if err = utils.AskForConfirm(c.Stdin, c.Printer, "delete all the images"); err != nil {
		return err
	}
	c.Printer.Verbose("Deleting all the images...")

	var multiErr error
	for _, img := range items {
		if id, ok := img.GetIdOk(); ok && id != nil {
			c.Printer.Verbose("Starting deleting image with id: %v...", *id)
			resp, err = c.CloudApiV6Services.Images().Delete(*id, resources.QueryParams{})
			if resp != nil && printer.GetId(resp) != "" {
				c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
			}
			if err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
				continue
			} else {
				_ = c.Printer.Warn(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
			}
			if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
				multiErr = multierr.Append(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
				continue
			}
		}
	}
	if multiErr != nil {
		return multiErr
	}
	return nil
}

// returns an ImageProperties object which reflects the currently set flags
func getDesiredImageAfterPatch(c *core.CommandConfig) resources.ImageProperties {
	input := resources.ImageProperties{}
	c.Command.Command.Flags().VisitAll(func(flag *pflag.Flag) {
		val := flag.Value.String()
		if val == "" {
			return
		}
		boolval, _ := strconv.ParseBool(val)
		switch flag.Name {
		case cloudapiv6.ArgName:
			input.SetName(val)
			break
		case cloudapiv6.ArgDescription:
			input.SetDescription(val)
			break
		case "cloud-init":
			input.SetCloudInit(val)
			break
		case cloudapiv6.ArgLicenceType:
			input.SetLicenceType(val)
			break
		case cloudapiv6.ArgCpuHotPlug:
			input.SetCpuHotPlug(boolval)
			break
		case cloudapiv6.ArgRamHotPlug:
			input.SetRamHotPlug(boolval)
			break
		case cloudapiv6.ArgNicHotPlug:
			input.SetNicHotPlug(boolval)
			break
		case cloudapiv6.ArgDiscVirtioHotPlug:
			input.SetDiscVirtioHotPlug(boolval)
			break
		case cloudapiv6.ArgDiscScsiHotPlug:
			input.SetDiscScsiHotPlug(boolval)
			break
		case cloudapiv6.ArgCpuHotUnplug:
			input.SetCpuHotUnplug(boolval)
			break
		case cloudapiv6.ArgRamHotUnplug:
			input.SetRamHotUnplug(boolval)
			break
		case cloudapiv6.ArgNicHotUnplug:
			input.SetNicHotUnplug(boolval)
			break
		case cloudapiv6.ArgDiscVirtioHotUnplug:
			input.SetDiscVirtioHotUnplug(boolval)
			break
		case cloudapiv6.ArgDiscScsiHotUnplug:
			input.SetDiscScsiHotUnplug(boolval)
			break
		default:
			// --image-id, verbose, filters, depth, etc
			break
		}
		c.Printer.Verbose(fmt.Sprintf("Property %s set: %s", flag.Name, flag.Value))
	})
	return input
}

func RunImageUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams

	input := getDesiredImageAfterPatch(c)
	img, resp, err := c.CloudApiV6Services.Images().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(constants.MessageRequestInfo, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, []resources.Image{*img}))
}

func PreRunImageUpload(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired("image")
	if err != nil {
		return err
	}

	ftpUser := viper.GetString(core.GetFlagName(c.NS, "ftp-user"))
	ftpPass := viper.GetString(core.GetFlagName(c.NS, "ftp-pass"))

	errHandler := func(err error) {
		errMsg := c.Printer.Warn("Warn: " +
			fmt.Errorf(
				"failed trying to use standard client credentials for FTP server: %w",
				err,
			).Error(),
		)
		die.Die(errMsg.Error())
	}
	if ftpUser == "" || ftpPass == "" {
		if client.Must(errHandler).IsTokenAuth() {
			c.Printer.Warn("Warn: You are using token authentication for your client. FTP connection does not support JWT token")
			c.Printer.Warn("Warn: You can use --ftp-user and --ftp-pass to override the client's credentials for the FTP server only")
		}
	}

	validExts := []string{".iso", ".img", ".vmdk", ".vhd", ".vhdx", ".cow", ".qcow", ".qcow2", ".raw", ".vpc", ".vdi"}
	images := viper.GetStringSlice(core.GetFlagName(c.NS, "image"))
	invalidImages := functional.Filter(
		functional.Map(images, func(s string) string {
			return filepath.Ext(s)
		}),
		func(ext string) bool {
			return !slices.Contains(
				validExts,
				ext,
			)
		},
	)
	if len(invalidImages) > 0 {
		return fmt.Errorf("%s is an invalid image extension. Valid extensions are: %s", strings.Join(invalidImages, ","), validExts)
	}

	// "Locations" flag only required if ftp-url custom flag contains a %s in which to add the location ID
	if strings.Contains(viper.GetString(core.GetFlagName(c.NS, "ftp-url")), "%s") {
		err = c.Command.Command.MarkFlagRequired(cloudapiv6.ArgLocation)
		if err != nil {
			return err
		}
	}

	validLocs := []string{"fra", "fkb", "txl", "lhr", "las", "ewr", "vit"}
	locs := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	invalidLocs := functional.Filter(
		locs,
		func(loc string) bool {
			return !slices.Contains(
				validLocs,
				loc,
			)
		},
	)
	if len(invalidLocs) > 0 {
		c.Printer.Verbose("WARN: %s is an invalid location. Valid IONOS locations are: %s", strings.Join(invalidLocs, ","), locs)
	}

	aliases := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
	if len(aliases) != 0 && len(aliases) != len(images) {
		return fmt.Errorf("slices of image files and image aliases are of different lengths. Uploading multiple images with the same alias is forbidden")
	}

	return nil
}

// Reads server certificate at given path.
// If path unset, returns nil.
// Otherwise, returns certificate pool containing server certificate
func getCertificate(path string) (*x509.CertPool, error) {
	if path == "" {
		return nil, nil
	}
	caCert, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return caCertPool, nil
}

func updateImagesAfterUpload(c *core.CommandConfig, diffImgs []ionoscloud.Image, properties resources.ImageProperties) ([]ionoscloud.Image, error) {
	// do a patch on the uploaded images
	var imgs []ionoscloud.Image
	for _, diffImg := range diffImgs {
		img, _, err := client.Must().CloudClient.ImagesApi.ImagesPatch(c.Context, *diffImg.GetId()).Image(properties.ImageProperties).Execute()
		imgs = append(imgs, img)
		if err != nil {
			return nil, err
		}
	}
	return imgs, nil
}

// getDiffUploadedImages will keep querying /images endpoint until the images with the given names and locations show up.
func getDiffUploadedImages(c *core.CommandConfig, names, locations []string) ([]ionoscloud.Image, error) {
	var diffImgs []ionoscloud.Image

	for {
		select {
		case <-c.Context.Done():
			return nil, fmt.Errorf("ran out of time: %w", c.Context.Err())
		default:
			req := client.Must().CloudClient.ImagesApi.ImagesGet(c.Context).Depth(1).Filter("public", "false")
			for _, n := range names {
				req.Filter("name", n)
			}
			for _, l := range locations {
				req.Filter("location", l)
			}

			imgs, _, err := req.Execute()
			if err != nil {
				return nil, fmt.Errorf("failed listing images")
			}
			c.Printer.Verbose("Got images by listing: %+v", *imgs.Items)

			diffImgs = append(diffImgs, *imgs.Items...)
			c.Printer.Verbose("Total images: %+v", diffImgs)

			if len(diffImgs) == len(names)*len(locations) {
				c.Printer.Verbose("Success! All images found via API: %+v", diffImgs)
				return diffImgs, nil
			}

			// New attempt...
			time.Sleep(10 * time.Second)
		}
	}
}

func RunImageUpload(c *core.CommandConfig) error {
	certPool, err := getCertificate(viper.GetString(core.GetFlagName(c.NS, "crt-path")))
	if err != nil {
		return err
	}

	images := viper.GetStringSlice(core.GetFlagName(c.NS, "image"))
	aliases := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
	locations := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	skipVerify := viper.GetBool(core.GetFlagName(c.NS, "skip-verify"))
	ftpUser := viper.GetString(core.GetFlagName(c.NS, "ftp-user"))
	ftpPass := viper.GetString(core.GetFlagName(c.NS, "ftp-pass"))

	errHandler := func(err error) {
		errMsg := fmt.Errorf(
			"failed trying to use standard client credentials for FTP server: %w",
			err,
		)
		die.Die(errMsg.Error())
	}
	if ftpUser == "" || ftpPass == "" {
		cfg := client.Must(errHandler).CloudClient.GetConfig()
		ftpUser = cfg.Username
		ftpPass = cfg.Password
	}

	ctx, cancel := context.WithTimeout(c.Context, time.Duration(viper.GetInt(core.GetFlagName(c.NS, constants.ArgTimeout)))*time.Second)
	defer cancel()
	c.Context = ctx

	var eg errgroup.Group
	for _, loc := range locations {
		for imgIdx, img := range images {
			url := viper.GetString(core.GetFlagName(c.NS, "ftp-url"))
			if strings.Contains(url, "%s") {
				url = fmt.Sprintf(url, loc) // Add the location modifier, if the URL supports it
			}
			c.Printer.Verbose("Uploading %s to %s", img, url)

			var isoOrHdd string
			if ext := filepath.Ext(img); ext == ".iso" || ext == ".img" {
				isoOrHdd = "iso"
			} else {
				isoOrHdd = "hdd"
			}

			serverFilePath := fmt.Sprintf("%s-images/", isoOrHdd) // iso-images / hdd-images
			if len(aliases) == 0 {
				serverFilePath += filepath.Base(img) // If no custom alias, use the filename
			} else {
				serverFilePath += aliases[imgIdx] + filepath.Ext(img) // Use custom alias
			}

			file, err := os.Open(img)
			if err != nil {
				return err
			}

			data := bufio.NewReader(file)
			// Catching error from goroutines. https://stackoverflow.com/questions/62387307/how-to-catch-errors-from-goroutines
			// Uploads each image to each location.
			eg.Go(func() error {

				err := c.CloudApiV6Services.Images().Upload(
					c.Context,
					resources.UploadProperties{
						FTPServerProperties: resources.FTPServerProperties{
							Url:               url,
							Port:              21,
							SkipVerify:        skipVerify,
							ServerCertificate: certPool,
							Username:          ftpUser,
							Password:          ftpPass,
						},
						ImageFileProperties: resources.ImageFileProperties{
							Path:       serverFilePath,
							DataBuffer: data,
						},
					},
				)
				if err != nil {
					return err
				}
				return file.Close()
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	if viper.GetBool(core.GetFlagName(c.NS, "skip-update")) {
		c.Printer.Verbose("Successfully uploaded images")
		return nil
	}

	names := images
	if len(aliases) != 0 {
		// Returns a slice containing `alias[i] + filepath.Ext(images[i])`
		// (i.e it gets the extensions from `images` flag, and appends them to each elem `image-alias`)
		// Resulting slice is the full image names, as returned by `ionosctl image list` on the Name column
		names = functional.MapIdx(aliases, func(k int, v string) string {
			return v + filepath.Ext(images[k])
		})
	}
	diffImgs, err := getDiffUploadedImages(c, names, locations) // Get UUIDs of uploaded images
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP sucessful: %w", err)
	}

	properties := getDesiredImageAfterPatch(c)
	imgs, err := updateImagesAfterUpload(c, diffImgs, properties)
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP sucessful: %w", err)
	}

	c.Printer.Verbose("Successfully uploaded and updated images")
	// oh my lord, we need to get rid of the `resources` wrappers...
	return c.Printer.Print(getImagePrint(nil, c, getImages(resources.Images{Images: ionoscloud.Images{Items: &imgs}})))
}

func PreRunImageList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ImagesFilters(), completer.ImagesFiltersUsage())
	}
	return nil
}

func PreRunImageId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgImageId)
}

func RunImageList(c *core.CommandConfig) error {
	_ = c.Printer.Warn("WARNING: The following flags are deprecated:" + c.Command.GetAnnotationsByKey(core.DeprecatedFlagsAnnotation) + ". Use --filters --order-by --max-results options instead!")
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	images, resp, err := c.CloudApiV6Services.Images().List(listQueryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)) {
		images = sortImagesByLocation(images, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLocation)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType)) {
		images = sortImagesByLicenceType(images, strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgLicenceType))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgType)) {
		images = sortImagesByType(images, strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgType))))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		images = sortImagesByAlias(images, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)))
	}
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)) {
		images = sortImagesByTime(images, viper.GetInt(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)))
	}
	if itemsOk, ok := images.GetItemsOk(); ok && itemsOk != nil {
		if len(*itemsOk) == 0 {
			return errors.New("error getting images based on given criteria")
		}
	}
	return c.Printer.Print(getImagePrint(nil, c, getImages(images)))
}

func RunImageGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams
	c.Printer.Verbose("Image with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))
	img, resp, err := c.CloudApiV6Services.Images().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)), queryParams)
	if resp != nil {
		c.Printer.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getImage(img)))
}

// Output Printing

var (
	defaultImageCols = []string{"ImageId", "Name", "ImageAliases", "Location", "LicenceType", "ImageType", "CloudInit", "CreatedDate"}
	allImageCols     = []string{"ImageId", "Name", "ImageAliases", "Location", "Size", "LicenceType", "ImageType", "Description", "Public", "CloudInit", "CreatedDate", "CreatedBy", "CreatedByUserId"}
)

type ImagePrint struct {
	ImageId         string    `json:"ImageId,omitempty"`
	Name            string    `json:"Name,omitempty"`
	Description     string    `json:"Description,omitempty"`
	Location        string    `json:"Location,omitempty"`
	Size            string    `json:"Size,omitempty"`
	LicenceType     string    `json:"LicenceType,omitempty"`
	ImageType       string    `json:"ImageType,omitempty"`
	Public          bool      `json:"Public,omitempty"`
	ImageAliases    []string  `json:"ImageAliases,omitempty"`
	CloudInit       string    `json:"CloudInit,omitempty"`
	CreatedBy       string    `json:"CreatedBy,omitempty"`
	CreatedByUserId string    `json:"CreatedByUserId,omitempty"`
	CreatedDate     time.Time `json:"CreatedDate,omitempty"`
}

func getImagePrint(resp *resources.Response, c *core.CommandConfig, imgs []resources.Image) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, constants.ArgWaitForRequest))
		}
		if imgs != nil {
			r.OutputJSON = imgs
			r.KeyValue = getImagesKVMaps(imgs)
			r.Columns = printer.GetHeaders(allImageCols, defaultImageCols, viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols)))
		}
	}
	return r
}

func getImages(images resources.Images) []resources.Image {
	imgs := make([]resources.Image, 0)
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, d := range *items {
			imgs = append(imgs, resources.Image{Image: d})
		}
	}
	return imgs
}

func getImage(image *resources.Image) []resources.Image {
	imgs := make([]resources.Image, 0)
	if image != nil {
		imgs = append(imgs, resources.Image{Image: image.Image})
	}
	return imgs
}

func getImagesKVMaps(imgs []resources.Image) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(imgs))
	for _, img := range imgs {
		o := getImageKVMap(img)
		out = append(out, o)
	}
	return out
}

func getImageKVMap(img resources.Image) map[string]interface{} {
	var imgPrint ImagePrint
	if idOk, ok := img.GetIdOk(); ok && idOk != nil {
		imgPrint.ImageId = *idOk
	}
	if propertiesOk, ok := img.GetPropertiesOk(); ok && propertiesOk != nil {
		if name, ok := propertiesOk.GetNameOk(); ok && name != nil {
			imgPrint.Name = *name
		}
		if description, ok := propertiesOk.GetDescriptionOk(); ok && description != nil {
			imgPrint.Description = *description
		}
		if loc, ok := propertiesOk.GetLocationOk(); ok && loc != nil {
			imgPrint.Location = *loc
		}
		if size, ok := propertiesOk.GetSizeOk(); ok && size != nil {
			imgPrint.Size = fmt.Sprintf("%vGB", *size)
		}
		if licType, ok := propertiesOk.GetLicenceTypeOk(); ok && licType != nil {
			imgPrint.LicenceType = *licType
		}
		if imgType, ok := propertiesOk.GetImageTypeOk(); ok && imgType != nil {
			imgPrint.ImageType = *imgType
		}
		if public, ok := propertiesOk.GetPublicOk(); ok && public != nil {
			imgPrint.Public = *public
		}
		if aliases, ok := propertiesOk.GetImageAliasesOk(); ok && aliases != nil {
			imgPrint.ImageAliases = *aliases
		}
		if cloudInit, ok := propertiesOk.GetCloudInitOk(); ok && cloudInit != nil {
			imgPrint.CloudInit = *cloudInit
		}
	}
	if metadataOk, ok := img.GetMetadataOk(); ok && metadataOk != nil {
		if createdDateOk, ok := metadataOk.GetCreatedDateOk(); ok && createdDateOk != nil {
			imgPrint.CreatedDate = *createdDateOk
		}
		if createdByOk, ok := metadataOk.GetCreatedByOk(); ok && createdByOk != nil {
			imgPrint.CreatedBy = *createdByOk
		}
		if createdByUserIdOk, ok := metadataOk.GetCreatedByUserIdOk(); ok && createdByUserIdOk != nil {
			imgPrint.CreatedByUserId = *createdByUserIdOk
		}
	}
	return structs.Map(imgPrint)
}

// Output Columns Sorting

func sortImagesByLocation(images resources.Images, location string) resources.Images {
	imgLocationItems := make([]ionoscloud.Image, 0)
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *items {
			properties := img.GetProperties()
			if loc, ok := properties.GetLocationOk(); ok && loc != nil {
				if *loc == location {
					imgLocationItems = append(imgLocationItems, img)
				}
			}
		}
	}
	images.Items = &imgLocationItems
	return images
}

func sortImagesByLicenceType(images resources.Images, licenceType string) resources.Images {
	imgLicenceTypeItems := make([]ionoscloud.Image, 0)
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *items {
			properties := img.GetProperties()
			if imgLicenceType, ok := properties.GetLicenceTypeOk(); ok && imgLicenceType != nil {
				if *imgLicenceType == licenceType {
					imgLicenceTypeItems = append(imgLicenceTypeItems, img)
				}
			}
		}
	}
	images.Items = &imgLicenceTypeItems
	return images
}

func sortImagesByType(images resources.Images, imgType string) resources.Images {
	imgTypeItems := make([]ionoscloud.Image, 0)
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *items {
			properties := img.GetProperties()
			if t, ok := properties.GetImageTypeOk(); ok && t != nil {
				if *t == imgType {
					imgTypeItems = append(imgTypeItems, img)
				}
			}
		}
	}
	images.Items = &imgTypeItems
	return images
}

func sortImagesByAlias(images resources.Images, alias string) resources.Images {
	imgTypeItems := make([]ionoscloud.Image, 0)
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *items {
			properties := img.GetProperties()
			if imageAliasesOk, ok := properties.GetImageAliasesOk(); ok && imageAliasesOk != nil {
				for _, imageAliaseOk := range *imageAliasesOk {
					if strings.Contains(imageAliaseOk, alias) {
						imgTypeItems = append(imgTypeItems, img)
					}
				}
			}
		}
	}
	images.Items = &imgTypeItems
	return images
}

func sortImagesByTime(images resources.Images, n int) resources.Images {
	if items, ok := images.GetItemsOk(); ok && items != nil {
		imageItems := *items
		if len(imageItems) > 0 {
			// Sort Requests using time.Time, in descending order
			sort.SliceStable(imageItems, func(i, j int) bool {
				return imageItems[i].Metadata.CreatedDate.Time.After(imageItems[j].Metadata.CreatedDate.Time)
			})
		}
		if len(imageItems) >= n {
			imageItems = imageItems[:n]
		}
		images.Items = &imageItems
	}
	return images
}
