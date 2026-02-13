package cluster

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/cilium/fake"
	cloudapiv6completer "github.com/ionos-cloud/ionosctl/v6/commands/cloudapi-v6/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/completer"
	"github.com/ionos-cloud/ionosctl/v6/commands/dbaas/mongo/templates"
	"github.com/ionos-cloud/ionosctl/v6/internal/client"
	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/json2table/resource2table"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/jsontabwriter"
	"github.com/ionos-cloud/ionosctl/v6/internal/printer/tabheaders"
	"github.com/ionos-cloud/ionosctl/v6/pkg/convbytes"
	"github.com/ionos-cloud/ionosctl/v6/pkg/functional"
	"github.com/ionos-cloud/ionosctl/v6/pkg/pointer"
	"github.com/ionos-cloud/sdk-go-bundle/products/dbaas/mongo/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
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
		fmt.Sprintf("ionosctl dbaas mongo cluster create --%s enterprise %s [%s] %s",
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
			/*
			 * For edition playground, only replica-set, =1 instance and playground template (33457e53-1f8b-4ed2-8a12-2d42355aa759, 1 core, 50 GB Storage, 2 GB RAM).
			 * For edition business, only replica-set, >1 instance and any template.
			 * For edition enterprise, type replica-set/sharded-cluster and
			 *  - CPU Cores: 1-8
			 *  - RAM Size (GB): <16 GB
			 *  - Storage Size:  >100GB for optimal perf. max 1048.576 GB.
			 *  - Shards: 2-32 shards. (Infer sharded-cluster if set. Else, replicaset). A sharded cluster can still have multiple replicas
			 *  - Instances: >3
			**/

			err := validateOrInferEditionByTemplate(c) // sets FlagEdition if unset and possible to infer
			if err != nil {
				return fmt.Errorf("failed inferring or validating edition: %w", err)
			}

			err = inferTypeForEnterprise(c) // sets FlagType if unset and possible to infer
			if err != nil {
				return fmt.Errorf("failed inferring type: %w", err)
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
			cluster := mongo.CreateClusterProperties{}

			if c.Command.Command.Flags().Changed(constants.FlagEdition) {
				edition, err := c.Command.Command.Flags().GetString(constants.FlagEdition)
				if err != nil {
					return err
				}
				cluster.Edition = pointer.From(edition)
			}

			if c.Command.Command.Flags().Changed(constants.FlagType) {
				clusterType, err := c.Command.Command.Flags().GetString(constants.FlagType)
				if err != nil {
					return err
				}
				cluster.Type = pointer.From(clusterType)
			}

			if c.Command.Command.Flags().Changed(constants.FlagTemplateId) {
				// Old flag kept for backwards compatibility. Behaviour fully included in --template flag
				templateId, err := c.Command.Command.Flags().GetString(constants.FlagTemplateId)
				if err != nil {
					return err
				}
				cluster.TemplateID = pointer.From(templateId)
			} else if c.Command.Command.Flags().Changed(constants.FlagTemplate) {
				template, err := c.Command.Command.Flags().GetString(constants.FlagTemplate)
				if err != nil {
					return err
				}
				tmplId, err := templates.Resolve(template)
				if err != nil {
					return err
				}
				cluster.TemplateID = pointer.From(tmplId)
			}

			if c.Command.Command.Flags().Changed(constants.FlagName) {
				name, err := c.Command.Command.Flags().GetString(constants.FlagName)
				if err != nil {
					return err
				}
				cluster.DisplayName = name
			}

			if c.Command.Command.Flags().Changed(constants.FlagVersion) {
				version, err := c.Command.Command.Flags().GetString(constants.FlagVersion)
				if err != nil {
					return err
				}
				cluster.MongoDBVersion = version
			}

			if c.Command.Command.Flags().Changed(constants.FlagLocation) {
				location, err := c.Command.Command.Flags().GetString(constants.FlagLocation)
				if err != nil {
					return err
				}
				cluster.Location = location
			}

			if c.Command.Command.Flags().Changed(constants.FlagInstances) {
				instances, err := c.Command.Command.Flags().GetInt32(constants.FlagInstances)
				if err != nil {
					return err
				}
				cluster.Instances = instances
			}

			if c.Command.Command.Flags().Changed(constants.FlagShards) {
				shards, err := c.Command.Command.Flags().GetInt32(constants.FlagShards)
				if err != nil {
					return err
				}
				cluster.Shards = pointer.From(shards)
			}

			cluster.MaintenanceWindow = &mongo.MaintenanceWindow{}
			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceDay) {
				day, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceDay)
				if err != nil {
					return err
				}
				cluster.MaintenanceWindow.DayOfTheWeek = mongo.DayOfTheWeek(day)
			}

			if c.Command.Command.Flags().Changed(constants.FlagMaintenanceTime) {
				maintenanceTime, err := c.Command.Command.Flags().GetString(constants.FlagMaintenanceTime)
				if err != nil {
					return err
				}
				cluster.MaintenanceWindow.Time = maintenanceTime
			}

			cluster.Connections = make([]mongo.Connection, 1)
			if c.Command.Command.Flags().Changed(constants.FlagCidr) {
				cidrList, err := c.Command.Command.Flags().GetStringSlice(constants.FlagCidr)
				if err != nil {
					return err
				}
				cluster.Connections[0].CidrList = cidrList
			}

			if c.Command.Command.Flags().Changed(constants.FlagDatacenterId) {
				datacenterId, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
				if err != nil {
					return err
				}
				cluster.Connections[0].DatacenterId = datacenterId
			}

			if c.Command.Command.Flags().Changed(constants.FlagLanId) {
				lanId, err := c.Command.Command.Flags().GetString(constants.FlagLanId)
				if err != nil {
					return err
				}
				cluster.Connections[0].LanId = lanId
			}

			// backup flags
			cluster.Backup = nil
			if c.Command.Command.Flags().Changed(flagBackupLocation) {
				backupLocation, err := c.Command.Command.Flags().GetString(flagBackupLocation)
				if err != nil {
					return err
				}
				if cluster.Backup == nil {
					cluster.Backup = &mongo.BackupProperties{}
				}
				cluster.Backup.Location = pointer.From(backupLocation)
			}

			cluster.BiConnector = nil
			if c.Command.Command.Flags().Changed(flagBiconnector) {
				hostAndPort, err := c.Command.Command.Flags().GetString(flagBiconnector)
				if err != nil {
					return err
				}
				if cluster.BiConnector == nil {
					cluster.BiConnector = &mongo.BiConnectorProperties{}
				}
				host, port, err := net.SplitHostPort(hostAndPort)
				if err != nil {
					return fmt.Errorf("failed splitting --%s %s into host and port: %w",
						flagBiconnector, hostAndPort, err)
				}
				cluster.BiConnector.Enabled = pointer.From(true)
				cluster.BiConnector.Host = pointer.From(host)
				cluster.BiConnector.Port = pointer.From(port)
			}

			if c.Command.Command.Flags().Changed(flagBiconnectorEnabled) {
				biconnectorEnabled, err := c.Command.Command.Flags().GetBool(flagBiconnectorEnabled)
				if err != nil {
					return err
				}
				if !biconnectorEnabled && cluster.BiConnector != nil {
					cluster.BiConnector.Enabled = pointer.From(false)
				}
			}

			// Enterprise flags
			if c.Command.Command.Flags().Changed(constants.FlagCores) {
				cores, err := c.Command.Command.Flags().GetInt32(constants.FlagCores)
				if err != nil {
					return err
				}
				cluster.Cores = pointer.From(cores)
			}

			if c.Command.Command.Flags().Changed(constants.FlagStorageType) {
				storageType, err := c.Command.Command.Flags().GetString(constants.FlagStorageType)
				if err != nil {
					return err
				}
				cluster.StorageType = (*mongo.StorageType)(pointer.From(storageType))
			}

			if c.Command.Command.Flags().Changed(constants.FlagStorageSize) {
				storageSizeStr, err := c.Command.Command.Flags().GetString(constants.FlagStorageSize)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(storageSizeStr, convbytes.MB)
				cluster.StorageSize = pointer.From(int32(sizeInt64))
			}

			if c.Command.Command.Flags().Changed(constants.FlagRam) {
				ramStr, err := c.Command.Command.Flags().GetString(constants.FlagRam)
				if err != nil {
					return err
				}
				sizeInt64 := convbytes.StrToUnit(ramStr, convbytes.MB)
				cluster.Ram = pointer.From(int32(sizeInt64))
			}

			createdCluster, _, err := client.Must().MongoClient.ClustersApi.ClustersPost(context.Background()).CreateClusterRequest(
				mongo.CreateClusterRequest{Properties: &cluster},
			).Execute()
			if err != nil {
				return fmt.Errorf("failed creating cluster: %w", err)
			}

			cols, _ := c.Command.Command.Flags().GetStringSlice(constants.ArgCols)

			clusterConverted, err := resource2table.ConvertDbaasMongoClusterToTable(createdCluster)
			if err != nil {
				return err
			}

			out, err := jsontabwriter.GenerateOutputPreconverted(createdCluster, clusterConverted,
				tabheaders.GetHeaders(allCols, defaultCols, cols))
			if err != nil {
				return err
			}

			fmt.Fprintf(c.Command.Command.OutOrStdout(), "%s", out)
			return nil
		},
		InitClient: true,
	})

	cmd.AddSetFlag(constants.FlagEdition, "e", "", enumEditions, "Cluster Edition", core.RequiredFlagOption())
	cmd.AddSetFlag(constants.FlagType, "", "replicaset", enumTypes, "Cluster Type. Required for enterprise clusters. Not required (inferred) if using --shards or --instances")

	cmd.AddStringFlag(constants.FlagTemplateId, "", "", "The ID of a Mongo Template. Please use --template instead")
	cmd.Command.Flags().MarkHidden(constants.FlagTemplateId)

	// Template
	cmd.AddStringFlag(constants.FlagTemplate, "", "", "The ID of a Mongo Template, or a word contained in the name of one. "+
		"Templates specify the number of cores, storage size, and memory. Business editions default to XS template. Playground editions default to playground template.")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagTemplate, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		ts, err := templates.List()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}
		names := functional.Fold(ts, func(acc []string, t mongo.TemplateResponse) []string {
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
	cmd.AddStringFlag(constants.FlagVersion, "", "7.0", "The MongoDB version of your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagVersion, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return completer.MongoDBVersions(), cobra.ShellCompDirectiveNoFileComp
	})

	cmd.AddStringFlag(constants.FlagLocation, constants.FlagLocationShort, "", "The physical location where the cluster will be created. (defaults to the location of the connected datacenter)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLocation, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.LocationIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagInstances, "", 1, "The total number of instances of the cluster (one primary and n-1 secondaries). Minimum of 3 for enterprise edition")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagInstances, func(cmdCobra *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1", "3", "5", "7"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddInt32Flag(constants.FlagShards, "", 1, "The total number of shards in the sharded_cluster cluster. Setting this flag is only possible for enterprise clusters and infers a sharded_cluster type. Possible values: 2 - 32. (required for sharded_cluster enterprise clusters)", core.RequiredFlagOption())

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
	// Enterprise-specific
	cmd.AddInt32Flag(constants.FlagCores, "", 1, "The total number of cores for the Server, e.g. 4. (required and only settable for enterprise edition)")
	cmd.AddStringFlag(constants.FlagRam, "", "2GB", "Custom RAM: multiples of 1024. e.g. --ram 1024 or --ram 1024MB or --ram 4GB (required and only settable for enterprise edition)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagRam, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"1024MB", "2GB", "4GB", "8GB", "12GB", "16GB"}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageType, "", "\"SSD Standard\"",
		"Custom Storage Type. (only settable for enterprise edition)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageType, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"HDD", "\"SSD Standard\"", "\"SSD Premium\""}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagStorageSize, "", "5GB", "Custom Storage: Minimum of 5GB, Greater performance for values greater than 100 GB. (only settable for enterprise edition)")
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagStorageSize, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"5GB", "10GB", "50GB", "100GB", "200GB", "400GB", "800GB", "1TB", "2TB"}, cobra.ShellCompDirectiveNoFileComp
	})

	// Connections
	cmd.AddStringFlag(constants.FlagDatacenterId, "", "", "The datacenter to which your cluster will be connected. Must be in the same location as the cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagDatacenterId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return cloudapiv6completer.DataCentersIds(), cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(constants.FlagLanId, "", "", "The numeric LAN ID with which you connect your cluster", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagLanId, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		datacenterId, _ := c.Flags().GetString(constants.FlagDatacenterId)
		return cloudapiv6completer.LansIds(datacenterId),
			cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringSliceFlag(constants.FlagCidr, "", nil, "The list of IPs and subnet for your cluster. All IPs must be in a /24 network", core.RequiredFlagOption())
	_ = cmd.Command.RegisterFlagCompletionFunc(constants.FlagCidr, func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var cidrs []string
		instances, _ := c.Flags().GetInt32(constants.FlagInstances)
		for i := 0; i < int(instances); i++ {
			cidrs = append(cidrs, fake.IP(fake.WithIPv4(), fake.WithIPCIDR("192.168.1.128/25"))+"/24")
		}

		return []string{strings.Join(cidrs, ",")}, cobra.ShellCompDirectiveNoFileComp
	})
	cmd.AddStringFlag(flagBackupLocation, "", "", "The location where the cluster backups will be stored. If not set, the backup is stored in the nearest location of the cluster")

	// From Backup: Snapshot_ID,Recovery_Target_Time

	// Biconnector
	cmd.AddStringFlag(flagBiconnector, "", "", "BI Connector host & port. The MongoDB Connector for Business Intelligence allows you to query a MongoDB database using SQL commands. Example: r1.m-abcdefgh1234.mongodb.de-fra.ionos.com:27015")
	cmd.AddBoolFlag(flagBiconnectorEnabled, "", true, fmt.Sprintf("Enable or disable the biconnector. To disable it, use --%s=false", flagBiconnectorEnabled))

	// Misc
	cmd.AddBoolFlag(constants.ArgWaitForRequest, constants.ArgWaitForRequestShort, constants.DefaultWait, "Wait for the Request to be executed")
	cmd.AddIntFlag(constants.ArgTimeout, constants.ArgTimeoutShort, constants.DefaultTimeoutSeconds, "Timeout option for Request [seconds]")

	// They do nothing... but we can't outright remove them in case some user already uses them in their scripts
	// would cause ('unknown flag: -w')
	cmd.Command.Flags().MarkHidden(constants.ArgWaitForRequest)
	cmd.Command.Flags().MarkHidden(constants.ArgTimeout)

	cmd.Command.SilenceUsage = true
	cmd.Command.Flags().SortFlags = false

	return cmd
}

// returns a slice of flags to be marked as required, depending on wanted edition and type
func getRequiredFlagsByEditionAndType(edition, cType string) ([]string, error) {
	alwaysRequired := []string{
		constants.FlagEdition, constants.FlagName,
		constants.FlagDatacenterId, constants.FlagLanId, constants.FlagCidr,
	}
	switch edition {
	case "playground":
		// Type inferred as replicaset. Template inferred as type playground. Instances inferred as 1
		return alwaysRequired, nil
	case "business":
		// Type inferred as replicaset.
		return append(alwaysRequired,
			// constants.FlagTemplate, // Decided to be defaulted to XS
			constants.FlagInstances,
		), nil
	case "enterprise":
		enterpriseBaseFlags := alwaysRequired
		// enterpriseBaseFlags := append(alwaysRequired, // Decided to be defaulted to lowest vals
		// constants.FlagCores, constants.FlagStorageType,
		// constants.FlagStorageSize, constants.FlagRam,
		// )
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

// SIDE EFFECT: sets FlagEdition if not set and can be inferred
func validateOrInferEditionByTemplate(c *core.PreCommandConfig) error {
	if c.Command.Command.Flags().Changed(constants.FlagTemplate) {
		templateStr, err := c.Command.Command.Flags().GetString(constants.FlagTemplate)
		if err != nil {
			return err
		}
		tmplId, err := templates.Resolve(templateStr)
		if err != nil {
			// Intentionally don't wrap this error since the deeper error would kind of say the same thing
			return err
		}
		template, err := templates.Find(func(x mongo.TemplateResponse) bool {
			return *x.Id == tmplId
		})
		if err != nil {
			return fmt.Errorf("failed finding template with ID %s: %w", tmplId, err)
		}

		if template.Properties == nil || template.Id == nil ||
			template.Properties.Edition == nil || template.Properties.Name == nil {
			return fmt.Errorf("found a template with some unset fields: %#v.\n Please use IONOS_LOG_LEVEL=trace and file a Github Issue", template)
		}

		if c.Command.Command.Flags().Changed(constants.FlagEdition) {
			edition, err := c.Command.Command.Flags().GetString(constants.FlagEdition)
			if err != nil {
				return err
			}
			// Check that template & edition aren't set to incompatible things

			if edition == "enterprise" && c.Command.Command.Flags().Changed(constants.FlagTemplate) {
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
	if !c.Command.Command.Flags().Changed(constants.FlagEdition) {
		return fmt.Errorf("set --%s or --%s (%s) to get a list of required flags",
			constants.FlagTemplate, constants.FlagEdition, strings.Join(enumEditions, " | "))
	}

	edition, err := c.Command.Command.Flags().GetString(constants.FlagEdition)
	if err != nil {
		return err
	}

	// Enterprise edition cannot have --template-id
	if edition == "enterprise" && c.Command.Command.Flags().Changed(constants.FlagTemplate) {
		return fmt.Errorf("for enterprise edition, setting --%s is forbidden. Use %s", constants.FlagTemplate,
			core.FlagsUsage(constants.FlagCores, constants.FlagRam, constants.FlagStorageType, constants.FlagStorageSize))
	}

	// Business edition: infer template as XS
	// Playground edition: infer template as Playground
	fn := core.GetFlagName(c.NS, constants.FlagTemplate)
	if edition == "business" && !c.Command.Command.Flags().Changed(constants.FlagTemplate) {
		viper.Set(fn, "XS")
	} else if edition == "playground" && !c.Command.Command.Flags().Changed(constants.FlagTemplate) {
		viper.Set(fn, "Playground")
	}

	// Special case for playground: infer that instances is 1
	flagInstances := core.GetFlagName(c.NS, constants.FlagInstances)
	if edition == "playground" && !c.Command.Command.Flags().Changed(constants.FlagInstances) {
		viper.Set(flagInstances, 1)
	}

	clusterType, _ := c.Command.Command.Flags().GetString(constants.FlagType)
	flags, err := getRequiredFlagsByEditionAndType(edition, clusterType)
	if err != nil {
		return fmt.Errorf("failed getting required flags for edition %s: %w", edition, err)
	}

	err = core.CheckRequiredFlags(c.Command, c.NS, flags...)
	if err != nil {
		return fmt.Errorf("not all %s edition flags are set: %w", edition, err)
	}

	return nil
}

// SIDE EFFECT: sets FlagType if not set and can be inferred via --instances or --sharded
func inferTypeForEnterprise(c *core.PreCommandConfig) error {
	fn := core.GetFlagName(c.NS, constants.FlagType)
	if c.Command.Command.Flags().Changed(constants.FlagType) {
		return nil
	}

	viper.Set(fn, "replicaset")
	if c.Command.Command.Flags().Changed(constants.FlagShards) {
		viper.Set(fn, "sharded-cluster")
	}

	return nil
}

// SIDE EFFECT: sets FlagLocation if not set and can be inferred
func inferLocationByDatacenter(c *core.PreCommandConfig) error {
	fn := core.GetFlagName(c.NS, constants.FlagLocation)
	if !c.Command.Command.Flags().Changed(constants.FlagLocation) {
		dcId, err := c.Command.Command.Flags().GetString(constants.FlagDatacenterId)
		if err != nil {
			return err
		}
		dc, _, err := client.Must().CloudClient.DataCentersApi.DatacentersFindById(context.Background(), dcId).Execute()
		if err != nil {
			return fmt.Errorf("failed inferring location via datacenter's ID: failed getting datacenter with ID %s: %w", dcId, err)
		}
		if dc.Properties == nil || dc.Properties.Location == nil {
			return fmt.Errorf("failed inferring location via datacenter's ID: datacenter %s location is nil: %w", dcId, err)
		}
		viper.Set(fn, *dc.Properties.Location)
	}
	return nil
}
