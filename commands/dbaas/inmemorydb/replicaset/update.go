package replicaset

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"regexp"

	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/viper"
)

func Update() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "update",
		Aliases:   []string{"u"},
		ShortDesc: "Partially modify a replicaset's properties. NOTE: Passwords cannot be modified! This command uses a combination of GET and PUT to simulate a PATCH operation",
		Example:   fmt.Sprintf("ionosctl dbaas inmemorydb replicaset update %s", core.FlagsUsage(constants.FlagReplicasetID, constants.FlagName, constants.FlagReplicas, constants.FlagMaintenanceDay, constants.FlagMaintenanceTime)), PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagReplicasetID); err != nil {
				return err
			}

			// if viper.IsSet(core.GetFlagName(c.NS, constants.ArgPassword)) {
			// 	return fmt.Errorf("changing passwords is not yet supported")
			// }

			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			id := viper.GetString(core.GetFlagName(c.NS, constants.FlagReplicasetID))

			rs, _, err := client.Must().InMemoryDBClient.
				ReplicaSetApi.
				ReplicasetsFindById(context.Background(), id).
				Execute()
			if err != nil {
				return fmt.Errorf("failed getting replicaset with id %q: %w", id, err)
			}

			input := rs.Properties

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.DisplayName = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				input.Version = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagReplicas); viper.IsSet(fn) {
				input.Replicas = int32(viper.GetInt(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(fn) {
				input.Resources.Cores = int32(viper.GetInt(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagRam); viper.IsSet(fn) && viper.GetString(fn) != "" {
				sizeInt64 := convbytes.StrToUnit(viper.GetString(fn), convbytes.GB)
				if sizeInt64 < math.MinInt32 || sizeInt64 > math.MaxInt32 {
					return fmt.Errorf("RAM size %d exceeds the range of int32", sizeInt64)
				}
				input.Resources.Ram = int32(sizeInt64)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagPersistenceMode); viper.IsSet(fn) {
				input.PersistenceMode = inmemorydb.PersistenceMode(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEvictionPolicy); viper.IsSet(fn) {
				input.EvictionPolicy = inmemorydb.EvictionPolicy(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagBackupLocation); viper.IsSet(fn) {
				input.Backup = &inmemorydb.BackupProperties{}
				input.Backup.Location = pointer.From(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				input.Connections[0].DatacenterId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				input.Connections[0].LanId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				input.Connections[0].Cidr = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.IsSet(fn) {
				input.MaintenanceWindow.Time = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.IsSet(fn) {
				input.MaintenanceWindow.DayOfTheWeek = inmemorydb.DayOfTheWeek(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) {
				input.Credentials.Username = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
				input.Credentials.Password = &inmemorydb.UserPassword{}
				password := viper.GetString(fn)
				hashFlag := viper.GetBool(core.GetFlagName(c.NS, constants.ArgHashPassword))

				isSHA256 := func(s string) bool {
					// Check if it's a 64-character hex string
					matched, _ := regexp.MatchString("^[a-fA-F0-9]{64}$", s)
					return matched
				}

				switch {
				case isSHA256(password):
					input.Credentials.Password.
						HashedPassword = &inmemorydb.HashedPassword{Hash: password, Algorithm: "SHA-256"}
				case !hashFlag:
					input.Credentials.Password = &inmemorydb.UserPassword{PlainTextPassword: pointer.From(password)}
				case hashFlag:
					hash := sha256.Sum256([]byte(password))
					input.Credentials.Password.HashedPassword = &inmemorydb.HashedPassword{
						Hash:      hex.EncodeToString(hash[:]),
						Algorithm: "SHA-256",
					}
				}
			}

			replica, _, err := client.Must().InMemoryDBClient.ReplicaSetApi.
				ReplicasetsPut(context.Background(), id).
				ReplicaSetEnsure(inmemorydb.ReplicaSetEnsure{Id: id, Properties: input}).Execute()
			if err != nil {
				return err
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			out, err := jsontabwriter.GenerateOutput("", jsonpaths.DbaasInMemoryDBReplicaSet, replica,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddStringFlag(constants.FlagReplicasetID, constants.FlagIdShort, "",
		"The ID of the Replica Set you want to delete",
		core.WithCompletion(utils.ReplicasetIDs, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	addPropertiesFlags(cmd)

	return cmd
}
