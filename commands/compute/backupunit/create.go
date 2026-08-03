package backupunit

import (
	"context"

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
		LongDesc: `Use this command to create a BackupUnit under your contract. A BackupUnit is the named storage container + login that the IONOS Managed Backup agent uses to store server backups.

You must supply --name, --email and --password.

Notes:

* --name becomes the backup login: it is concatenated with your contract number as CONTRACT_NUMBER-NAME, so it must be GLOBALLY UNIQUE across all IONOS contracts. It CANNOT be changed after creation (to rename, delete and recreate).
* --password is the login secret used to register the backup agent. It is WRITE-ONLY: the Cloud API never returns it, so record it now. It can be changed later with ` + "`backupunit update`" + `.
* --email receives service reports from the backup system and does NOT need to match your Cloud API username. It can be changed later.
* After creation, log in to the backup console at https://backup.ionos.com (or via DCD, https://dcd.ionos.com/latest/). Use ` + "`backupunit get-sso-url`" + ` for a one-click SSO link.

Required values to run a command:

* Name
* Email
* Password`,
		Example: `# Create a BackupUnit
ionosctl compute backupunit create --name mybackups --email ops@example.com --password 'S3cretPass!'

# Then open the backup console via SSO (grab the id from the create output)
ionosctl compute backupunit get-sso-url --backupunit-id BACKUPUNIT_ID`,
		PreCmdRun:  PreRunBackupUnitNameEmailPwd,
		CmdRun:     RunBackupUnitCreate,
		InitClient: true,
	})
	cmd.AddStringFlag(cloudapiv6.ArgName, cloudapiv6.ArgNameShort, "", "Alphanumeric name for the BackupUnit. Combined with your contract number it forms the backup login (CONTRACT_NUMBER-NAME), so it must be globally unique and CANNOT be changed after creation", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgEmail, cloudapiv6.ArgEmailShort, "", "E-mail address that will receive backup service reports. Does not need to match your Cloud API username", core.RequiredFlagOption())
	cmd.AddStringFlag(cloudapiv6.ArgPassword, cloudapiv6.ArgPasswordShort, "", "Login secret used to register the backup agent. Write-only: it is never returned by the API, so record it now", core.RequiredFlagOption())

	return cmd
}
