package commands

import (
	"bufio"
	"context"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/internal/waitfor"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"golang.org/x/exp/slices"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/query"
	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/waiter"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

var (
	defaultImageCols = []string{"ImageId", "Name", "ImageAliases", "Location", "LicenceType", "ImageType", "CloudInit", "CreatedDate"}
	allImageCols     = []string{"ImageId", "Name", "ImageAliases", "Location", "Size", "LicenceType", "ImageType", "Description", "Public", "CloudInit", "CreatedDate", "CreatedBy", "CreatedByUserId"}
)

const (
	FlagRenameImages    = "rename"
	FlagImage           = "image"
	FlagSkipUpdate      = "skip-update"
	FlagSkipVerify      = "skip-verify"
	FlagFtpUrl          = "ftp-url"
	FlagCertificatePath = "crt-path"
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
	globalFlags.StringSliceP(constants.ArgCols, "", defaultImageCols, tabheaders.ColsMessage(allImageCols))
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

	deprecatedMessage := "incompatible with --max-results. Use --filters --order-by --max-results options instead!"

	list.AddStringFlag(constants.FlagType, "", "", "The type of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = list.Command.RegisterFlagCompletionFunc(constants.FlagType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"CDROM", "HDD"}, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgLicenceType, "", "", "The licence type of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLicenceType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return constants.EnumLicenceType, cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, "", "The location of the Image", core.DeprecatedFlagOption(deprecatedMessage))
	_ = list.Command.RegisterFlagCompletionFunc(cloudapiv6.ArgLocation, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	list.AddStringFlag(cloudapiv6.ArgImageAlias, "", "", "Image Alias or part of Image Alias to sort Images by", core.DeprecatedFlagOption(deprecatedMessage))
	list.AddIntFlag(cloudapiv6.ArgLatest, "", 0, "Show the latest N Images, based on creation date, starting from now in descending order. If it is not set, all Images will be printed", core.DeprecatedFlagOption("Use --filters --order-by --max-results options instead!"))
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
		return completer.ImageIds(), cobra.ShellCompDirectiveNoFileComp
	})
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
		return completer.ImageIds(func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			return request.Filter("public", "false")
		}), cobra.ShellCompDirectiveNoFileComp
	})

	update.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end

	update.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	update.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")
	update.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	addPropertiesFlags := func(command *core.Command) {
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
		return completer.ImageIds(func(request ionoscloud.ApiImagesGetRequest) ionoscloud.ApiImagesGetRequest {
			return request.Filter("public", "false")
		}), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(constants.ArgAll, constants.ArgAllShort, false, "Delete all non-public images")

	deleteCmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request for Image update to be executed")
	deleteCmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for Image update [seconds]")
	deleteCmd.AddInt32Flag(cloudapiv6.ArgDepth, cloudapiv6.ArgDepthShort, cloudapiv6.DefaultGetDepth, cloudapiv6.ArgDepthDescription)

	/*
		upload Command
		Perform upload to given FTP server. All locations from `location list` + karlsruhe (fkb)
		- ftp://ftp-fkb.ionos.com/hdd-images
		- ftp://ftp-fkb.ionos.com/iso-images
		https://docs.ionos.com/cloud/compute-engine/block-storage/block-storage-faq#how-do-i-upload-my-own-images-with-ftp
	*/
	upload := core.NewCommand(ctx, imageCmd, core.CommandBuilder{
		Namespace: "image",
		Resource:  "image",
		Verb:      "upload",
		Aliases:   []string{"ftp-upload", "ftp", "upl"},
		ShortDesc: "Upload an image to FTP server using FTP over TLS (FTPS)",
		LongDesc: fmt.Sprintf(`WARNING: 
This command can only be used if 2-Factor Authentication is disabled on your account and you're logged in using IONOS_USERNAME and IONOS_PASSWORD environment variables (see "Authenticating with Ionos Cloud" at https://docs.ionos.com/cli-ionosctl).

OVERVIEW:
  Use this command to securely upload one or more HDD or ISO images to the specified FTP server using FTP over TLS (FTPS). This command supports a variety of options to provide flexibility during the upload process:
  - The command supports renaming the uploaded images with the '--%s' flag. If uploading multiple images, you must provide an alias for each image.
  - Specify the context deadline for the FTP connection using the '--%s' flag. The operation as a whole will terminate after the specified number of seconds, i.e. if the FTP upload had finished but your PATCH operation did not, only the PATCH operation will be intrerrupted.
POST-UPLOAD OPERATIONS:
  By default, this command will query 'GET /images' endpoint for your uploaded images, then try to use 'PATCH /images/<UUID>' to update the uploaded images with the given property fields.
  - It is necessary to use valid API credentials for this.
  - To skip this API behaviour, you can use '--%s'.
CUSTOM URLs:
  This command supports usage of other FTP servers too, not just the IONOS ones.
  - The '--%s' flag is only required if your '--%s' contains a placeholder variable (i.e. %%s).
  In this case, for every location in that slice, an attempt of FTP upload would be made at the URL computed by embedding it into the placeholder variable
  - Use the '--%s' flag to skip the verification of the server certificate. This can be useful when using a custom ftp-url,
  but be warned that this could expose you to a man-in-the-middle attack.
  - If you're using a self-signed FTP server, you can provide the path to the server certificate file in base64 PEM format using the '--%s' flag.
`, cloudapiv6.ArgImageAlias, constants.ArgTimeout, FlagSkipUpdate, cloudapiv6.ArgLocation, FlagFtpUrl, FlagSkipVerify, FlagCertificatePath),
		Example: `- 'ionosctl img upload -i kolibri.iso -l fkb,fra,vit --skip-update': Simply upload the image 'kolibri.iso' from the current directory to IONOS FTP servers 'ftp://ftp-fkb.ionos.com/iso-images', 'ftp://ftp-fra.ionos.com/iso-images', 'ftp://ftp-vit.ionos.com/iso-images'.
- 'ionosctl img upload -i kolibri.iso -l fra': Upload the image 'kolibri.iso' from the current directory to IONOS FTP server 'ftp://ftp-fra.ionos.com/iso-images'. Once the upload has finished, start querying 'GET /images' with a filter for 'kolibri', to get the UUID of the image as seen by the Images API. When UUID is found, perform a 'PATCH /images/<UUID>' to set the default flag values.
- 'ionosctl img upload -i kolibri.iso --skip-update --skip-verify --ftp-url ftp://12.34.56.78': Use your own custom server. Use skip verify to skip checking server's identity
- 'ionosctl img upload -i kolibri.iso -l fra --ftp-url ftp://myComplexFTPServer/locations/%s --crt-path certificates/my-servers-cert.crt --location Paris,Berlin,LA,ZZZ --skip-update': Upload the image to multiple FTP servers, with location embedding into URL.`,
		PreCmdRun: core.PreRunWithDeprecatedFlags(PreRunImageUpload,
			functional.Tuple[string]{First: FlagRenameImages, Second: cloudapiv6.ArgImageAlias}),
		CmdRun:     RunImageUpload,
		InitClient: true,
	})

	upload.AddStringSliceFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, nil, fmt.Sprintf("Location to upload to. Must be an array containing only fra, fkb, txl, lhr, las, ewr, vit if not using --%s", FlagFtpUrl), core.RequiredFlagOption())
	upload.AddStringSliceFlag(FlagRenameImages, "", nil, "Rename the uploaded images before trying to upload. These names should not contain any extension. By default, this is the base of the image path")
	upload.AddStringSliceFlag(FlagImage, "i", nil, "Slice of paths to images, can be absolute path or relative to current working directory", core.RequiredFlagOption())
	upload.AddStringFlag(FlagFtpUrl, "", "ftp-%s.ionos.com", "URL of FTP server, with %s flag if location is embedded into url")
	upload.AddBoolFlag(FlagSkipVerify, "", false, "Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!")
	upload.AddBoolFlag(FlagSkipUpdate, "", false, "Skip setting image properties after it has been uploaded. Normal behavior is to send a PATCH to the API, after the image has been uploaded, with the contents of the image properties flags and emulate a \"create\" command.")
	upload.AddStringFlag(FlagCertificatePath, "", "", "(Not needed for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.")
	upload.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, 300, "(seconds) Context Deadline. FTP connection will time out after this many seconds")

	upload.AddStringSliceFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, nil, "")
	upload.Command.Flags().MarkHidden(cloudapiv6.ArgImageAlias)

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

		return nil
	}

	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete image", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	imgId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deletion on image with ID: %v...", imgId))

	resp, err := c.CloudApiV6Services.Images().Delete(imgId, queryParams)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
	}
	if err != nil {
		return err
	}

	if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), jsontabwriter.GenerateLogOutput("Image deleted successfully"))

	return nil
}
func PreRunImageUpload(c *core.PreCommandConfig) error {
	err := c.Command.Command.MarkFlagRequired(FlagImage)
	if err != nil {
		return err
	}

	validExts := []string{".iso", ".img", ".vmdk", ".vhd", ".vhdx", ".cow", ".qcow", ".qcow2", ".raw", ".vpc", ".vdi"}
	images := viper.GetStringSlice(core.GetFlagName(c.NS, FlagImage))
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
	if strings.Contains(viper.GetString(core.GetFlagName(c.NS, FlagFtpUrl)), "%s") {
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
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"WARN: %s is an invalid location. Valid IONOS locations are: %s", strings.Join(invalidLocs, ","), locs))
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
func RunImageUpload(c *core.CommandConfig) error {
	certPool, err := getCertificate(viper.GetString(core.GetFlagName(c.NS, FlagCertificatePath)))
	if err != nil {
		return err
	}

	url := viper.GetString(core.GetFlagName(c.NS, FlagFtpUrl))
	images := viper.GetStringSlice(core.GetFlagName(c.NS, FlagImage))
	aliases := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgImageAlias))
	locations := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	skipVerify := viper.GetBool(core.GetFlagName(c.NS, FlagSkipVerify))

	ctx, cancel := context.WithTimeout(c.Context, time.Duration(viper.GetInt(core.GetFlagName(c.NS, constants.ArgTimeout)))*time.Second)
	defer cancel()
	c.Context = ctx

	// just a simple patch to force entry into the `for` loop below if no locations are provided
	if !strings.Contains(url, "%s") &&
		(locations == nil || len(locations) == 0) {
		sentinel := []string{""}
		locations = sentinel
	}

	var eg errgroup.Group

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
		"Uploading %+v to %+v", images, locations))

	originalURL := url
	for _, loc := range locations {
		for imgIdx, img := range images {
			if strings.Contains(originalURL, "%s") {
				url = fmt.Sprintf(originalURL, loc) // Add the location modifier, if the URL supports it
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
				"Uploading %s to %s", img, url))

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
			eg.Go(func() error {
				err := resources.FtpUpload(
					c.Context,
					resources.UploadProperties{
						FTPServerProperties: resources.FTPServerProperties{
							Url:               url,
							Port:              21,
							SkipVerify:        skipVerify,
							ServerCertificate: certPool,
							Username:          client.Must().CloudClient.GetConfig().Username,
							Password:          client.Must().CloudClient.GetConfig().Password,
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

	// If --skip-update is set, we are done
	if viper.GetBool(core.GetFlagName(c.NS, FlagSkipUpdate)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(
			"Successfully uploaded images"))
		return nil
	}

	// Below, we query that the images have been uploaded, and then PATCH them with the given properties
	names := images
	if len(aliases) != 0 {
		// Returns a slice containing `alias[i] + filepath.Ext(images[i])`
		names = functional.MapIdx(aliases, func(k int, v string) string {
			return v + filepath.Ext(images[k])
		})
	}
	diffImgs, err := getDiffUploadedImages(c, names, locations) // Get UUIDs of uploaded images
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP sucessful: %w", err)
	}

	properties := getDesiredImageAfterPatch(c, true)
	imgs, err := updateImagesAfterUpload(c, diffImgs, properties)
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP sucessful: %w", err)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Image, imgs,
		tabheaders.GetHeaders(allImageCols, defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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
			j, err := json.Marshal(*imgs.Items)
			if err == nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Got images by listing: %s", string(j)))
			}

			diffImgs = append(diffImgs, *imgs.Items...)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Total images: %+v", diffImgs))

			if len(diffImgs) == len(names)*len(locations) {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Success! All images found via API: %+v", diffImgs))
				return diffImgs, nil
			}

			// New attempt...
			time.Sleep(10 * time.Second)
		}
	}
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

	items, err := getNonPublicImages(*allItems, c.Command.Command.ErrOrStderr())
	if err != nil {
		return err
	}
	if len(items) < 1 {
		return errors.New("no non-public images found")
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput("Images to be deleted:"))
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

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(delIdAndName))
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete all the images", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Deleting all the images..."))

	var multiErr error
	for _, img := range items {
		if id, ok := img.GetIdOk(); ok && id != nil {
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Starting deleting image with id: %v...", *id))

			resp, err = c.CloudApiV6Services.Images().Delete(*id, resources.QueryParams{})
			if resp != nil && request.GetId(resp) != "" {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
			}
			if err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
				continue
			} else {
				_ = jsontabwriter.GenerateLogOutput(fmt.Sprintf(constants.MessageDeletingAll, c.Resource, *id))
			}

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateLogOutput(constants.MessageDeletingAll, c.Resource, *id))

			if err = waitfor.WaitForRequest(c, waiter.RequestInterrogator, request.GetId(resp)); err != nil {
				multiErr = errors.Join(multiErr, fmt.Errorf(constants.ErrDeleteAll, c.Resource, *id, err))
				continue
			}
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
			fmt.Fprintf(verboseOut, jsontabwriter.GenerateVerboseOutput("skipping %s: properties are nil\n", *i.GetId()))
			continue
		}

		isPublic, ok := properties.GetPublicOk()
		if !ok {
			fmt.Fprintf(verboseOut, jsontabwriter.GenerateVerboseOutput("skipping %s: field `public` is nil\n", *i.GetId()))
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
		default:
			// --image-id, verbose, filters, depth, etc
			break
		}

		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Property %s set: %s", flag.Name, flag.Value))
	})
	return input
}

func RunImageUpdate(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}
	queryParams := listQueryParams.QueryParams

	input := getDesiredImageAfterPatch(c, false)
	img, resp, err := c.CloudApiV6Services.Images().Update(
		viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)),
		input,
		queryParams,
	)
	if resp != nil && request.GetId(resp) != "" {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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
	// Add Query Parameters for GET Requests
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	images, resp, err := c.CloudApiV6Services.Images().List(listQueryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
}

func RunImageGet(c *core.CommandConfig) error {
	listQueryParams, err := query.GetListQueryParams(c)
	if err != nil {
		return err
	}

	queryParams := listQueryParams.QueryParams
	fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput("Image with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId))))

	img, resp, err := c.CloudApiV6Services.Images().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgImageId)), queryParams)
	if resp != nil {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), jsontabwriter.GenerateVerboseOutput(constants.MessageRequestTime, resp.RequestTime))
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

	fmt.Fprintf(c.Command.Command.OutOrStdout(), out)

	return nil
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
