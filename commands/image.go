package commands

import (
	"context"
	"errors"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/builder"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
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
	viper.BindPFlag(builder.GetGlobalFlagName(imageCmd.Command.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))

	/*
		List Command
	*/
	listCmd := builder.NewCommand(context.TODO(), imageCmd, noPreRun, RunImageList, "list", "List images",
		"Use this command to get a list of available public images. Use flags to retrieve a list of sorted images by location, licence type, type or size.",
		"", true)
	listCmd.AddFloat32Flag(config.ArgImageSize, "", 0, "")
	listCmd.AddStringFlag(config.ArgImageType, "", "", "")
	listCmd.AddStringFlag(config.ArgImageLicenceType, "", "", "")
	listCmd.AddStringFlag(config.ArgImageLocation, "", "", "")
	listCmd.Command.RegisterFlagCompletionFunc(config.ArgImageLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getLocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Get Command
	*/
	getCmd := builder.NewCommand(context.TODO(), imageCmd, PreRunImageIdValidate, RunImageGet, "get", "Get a specified Image",
		"Use this command to get information about a specified Image.\n\nRequired values to run command:\n\n* Image Id",
		"", true)
	getCmd.AddStringFlag(config.ArgImageId, "", "", "The unique Image Id. [Required flag]")
	getCmd.Command.RegisterFlagCompletionFunc(config.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getImageIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Delete Command
	*/
	deleteCmd := builder.NewCommand(context.TODO(), imageCmd, PreRunImageIdValidate, RunImageDelete, "delete", "Delete a private Image",
		"Use this command to delete the specified private image. This only applies to private images that you have uploaded.\n\nRequired values to run command:\n\n* Image Id",
		"", true)
	deleteCmd.AddStringFlag(config.ArgImageId, "", "", "The unique Image Id. [Required flag]")
	deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgImageId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
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
	return c.Printer.Print(printer.Result{
		OutputJSON: images,
		KeyValue:   getImagesKVMaps(getImages(images)),
		Columns:    getImageCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
}

func RunImageGet(c *builder.CommandConfig) error {
	img, _, err := c.Images().Get(viper.GetString(builder.GetFlagName(c.ParentName, c.Name, config.ArgImageId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(printer.Result{
		OutputJSON: img,
		KeyValue:   getImagesKVMaps([]resources.Image{*img}),
		Columns:    getImageCols(builder.GetGlobalFlagName(c.ParentName, config.ArgCols), c.Printer.GetStderr()),
	})
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
	return c.Printer.Print(printer.Result{
		ApiResponse: resp,
		Resource:    "image",
		Verb:        "delete",
		WaitFlag:    viper.GetBool(builder.GetFlagName(c.ParentName, c.Name, config.ArgWait)),
	})
}

var defaultImageCols = []string{"ImageId", "Name", "Location", "Size", "LicenceType", "ImageType"}

type imagePrint struct {
	ImageId     string  `json:"ImageId,omitempty"`
	Name        string  `json:"Name,omitempty"`
	Description string  `json:"Description,omitempty"`
	Location    string  `json:"Location,omitempty"`
	Size        float32 `json:"Size,omitempty"`
	LicenceType string  `json:"LicenceType,omitempty"`
	ImageType   string  `json:"ImageType,omitempty"`
	Public      bool    `json:"Public,omitempty"`
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

func getImages(datacenters resources.Images) []resources.Image {
	dc := make([]resources.Image, 0)
	for _, d := range *datacenters.Items {
		dc = append(dc, resources.Image{d})
	}
	return dc
}

func getImagesKVMaps(imgs []resources.Image) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(imgs))
	for _, img := range imgs {
		properties := img.GetProperties()
		var imgPrint imagePrint
		if imgId, ok := img.GetIdOk(); ok && imgId != nil {
			imgPrint.ImageId = *imgId
		}
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
		o := structs.Map(imgPrint)
		out = append(out, o)
	}
	return out
}

func getImageIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	imageSvc := resources.NewImageService(clientSvc.Get(), context.TODO())
	images, _, err := imageSvc.List()
	clierror.CheckError(err, outErr)
	imgsIds := make([]string, 0)
	if images.Images.Items != nil {
		for _, d := range *images.Images.Items {
			imgsIds = append(imgsIds, *d.GetId())
		}
	} else {
		return nil
	}
	return imgsIds
}

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
