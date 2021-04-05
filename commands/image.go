package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func image() *builder.Command {
	imageCmd := &builder.Command{
		Command: &cobra.Command{
			Use:              "image",
			Aliases:          []string{"images", "img"},
			Short:            "Image Operations",
			Long:             `The sub-command of ` + "`" + `ionosctl image` + "`" + ` allows you to see information about images available.`,
			TraverseChildren: true,
		},
	}
	globalFlags := imageCmd.Command.PersistentFlags()
	globalFlags.StringSlice(config.ArgCols, defaultImageCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(builder.GetGlobalFlagName(imageCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	list := builder.NewCommand(context.TODO(), imageCmd, noPreRun, RunImageList, "list", "List Images",
		"Use this command to get a list of available public Images. Use flags to retrieve a list of sorted images by location, licence type, type or size.",
		listImagesExample, true)
	list.AddFloat32Flag(config.ArgImageSize, "", 0, "The size of the Image")
	list.AddStringFlag(config.ArgImageType, "", "", "The type of the Image")
	list.AddStringFlag(config.ArgImageLicenceType, "", "", "The licence type of the Image")
	list.AddStringFlag(config.ArgImageLocation, "", "", "The location of the Image")
	_ = list.Command.RegisterFlagCompletionFunc(config.ArgImageLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	get := builder.NewCommand(context.TODO(), imageCmd, PreRunImageIdValidate, RunImageGet, "get", "Get a specified Image",
		"Use this command to get information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		getImageExample, true)
	get.AddStringFlag(config.ArgImageId, "", "", config.RequiredFlagImageId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), imageCmd, PreRunImageIdValidate, RunImageDelete, "delete", "Delete a private Image",
		"Use this command to delete the specified private image. This only applies to private images that you have uploaded.\n\nRequired values to run command:\n\n* Image Id",
		"", true)
	deleteCmd.AddStringFlag(config.ArgImageId, "", "", config.RequiredFlagImageId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	return imageCmd
}

func PreRunImageIdValidate(c *builder.PreCommandConfig) error {
	return builder.CheckRequiredFlags(c.ParentName, c.Name, config.ArgImageId)
}

func RunImageList(c *builder.CommandConfig) error {
	images, _, err := c.Images().List()
	if err != nil {
		return err
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageLocation)) {
		images = sortImagesByLocation(images, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageLocation)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageLicenceType)) {
		images = sortImagesByLicenceType(images, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageLicenceType)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageType)) {
		images = sortImagesByType(images, viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageType)))
	}
	if viper.IsSet(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageSize)) {
		images = sortImagesBySize(images, float32(viper.GetFloat64(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageSize))))
	}
	return c.Printer.Print(getImagePrint(nil, c, getImages(images)))
}

func RunImageGet(c *builder.CommandConfig) error {
	img, _, err := c.Images().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(nil, c, getImage(img)))
}

func RunImageDelete(c *builder.CommandConfig) error {
	err := utils.AskForConfirm(c.Stdin, c.Printer, "delete private image")
	if err != nil {
		return err
	}
	resp, err := c.Images().Delete(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageId)))
	if err != nil {
		return err
	}
	err = waitForAction(c, printer.GetRequestPath(resp))
	if err != nil {
		return err
	}
	return c.Printer.Print(getImagePrint(resp, c, nil))
}

// Output Printing

var defaultImageCols = []string{"ImageId", "Name", "Location", "Size", "LicenceType", "ImageType"}

type ImagePrint struct {
	ImageId     string  `json:"ImageId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	Description string  `json:"Description,omitempty"`
	Location    string  `json:"Location,omitempty"`
	Size        float32 `json:"Size,omitempty"`
	LicenceType string  `json:"LicenceType,omitempty"`
	ImageType   string  `json:"ImageType,omitempty"`
	Public      bool    `json:"Public,omitempty"`
}

func getImagePrint(resp *resources.Response, c *builder.CommandConfig, imgs []resources.Image) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.ParentName
			r.Verb = c.Name
			r.WaitFlag = viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait))
		}
		if imgs != nil {
			r.OutputJSON = imgs
			r.KeyValue = getImagesKVMaps(imgs)
			r.Columns = getImageCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr())
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
		"ImageId":     "ImageId",
		"Name":        "Name",
		"Description": "Description",
		"Location":    "Location",
		"Size":        "Size",
		"LicenceType": "LicenceType",
		"ImageType":   "ImageType",
		"Public":      "Public",
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
	if imgId, ok := img.GetIdOk(); ok && imgId != nil {
		imgPrint.ImageId = *imgId
	}
	if properties, ok := img.GetPropertiesOk(); ok && properties != nil {
		if name, ok := properties.GetNameOk(); ok && name != nil {
			imgPrint.Name = *name
		}
		if description, ok := properties.GetDescriptionOk(); ok && description != nil {
			imgPrint.Description = *description
		}
		if loc, ok := properties.GetLocationOk(); ok && loc != nil {
			imgPrint.Location = *loc
		}
		if size, ok := properties.GetSizeOk(); ok && size != nil {
			imgPrint.Size = *size
		}
		if licType, ok := properties.GetLicenceTypeOk(); ok && licType != nil {
			imgPrint.LicenceType = *licType
		}
		if imgType, ok := properties.GetImageTypeOk(); ok && imgType != nil {
			imgPrint.ImageType = *imgType
		}
		if public, ok := properties.GetPublicOk(); ok && public != nil {
			imgPrint.Public = *public
		}
	}
	return structs.Map(imgPrint)
}

func getImageIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	imageSvc := resources.NewImageService(clientSvc.Get(), context.TODO())
	images, _, err := imageSvc.List()
	clierror.CheckError(err, outErr)
	imgsIds := make([]string, 0)
	if items, ok := images.Images.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				imgsIds = append(imgsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return imgsIds
}

// Output Columns Sorting

func sortImagesByLocation(images resources.Images, location string) resources.Images {
	var imgLocationItems []ionoscloud.Image
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *images.Items {
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
	var imgLicenceTypeItems []ionoscloud.Image
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *images.Items {
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
	var imgTypeItems []ionoscloud.Image
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *images.Items {
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

func sortImagesBySize(images resources.Images, size float32) resources.Images {
	var imgTypeItems []ionoscloud.Image
	if items, ok := images.GetItemsOk(); ok && items != nil {
		for _, img := range *images.Items {
			properties := img.GetProperties()
			if imgSize, ok := properties.GetSizeOk(); ok && imgSize != nil {
				if *imgSize == size {
					imgTypeItems = append(imgTypeItems, img)
				}
			}
		}
	}
	images.Items = &imgTypeItems
	return images
}
