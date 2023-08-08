package cluster

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/spf13/cobra"
)

func ClusterCreateCmd() *core.Command {

	playgroundRequired, _ := getRequiredFlagsByEditionAndType("playground", "")
	businessRequired, _ := getRequiredFlagsByEditionAndType("business", "")
	enterpriseReplicasetRequired, _ := getRequiredFlagsByEditionAndType("enterprise", "replicaset")
	enterpriseShardedRequired, _ := getRequiredFlagsByEditionAndType("enterprise", "sharded-cluster")

	examples := []string{
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s playground %s",
			constants.FlagEdition, core.FlagsUsage(playgroundRequired[1:]...)),
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s business %s",
			constants.FlagEdition, core.FlagsUsage(businessRequired[1:]...)),
		// Example where type is inferred by using --shards or --instances
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s enterprise (%s | %s) %s",
			constants.FlagEdition, core.FlagUsage(constants.FlagInstances), core.FlagUsage(constants.FlagShards), core.FlagsUsage(enterpriseReplicasetRequired[1:len(enterpriseReplicasetRequired)-1]...)),
		// Example where using --type creates a requirement for --instances
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s enterprise --%s replicaset %s",
			constants.FlagEdition, constants.FlagType, core.FlagsUsage(enterpriseReplicasetRequired[1:]...)),
		// Example where using --type creates a requirement for --shards
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s enterprise --%s sharded-cluster %s",
			constants.FlagEdition, constants.FlagType, core.FlagsUsage(enterpriseShardedRequired[1:]...)),
	}

	cmd := core.NewCommand(context.TODO(), nil, core.CommandBuilder{
		Namespace: "dbaas-mongo",
		Resource:  "cluster",
		Verb:      "create", // used in AVAILABLE COMMANDS in help
		Aliases:   []string{"c"},
		ShortDesc: "Create DBaaS Mongo Replicaset or Sharded Clusters for your chosen edition",
		Example:   strings.Join(examples, "\n\n"),
		PreCmdRun: func(c *core.PreCommandConfig) error {
			// old: "cidr", "datacenter-id", "instances", "lan-id", "maintenance-day", "maintenance-time", "name", "template-id"

			/* Supermongo:
			 * For edition playground, only replica-set, =1 instance and playground template (33457e53-1f8b-4ed2-8a12-2d42355aa759, 1 core, 50 GB Storage, 2 GB RAM).
			 * For edition business, only replica-set, >1 instance and any template.
			 * For edition enterprise, type replica-set/sharded-cluster and
			 *  - CPU Cores: 1-8
			 *  - RAM Size (GB): <16 GB
			 *  - Storage Size:  >100GB for optimal perf. max 1048.576 GB.
			 *  - Shards: 2-32 shards.
			 *  - Instances: >3
			**/

			if fn := core.GetFlagName(c.NS, constants.FlagEdition); !viper.IsSet(fn) {
				return fmt.Errorf("--%s (%s) is a required option",
					constants.FlagEdition, strings.Join(enumEditions, " | "))
			} else {
				edition := viper.GetString(fn)
				// type is a reserved keyword
				cType := viper.GetString(core.GetFlagName(c.NS, constants.FlagType))
				flags, err := getRequiredFlagsByEditionAndType(edition, cType)
				if err != nil {
					return fmt.Errorf("failed getting required flags for edition %s: %w", edition, err)
				}

				err = core.CheckRequiredFlags(c.Command, c.NS, flags...)
				if err != nil {
					return fmt.Errorf("not all %s edition flags are set: %w", edition, err)
				}
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			// TODO
			return nil
			/*
				// Extra CLI helpers
				if *input.Properties.Location == "" {
					// If location isn't set to Datacenter's Location, Mongo API throws an error. Location property is also marked as required
					// To improve user experience we mark it as optional, and set it to the datacenter's location implicitly via connections.datacenterID.
					dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(c.Context, *createConn.DatacenterId).Execute()
					if err != nil {
						return err
					}
					input.Properties.Location = dc.Properties.Location
				}

			*/
		},
		InitClient: true,
	})

	cmd.AddSetFlag(constants.FlagEdition, "e", "", enumEditions, "Cluster Edition", core.RequiredFlagOption())
	cmd.AddSetFlag(constants.FlagType, "", "replicaset", enumTypes, "Cluster Type. Required for enterprise clusters")

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster")
	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "The physical location where the cluster will be created. (defaults to the location of the connected datacenter)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagTemplateId, "", "", "The unique ID of the template, which specifies the number of cores, storage size, and memory")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTemplateId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoTemplateIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances in the cluster (one primary and n-1 secondaries). (required for non-playground clusters)")
	cmd.AddStringFlag(constants.FlagVersion, "", "6.0", "The MongoDB version of your cluster")

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request [seconds]")

	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LansIds(os.Stderr, viper.GetString(core.GetFlagName(cmd.NS, constants.FlagDatacenterId))),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.FlagCidr, "", nil, "The list of IPs and subnet for your cluster. All IPs must be in a /24 network. Note the following unavailable IP range: 10.233.114.0/24", core.RequiredFlagOption())

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

// returns a slice of flags to be marked as required, depending on wanted edition and type
func getRequiredFlagsByEditionAndType(edition, cType string) ([]string, error) {
	alwaysRequired := []string{
		constants.FlagEdition, constants.FlagName,
		constants.FlagMaintenanceDay, constants.FlagMaintenanceTime,
		constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
	}
	switch edition {
	case "playground":
		// Type inferred as replicaset. Template inferred as type playground. Instances inferred as 1
		return alwaysRequired, nil
	case "business":
		// Type inferred as replicaset.
		return append(alwaysRequired,
			constants.FlagTemplateId, constants.FlagInstances,
		), nil
	case "enterprise":
		enterpriseBaseFlags := append(alwaysRequired,
			constants.FlagCores, constants.FlagStorageType,
			constants.FlagStorageSize, constants.FlagRam,
		)
		switch cType {
		case "replicaset":
			return append(enterpriseBaseFlags, constants.FlagInstances), nil
		case "sharded-cluster":
			return append(enterpriseBaseFlags, constants.FlagShards), nil
		case "":
			return nil, fmt.Errorf("--%s is required for enterprise clusters (valid: <%s>)",
				constants.FlagType, strings.Join(enumTypes, " | "))
		default:
			return nil, fmt.Errorf("unknown type %s (valid: <%s>)",
				cType, strings.Join(enumTypes, " | "))
		}
	default:
		return nil, fmt.Errorf("unknown edition %s (valid: <%s>)",
			edition, strings.Join(enumEditions, " | "))
	}
}
