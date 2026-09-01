package image

import (
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/table"
	"github.com/ionos-cloud/sdk-go-bundle/shared/fileconfiguration"
	"github.com/spf13/cobra"
)

var allImageCols = []table.Column{
	{Name: "ImageId", JSONPath: "id", Default: true},
	{Name: "Name", JSONPath: "properties.name", Default: true},
	{Name: "ImageAliases", JSONPath: "properties.imageAliases", Default: true},
	{Name: "Location", JSONPath: "properties.location", Default: true},
	{Name: "LicenceType", JSONPath: "properties.licenceType", Default: true},
	{Name: "ImageType", JSONPath: "properties.imageType", Default: true},
	{Name: "CloudInit", JSONPath: "properties.cloudInit", Default: true},
	{Name: "CreatedDate", JSONPath: "metadata.createdDate", Default: true},
	{Name: "Size", JSONPath: "properties.size"},
	{Name: "Description", JSONPath: "properties.description"},
	{Name: "Public", JSONPath: "properties.public"},
	{Name: "CreatedBy", JSONPath: "metadata.createdBy"},
	{Name: "CreatedByUserId", JSONPath: "metadata.createdByUserId"},
	{Name: "ExposeSerial", JSONPath: "properties.exposeSerial"},
	{Name: "RequireLegacyBios", JSONPath: "properties.requireLegacyBios"},
	{Name: "ApplicationType", JSONPath: "properties.applicationType"},
	{Name: "RequiredFeatures", JSONPath: "properties.requiredFeatures"},
}

func ImageCmd() *core.Command {
	imageCmd := &core.Command{
		Command: &cobra.Command{
			Use:     "image",
			Aliases: []string{"img"},
			Short:   "Manage bootable OS images and ISOs used to create volumes and servers",
			Long: `An Image is a bootable disk template that a Volume is cloned from when you create a server. Every image is scoped to a single physical location (e.g. de/fra) and has a licence-type (LINUX, WINDOWS2022, ...) that tells the platform how to bill and configure the guest OS.

There are two kinds of image:
  * PUBLIC  - IONOS-maintained OS distributions and ISOs, shared across all contracts and available in every location. You cannot modify or delete them. They usually carry one or more image-aliases (e.g. "ubuntu:latest", "debian:12") so you can reference the newest build by a stable name instead of a UUID.
  * PRIVATE - images you own, created either by uploading a disk file over FTP (see 'image upload') or by snapshotting an existing volume. Only these can be updated ('image update') and, in a limited way, deleted.

An image also advertises hardware capabilities that servers booted from it inherit: hot-plug/hot-unplug support for CPU, RAM, NIC and disks (add/remove the resource on a running VM with no reboot), cloud-init support for first-boot provisioning, and Confidential Computing (SEV-SNP) support surfaced via the RequiredFeatures column.

The sub-commands below let you list and inspect PUBLIC and PRIVATE images, upload your own image files, and set an uploaded image's properties.`,
			TraverseChildren: true,
		},
	}
	imageCmd.AddColsFlag(allImageCols)

	imageCmd.AddCommand(ImageListCmd())
	imageCmd.AddCommand(ImageGetCmd())
	imageCmd.AddCommand(ImageUpdateCmd())
	imageCmd.AddCommand(ImageDeleteCmd())
	imageCmd.AddCommand(Upload())

	return core.WithConfigOverride(imageCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
