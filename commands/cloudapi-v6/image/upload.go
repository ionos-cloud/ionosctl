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
		LongDesc: fmt.Sprintf(`WARNING:
This command can only be used if 2-Factor Authentication is disabled on your account and you're logged in using IONOS_USERNAME and IONOS_PASSWORD environment variables (see "Authenticating with Ionos Cloud" at https://docs.ionos.com/cli-ionosctl).

OVERVIEW:
  Use this command to securely upload one or more HDD or ISO images to the specified FTP server using FTP over TLS (FTPS). This command supports a variety of options to provide flexibility during the upload process:
  - The command supports renaming the uploaded images with the '--%s' flag. If uploading multiple images, you must provide a new name for each image.
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
`, FlagRenameImages, constants.ArgTimeout, FlagSkipUpdate, cloudapiv6.ArgLocation, FlagFtpUrl, FlagSkipVerify, FlagCertificatePath),
		Example: `- 'ionosctl img upload -i kolibri.iso -l fkb,fra,vit --skip-update': Simply upload the image 'kolibri.iso' from the current directory to IONOS FTP servers 'ftp://ftp-fkb.ionos.com/iso-images', 'ftp://ftp-fra.ionos.com/iso-images', 'ftp://ftp-vit.ionos.com/iso-images'.
- 'ionosctl img upload -i kolibri.iso -l fra': Upload the image 'kolibri.iso' from the current directory to IONOS FTP server 'ftp://ftp-fra.ionos.com/iso-images'. Once the upload has finished, start querying 'GET /images' with a filter for 'kolibri', to get the UUID of the image as seen by the Images API. When UUID is found, perform a 'PATCH /images/<UUID>' to set the default flag values.
- 'ionosctl img upload -i kolibri.iso --skip-update --skip-verify --ftp-url ftp://12.34.56.78': Use your own custom server. Use skip verify to skip checking server's identity
- 'ionosctl img upload -i kolibri.iso -l fra --ftp-url ftp://myComplexFTPServer/locations/%s --crt-path certificates/my-servers-cert.crt --location Paris,Berlin,LA,ZZZ --skip-update': Upload the image to multiple FTP servers, with location embedding into URL.`,
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
	upload.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, 300, "(seconds) Context Deadline. FTP connection will time out after this many seconds")

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

			imgs, _, err := req.Execute()
			if err != nil {
				return nil, fmt.Errorf("failed listing images")
			}
			j, err := json.Marshal(*imgs.Items)
			if err == nil {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Got images by listing: %s", string(j)))
			}

			diffImgs = append(diffImgs, *imgs.Items...)
			fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Total images: %+v", diffImgs))

			if len(diffImgs) == len(names)*len(locations) {
				fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput("Success! All images found via API: %+v", diffImgs))
				return diffImgs, nil
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

	// Accept short tokens and API-style tokens. Validate by the canonical region token.
	validRegions := []string{"fra", "fkb", "txl", "lhr", "bhx", "las", "ewr", "vit", "par", "mci"}
	locs := viper.GetStringSlice(core.GetFlagName(c.NS, cloudapiv6.ArgLocation))
	invalidLocs := functional.Filter(
		locs,
		func(loc string) bool {
			r := regionPart(loc)
			return !slices.Contains(validRegions, r)
		},
	)
	if len(invalidLocs) > 0 {
		fmt.Fprintf(c.Command.Command.ErrOrStderr(), "%s", jsontabwriter.GenerateVerboseOutput(
			"WARN: %s is an invalid location. Valid IONOS regions are: %s", strings.Join(invalidLocs, ","), validRegions))
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
