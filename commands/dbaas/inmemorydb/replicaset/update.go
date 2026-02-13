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
			id, err := c.Command.Command.Flags().GetString(constants.FlagReplicasetID)
			if err != nil {
				return err
			}

			rs, _, err := client.Must().InMemoryDBClient.
				ReplicaSetApi.
				ReplicasetsFindById(context.Background(), id).
				Execute()
			if err != nil {
				return fmt.Errorf("failed getting replicaset with id %q: %w", id, err)
			}

			input := rs.Properties

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				input.DisplayName = name
			}

			if c.Command.Command.Flags().Changed(constants.FlagVersion) {
				version, err := c.Command.Command.Flags().GetString(constants.FlagVersion)
				if err != nil {
					return err
				}
				input.Version = version
			}

			if c.Command.Command.Flags().Changed(constants.FlagReplicas) {
				replicas, err := c.Command.Command.Flags().GetInt(constants.FlagReplicas)
				if err != nil {
					return err
				}
				input.Replicas = int32(replicas)
			}
			if c.Command.Command.Flags().Changed(constants.FlagCores) {
				cores, err := c.Command.Command.Flags().GetInt(constants.FlagCores)
				if err != nil {
					return err
				}
				input.Resources.Cores = int32(cores)
			}
			if c.Command.Command.Flags().Changed(constants.FlagRam) {
				ram, err := c.Command.Command.Flags().GetString(constants.FlagRam)
				if err != nil {
					return err
				}
				if ram != "" {
					sizeInt64 := convbytes.StrToUnit(ram, convbytes.GB)
					if sizeInt64 < math.MinInt32 || sizeInt64 > math.MaxInt32 {
						return fmt.Errorf("RAM size %d exceeds the range of int32", sizeInt64)
					}
					input.Resources.Ram = int32(sizeInt64)
				}
			}

			if c.Command.Command.Flags().Changed(constants.FlagPersistenceMode) {
				persistenceMode, err := c.Command.Command.Flags().GetString(constants.FlagPersistenceMode)
				if err != nil {
					return err
				}
				input.PersistenceMode = inmemorydb.PersistenceMode(persistenceMode)
			}
			if c.Command.Command.Flags().Changed(constants.FlagEvictionPolicy) {
				evictionPolicy, err := c.Command.Command.Flags().GetString(constants.FlagEvictionPolicy)
				if err != nil {
					return err
				}
				input.EvictionPolicy = inmemorydb.EvictionPolicy(evictionPolicy)
			}

			if c.Command.Command.Flags().Changed(constants.FlagBackupLocation) {
				backupLocation, err := c.Command.Command.Flags().GetString(constants.FlagBackupLocation)
				if err != nil {
					return err
				}
				input.Backup = &inmemorydb.BackupProperties{}
				input.Backup.Location = pointer.From(backupLocation)
			}

			if c.Command.Command.Flags().Changed(constants.FlagDatacenterId) {
				datacenterId, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
				if err != nil {
					return err
				}
				input.Connections[0].DatacenterId = datacenterId
			}
			if c.Command.Command.Flags().Changed(constants.FlagLanId) {
				lanId, err := c.Command.Command.Flags().GetString(constants.FlagLanId)
				if err != nil {
					return err
				}
				input.Connections[0].LanId = lanId
			}
			if c.Command.Command.Flags().Changed(constants.FlagCidr) {
				cidr, err := c.Command.Command.Flags().GetString(constants.FlagCidr)
				if err != nil {
					return err
				}
				input.Connections[0].Cidr = cidr
			}

			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceTime) {
				maintenanceTime, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceTime)
				if err != nil {
					return err
				}
				input.MaintenanceWindow.Time = maintenanceTime
			}
			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceDay) {
				maintenanceDay, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceDay)
				if err != nil {
					return err
				}
				input.MaintenanceWindow.DayOfTheWeek = inmemorydb.DayOfTheWeek(maintenanceDay)
			}

			if c.Command.Command.Flags().Changed(constants.ArgUser) {
				user, err := c.Command.Command.Flags().GetString(constants.ArgUser)
				if err != nil {
					return err
				}
				input.Credentials.Username = user
			}

			if c.Command.Command.Flags().Changed(constants.ArgPassword) {
				input.Credentials.Password = &inmemorydb.UserPassword{}
				password, err := c.Command.Command.Flags().GetString(constants.ArgPassword)
				if err != nil {
					return err
				}
				hashFlag, err := c.Command.Command.Flags().GetBool(constants.ArgHashPassword)
				if err != nil {
					return err
				}

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
