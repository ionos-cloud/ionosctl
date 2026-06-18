package group

import (
	"fmt"

	"github.com/ionos-cloud/ionosctl/v6/internal/constants"
	"github.com/ionos-cloud/ionosctl/v6/internal/core"
	"github.com/ionos-cloud/ionosctl/v6/internal/request"
	"github.com/ionos-cloud/ionosctl/v6/pkg/confirm"
	cloudapiv6 "github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
	ionoscloud "github.com/ionos-cloud/sdk-go/v6"
	"github.com/spf13/viper"
)

func PreRunGroupId(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgGroupId)
}

func PreRunGroupDelete(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlagsSets(c.Command, c.NS,
		[]string{cloudapiv6.ArgGroupId},
		[]string{cloudapiv6.ArgAll},
	)
}

func PreRunGroupUserIds(c *core.PreCommandConfig) error {
	return core.CheckRequiredFlags(c.Command, c.NS, cloudapiv6.ArgGroupId, cloudapiv6.ArgUserId)
}

func RunGroupList(c *core.CommandConfig) error {
	groups, resp, err := c.CloudApiV6Services.Groups().List()
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupCols).Prefix("items").Print(groups.Groups)
}

func RunGroupGet(c *core.CommandConfig) error {
	c.Verbose("Group with id: %v is getting...", viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))

	u, resp, err := c.CloudApiV6Services.Groups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if resp != nil {
		c.Verbose(constants.MessageRequestTime, resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupCols).Print(u.Group)
}

func RunGroupCreate(c *core.CommandConfig) error {
	properties := getGroupCreateInfo(c)

	newGroup := resources.Group{
		Group: ionoscloud.Group{
			Properties: &properties.GroupProperties,
		},
	}
	u, resp, err := c.CloudApiV6Services.Groups().Create(newGroup)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupCols).Print(u.Group)
}

func RunGroupUpdate(c *core.CommandConfig) error {
	u, resp, err := c.CloudApiV6Services.Groups().Get(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)))
	if err != nil {
		return err
	}

	properties := getGroupUpdateInfo(u, c)
	newGroup := resources.Group{
		Group: ionoscloud.Group{
			Properties: &properties.GroupProperties,
		},
	}
	groupUpd, resp, err := c.CloudApiV6Services.Groups().Update(viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId)), newGroup)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	return c.Printer(allGroupCols).Print(groupUpd.Group)
}

func RunGroupDelete(c *core.CommandConfig) error {
	groupId := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgGroupId))

	if viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAll)) {
		if err := DeleteAllGroups(c); err != nil {
			return err
		}

		return nil
	}

	if !confirm.FAsk(c.Command.Command.InOrStdin(), "delete group", viper.GetBool(constants.ArgForce)) {
		return fmt.Errorf(confirm.UserDenied)
	}

	c.Verbose("Starting deleting Group with id: %v...", groupId)

	resp, err := c.CloudApiV6Services.Groups().Delete(groupId)
	if resp != nil && request.GetId(resp) != "" {
		c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
	}
	if err != nil {
		return err
	}

	c.Msg("Group successfully deleted")

	return nil

}

func getGroupCreateInfo(c *core.CommandConfig) *resources.GroupProperties {
	name := viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))
	createDc := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateDc))
	createSnap := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateSnapshot))
	reserveIp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgReserveIp))
	accessLog := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessLog))
	createBackUp := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateBackUpUnit))
	createPcc := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreatePcc))
	createNic := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateNic))
	createK8s := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateK8s))
	s3 := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgS3Privilege))
	createFlowLog := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateFlowLog))
	monitoring := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessMonitoring))
	certs := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessCerts))
	dns := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessDNS))
	manageDb := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgManageDbaas))
	manageReg := viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgManageRegistry))

	c.Verbose("Properties set for creating the group: Name: %v, CreateDatacenter: %v, CreateSnapshot: %v, "+
		"ReserveIp: %v, AccessActivityLog: %v, CreateBackupUnit: %v, CreatePcc: %v, CreateInternetAccess: %v, CreateK8sCluster: %v, "+
		"S3Privilege: %v, CreateFlowLog: %v, AccessAndManageMonitoring: %v, AccessAndManageCertificates: %v, AccessAndManageDns: %v,"+
		"ManageDBaaS: %v, ManageRegistry: %v",
		name, createDc, createSnap, reserveIp, accessLog, createBackUp, createPcc, createNic, createK8s, s3, createFlowLog, monitoring, certs,
		dns, manageDb, manageReg)

	return &resources.GroupProperties{
		GroupProperties: ionoscloud.GroupProperties{
			Name:                        &name,
			CreateDataCenter:            &createDc,
			CreateSnapshot:              &createSnap,
			ReserveIp:                   &reserveIp,
			AccessActivityLog:           &accessLog,
			CreatePcc:                   &createPcc,
			S3Privilege:                 &s3,
			CreateBackupUnit:            &createBackUp,
			CreateInternetAccess:        &createNic,
			CreateK8sCluster:            &createK8s,
			CreateFlowLog:               &createFlowLog,
			AccessAndManageMonitoring:   &monitoring,
			AccessAndManageCertificates: &certs,
			AccessAndManageDns:          &dns,
			ManageDBaaS:                 &manageDb,
			ManageRegistry:              &manageReg,
		},
	}
}

func getGroupUpdateInfo(oldGroup *resources.Group, c *core.CommandConfig) *resources.GroupProperties {
	var (
		groupName                                                           string
		createDc, createSnap, createPcc, createBackUp, createNic, createK8s bool
		reserveIp, accessLog, s3, createFlowLog, monitoring, certs, dns     bool
		manageReg, manageDb                                                 bool
	)

	if properties, ok := oldGroup.GetPropertiesOk(); ok && properties != nil {
		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgName)) {
			groupName = viper.GetString(core.GetFlagName(c.NS, cloudapiv6.ArgName))

			c.Verbose("Property Name set: %v", groupName)
		} else {
			if name, ok := properties.GetNameOk(); ok && name != nil {
				groupName = *name
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateDc)) {
			createDc = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateDc))

			c.Verbose("Property CreateDataCenter set: %v", createDc)
		} else {
			if dc, ok := properties.GetCreateDataCenterOk(); ok && dc != nil {
				createDc = *dc
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateSnapshot)) {
			createSnap = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateSnapshot))

			c.Verbose("Property CreateSnapshot set: %v", createSnap)
		} else {
			if s, ok := properties.GetCreateSnapshotOk(); ok && s != nil {
				createSnap = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreatePcc)) {
			createPcc = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreatePcc))

			c.Verbose("Property CreatePcc set: %v", createPcc)
		} else {
			if s, ok := properties.GetCreatePccOk(); ok && s != nil {
				createPcc = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateK8s)) {
			createK8s = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateK8s))

			c.Verbose("Property CreateK8sCluster set: %v", createK8s)
		} else {
			if s, ok := properties.GetCreateK8sClusterOk(); ok && s != nil {
				createK8s = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateNic)) {
			createNic = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateNic))

			c.Verbose("Property CreateInternetAccess set: %v", createNic)
		} else {
			if s, ok := properties.GetCreateInternetAccessOk(); ok && s != nil {
				createNic = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateBackUpUnit)) {
			createBackUp = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateBackUpUnit))

			c.Verbose("Property CreateBackupUnit set: %v", createBackUp)
		} else {
			if s, ok := properties.GetCreateBackupUnitOk(); ok && s != nil {
				createBackUp = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgReserveIp)) {
			reserveIp = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgReserveIp))

			c.Verbose("Property ReserveIp set: %v", reserveIp)
		} else {
			if ip, ok := properties.GetReserveIpOk(); ok && ip != nil {
				reserveIp = *ip
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAccessLog)) {
			accessLog = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessLog))

			c.Verbose("Property AccessActivityLog set: %v", accessLog)
		} else {
			if log, ok := properties.GetAccessActivityLogOk(); ok && log != nil {
				accessLog = *log
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgS3Privilege)) {
			s3 = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgS3Privilege))

			c.Verbose("Property S3Privilege set: %v", s3)
		} else {
			if s, ok := properties.GetS3PrivilegeOk(); ok && s != nil {
				s3 = *s
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgCreateFlowLog)) {
			createFlowLog = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgCreateFlowLog))

			c.Verbose("Property CreateFlowLog set: %v", createFlowLog)
		} else {
			if f, ok := properties.GetCreateFlowLogOk(); ok && f != nil {
				createFlowLog = *f
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAccessMonitoring)) {
			monitoring = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessMonitoring))

			c.Verbose("Property AccessAndManageMonitoring set: %v", monitoring)
		} else {
			if m, ok := properties.GetAccessAndManageMonitoringOk(); ok && m != nil {
				monitoring = *m
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAccessCerts)) {
			certs = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessCerts))

			c.Verbose("Property AccessAndManageCertificates set: %v", certs)
		} else {
			if accessCerts, ok := properties.GetAccessAndManageCertificatesOk(); ok && accessCerts != nil {
				certs = *accessCerts
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgAccessDNS)) {
			dns = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgAccessDNS))

			c.Verbose("Property AccessAndManageDNS set: %v", dns)
		} else {
			if accessDns, ok := properties.GetAccessAndManageDnsOk(); ok && accessDns != nil {
				dns = *accessDns
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgManageDbaas)) {
			manageDb = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgManageDbaas))

			c.Verbose("Property ManageDBaaS set: %v", manageDb)
		} else {
			if manageDBaaSOk, ok := properties.GetManageDBaaSOk(); ok && manageDBaaSOk != nil {
				manageDb = *manageDBaaSOk
			}
		}

		if viper.IsSet(core.GetFlagName(c.NS, cloudapiv6.ArgManageRegistry)) {
			manageReg = viper.GetBool(core.GetFlagName(c.NS, cloudapiv6.ArgManageRegistry))

			c.Verbose("Property ManageRegistry set: %v", manageReg)
		} else {
			if manageRegistry, ok := properties.GetManageRegistryOk(); ok && manageRegistry != nil {
				manageReg = *manageRegistry
			}
		}
	}

	return &resources.GroupProperties{
		GroupProperties: ionoscloud.GroupProperties{
			Name:                        &groupName,
			CreateDataCenter:            &createDc,
			CreateSnapshot:              &createSnap,
			ReserveIp:                   &reserveIp,
			AccessActivityLog:           &accessLog,
			CreatePcc:                   &createPcc,
			S3Privilege:                 &s3,
			CreateBackupUnit:            &createBackUp,
			CreateInternetAccess:        &createNic,
			CreateK8sCluster:            &createK8s,
			CreateFlowLog:               &createFlowLog,
			AccessAndManageMonitoring:   &monitoring,
			AccessAndManageCertificates: &certs,
			AccessAndManageDns:          &dns,
			ManageDBaaS:                 &manageDb,
			ManageRegistry:              &manageReg,
		},
	}
}

func DeleteAllGroups(c *core.CommandConfig) error {
	return core.DeleteAll(c, core.DeleteAllOptions[ionoscloud.Group]{
		Resource: "Group",
		List: func() ([]ionoscloud.Group, error) {
			groups, _, err := c.CloudApiV6Services.Groups().List()
			if err != nil {
				return nil, err
			}

			items, ok := groups.GetItemsOk()
			if !ok || items == nil {
				return nil, fmt.Errorf("could not get items of Groups")
			}

			return *items, nil
		},
		Summary: func(group ionoscloud.Group) string {
			summary := fmt.Sprintf("id: %s", *group.GetId())
			if props, ok := group.GetPropertiesOk(); ok && props != nil {
				if name, ok := props.GetNameOk(); ok && name != nil && *name != "" {
					summary = fmt.Sprintf("%s (name: %s)", summary, *name)
				}
			}
			return summary
		},
		ID: func(group ionoscloud.Group) string {
			return *group.GetId()
		},
		Delete: func(group ionoscloud.Group) error {
			resp, err := c.CloudApiV6Services.Groups().Delete(*group.GetId())
			if resp != nil && request.GetId(resp) != "" {
				c.Verbose(constants.MessageRequestInfo, request.GetId(resp), resp.RequestTime)
			}
			return err
		},
	})
}
