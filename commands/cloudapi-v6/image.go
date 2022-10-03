package commands

import (
	"context"
	"errors"
	"fmt"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/spf13/pflag"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/printer"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	globalFlags.StringSliceP(config.ArgCols, "", defaultImageCols, printer.ColsMessage(allImageCols))
	_ = viper.BindPFlag(core.GetGlobalFlagName(imageCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = imageCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	list.AddInt32Flag(cloudapiv6.ArgMaxResults, cloudapiv6.ArgMaxResultsShort, cloudapiv6.DefaultMaxResults, cloudapiv6.ArgMaxResultsDescription)
	list.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultListDepth, cloudapiv6.ArgDepthDescription)
	list.AddStringFlag(cloudapiv6.ArgOrderBy, "", "", cloudapiv6.ArgOrderByDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgOrderBy, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringSliceFlag(cloudapiv6.ArgFilters, cloudapiv6.ArgFiltersShort, []string{""}, cloudapiv6.ArgFiltersDescription)
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgFilters, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImagesFilters(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)

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
	get.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	get.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	update := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace:  "image",
		Resource:   "image",
		Verb:       "update",
		Aliases:    []string{"g"},
		ShortDesc:  "Update a specified Image",
		LongDesc:   "Use this command to get information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		Example:    "placeholder", // TODO
		PreCmdRun:  PreRunImageId,
		CmdRun:     RunImageUpdate,
		InitClient: true,
	})
	update.AddUUIDFlag(cloudapiv6.ArgImageId, cloudapiv6.ArgIdShort, "", cloudapiv6.ImageId, core.RequiredFlagOption())
	_ = update.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.ImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	update.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end

	update.AddBoolFlag(config.ArgWaitForRequest, config.ArgWaitForRequestShort, config.DefaultWait, "Wait for the Request for Image update to be executed")
	update.AddIntFlag(config.ArgTimeout, config.ArgTimeoutShort, config.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")
	update.AddBoolFlag(config.ArgNoHeaders, "", false, cloudapiv6.ArgNoHeadersDescription)
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	update.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Image")
	update.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Image")
	update.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", []string{"UNKNOWN", "WINDOWS", "WINDOWS2016", "WINDOWS2022", "LINUX", "OTHER"}, "The OS type of this image")
	update.AddSetFlag("cloud-init", "", "V1", []string{"V1", "NONE"}, "Cloud init compatibility")
	update.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", true, "'Hot-Plug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
	update.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", true, "'Hot-Plug' RAM")
	update.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", true, "'Hot-Plug' NIC")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", true, "'Hot-Plug' Virt-IO drive")
	update.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", true, "'Hot-Plug' SCSI drive")
	update.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "'Hot-Unplug' CPU. It is not possible to have a hot-unplug CPU which you previously did not hot-plug")
	update.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "'Hot-Unplug' RAM")
	update.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "'Hot-Unplug' NIC")
	update.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "'Hot-Unplug' Virt-IO drive")
	update.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "'Hot-Unplug' SCSI drive")

	return imageCmd
}

func RunImageUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams

	input := resources.ImageProperties{}
	c.Command.Command.Flags().Visit(func(flag *pflag.Flag) {
		val := flag.Value.String()
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
	img, resp, err := c.CloudApiV6Services.Images().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)),
		input,
		queryParams,
	)
	if resp != nil && printer.GetId(resp) != "" {
		c.Printer.Verbose(config.RequestInfoMessage, printer.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, waiter.RequestInterrogator, printer.GetId(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, []resources.Image{*img}))
}

func PreRunImageList(c *core.PreCommandConfig) error {
	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgFilters)) {
		return query.ValidateFilters(c, completer.ImagesFilters(), completer.ImagesFiltersUsage())
	}
	return nil
}

func PreRunImageId(c *core.PreCommandConfig) error {
	if err := c.Command.Command.MarkFlagRequired(cloudapiv6.ArgImageId); err != nil {
		return err
	}

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
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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
		c.Printer.Verbose(config.RequestTimeMessage, resp.RequestTime)
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
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if imgs != nil {
			r.OutputJSON = imgs
			r.KeyValue = getImagesKVMaps(imgs)
			r.Columns = getImageCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getImageCols(flagName string, outErr io.Writer) []string {
	var cols []string
	if viper.IsSet(flagName) {
		cols = viper.GetStringSlice(flagName)
	} else {
		return defaultImageCols
	}

	columnsMap := map[string]string{
		"ImageId":         "ImageId",
		"Name":            "Name",
		"Description":     "Description",
		"Location":        "Location",
		"Size":            "Size",
		"LicenceType":     "LicenceType",
		"ImageType":       "ImageType",
		"Public":          "Public",
		"ImageAliases":    "ImageAliases",
		"CloudInit":       "CloudInit",
		"CreatedDate":     "CreatedDate",
		"CreatedBy":       "CreatedBy",
		"CreatedByUserId": "CreatedByUserId",
	}
	var datacenterCols []string
	for _, k := range cols {
		col := columnsMap[k]
		if col != "" {
			datacenterCols = append(datacenterCols, col)
		} else {
			clierror.CheckError(errors.New("unknown column "+k), outErr)
		}
	}
	return datacenterCols
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
