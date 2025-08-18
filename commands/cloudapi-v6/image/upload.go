package image

import (
	"context"
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

/*
	Perform upload to given FTP server.
	- ftp://ftp-fkb.ionos.com/hdd-images
	- ftp://ftp-fkb.ionos.com/iso-images
	https://docs.ionos.com/cloud/compute-engine/block-storage/block-storage-faq#how-do-i-upload-my-own-images-with-ftp
*/

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

	return upload
}
