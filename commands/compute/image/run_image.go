package image

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	FlagRenameImages    = "rename"
	FlagImage           = "image"
	FlagSkipUpdate      = "skip-update"
	FlagSkipVerify      = "skip-verify"
	FlagFtpUrl          = "ftp-url"
	FlagCertificatePath = "crt-path"
)

func addPropertiesFlags(command *core.Command) {
	command.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Name of the Image")
	command.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Description of the Image")
	command.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", constants.EnumLicenceType, "The OS type of this image")
	command.AddSetFlag(constants.FlagCloudInit, "", "V1", []string{"V1", "NONE"}, "Cloud init compatibility")
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
	command.AddBoolFlag(cloudapiv6.ArgExposeSerial, "", false, "If set to `true` will expose the serial id of the disk attached to the server")
	command.AddBoolFlag(cloudapiv6.ArgRequireLegacyBios, "", true, "Indicates if the image requires the legacy BIOS for compatibility or specific needs.")
	command.AddSetFlag(cloudapiv6.ArgApplicationType, "", "UNKNOWN", constants.EnumApplicationType, "The type of application that is hosted on this resource")
}

func PreRunImageDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS, []string{cloudapiv6.ArgImageId}, []string{cloudapiv6.ArgAll})
}

func RunImageDelete(c *core.CommandConfig) error {
	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllNonPublicImages(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete image", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	imgId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Starting deletion on image with ID: %v...", imgId))

	resp, err := c.CloudApiV6Services.Images().Delete(imgId)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", jsontabwriter.GenerateLogOutput("Image deleted successfully"))

	return nil
}

// DeleteAllNonPublicImages deletes non-public images, as deleting public images is forbidden by the API.
func DeleteAllNonPublicImages(c *core.CommandConfig) error {
	images, resp, err := c.CloudApiV6Services.Images().List()
	if err != nil {
		return err
	}
	allItems, ok := images.GetItemsOk()
	if !(ok && len(*allItems) > 0 && allItems != nil) {
		return errors.New("could not retrieve images")
	}

	items, err := getNonPublicImages(*allItems, c.Command.Command.ErrOrStderr())
	if err != nil {
		return err
	}
	if len(items) < 1 {
		return errors.New("no non-public images found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateLogOutput("Images to be deleted:"))
	// TODO: this is duplicated across all resources - refactor this (across all resources)
	var multiErr error
	for _, img := range items {
		id := img.GetId()
		name := img.GetProperties().Name

		if !confirm.FAsk(c.Command.Command.InOrStdin(), fmt.Sprintf("Delete the image with Id: %s , Name: %s", *id, *name), viper.GetBool(constants.ArgForce)) {
			return fmt.Errorf(confirm.UserDenied)
		}

		resp, err = c.CloudApiV6Services.Images().Delete(*id)
		if resp != nil && request.GetId(resp) != "" {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
		}
		if err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		} else {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s",
				jsontabwriter.GenerateLogOutput("%s", fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id)))
		}

		if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
			multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
			continue
		}
	}

	if multiErr != nil {
		return multiErr
	}

	return nil
}

// Util func - Given a slice of public & non-public images, return only those images that are non-public.
// If any image in the slice has null properties, or "Properties.Public" field is nil, the image is skipped (and a verbose message is shown)
func getNonPublicImages(imgs []ionoscloud.Image, verboseOut io.Writer) ([]ionoscloud.Image, error) {
	var nonPublicImgs []ionoscloud.Image
	for _, i := range imgs {
		properties, ok := i.GetPropertiesOk()
		if !ok {
			fmt.Fprintf(verboseOut, "%s", jsontabwriter.GenerateVerboseOutput("skipping %s: properties are nil\n", *i.GetId()))
			continue
		}

		isPublic, ok := properties.GetPublicOk()
		if !ok {
			fmt.Fprintf(verboseOut, "%s", jsontabwriter.GenerateVerboseOutput("skipping %s: field `public` is nil\n", *i.GetId()))
			continue
		}

		if !*isPublic {
			nonPublicImgs = append(nonPublicImgs, i)
		}
	}
	return nonPublicImgs, nil
}

// returns an ImageProperties object which reflects the currently set flags
func getDesiredImageAfterPatch(c *core.CommandConfig, useUnsetFlags bool) resources.ImageProperties {
	input := resources.ImageProperties{}

	// flagTraverser is a reference to the pflag function that traverses the flags.
	// The specific function (either `Visit` or `VisitAll`) is determined by the `useUnsetFlags` argument.
	flagTraverser := c.Command.Command.Flags().Visit
	if useUnsetFlags {
		flagTraverser = c.Command.Command.Flags().VisitAll
	}

	flagTraverser(func(flag *pflag.Flag) {
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
		case cloudapiv6.ArgExposeSerial:
			input.SetExposeSerial(boolval)
			break
		case cloudapiv6.ArgRequireLegacyBios:
			if flag.Changed == true {
				input.SetRequireLegacyBios(boolval)
			}
			break
		case cloudapiv6.ArgApplicationType:
			input.SetApplicationType(val)
			break
		default:
			// --image-id, verbose, filters, depth, etc
			break
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Property %s set: %s", flag.Name, flag.Value))
	})
	return input
}

func RunImageUpdate(c *core.CommandConfig) error {
	input := getDesiredImageAfterPatch(c, false)
	img, resp, err := c.CloudApiV6Services.Images().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)),
		input,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Image, img.Image,
		tabheaders.GetHeaders(allImageCols, defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func PreRunImageId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgImageId)
}

func RunImageList(c *core.CommandConfig) error {

	images, resp, err := c.CloudApiV6Services.Images().List()
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	if viper.IsSet(core.GetFlagName(c.NS, constants.FlagType)) {
		images = sortImagesByType(images, strings.ToUpper(viper.GetString(core.GetFlagName(c.NS, constants.FlagType))))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)) {
		images = sortImagesByAlias(images, viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias)))
	}

	if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)) {
		images = sortImagesByTime(images, viper.GetInt(core.GetFlagName(c.NS, cloudapiv6.ArgLatest)))
	}

	if itemsOk, ok := images.GetItemsOk(); !ok || itemsOk == nil {
		return nil
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("items", jsonpaths.Image, images.Images,
		tabheaders.GetHeaders(allImageCols, defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

func RunImageGet(c *core.CommandConfig) error {
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Image with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))))

	img, resp, err := c.CloudApiV6Services.Images().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
	}
	if err != nil {
		return err
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Image, img.Image,
		tabheaders.GetHeaders(allImageCols, defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

	return nil
}

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

const listImagesExample = `ionosctl compute image list

ionosctl compute image list --location us/las --type HDD --licence-type LINUX`

const getImageExample = `ionosctl compute image get --image-id IMAGE_ID`
