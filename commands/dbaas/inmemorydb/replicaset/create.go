package replicaset

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"time"

	"github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	completer2 "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/jsonpaths"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/uuidgen"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/inmemorydb/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Create() *core.Command {
	cmd := core.NewCommand(context.Background(), nil, core.CommandBuilder{
		Namespace: "dbaas inmemorydb",
		Resource:  "replicaset",
		Verb:      "create",
		Aliases:   []string{"post", "c"},
		ShortDesc: "Create a replica set",
		LongDesc: `Create a replica set. In-Memory DB replica set with support for a single instance or a In-Memory DB replication in leader follower mode. The mode is determined by the number of replicas. One replica is standalone, everything else an In-Memory DB replication as leader follower mode with one active and n-1 passive replicas.

PersistenceMode:
None: Data is inMemory only and will not be persisted. Useful for cache only applications.
AOF (Append Only File): AOF persistence logs every write operation received by the server. These operations can then be replayed again at server startup, reconstructing the original dataset. Commands are logged using the same format as the In-Memory DB protocol itself.
RDB: RDB persistence performs snapshots of the current in memory state.
RDB_AOF: Both RDB and AOF persistence are enabled.

EvictionPolicy:
noeviction: No eviction policy is used. In-Memory DB will never remove any data. If the memory limit is reached, an error will be returned on write operations.
allkeys-lru: The least recently used keys will be removed first.
allkeys-lfu: The least frequently used keys will be removed first.
allkeys-random: Random keys will be removed.
volatile-lru: The least recently used keys will be removed first, but only among keys with the expire field set to true.
volatile-lfu: The least frequently used keys will be removed first, but only among keys with the expire field set to true.
volatile-random: Random keys will be removed, but only among keys with the expire field set to true.
volatile-ttl: The key with the nearest time to live will be removed first, but only among keys with the expire field set to true.`,
		Example: "ionosctl dbaas inmemorydb replicaset create " + core.FlagsUsage(constants.FlagLocation, constants.FlagName,
			constants.FlagReplicas, constants.FlagCores, constants.FlagRam, constants.ArgUser, constants.ArgPassword,
			constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			if err := core.CheckRequiredFlags(c.Command, c.NS,
				constants.FlagName, constants.FlagReplicas,
				constants.FlagCores, constants.FlagRam,
				constants.ArgUser, constants.ArgPassword,
				constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr); err != nil {
				return err
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			input := inmemorydb.ReplicaSet{}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				input.DisplayName = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagVersion); true {
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

			if fn := core.GetFlagName(c.NS, constants.FlagPersistenceMode); true {
				input.PersistenceMode = inmemorydb.PersistenceMode(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagEvictionPolicy); true {
				input.EvictionPolicy = inmemorydb.EvictionPolicy(viper.GetString(fn))
			}

			if fn := core.GetFlagName(c.NS, constants.FlagBackupLocation); viper.IsSet(fn) {
				input.Backup = &inmemorydb.BackupProperties{}
				input.Backup.Location = pointer.From(viper.GetString(fn))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagSnapshotId); viper.IsSet(fn) {
				input.InitialSnapshotId = pointer.From(viper.GetString(fn))
			}

			input.Connections = make([]inmemorydb.Connection, 1)
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				input.Connections[0].DatacenterId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				input.Connections[0].LanId = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				input.Connections[0].Cidr = viper.GetString(fn)
			}

			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); true {
				if input.MaintenanceWindow == nil {
					input.MaintenanceWindow = &inmemorydb.MaintenanceWindow{}
				}
				input.MaintenanceWindow.Time = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); true {
				if input.MaintenanceWindow == nil {
					input.MaintenanceWindow = &inmemorydb.MaintenanceWindow{}
				}
				input.MaintenanceWindow.DayOfTheWeek = inmemorydb.DayOfTheWeek(viper.GetString(fn))
			}

			input.Credentials = inmemorydb.User{Password: &inmemorydb.UserPassword{}}
			if fn := core.GetFlagName(c.NS, constants.ArgUser); viper.IsSet(fn) {
				input.Credentials.Username = viper.GetString(fn)
			}
			if fn := core.GetFlagName(c.NS, constants.ArgPassword); viper.IsSet(fn) {
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

			id := uuidgen.Must()
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

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	addPropertiesFlags(cmd)

	return cmd
}

func addPropertiesFlags(cmd *core.Command) {
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of the Replica Set", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "7.2", "The In-Memory DB version of your Replica Set", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagReplicas, "", 1,
		"The total number of replicas in the Replica Set (one active and n-1 passive)."+
			" In case of a standalone instance, the value is 1. In all other cases, the value is >1. "+
			"The replicas will not be available as read replicas, they are only standby for a failure of the active instance", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagCores, "", 1, "The number of CPU cores per instance", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "12", "16", "24", "31"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance in gigabytes (GB)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB", "64GB", "128GB", "256GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagPersistenceMode, "", "RDB",
		[]string{"None", "AOF", "RDB", "RDB_AOF"}, "Specifies how and if data is persisted (refer to the long description for more details)")
	cmd.AddSetFlag(constants.FlagEvictionPolicy, "", "allkeys-lru",
		[]string{"noeviction", "allkeys-lru", "allkeys-lfu", "allkeys-random", "volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl"}, "The eviction policy for the replica set (refer to the long description for more details)")

	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to connect your instance to",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.DataCentersIds()
		}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric Private LAN ID to connect your instance to",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId)))
		}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)
	cmd.AddStringFlag(constants.FlagCidr, "", "", "The IP and subnet for your instance."+
		" Note the following unavailable IP ranges: 10.210.0.0/16 10.212.0.0/14", core.RequiredFlagOption(),
		core.WithCompletionComplex(completer2.GetCidrCompletionFunc(cmd), constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)

	// Maintenance
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7) // Random hour 10-16
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", fmt.Sprintf("%02d:00:00", hour),
		"Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: 16:30:59. "+
			"Defaults to a random day during Mon-Fri, during the hours 10:00-16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))],
		"Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. "+
			"Defaults to a random day during Mon-Fri, during the hours 10:00-16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})

	// credentials
	cmd.AddStringFlag(constants.ArgUser, "", "", "The initial username", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.ArgPassword, "", "", "Password (plaintext or SHA-256). If plaintext, it’s hashed when --hash-password is true; otherwise sent as-is", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgHashPassword, "", true, "Hash plaintext passwords before sending. Use '--hash-password=false' to send plaintext passwords as-is")

	cmd.AddStringFlag(constants.FlagBackupLocation, "", "", "The S3 location where the backups will be stored")
	cmd.AddStringFlag(constants.FlagSnapshotId, "", "",
		"If set, will create the replicaset from the specified snapshot",
		core.WithCompletion(
			func() []string {
				// for each snapshot
				return utils.SnapshotProperty(func(snapshot inmemorydb.SnapshotRead) string {
					// return its ID
					return snapshot.Id + "\t" + snapshot.Metadata.SnapshotTime.Format("2006-01-02 15:04:05")
				})
			}, constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations,
		),
	)
}
