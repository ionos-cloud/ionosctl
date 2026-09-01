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

	"github.com/ionos-cloud/ionosctl/v6/commands/compute/completer"
	completer2 "github.com/ionos-cloud/ionosctl/v6/commands/dbaas/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/inmemorydb/utils"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
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
		ShortDesc: "Create an In-Memory DB Replica Set",
		LongDesc: `Create a new In-Memory DB Replica Set (a Redis-compatible database).

The mode is set by --replicas: 1 replica is a standalone instance; more than 1 is a leader-follower replication with one active and n-1 passive standby replicas (passive replicas fail over but do not serve reads). Up to 5 replicas are allowed.

Each replica gets its own --cores (1-16) and --ram (4-32 GB); storage is derived automatically from the RAM and persistence mode and cannot be set. The replica set attaches to exactly one network connection: --datacenter-id + --lan-id + --cidr. An initial user (--user / --password) is created for you.

There are two ways to create a replica set:
  1. Empty: pass the sizing, connection and credential flags (as in the basic example below).
  2. From a snapshot: additionally pass --snapshot-id to restore an existing point-in-time snapshot into the new replica set. The snapshot must belong to the same location.

PersistenceMode (--persistence-mode, controls how data survives restarts):
  None:    In-memory only, nothing is persisted. Best for pure caches.
  AOF:     Append Only File - every write is logged and replayed on restart, reconstructing the dataset.
  RDB:     Periodic point-in-time dumps of the in-memory state.
  RDB_AOF: Both RDB and AOF are enabled.

EvictionPolicy (--eviction-policy, what happens when the memory limit is hit):
  noeviction:      Never evict; write operations return an error once memory is full.
  allkeys-lru:     Evict the least recently used keys first.
  allkeys-lfu:     Evict the least frequently used keys first.
  allkeys-random:  Evict random keys.
  volatile-lru:    As allkeys-lru, but only among keys with a TTL (expire) set.
  volatile-lfu:    As allkeys-lfu, but only among keys with a TTL (expire) set.
  volatile-random: Evict random keys, but only among keys with a TTL (expire) set.
  volatile-ttl:    Evict the key with the nearest TTL first, among keys with a TTL set.`,
		Example: `# Create a standalone (single-instance) replica set
ionosctl dbaas in-memory-db replicaset create --location de/fra --name mycache --replicas 1 --cores 1 --ram 4GB --datacenter-id DATACENTER_ID --lan-id 1 --cidr 192.168.1.100/24 --user dbadmin --password MyStrongPass1

# Advanced: a 3-node leader-follower set with AOF persistence, a custom eviction policy and maintenance window
ionosctl dbaas in-memory-db replicaset create --location de/fra --name prod-cache --replicas 3 --cores 4 --ram 16GB --persistence-mode RDB_AOF --eviction-policy volatile-lru --datacenter-id DATACENTER_ID --lan-id 1 --cidr 192.168.1.100/24 --user dbadmin --password MyStrongPass1 --maintenance-day Saturday --maintenance-time 03:00:00`,
		PreCmdRun: func(c *core.PreCommandConfig) error {
			return c.CheckRequiredFlagsAndLocation(
				constants.FlagName, constants.FlagReplicas,
				constants.FlagCores, constants.FlagRam,
				constants.ArgUser, constants.ArgPassword,
				constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr)
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

			return c.Printer(allCols).Print(replica)
		},
		InitClient: true,
	})

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	addPropertiesFlags(cmd)

	return cmd
}

func addPropertiesFlags(cmd *core.Command) {
	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The human-readable display name of the Replica Set", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "7.2", "The In-Memory DB (Redis-compatible) engine version. Supported values: 7.0, 7.2", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagReplicas, "", 1,
		"Total number of replicas (1-5): one active plus n-1 passive. Set 1 for a standalone instance; >1 enables leader-follower replication. "+
			"Passive replicas are hot standbys for failover of the active instance only - they are NOT read replicas and do not serve reads", core.RequiredFlagOption())
	cmd.AddIntFlag(constants.FlagCores, "", 1, "The number of CPU cores per instance (1-16)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCores, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "2", "4", "8", "12", "16"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagRam, "", "4GB", "The amount of memory per instance, 4-32 GB (e.g. --ram 8 or --ram 8GB). Storage size is derived automatically from RAM and persistence mode and is not configurable", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"4GB", "8GB", "16GB", "32GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagPersistenceMode, "", "RDB",
		[]string{"None", "AOF", "RDB", "RDB_AOF"}, "How data is persisted across restarts: None (cache only), AOF (write log), RDB (periodic dumps), or RDB_AOF (both). See the long description for details")
	cmd.AddSetFlag(constants.FlagEvictionPolicy, "", "allkeys-lru",
		[]string{"noeviction", "allkeys-lru", "allkeys-lfu", "allkeys-random", "volatile-lru", "volatile-lfu", "volatile-random", "volatile-ttl"}, "What to evict when the memory limit is reached. 'volatile-*' policies only touch keys with a TTL; 'noeviction' errors on writes when full. See the long description for details")

	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The ID of the Virtual Datacenter the replica set connects into. Must be in the same location as the replica set",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.DataCentersIds()
		}, "api.ionos.com", []string{}),
	)
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric ID of a private LAN (within the datacenter) to attach the replica set to",
		core.RequiredFlagOption(),
		core.WithCompletion(func() []string {
			return completer.LansIds(viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId)))
		}, "api.ionos.com", []string{}),
	)
	cmd.AddStringFlag(constants.FlagCidr, "", "", "The private IP and subnet assigned to the replica set on the LAN, in CIDR notation (e.g. 192.168.1.100/24)."+
		" These ranges are reserved and cannot be used: 10.210.0.0/16, 10.212.0.0/14", core.RequiredFlagOption(),
		core.WithCompletionComplex(completer2.GetCidrCompletionFunc(cmd), constants.InMemoryDBApiRegionalURL, constants.InMemoryDBLocations),
	)

	// Maintenance
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	hour := 10 + r.Intn(7) // Random hour 10-16
	workingDaysOfWeek := []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday"}

	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", fmt.Sprintf("%02d:00:00", hour),
		"Start time (UTC, HH:MM:SS) of the weekly 4-hour maintenance window during which upgrades and patches may occur, e.g. 16:30:59. "+
			"Defaults to a random time between 10:00 and 16:00")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", workingDaysOfWeek[rand.Intn(len(workingDaysOfWeek))],
		"Day of the week for the weekly 4-hour maintenance window (Monday-Sunday). "+
			"Defaults to a random weekday (Mon-Fri)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return append(workingDaysOfWeek, "Saturday", "Sunday"), cobra.ShellCompDirectiveNoFileComp
	})

	// credentials
	cmd.AddStringFlag(constants.ArgUser, "", "", "Username for the initial database user. 1-16 chars, letters/digits/underscore only. Some names are reserved (e.g. 'admin', 'standby')", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.ArgPassword, "", "", "Password for the initial user. Plaintext must be 10-63 chars; by default it is hashed (SHA-256) client-side before sending. A value that is already a 64-char SHA-256 hash is sent as-is", core.RequiredFlagOption())
	cmd.AddBoolFlag(constants.ArgHashPassword, "", true, "Hash a plaintext --password (SHA-256) before sending. Set --hash-password=false to send the plaintext password to the API as-is")

	cmd.AddStringFlag(constants.FlagBackupLocation, "", "", "The S3 (Object Storage) location where automatic backups/snapshots are stored, e.g. 'de'. Defaults to a location near the cluster")
	cmd.AddStringFlag(constants.FlagSnapshotId, "", "",
		"Create the replica set restored from this existing snapshot instead of empty. The snapshot must be in the same location; all sizing/connection/credential flags still apply",
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
