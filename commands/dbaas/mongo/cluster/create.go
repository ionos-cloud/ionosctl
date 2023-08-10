package cluster

import (
	"context"
	"fmt"
	"os"
	"strings"

	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/functional"
	"github.com/ionos-cloud/ionosctl/v6/internal/pointer"
	"github.com/ionos-cloud/ionosctl/v6/pkg/constants"
	"github.com/ionos-cloud/ionosctl/v6/pkg/core"
	"github.com/ionos-cloud/ionosctl/v6/pkg/utils"
	ionoscloud "github.com/ionos-cloud/sdk-go-dbaas-mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

// SIDE EFFECT: sets FlagEdition if not set and can be inferred
func validateOrInferEditionByTemplate(c *core.PreCommandConfig) error {
	if fn := core.GetFlagName(c.NS, constants.FlagTemplate); viper.IsSet(fn) {
		tmplId, err := templates.Resolve(viper.GetString(fn))
		if err != nil {
			return fmt.Errorf("failed resolving %s to an ID of template: %w", viper.GetString(fn), err)
		}
		template, err := templates.Find(func(x ionoscloud.TemplateResponse) bool {
			return *x.Id == tmplId
		})
		if err != nil {
			return fmt.Errorf("failed finding template with ID %s: %w", tmplId, err)
		}

		if template.Properties == nil || template.Id == nil ||
			template.Properties.Edition == nil || template.Properties.Name == nil {
			return fmt.Errorf("found a template with some unset fields: %#v.\n Please use IONOS_LOG_LEVEL=trace and file a Github Issue", template)
		}

		if fnEd := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fnEd) {
			edition := viper.GetString(fnEd)
			// Check that template & edition aren't set to incompatible things

			if edition == "enterprise" && viper.IsSet(core.GetFlagName(c.NS, constants.FlagTemplate)) {
				return fmt.Errorf("for enterprise edition, setting --%s is forbidden. Use %s", constants.FlagTemplate,
					core.FlagsUsage(constants.FlagCores, constants.FlagRam, constants.FlagStorageType, constants.FlagStorageSize))
			}

			if *template.Properties.Edition != edition {
				return fmt.Errorf("the edition %s is not compatible with template %s (must be %s). Unset the flag --%s to use that edition instead",
					edition, *template.Properties.Name, *template.Properties.Edition, constants.FlagEdition)
			}
		} else {
			// Fallback edition to inferred one via template ID, if not explicitly set
			if slices.Contains(enumEditions, *template.Properties.Edition) {
				viper.Set(core.GetFlagName(c.NS, constants.FlagEdition), *template.Properties.Edition)
			}
		}
	}

	return nil
}

// validateEdition validates edition settings
func validateEdition(c *core.PreCommandConfig) error {
	fnEd := core.GetFlagName(c.NS, constants.FlagEdition)
	if !viper.IsSet(fnEd) {
		return fmt.Errorf("set --%s or --%s (%s) to get a list of required flags",
			constants.FlagTemplate, constants.FlagEdition, strings.Join(enumEditions, " | "))
	}

	edition := viper.GetString(fnEd)
	// Enterprise edition cannot have --template-id
	if edition == "enterprise" && viper.IsSet(core.GetFlagName(c.NS, constants.FlagTemplate)) {
		return fmt.Errorf("for enterprise edition, setting --%s is forbidden. Use %s", constants.FlagTemplate,
			core.FlagsUsage(constants.FlagCores, constants.FlagRam, constants.FlagStorageType, constants.FlagStorageSize))
	}

	flags, err := getRequiredFlagsByEditionAndType(edition, viper.GetString(core.GetFlagName(c.NS, constants.FlagType)))
	if err != nil {
		return fmt.Errorf("failed getting required flags for edition %s: %w", edition, err)
	}

	err = core.CheckRequiredFlags(c.Command, c.NS, flags...)
	if err != nil {
		return fmt.Errorf("not all %s edition flags are set: %w", edition, err)
	}

	return nil
}

// SIDE EFFECT: sets FlagLocation if not set and can be inferred
func inferLocationByDatacenter(c *core.PreCommandConfig) error {
	if fn := core.GetFlagName(c.NS, constants.FlagLocation); !viper.IsSet(fn) {
		dcId := viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
		dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
		if err != nil {
			return fmt.Errorf("failed inferring location via datacenter's ID: failed getting datacenter with ID %s: %w", dcId, err)
		}
		if dc.Properties == nil || dc.Properties.Location == nil {
			return fmt.Errorf("failed inferring location via datacenter's ID: datacenter %s location is nil: %w", dcId, err)
		}
		viper.Set(fn, dcId)
	}
	return nil
}

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

			err := validateOrInferEditionByTemplate(c) // sets FlagEdition if unset and possible to infer
			if err != nil {
				return fmt.Errorf("failed inferring or validating edition: %w", err)
			}

			err = validateEdition(c)
			if err != nil {
				return fmt.Errorf("failed validating edition specific flags: %w", err)
			}

			err = inferLocationByDatacenter(c)
			if err != nil {
				return fmt.Errorf("failed inferring location: %w", err)
			}
			return nil
		},
		CmdRun: func(c *core.CommandConfig) error {
			cluster := ionoscloud.CreateClusterProperties{}
			if fn := core.GetFlagName(c.NS, constants.FlagEdition); viper.IsSet(fn) {
				cluster.Edition = pointer.From(viper.GetString(constants.FlagEdition))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagType); viper.IsSet(fn) {
				cluster.Type = pointer.From(viper.GetString(constants.FlagType))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagTemplateId); viper.IsSet(fn) {
				// Old flag kept for backwards compatibility. Behaviour fully included in --template flag
				cluster.TemplateID = pointer.From(viper.GetString(constants.FlagTemplateId))
			} else {
				if fn := core.GetFlagName(c.NS, constants.FlagTemplate); viper.IsSet(fn) {
					tmplId, err := templates.Resolve(viper.GetString(constants.FlagTemplate))
					if err != nil {
						return err
					}
					cluster.TemplateID = pointer.From(tmplId)
				}
			}

			if fn := core.GetFlagName(c.NS, constants.FlagName); viper.IsSet(fn) {
				cluster.DisplayName = pointer.From(viper.GetString(constants.FlagName))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagVersion); viper.IsSet(fn) {
				cluster.MongoDBVersion = pointer.From(viper.GetString(constants.FlagVersion))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLocation); viper.IsSet(fn) {
				cluster.Location = pointer.From(viper.GetString(constants.FlagLocation))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagInstances); viper.IsSet(fn) {
				cluster.Instances = pointer.From(viper.GetInt32(constants.FlagInstances))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagShards); viper.IsSet(fn) {
				cluster.Shards = pointer.From(viper.GetInt32(constants.FlagShards))
			}

			cluster.MaintenanceWindow = &ionoscloud.MaintenanceWindow{}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceDay); viper.IsSet(fn) {
				cluster.MaintenanceWindow.DayOfTheWeek = (*ionoscloud.DayOfTheWeek)(pointer.From(
					viper.GetString(constants.FlagMaintenanceDay)))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagMaintenanceTime); viper.IsSet(fn) {
				cluster.MaintenanceWindow.Time = pointer.From(
					viper.GetString(constants.FlagMaintenanceTime))
			}

			cluster.Connections = pointer.From(make([]ionoscloud.Connection, 1))
			if fn := core.GetFlagName(c.NS, constants.FlagCidr); viper.IsSet(fn) {
				(*cluster.Connections)[0].CidrList = pointer.From(
					viper.GetStringSlice(constants.FlagCidr))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagDatacenterId); viper.IsSet(fn) {
				(*cluster.Connections)[0].DatacenterId = pointer.From(
					viper.GetString(constants.FlagDatacenterId))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagLanId); viper.IsSet(fn) {
				(*cluster.Connections)[0].LanId = pointer.From(
					viper.GetString(constants.FlagLanId))
			}

			// Enterprise flags
			if fn := core.GetFlagName(c.NS, constants.FlagCores); viper.IsSet(fn) {
				cluster.Cores = pointer.From(viper.GetInt32(constants.FlagCores))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageType); viper.IsSet(fn) {
				cluster.StorageType = (*ionoscloud.StorageType)(pointer.From(viper.GetString(constants.FlagStorageType)))
			}
			if fn := core.GetFlagName(c.NS, constants.FlagStorageSize); viper.IsSet(fn) {
				storageHuman := viper.GetString(constants.FlagStorageSize)
				utils.ConvertToGB()
				cluster.StorageSize = pointer.From()
			}
			return nil
		},
		InitClient: true,
	})

	cmd.AddSetFlag(constants.FlagEdition, "e", "", enumEditions, "Cluster Edition", core.RequiredFlagOption())
	cmd.AddSetFlag(constants.FlagType, "", "replicaset", enumTypes, "Cluster Type. Required for enterprise clusters. Not required (inferred) if using --shards or --instances")

	// Template
	cmd.AddStringFlag(constants.FlagTemplate, "", "", "The ID of a Mongo Template, or a word contained in the name of one. Templates specify the number of cores, storage size, and memory. (Required only for business edition)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTemplate, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ts, err := templates.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := functional.Fold(ts, func(acc []string, t ionoscloud.TemplateResponse) []string {
			if t.Properties == nil || t.Properties.Name == nil {
				return acc
			}
			wordsInTemplateName := strings.Split(*t.Properties.Name, " ")
			// Add only the last words of templates (e.g. 4XL, L, S, XS, Playground) since completions dont support spaces and they have multiple words in their names
			return append(acc, wordsInTemplateName[len(wordsInTemplateName)-1])
		}, nil)

		return names, cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagName, constants.FlagNameShort, "", "The name of your cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagVersion, "", "6.0", "The MongoDB version of your cluster", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "The physical location where the cluster will be created. (defaults to the location of the connected datacenter)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances in the replicaset cluster (one primary and n-1 secondaries). Setting this flag infers a replicaset type. Limited to at least 3 for business edition. (required for non-playground replicaset clusters)", core.RequiredFlagOption())
	cmd.AddInt32Flag(constants.FlagShards, "", 1, "The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters)", core.RequiredFlagOption())

	// Maintenance
	cmd.AddStringFlag(constants.FlagMaintenanceTime, "", "", "Time for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long window, during which maintenance might occur. e.g.: 16:30:59", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceTime, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"00:00:00", "04:00:00", "08:00:00", "10:00:00", "12:00:00", "16:00:00", "20:00:00"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagMaintenanceDay, "", "", "Day Of the Week for the MaintenanceWindows. The MaintenanceWindow is a weekly 4 hour-long windows, during which maintenance might occur. e.g.: Saturday", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagMaintenanceDay, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Enterprise-specific
	cmd.AddIntFlag(constants.FlagCores, "", 0, "The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagRam, "", "", "Custom RAM: multiples of 256. e.g. --ram 256 or --ram 256MB or --ram 4GB (required and only settable for enterprise edition)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"256MB", "512MB", "1024MB", "2GB", "4GB", "8GB", "12GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddSetFlag(constants.FlagStorageType, "", "", []string{"HDD", "SSD", "SSD Premium"},
		"Custom Storage Type. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	cmd.AddStringFlag(constants.FlagStorageSize, "", "", "Custom Storage: Greater performance for values greater than 100 GB. (required and only settable for enterprise edition)", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"50GB", "100GB", "200GB", "1TB", "10TB", "100TB", "1000TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Connections
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

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request [seconds]")

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
			constants.FlagTemplate, constants.FlagInstances,
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
