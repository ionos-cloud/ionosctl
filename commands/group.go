package commands

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/fatih/structs"
	"github.com/ionos-cloud/ionosctl/pkg/config"
	"github.com/ionos-cloud/ionosctl/pkg/core"
	"github.com/ionos-cloud/ionosctl/pkg/resources"
	"github.com/ionos-cloud/ionosctl/pkg/utils"
	"github.com/ionos-cloud/ionosctl/pkg/utils/clierror"
	"github.com/ionos-cloud/ionosctl/pkg/utils/printer"
	ionoscloud "github.com/ionos-cloud/sdk-go/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func group() *core.Command {
	ctx := context.TODO()
	groupCmd := &core.Command{
		Command: &cobra.Command{
			Use:              "group",
			Short:            "Group Operations",
			Long:             `The sub-commands of ` + "`" + `ionosctl group` + "`" + ` allow you to list, get, create, update, delete Groups, but also operations: add/remove/list/update User from the Group.`,
			TraverseChildren: true,
		},
	}
	globalFlags := groupCmd.GlobalFlags()
	globalFlags.StringSlice(config.ArgCols, defaultGroupCols, "Columns to be printed in the standard output")
	_ = viper.BindPFlag(core.GetGlobalFlagName(groupCmd.Name(), config.ArgCols), globalFlags.Lookup(config.ArgCols))
	_ = groupCmd.Command.RegisterFlagCompletionFunc(config.ArgCols, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return defaultGroupCols, cobra.ShellCompDirectiveNoFileComp
	})

	/*
		List Command
	*/
	core.NewCommand(ctx, groupCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "group",
		Verb:       "list",
		ShortDesc:  "List Groups",
		LongDesc:   "Use this command to get a list of available Groups available on your account.",
		Example:    listGroupExample,
		PreCmdRun:  noPreRun,
		CmdRun:     RunGroupList,
		InitClient: true,
	})

	/*
		Get Command
	*/
	get := core.NewCommand(ctx, groupCmd, core.CommandBuilder{
		Namespace:  "group",
		Resource:   "group",
		Verb:       "get",
		ShortDesc:  "Get a Group",
		LongDesc:   "Use this command to retrieve details about a specific Group.\n\nRequired values to run command:\n\n* Group Id",
		Example:    getGroupExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupGet,
		InitClient: true,
	})
	get.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = get.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})

	/*
		Create Command
	*/
	create := core.NewCommand(ctx, groupCmd, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "create",
		ShortDesc: "Create a Group",
		LongDesc: `Use this command to create a new Group and set Group privileges. You need to specify the name for the new Group. By default, all privileges will be set to false. You need to use flags privileges to be set to true.

Required values to run a command:

* Group Name`,
		Example:    createGroupExample,
		PreCmdRun:  PreRunGroupName,
		CmdRun:     RunGroupCreate,
		InitClient: true,
	})
	create.AddStringFlag(config.ArgGroupName, "", "", "Name for the Group "+config.RequiredFlag)
	create.AddBoolFlag(config.ArgGroupCreateDc, "", false, "The group will be allowed to create Data Centers")
	create.AddBoolFlag(config.ArgGroupCreateSnapshot, "", false, "The group will be allowed to create Snapshots")
	create.AddBoolFlag(config.ArgGroupReserveIp, "", false, "The group will be allowed to reserve IP addresses")
	create.AddBoolFlag(config.ArgGroupAccessLog, "", false, "The group will be allowed to access the activity log")
	create.AddBoolFlag(config.ArgGroupCreatePcc, "", false, "The group will be allowed to create PCCs")
	create.AddBoolFlag(config.ArgGroupS3Privilege, "", false, "The group will be allowed to manage S3")
	create.AddBoolFlag(config.ArgGroupCreateBackUpUnit, "", false, "The group will be able to manage Backup Units")
	create.AddBoolFlag(config.ArgGroupCreateNic, "", false, "The group will be allowed to create NICs")
	create.AddBoolFlag(config.ArgGroupCreateK8s, "", false, "The group will be allowed to create K8s Clusters")
	create.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Request for Group creation to be executed")
	create.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Group creation [seconds]")

	/*
		Update Command
	*/
	update := core.NewCommand(ctx, groupCmd, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "update",
		ShortDesc: "Update a Group",
		LongDesc: `Use this command to update details about a specific Group.

You can wait for the Request to be executed using ` + "`" + `--wait-for-request` + "`" + ` option.

Required values to run command:

* Group Id`,
		Example:    updateGroupExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupUpdate,
		InitClient: true,
	})
	update.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = update.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	update.AddStringFlag(config.ArgGroupName, "", "", "Name for the Group "+config.RequiredFlag)
	update.AddBoolFlag(config.ArgGroupCreateDc, "", false, "The group will be allowed to create Data Centers")
	update.AddBoolFlag(config.ArgGroupCreateSnapshot, "", false, "The group will be allowed to create Snapshots")
	update.AddBoolFlag(config.ArgGroupReserveIp, "", false, "The group will be allowed to reserve IP addresses")
	update.AddBoolFlag(config.ArgGroupAccessLog, "", false, "The group will be allowed to access the activity log")
	update.AddBoolFlag(config.ArgGroupCreatePcc, "", false, "The group will be allowed to create PCCs")
	update.AddBoolFlag(config.ArgGroupS3Privilege, "", false, "The group will be allowed to manage S3")
	update.AddBoolFlag(config.ArgGroupCreateBackUpUnit, "", false, "The group will be able to manage Backup Units")
	update.AddBoolFlag(config.ArgGroupCreateNic, "", false, "The group will be allowed to create NICs")
	update.AddBoolFlag(config.ArgGroupCreateK8s, "", false, "The group will be allowed to create K8s Clusters")
	update.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Request for Group update to be executed")
	update.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Group update [seconds]")

	/*
		Delete Command
	*/
	deleteCmd := core.NewCommand(ctx, groupCmd, core.CommandBuilder{
		Namespace: "group",
		Resource:  "group",
		Verb:      "delete",
		ShortDesc: "Delete a Group",
		LongDesc: `Use this operation to delete a single Group. Resources that are assigned to the Group are NOT deleted, but are no longer accessible to the Group members unless the member is a Contract Owner, Admin, or Resource Owner.

Required values to run command:

* Group Id`,
		Example:    deleteGroupExample,
		PreCmdRun:  PreRunGroupId,
		CmdRun:     RunGroupDelete,
		InitClient: true,
	})
	deleteCmd.AddStringFlag(config.ArgGroupId, "", "", config.RequiredFlagGroupId)
	_ = deleteCmd.Command.RegisterFlagCompletionFunc(config.ArgGroupId, func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return getGroupsIds(os.Stderr), cobra.ShellCompDirectiveNoFileComp
	})
	deleteCmd.AddBoolFlag(config.ArgWaitForRequest, "", config.DefaultWait, "Wait for Request for Group deletion to be executed")
	deleteCmd.AddIntFlag(config.ArgTimeout, "", config.DefaultTimeoutSeconds, "Timeout option for Request for Group deletion [seconds]")

	groupCmd.AddCommand(groupResource())
	groupCmd.AddCommand(groupUser())
	return groupCmd
}

func PreRunGroupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgGroupId)
}

func PreRunGroupUserIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgGroupId, config.ArgUserId)
}

func PreRunGroupName(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.NS, config.ArgGroupName)
}

func RunGroupList(c *core.CommandConfig) error {
	groups, _, err := c.Groups().List()
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(nil, c, getGroups(groups)))
}

func RunGroupGet(c *core.CommandConfig) error {
	u, _, err := c.Groups().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(nil, c, getGroup(u)))
}

func RunGroupCreate(c *core.CommandConfig) error {
	properties := getGroupCreateInfo(c)
	newGroup := resources.Group{
		Group: ionoscloud.Group{
			Properties: &properties.GroupProperties,
		},
	}
	u, resp, err := c.Groups().Create(newGroup)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(nil, c, getGroup(u)))
}

func RunGroupUpdate(c *core.CommandConfig) error {
	u, resp, err := c.Groups().Get(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if err != nil {
		return err
	}
	properties := getGroupUpdateInfo(u, c)
	newGroup := resources.Group{
		Group: ionoscloud.Group{
			Properties: &properties.GroupProperties,
		},
	}
	groupUpd, resp, err := c.Groups().Update(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)), newGroup)
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(resp, c, getGroup(groupUpd)))
}

func RunGroupDelete(c *core.CommandConfig) error {
	if err := utils.AskForConfirm(c.Stdin, c.Printer, "delete group"); err != nil {
		return err
	}
	resp, err := c.Groups().Delete(viper.GetString(core.GetFlagName(c.NS, config.ArgGroupId)))
	if err != nil {
		return err
	}
	if err = utils.WaitForRequest(c, printer.GetRequestPath(resp)); err != nil {
		return err
	}
	return c.Printer.Print(getGroupPrint(resp, c, nil))
}

func getGroupCreateInfo(c *core.CommandConfig) *resources.GroupProperties {
	name := viper.GetString(core.GetFlagName(c.NS, config.ArgGroupName))
	createDc := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateDc))
	createSnap := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateSnapshot))
	reserveIp := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupReserveIp))
	accessLog := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupAccessLog))
	createBackUp := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateBackUpUnit))
	createPcc := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreatePcc))
	createNic := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateNic))
	createK8s := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateK8s))
	s3 := viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupS3Privilege))
	return &resources.GroupProperties{
		GroupProperties: ionoscloud.GroupProperties{
			Name:                 &name,
			CreateDataCenter:     &createDc,
			CreateSnapshot:       &createSnap,
			ReserveIp:            &reserveIp,
			AccessActivityLog:    &accessLog,
			CreatePcc:            &createPcc,
			S3Privilege:          &s3,
			CreateBackupUnit:     &createBackUp,
			CreateInternetAccess: &createNic,
			CreateK8sCluster:     &createK8s,
		},
	}
}

func getGroupUpdateInfo(oldGroup *resources.Group, c *core.CommandConfig) *resources.GroupProperties {
	var (
		groupName                                                           string
		createDc, createSnap, createPcc, createBackUp, createNic, createK8s bool
		reserveIp, accessLog, s3                                            bool
	)
	if properties, ok := oldGroup.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupName)) {
			groupName = viper.GetString(core.GetFlagName(c.NS, config.ArgGroupName))
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				groupName = *name
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreateDc)) {
			createDc = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateDc))
		} else {
			if dc, ok := properties.GetCreateDataCenterOk(); ok && dc != nil {
				createDc = *dc
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreateSnapshot)) {
			createSnap = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateSnapshot))
		} else {
			if s, ok := properties.GetCreateSnapshotOk(); ok && s != nil {
				createSnap = *s
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreatePcc)) {
			createPcc = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreatePcc))
		} else {
			if s, ok := properties.GetCreatePccOk(); ok && s != nil {
				createPcc = *s
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreateK8s)) {
			createK8s = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateK8s))
		} else {
			if s, ok := properties.GetCreateK8sClusterOk(); ok && s != nil {
				createK8s = *s
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreateNic)) {
			createNic = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateNic))
		} else {
			if s, ok := properties.GetCreateInternetAccessOk(); ok && s != nil {
				createNic = *s
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupCreateBackUpUnit)) {
			createBackUp = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupCreateBackUpUnit))
		} else {
			if s, ok := properties.GetCreateBackupUnitOk(); ok && s != nil {
				createBackUp = *s
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupReserveIp)) {
			reserveIp = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupReserveIp))
		} else {
			if ip, ok := properties.GetReserveIpOk(); ok && ip != nil {
				reserveIp = *ip
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupAccessLog)) {
			accessLog = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupAccessLog))
		} else {
			if log, ok := properties.GetAccessActivityLogOk(); ok && log != nil {
				accessLog = *log
			}
		}
		if viper.IsSet(core.GetFlagName(c.NS, config.ArgGroupS3Privilege)) {
			s3 = viper.GetBool(core.GetFlagName(c.NS, config.ArgGroupS3Privilege))
		} else {
			if s, ok := properties.GetS3PrivilegeOk(); ok && s != nil {
				s3 = *s
			}
		}
	}
	return &resources.GroupProperties{
		GroupProperties: ionoscloud.GroupProperties{
			Name:                 &groupName,
			CreateDataCenter:     &createDc,
			CreateSnapshot:       &createSnap,
			ReserveIp:            &reserveIp,
			AccessActivityLog:    &accessLog,
			CreatePcc:            &createPcc,
			S3Privilege:          &s3,
			CreateBackupUnit:     &createBackUp,
			CreateInternetAccess: &createNic,
			CreateK8sCluster:     &createK8s,
		},
	}
}

// Output Printing

var defaultGroupCols = []string{"GroupId", "Name", "CreateDataCenter", "CreateSnapshot", "ReserveIp", "AccessActivityLog", "CreatePcc", "S3Privilege", "CreateBackupUnit", "CreateInternetAccess", "CreateK8s"}

type groupPrint struct {
	GroupId              string `json:"GroupId,omitempty"`
	Name                 string `json:"Name,omitempty"`
	CreateDataCenter     bool   `json:"CreateDataCenter,omitempty"`
	CreateSnapshot       bool   `json:"CreateSnapshot,omitempty"`
	ReserveIp            bool   `json:"ReserveIp,omitempty"`
	AccessActivityLog    bool   `json:"AccessActivityLog,omitempty"`
	CreatePcc            bool   `json:"CreatePcc,omitempty"`
	S3Privilege          bool   `json:"S3Privilege,omitempty"`
	CreateBackupUnit     bool   `json:"CreateBackupUnit,omitempty"`
	CreateInternetAccess bool   `json:"CreateInternetAccess,omitempty"`
	CreateK8s            bool   `json:"CreateK8s,omitempty"`
}

func getGroupPrint(resp *resources.Response, c *core.CommandConfig, groups []resources.Group) printer.Result {
	r := printer.Result{}
	if c != nil {
		if resp != nil {
			r.ApiResponse = resp
			r.Resource = c.Resource
			r.Verb = c.Verb
			r.WaitForRequest = viper.GetBool(core.GetFlagName(c.NS, config.ArgWaitForRequest))
		}
		if groups != nil {
			r.OutputJSON = groups
			r.KeyValue = getGroupsKVMaps(groups)
			r.Columns = getGroupCols(core.GetGlobalFlagName(c.Resource, config.ArgCols), c.Printer.GetStderr())
		}
	}
	return r
}

func getGroupCols(flagName string, outErr io.Writer) []string {
	if viper.IsSet(flagName) {
		var groupCols []string
		columnsMap := map[string]string{
			"GroupId":              "GroupId",
			"Name":                 "Name",
			"CreateDataCenter":     "CreateDataCenter",
			"CreateSnapshot":       "CreateSnapshot",
			"ReserveIp":            "ReserveIp",
			"AccessActivityLog":    "AccessActivityLog",
			"CreatePcc":            "CreatePcc",
			"S3Privilege":          "S3Privilege",
			"CreateBackupUnit":     "CreateBackupUnit",
			"CreateInternetAccess": "CreateInternetAccess",
			"CreateK8s":            "CreateK8s",
		}
		for _, k := range viper.GetStringSlice(flagName) {
			col := columnsMap[k]
			if col != "" {
				groupCols = append(groupCols, col)
			} else {
				clierror.CheckError(errors.New("unknown column "+k), outErr)
			}
		}
		return groupCols
	} else {
		return defaultGroupCols
	}
}

func getGroups(groups resources.Groups) []resources.Group {
	u := make([]resources.Group, 0)
	if items, ok := groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			u = append(u, resources.Group{Group: item})
		}
	}
	return u
}

func getGroup(u *resources.Group) []resources.Group {
	groups := make([]resources.Group, 0)
	if u != nil {
		groups = append(groups, resources.Group{Group: u.Group})
	}
	return groups
}

func getGroupsKVMaps(gs []resources.Group) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(gs))
	for _, g := range gs {
		var gPrint groupPrint
		if id, ok := g.GetIdOk(); ok && id != nil {
			gPrint.GroupId = *id
		}
		if properties, ok := g.GetPropertiesOk(); ok && properties != nil {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				gPrint.Name = *name
			}
			if createDc, ok := properties.GetCreateDataCenterOk(); ok && createDc != nil {
				gPrint.CreateDataCenter = *createDc
			}
			if createSnapshot, ok := properties.GetCreateSnapshotOk(); ok && createSnapshot != nil {
				gPrint.CreateSnapshot = *createSnapshot
			}
			if reserveIp, ok := properties.GetReserveIpOk(); ok && reserveIp != nil {
				gPrint.ReserveIp = *reserveIp
			}
			if accessLog, ok := properties.GetAccessActivityLogOk(); ok && accessLog != nil {
				gPrint.AccessActivityLog = *accessLog
			}
			if createPcc, ok := properties.GetCreatePccOk(); ok && createPcc != nil {
				gPrint.CreatePcc = *createPcc
			}
			if s3, ok := properties.GetS3PrivilegeOk(); ok && s3 != nil {
				gPrint.S3Privilege = *s3
			}
			if createBackup, ok := properties.GetCreateBackupUnitOk(); ok && createBackup != nil {
				gPrint.CreateBackupUnit = *createBackup
			}
			if createNic, ok := properties.GetCreateInternetAccessOk(); ok && createNic != nil {
				gPrint.CreateInternetAccess = *createNic
			}
			if createK8s, ok := properties.GetCreateK8sClusterOk(); ok && createK8s != nil {
				gPrint.CreateK8s = *createK8s
			}
		}
		o := structs.Map(gPrint)
		out = append(out, o)
	}
	return out
}

func getGroupsIds(outErr io.Writer) []string {
	err := config.Load()
	clierror.CheckError(err, outErr)
	clientSvc, err := resources.NewClientService(
		viper.GetString(config.Username),
		viper.GetString(config.Password),
		viper.GetString(config.Token),
		viper.GetString(config.ArgServerUrl),
	)
	clierror.CheckError(err, outErr)
	groupSvc := resources.NewGroupService(clientSvc.Get(), context.TODO())
	groups, _, err := groupSvc.List()
	clierror.CheckError(err, outErr)
	groupsIds := make([]string, 0)
	if items, ok := groups.Groups.GetItemsOk(); ok && items != nil {
		for _, item := range *items {
			if itemId, ok := item.GetIdOk(); ok && itemId != nil {
				groupsIds = append(groupsIds, *itemId)
			}
		}
	} else {
		return nil
	}
	return groupsIds
}
