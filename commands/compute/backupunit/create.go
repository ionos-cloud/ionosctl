package backupunit

import (
	"context"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
)

func BackupUnitCreateCmd() *core.Command {
	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "backupunit",
		Resource:  "backupunit",
		Verb:      "create",
		Aliases:   []string{"c"},
		ShortDesc: "Create a BackupUnit",
		LongDesc: `Use this command to create a BackupUnit under a particular contract. You need to specify the name, email and password for the new BackupUnit.

Notes:

* The name assigned to the BackupUnit will be concatenated with the contract number to create the login name for the backup system. The name may NOT be changed after creation.
* The password set here is used along with the login name described above to register the backup agent with the backup system. When setting the password, please make a note of it, as the value cannot be retrieved using the Cloud API.
* The e-mail address supplied here does NOT have to be the same as your Cloud API username. This e-mail address will receive service reports from the backup system.
* To login to backup agent, please use [https://dcd.ionos.com/latest/](https://dcd.ionos.com/latest/) and access BackupUnit Console or use [https://backup.ionos.com](https://backup.ionos.com)

Required values to run a command:

* Name
* Email
* Password`,
		Example:    `ionosctl compute backupunit create --name NAME --email EMAIL --password PASSWORD`,
		PreCmdRun:  PreRunBackupUnitNameEmailPwd,
		CmdRun:     RunBackupUnitCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Alphanumeric name you want to assign to the BackupUnit", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "The e-mail address you want to assign to the BackupUnit", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Alphanumeric password you want to assign to the BackupUnit", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait,
		"Wait for the Request for BackupUnit creation to be executed")
	cmd.Command.Flags().MarkHidden(constants.ArgWaitForRequest) // Backupunit resources are not tracked by /requests endpoint yet - but keep the flag for backward compatibility
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request for BackupUnit creation [seconds]")

	return cmd
}
