package cdrom

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
}

func ServerCdromCmd() *core.Command {
	serverCdromCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "cdrom",
			Aliases:          []string{"cd"},
			Short:            "Server CD-ROM Operations",
			Long:             "The sub-commands of `ionosctl compute server cdrom` allow you to attach, get, list, detach CD-ROMs from Servers.",
			TraverseChildren: true,
		},
	}
	serverCdromCmd.AddColsFlag(allImageCols)

	serverCdromCmd.AddCommand(ServerCdromAttachCmd())
	serverCdromCmd.AddCommand(ServerCdromListCmd())
	serverCdromCmd.AddCommand(ServerCdromGetCmd())
	serverCdromCmd.AddCommand(ServerCdromDetachCmd())

	return core.WithConfigOverride(serverCdromCmd, []string{fileconfiguration.Cloud, "compute"}, "")
}
