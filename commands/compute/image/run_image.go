package image

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
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
	FlagFtpPort         = "ftp-port"
	FlagCertificatePath = "crt-path"
)

func addPropertiesFlags(command *core.Command) {
	command.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Human-friendly display name to give the uploaded image (does not have to be unique)")
	command.AddStringFlag(cloudapiv6.ArgDescription, cloudapiv6.ArgDescriptionShort, "", "Free-text description to set on the uploaded image")
	command.AddSetFlag(cloudapiv6.ArgLicenceType, "", "UNKNOWN", constants.EnumLicenceType, "OS licence type. Determines how IONOS bills and configures the guest (e.g. Windows editions are licensed). Use LINUX/RHEL for Linux, WINDOWS2016/2019/2022/2025 for the matching Windows Server, OTHER/UNKNOWN otherwise")
	command.AddSetFlag(constants.FlagCloudInit, "", "V1", []string{"V1", "NONE"}, "Whether servers built from this image accept cloud-init user-data for first-boot provisioning. V1 enables the cloud-init datasource; NONE disables it. Confidential Computing images are forced to NONE")
	command.AddBoolFlag(cloudapiv6.ArgCpuHotPlug, "", true, "Guest supports adding CPU cores while running (no reboot)")
	command.AddBoolFlag(cloudapiv6.ArgRamHotPlug, "", true, "Guest supports adding RAM while running (no reboot)")
	command.AddBoolFlag(cloudapiv6.ArgNicHotPlug, "", true, "Guest supports attaching a NIC while running (no reboot)")
	command.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotPlug, "", true, "Guest supports attaching a Virt-IO disk while running (no reboot)")
	command.AddBoolFlag(cloudapiv6.ArgDiscScsiHotPlug, "", true, "Guest supports attaching a SCSI disk while running (no reboot)")
	command.AddBoolFlag(cloudapiv6.ArgCpuHotUnplug, "", false, "Guest supports removing CPU cores while running. Only valid if CPU hot-plug is also enabled")
	command.AddBoolFlag(cloudapiv6.ArgRamHotUnplug, "", false, "Guest supports removing RAM while running. Only valid if RAM hot-plug is also enabled")
	command.AddBoolFlag(cloudapiv6.ArgNicHotUnplug, "", false, "Guest supports detaching a NIC while running. Only valid if NIC hot-plug is also enabled")
	command.AddBoolFlag(cloudapiv6.ArgDiscVirtioHotUnplug, "", false, "Guest supports detaching a Virt-IO disk while running. Not supported on Windows guests")
	command.AddBoolFlag(cloudapiv6.ArgDiscScsiHotUnplug, "", false, "Guest supports detaching a SCSI disk while running. Not supported on Windows guests")
	command.AddBoolFlag(cloudapiv6.ArgExposeSerial, "", false, "Expose the attached disk's serial id to the guest. Some OSes/software need it; note it can influence licensed-software (e.g. Windows) behavior")
	command.AddBoolFlag(cloudapiv6.ArgRequireLegacyBios, "", true, "Boot the image in legacy BIOS mode instead of UEFI. Set false for images that require/expect UEFI")
	command.AddSetFlag(cloudapiv6.ArgApplicationType, "", "UNKNOWN", constants.EnumApplicationType, "Application pre-installed on the image (e.g. an MSSQL edition). Only PUBLIC images may set a value other than UNKNOWN, so this is a no-op on uploaded (private) images")
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
	c.Verbose("Starting deletion on image with ID: %v...", imgId)

	resp, err := c.CloudApiV6Services.Images().Delete(imgId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Image deleted successfully")

	return nil
}

// DeleteAllNonPublicImages deletes non-public images, as deleting public images is forbidden by the API.
func DeleteAllNonPublicImages(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.Image]{
		Resource: "image",
		List: func() ([]ionoscloud.Image, error) {
			images, _, err := c.CloudApiV6Services.Images().List()
			if err != nil {
				return nil, err
			}

			allItems, ok := images.GetItemsOk()
			if !ok || allItems == nil {
				return nil, errors.New("could not retrieve images")
			}

			// deleting public images is forbidden by the API, so only consider non-public ones
			nonPublic, err := getNonPublicImages(*allItems, c.Command.Command.ErrOrStderr())
			if err != nil {
				return nil, err
			}
			if len(nonPublic) == 0 && len(*allItems) > 0 {
				return nil, errors.New("no deletable images found: all images are public (deleting public images is not allowed)")
			}
			return nonPublic, nil
		},
		Summary: func(img ionoscloud.Image) string {
			var id, name, location, description string
			if img.Id != nil {
				id = *img.Id
			}
			if p := img.Properties; p != nil {
				if p.Name != nil {
					name = *p.Name
				}
				if p.Location != nil {
					location = *p.Location
				}
				if p.Description != nil {
					description = *p.Description
				}
			}

			s := fmt.Sprintf("%s (id: %s, location: %s)", name, id, location)
			if description != "" {
				s = fmt.Sprintf("%s (id: %s, location: %s, desc: %s)", name, id, location, description)
			}
			return s
		},
		ID: func(img ionoscloud.Image) string {
			if img.Id != nil {
				return *img.Id
			}
			return ""
		},
		Delete: func(img ionoscloud.Image) error {
			resp, err := c.CloudApiV6Services.Images().Delete(*img.Id)
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}

// Util func - Given a slice of public & non-public images, return only those images that are non-public.
// If any image in the slice has null properties, or "Properties.Public" field is nil, the image is skipped (and a verbose message is shown)
func getNonPublicImages(imgs []ionoscloud.Image, verboseOut io.Writer) ([]ionoscloud.Image, error) {
	var nonPublicImgs []ionoscloud.Image
	for _, i := range imgs {
		properties, ok := i.GetPropertiesOk()
		if !ok {
			fmt.Fprintf(verboseOut, "[INFO] skipping %s: properties are nil\n", *i.GetId())
			continue
		}

		isPublic, ok := properties.GetPublicOk()
		if !ok {
			fmt.Fprintf(verboseOut, "[INFO] skipping %s: field `public` is nil\n", *i.GetId())
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
		case cloudapiv6.ArgDescription:
			input.SetDescription(val)
		case "cloud-init":
			input.SetCloudInit(val)
		case cloudapiv6.ArgLicenceType:
			input.SetLicenceType(val)
		case cloudapiv6.ArgCpuHotPlug:
			input.SetCpuHotPlug(boolval)
		case cloudapiv6.ArgRamHotPlug:
			input.SetRamHotPlug(boolval)
		case cloudapiv6.ArgNicHotPlug:
			input.SetNicHotPlug(boolval)
		case cloudapiv6.ArgDiscVirtioHotPlug:
			input.SetDiscVirtioHotPlug(boolval)
		case cloudapiv6.ArgDiscScsiHotPlug:
			input.SetDiscScsiHotPlug(boolval)
		case cloudapiv6.ArgCpuHotUnplug:
			input.SetCpuHotUnplug(boolval)
		case cloudapiv6.ArgRamHotUnplug:
			input.SetRamHotUnplug(boolval)
		case cloudapiv6.ArgNicHotUnplug:
			input.SetNicHotUnplug(boolval)
		case cloudapiv6.ArgDiscVirtioHotUnplug:
			input.SetDiscVirtioHotUnplug(boolval)
		case cloudapiv6.ArgDiscScsiHotUnplug:
			input.SetDiscScsiHotUnplug(boolval)
		case cloudapiv6.ArgExposeSerial:
			input.SetExposeSerial(boolval)
		case cloudapiv6.ArgRequireLegacyBios:
			if flag.Changed {
				input.SetRequireLegacyBios(boolval)
			}
		case cloudapiv6.ArgApplicationType:
			input.SetApplicationType(val)
		default:
			// --image-id, verbose, filters, depth, etc
		}

		c.Verbose("Property %s set: %s", flag.Name, flag.Value)
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
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allImageCols).Print(img.Image)
}

func PreRunImageId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgImageId)
}

func RunImageList(c *core.CommandConfig) error {
	images, resp, err := c.CloudApiV6Services.Images().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
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

	return c.Printer(allImageCols).Prefix("items").Print(images.Images)
}

func RunImageGet(c *core.CommandConfig) error {
	c.Verbose("Image with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))

	img, resp, err := c.CloudApiV6Services.Images().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allImageCols).Print(img.Image)
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

const listImagesExample = `# List every image visible to your contract (public + private)
ionosctl compute image list

# List only your own uploaded/snapshotted images
ionosctl compute image list --filters public=false

# Find public Ubuntu HDD images in Frankfurt via server-side filters
ionosctl compute image list --filters public=true,imageAliases=ubuntu:latest,location=de/fra`

const getImageExample = `ionosctl compute image get --image-id IMAGE_ID`
