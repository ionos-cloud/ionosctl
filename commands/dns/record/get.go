package record

import (
	"context"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
)

func ZonesRecordsFindByIdCmd() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dns",
		Resource:  "record",
		Verb:      "get",
		Aliases:   []string{"g"},
		ShortDesc: "Retrieve a record",
		Example:   "ionosctl dns record get --zoneId ZONE_ID --recordId RECORD_ID",
		PreCmdRun: func(c *core.PreCommandConfig) error {
			err := c.Command.Command.MarkFlagRequired("zoneId")
			if err != nil {
				return err
			}
			err = c.Command.Command.MarkFlagRequired("recordId")
			if err != nil {
				return err
			}

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			client.Must().DNSClient
		},
		InitClient: true,
	})

	cmd.AddStringFlag(zoneId, "", "", "The ID (UUID) of the DNS zone", core.RequiredFlagOption())
	cmd.AddStringFlag(recordId, "", "", "The ID (UUID) of the record", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagId, "", "", "")
	cmd.AddStringFlag(constants.FlagFqdn, "", "", "A fully qualified domain name. FQDN consists of two parts - the hostname and the domain name")
	cmd.AddStringFlag(constants.FlagLastModifiedDate, "", "", "The date of the last change formatted as yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")
	cmd.AddStringFlag(constants.FlagState, "", "", "The list of possible provisioning states in which DNS resource could be at the specific time")
	cmd.AddStringFlag(constants.FlagZoneId, "", "", "The ID (UUID) of the DNS zone of which record belongs to")
	cmd.AddStringFlag(constants.FlagCreatedDate, "", "", "The date of the record creation formatted as yyyy-MM-dd'T'HH:mm:ss.SSS'Z'")
	cmd.AddStringFlag(constants.FlagType, "", "", "Holds supported DNS resource record types. In the DNS context a record is a DNS resource record")
	cmd.AddStringFlag(constants.FlagContent, "", "", "")
	cmd.AddBoolFlag(constants.FlagEnabled, "", false, "When true - the record is visible for lookup")
	cmd.AddStringFlag(constants.FlagName, "", "", "")
	cmd.AddIntFlag(constants.FlagPriority, "", 0, "Priority value is between 0 and 65535. Priority is mandatory for MX, SRV and URI record types and ignored for all other types")
	cmd.AddIntFlag(constants.FlagTtl, "", 0, "Time to live for the record, recommended 3600")

	return cmd
}
