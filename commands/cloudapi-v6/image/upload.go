package image

import (
	"bufio"
	"context"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
)

/*
	Perform upload to given FTP server.
	- ftp://ftp-fkb.ionos.com/hdd-images
	- ftp://ftp-fkb.ionos.com/iso-images
	https://docs.ionos.com/cloud/compute-engine/block-storage/block-storage-faq#how-do-i-upload-my-own-images-with-ftp
*/

type locationInfo struct {
	FTP string // fragment to embed into ftp-%s.ionos.com ('fra')
	API string // canonical CloudAPI location ('de/fra')
}

var knownLocations = map[string]locationInfo{
	// short forms -> ftp fragment, API location
	"fra":   {FTP: "fra", API: "de/fra"},
	"fra/2": {FTP: "fra-2", API: "de/fra/2"},
	"fkb":   {FTP: "fkb", API: "de/fkb"},
	"txl":   {FTP: "txl", API: "de/txl"},
	"lhr":   {FTP: "lhr", API: "gb/lhr"},
	"bhx":   {FTP: "bhx", API: "gb/bhx"},
	"las":   {FTP: "las", API: "us/las"},
	"ewr":   {FTP: "ewr", API: "us/ewr"},
	"vit":   {FTP: "vit", API: "es/vit"},
	"par":   {FTP: "par", API: "fr/par"},
	"mci":   {FTP: "mci", API: "us/mci"},

	// API-style and suffixed forms
	"de/fra":   {FTP: "fra", API: "de/fra"},
	"de/fra/2": {FTP: "fra-2", API: "de/fra/2"},
	"es/vit":   {FTP: "vit", API: "es/vit"},
	"gb/lhr":   {FTP: "lhr", API: "gb/lhr"},
	"gb/bhx":   {FTP: "bhx", API: "gb/bhx"},
	"fr/par":   {FTP: "par", API: "fr/par"},
	"us/las":   {FTP: "las", API: "us/las"},
	"us/ewr":   {FTP: "ewr", API: "us/ewr"},
	"us/mci":   {FTP: "mci", API: "us/mci"},
	"de/txl":   {FTP: "txl", API: "de/txl"},
	"de/fkb":   {FTP: "fkb", API: "de/fkb"},
}

// fallback country guesses for short regions not explicitly present in knownLocations
var baseRegionToCountry = map[string]string{
	"fra": "de",
	"fkb": "de",
	"txl": "de",
	"lhr": "gb",
	"bhx": "gb",
	"las": "us",
	"ewr": "us",
	"vit": "es",
	"par": "fr",
	"mci": "us",
}

// lookupFTP returns the FTP fragment that should be embedded into ftp-%s.ionos.com.
//
//	"fra"       -> "fra"
//	"de/fra"    -> "fra"
//	"de/fra/2"  -> "fra-2"
func lookupFTP(loc string) string {
	if loc == "" {
		return ""
	}
	if info, ok := knownLocations[loc]; ok {
		return info.FTP
	}

	// heuristic fallback:
	// - if input contains '/', assume region is parts[1] and optional suffix parts[2]
	// - else assume input is short region token and use it
	parts := strings.Split(loc, "/")
	if len(parts) == 1 {
		return parts[0]
	}
	region := parts[1]
	if len(parts) >= 3 && parts[2] != "" {
		return fmt.Sprintf("%s-%s", region, parts[2])
	}
	return region
}

// lookupAPI returns canonical API location for polling Images API.
//
//	"fra"       -> "de/fra"
//	"de/fra"    -> "de/fra"
//	"de/fra/2"  -> "de/fra/2"
func lookupAPI(loc string) string {
	if loc == "" {
		return ""
	}
	if info, ok := knownLocations[loc]; ok {
		return info.API
	}
	// If looks like api form, return it verbatim
	if strings.Contains(loc, "/") {
		return loc
	}
	// fallback guess using baseRegionToCountry
	if c, ok := baseRegionToCountry[loc]; ok {
		return fmt.Sprintf("%s/%s", c, loc)
	}
	// last resort: return region token
	return loc
}

// regionPart returns the canonical region token used for validation
//
//	"fra"         -> "fra"
//	"de/fra"      -> "fra"
//	"de/fra/2"    -> "fra"
func regionPart(loc string) string {
	if loc == "" {
		return ""
	}
	parts := strings.Split(loc, "/")
	if len(parts) == 1 {
		return parts[0]
	}
	// for api-form like de/fra or de/fra/2, region is parts[1]
	return parts[1]
}

func Upload() *core.Command {
	upload := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "image",
		Resource:  "image",
		Verb:      "upload",
		Aliases:   []string{"ftp-upload", "ftp", "upl"},
		ShortDesc: "Upload an image to FTP server using FTP over TLS (FTPS)",
		LongDesc: `This command uploads one or more disk images to an FTP server using FTP over TLS (FTPS), then optionally updates the uploaded images via the Images API to set properties you passed as flags.
This command requires that you are logged in using IONOS_USERNAME and IONOS_PASSWORD environment variables.

High level steps:
  1. Upload file(s) concurrently to the target FTP server(s).
  2. If you do not use --skip-update, poll the Images API for the uploaded image(s) to appear.
  3. When the API shows the uploaded image(s), perform PATCH /images/<UUID> to apply the requested image properties.
  4. Print the resulting image objects to stdout in the chosen table or JSON format.

AUTH AND SAFETY
  - The FTP server relies on API credentials via environment variables IONOS_USERNAME and IONOS_PASSWORD. You can debug your current setup with "ionosctl whoami --provenance".
  - Use --skip-update to skip the API PATCH step if you only want to perform an FTP upload and not modify images through the API.
  - Use --skip-verify to skip verifying the FTP server certificate. Only use that for trusted servers. Skipping certificate verification can expose you to man-in-the-middle attacks.
  - If using a custom FTP server it is advised to use a self-signed certificate instead of --skip-verify. Provide its PEM file via --crt-path. The file should contain the server certificate in base64 PEM format.

FTP URLs
  - Default IONOS FTP servers are of the form ftp-<region>.ionos.com (for example ftp-fra.ionos.com).
  - If uploading to default IONOS FTP servers, --ftp-url is optional. The command will construct the URL automatically from the locations you provide via --location (i.e. 'de/fra' or 'fra').
  - The command chooses the remote path automatically:
      * Files ending in .iso or .img are uploaded to the iso-images/ directory.
      * All other supported image extensions are uploaded to the hdd-images/ directory.
  - If you supply a custom --ftp-url that contains a placeholder, for example ftp://myftp.example/locations/%s, you must also supply one or more --location values. The command will replace %s with the location-specific fragment for each location. Example: --ftp-url ftp://myftp.example/locations/%s --location fra,fkb
  - If you supply a custom --ftp-url without a placeholder, you may provide multiple --ftp-url values to try multiple servers.

POLLING AND TIMEOUTS
  - After upload, unless you use --skip-update, the command repeatedly queries GET /images with filters for the uploaded file names and locations.
  - Polling runs until either all expected images appear, or the command context deadline expires.
  - The context deadline is controlled with --timeout (seconds). The FTP connection and the subsequent API operations share the same context. If a timeout occurs after FTP finished but before the PATCH completed, the PATCH will be cancelled.

NOTES
  - Uploading multiple images with the same name to the same location is forbidden.
  - The command does not delete or overwrite existing images on the FTP server. If an image with the same name already exists on the server, the upload will fail.
  - The command does not check if the uploaded image is valid or bootable. It only checks the file extension.
  - You can use 'ionosctl image list --filter public=false' to see your uploaded images.
  - You must contact support to delete images you uploaded via FTP. Deleting them via API will only set their size to 0B.

EXAMPLES
  - Simple upload to IONOS servers:
    ionosctl img upload -i image.iso -l de/fra,de/fkb,es/vit --skip-update
    Uploads image.iso to ftp://ftp-fkb.ionos.com/iso-images, ftp://ftp-fra.ionos.com/iso-images and ftp://ftp-vit.ionos.com/iso-images, then exits without calling the Images API.

  - Upload and let the CLI set properties via API:
    ionosctl img upload -i image.iso -l de/fra
    Uploads to ftp://ftp-fra.ionos.com/iso-images, polls GET /images until the image appears, then PATCHes that image with the properties you supplied via flags.

  - Use a custom FTP server:
    ionosctl img upload -i image.iso --ftp-url "ftp://myftp.example" --crt-path certificates/my-server-crt.pem --skip-update`,
		PreCmdRun: core.PreRunWithDeprecatedFlags(PreRunImageUpload,
			functional.Tuple[string]{First: FlagRenameImages, Second: cloudapiv6.ArgImageAlias}),
		CmdRun:     RunImageUpload,
		InitClient: true,
	})

	validLocations := []string{"de/fra", "de/fra/2", "es/vit", "gb/lhr", "gb/bhx", "fr/par", "us/las", "us/ewr", "us/mci", "de/txl", "de/fkb"}
	validLocationsStr := strings.Join(validLocations, ", ")

	upload.AddStringSliceFlag(cloudapiv6.ArgLocation, cloudapiv6.ArgLocationShort, nil,
		fmt.Sprintf("Location to upload to. Can be one of %s if not using --%s",
			validLocationsStr, FlagFtpUrl), core.RequiredFlagOption())

	upload.AddStringSliceFlag(FlagRenameImages, "", nil, "Rename the uploaded images before trying to upload. These names should not contain any extension. By default, this is the base of the image path")
	upload.AddStringSliceFlag(FlagImage, "i", nil, "Slice of paths to images, can be absolute path or relative to current working directory", core.RequiredFlagOption())
	upload.AddStringFlag(FlagFtpUrl, "", "ftp-%s.ionos.com", "URL of FTP server, with %s flag if location is embedded into url")
	upload.AddBoolFlag(FlagSkipVerify, "", false, "Skip verification of server certificate, useful if using a custom ftp-url. WARNING: You can be the target of a man-in-the-middle attack!")
	upload.AddBoolFlag(FlagSkipUpdate, "", false, "Skip setting image properties after it has been uploaded. Normal behavior is to send a PATCH to the API, after the image has been uploaded, with the contents of the image properties flags and emulate a \"create\" command.")
	upload.AddStringFlag(FlagCertificatePath, "", "", "(Not needed for IONOS FTP Servers) Path to file containing server certificate. If your FTP server is self-signed, you need to add the server certificate to the list of certificate authorities trusted by the client.")
	upload.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, 1200, "(seconds) Context Deadline. FTP connection will time out after this many seconds")

	upload.AddStringSliceFlag(cloudapiv6.ArgImageAlias, cloudapiv6.ArgImageAliasShort, nil, "")
	upload.Command.Flags().MarkHidden(cloudapiv6.ArgImageAlias)

	addPropertiesFlags(upload)

	upload.Command.Flags().SortFlags = false // Hot Plugs generate a lot of flags to scroll through, put them at the end
	upload.Command.SilenceUsage = true       // Don't print help if setting only 1 out of 2 required flags - too many flags. Help must be invoked manually via --help

	return upload
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

	fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
		"Uploading %+v to %+v", images, locations))

	originalURL := url
	for _, loc := range locations {
		for imgIdx, img := range images {
			// build FTP URL: replace %s with the ftp fragment mapped for the location
			if strings.Contains(originalURL, "%s") {
				ftpSub := lookupFTP(loc)
				url = fmt.Sprintf(originalURL, ftpSub) // Add the location modifier for FTP URL
			}
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
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
				defer file.Close()
				return resources.FtpUpload(
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
			})
		}
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	// If --skip-update is set, we are done
	if viper.GetBool(core.GetFlagName(c.NS, FlagSkipUpdate)) {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
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

	// Map user-provided locations to API locations for polling (e.g. "fra" -> "de/fra")
	apiLocations := make([]string, len(locations))
	for i, l := range locations {
		apiLocations[i] = lookupAPI(l)
	}

	diffImgs, err := getDiffUploadedImages(c, names, apiLocations) // Get UUIDs of uploaded images
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP successful: %w", err)
	}

	properties := getDesiredImageAfterPatch(c, true)
	imgs, err := updateImagesAfterUpload(c, diffImgs, properties)
	if err != nil {
		return fmt.Errorf("failed updating image with given properties, but uploading to FTP successful: %w", err)
	}

	cols := viper.GetStringSlice(core.GetFlagName(c.Resource, constants.ArgCols))

	out, err := jsontabwriter.GenerateOutput("", jsonpaths.Image, imgs,
		tabheaders.GetHeaders(allImageCols, defaultImageCols, cols))
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)

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

			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
				"Looking for images with names '%s' in locations '%s'...", strings.Join(names, ","), strings.Join(locations, ",")))

			imgs, _, err := req.Execute()
			if err != nil {
				return nil, fmt.Errorf("failed listing images")
			}
			j, err := json.Marshal(*imgs.Items)
			if err == nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Got images by listing: %s", string(j)))
			}

			diffImgs = append(diffImgs, *imgs.Items...)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Total images: %+v", len(diffImgs)))

			if len(diffImgs) == len(names)*len(locations) {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Success! All images found via API: %+v", len(diffImgs)))
				return diffImgs, nil
			}

			if len(diffImgs) > len(names)*len(locations) {
				return nil, fmt.Errorf("more images found (%d) than expected (%d). "+
					"Something went terribly wrong. Please open an issue at github.com/ionos-cloud/ionosctl/issues/new", len(diffImgs), len(names)*len(locations))
			}

			// New attempt...
			time.Sleep(10 * time.Second)
		}
	}
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

	validRegions := []string{"de/fra", "de/fra/2", "es/vit", "gb/lhr", "gb/bhx", "fr/par", "us/las", "us/ewr", "us/mci", "de/txl", "de/fkb"}
	locs := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	invalidLocs := functional.Filter(
		locs,
		func(loc string) bool {
			return !slices.Contains(validRegions, loc)
		},
	)
	if len(invalidLocs) > 0 {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"WARN: '%s' is an invalid location. Valid IONOS regions are: '%s'", strings.Join(invalidLocs, ", "), strings.Join(validRegions, ", ")))
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
