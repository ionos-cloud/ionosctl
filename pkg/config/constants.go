package config

// Resources
const (
	DatacenterResource = "datacenter"
	ServerResource     = "server"
	VolumeResource     = "volume"
	IpBlockResource    = "ipblock"
	SnapshotResource   = "snapshot"
)

// Flags
const (
	ArgConfig                      = "config"
	ArgOutput                      = "output"
	ArgQuiet                       = "quiet"
	ArgWaitForRequest              = "wait-for-request"
	ArgWaitForState                = "wait-for-state"
	ArgTimeout                     = "timeout"
	ArgForce                       = "force"
	ArgServerUrl                   = "api-url"
	ArgCols                        = "cols"
	ArgDataCenterId                = "datacenter-id"
	ArgDataCenterName              = "datacenter-name"
	ArgDataCenterDescription       = "datacenter-description"
	ArgDataCenterRegion            = "datacenter-location"
	ArgServerId                    = "server-id"
	ArgServerName                  = "server-name"
	ArgServerZone                  = "server-zone"
	ArgServerCores                 = "server-cores"
	ArgServerRAM                   = "server-ram"
	ArgServerCPUFamily             = "server-cpu-family"
	ArgVolumeId                    = "volume-id"
	ArgVolumeSize                  = "volume-size"
	ArgVolumeBus                   = "volume-bus"
	ArgVolumeLicenceType           = "volume-licence-type"
	ArgVolumeType                  = "volume-type"
	ArgVolumeName                  = "volume-name"
	ArgVolumeZone                  = "volume-zone"
	ArgVolumeSshKey                = "volume-ssh-keys"
	ArgLanId                       = "lan-id"
	ArgLanName                     = "lan-name"
	ArgLanPublic                   = "lan-public"
	ArgNicId                       = "nic-id"
	ArgNicName                     = "nic-name"
	ArgNicIps                      = "nic-ips"
	ArgNicDhcp                     = "nic-dhcp"
	ArgLoadBalancerId              = "loadbalancer-id"
	ArgLoadBalancerName            = "loadbalancer-name"
	ArgLoadBalancerIp              = "loadbalancer-ip"
	ArgLoadBalancerDhcp            = "loadbalancer-dhcp"
	ArgRequestId                   = "request-id"
	ArgSnapshotName                = "snapshot-name"
	ArgSnapshotDescription         = "snapshot-description"
	ArgSnapshotLicenceType         = "snapshot-licence-type"
	ArgSnapshotCpuHotPlug          = "snapshot-cpu-hot-plug"
	ArgSnapshotCpuHotUnplug        = "snapshot-cpu-hot-unplug"
	ArgSnapshotRamHotPlug          = "snapshot-ram-hot-plug"
	ArgSnapshotRamHotUnplug        = "snapshot-ram-hot-unplug"
	ArgSnapshotNicHotPlug          = "snapshot-nic-hot-plug"
	ArgSnapshotNicHotUnplug        = "snapshot-nic-hot-unplug"
	ArgSnapshotDiscVirtioHotPlug   = "snapshot-disc-virtio-hot-plug"
	ArgSnapshotDiscVirtioHotUnplug = "snapshot-disc-virtio-hot-unplug"
	ArgSnapshotDiscScsiHotPlug     = "snapshot-disc-scsi-hot-plug"
	ArgSnapshotDiscScsiHotUnplug   = "snapshot-disc-scsi-hot-unplug"
	ArgSnapshotSecAuthProtection   = "snapshot-sec-auth-protection"
	ArgSnapshotId                  = "snapshot-id"
	ArgImageId                     = "image-id"
	ArgImageLocation               = "image-location"
	ArgImageLicenceType            = "image-licence-type"
	ArgImageType                   = "image-type"
	ArgImageSize                   = "image-size"
	ArgIpBlockId                   = "ipblock-id"
	ArgIpBlockName                 = "ipblock-name"
	ArgIpBlockLocation             = "ipblock-location"
	ArgIpBlockSize                 = "ipblock-size"
	ArgFirewallRuleId              = "firewallrule-id"
	ArgFirewallRuleName            = "firewallrule-name"
	ArgFirewallRuleProtocol        = "firewallrule-protocol"
	ArgFirewallRuleSourceMac       = "firewallrule-source-mac"
	ArgFirewallRuleSourceIp        = "firewallrule-source-ip"
	ArgFirewallRuleTargetIp        = "firewallrule-target-ip"
	ArgFirewallRuleIcmpCode        = "firewallrule-icmp-code"
	ArgFirewallRuleIcmpType        = "firewallrule-icmp-type"
	ArgFirewallRulePortRangeStart  = "firewallrule-port-range-start"
	ArgFirewallRulePortRangeStop   = "firewallrule-port-range-end"
	ArgLabelUrn                    = "label-urn"
	ArgLabelKey                    = "label-key"
	ArgLabelValue                  = "label-value"
	ArgResourceLimits              = "resource-limits"
	ArgUserId                      = "user-id"
	ArgUserFirstName               = "user-first-name"
	ArgUserLastName                = "user-last-name"
	ArgUserEmail                   = "user-email"
	ArgUserPassword                = "user-password"
	ArgUserAdministrator           = "user-administrator"
	ArgUserForceSecAuth            = "user-force-secure"
	ArgGroupId                     = "group-id"
	ArgGroupName                   = "group-name"
	ArgGroupCreateDc               = "group-create-dc"
	ArgGroupCreateSnapshot         = "group-create-snapshot"
	ArgGroupReserveIp              = "group-reserve-ip"
	ArgGroupAccessLog              = "group-access-logs"
	ArgGroupS3Privilege            = "group-s3privilege"
	ArgGroupCreateBackUpUnit       = "group-create-backup"
	ArgGroupCreatePcc              = "group-create-pcc"
	ArgGroupCreateNic              = "group-create-nic"
	ArgGroupCreateK8s              = "group-create-k8s"
	ArgResourceId                  = "resource-id"
	ArgResourceType                = "resource-type"
	ArgEditPrivilege               = "edit-privilege"
	ArgSharePrivilege              = "share-privilege"
	ArgS3KeyId                     = "s3key-id"
	ArgS3KeyActive                 = "s3key-active"
	ArgBackupUnitId                = "backupunit-id"
	ArgBackupUnitName              = "backupunit-name"
	ArgBackupUnitPassword          = "backupunit-password"
	ArgBackupUnitEmail             = "backupunit-email"
	ArgPccId                       = "pcc-id"
	ArgPccName                     = "pcc-name"
	ArgPccDescription              = "pcc-description"
	ArgK8sClusterId                = "cluster-id"
	ArgK8sClusterName              = "cluster-name"
	ArgK8sClusterVersion           = "cluster-version"
	ArgK8sNodePoolId               = "nodepool-id"
	ArgK8sNodePoolName             = "nodepool-name"
	ArgK8sNodePoolVersion          = "nodepool-version"
	ArgK8sNodeCount                = "node-count"
	ArgCoresCount                  = "cores-count"
	ArgCpuFamily                   = "cpu-family"
	ArgRamSize                     = "ram-size"
	ArgK8sNodeZone                 = "node-zone"
	ArgStorageType                 = "storage-type"
	ArgStorageSize                 = "storage-size"
	ArgK8sMinNodeCount             = "min-node-count"
	ArgK8sMaxNodeCount             = "max-node-count"
	ArgK8sNodeId                   = "node-id"
	ArgK8sMaintenanceDay           = "maintenance-day"
	ArgK8sMaintenanceTime          = "maintenance-time"
	ArgK8sAnnotationKey            = "annotation-key"
	ArgK8sAnnotationValue          = "annotation-value"
)

// Default values
const (
	DefaultApiURL           = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName   = "/config.json"
	DefaultOutputFormat     = "text"
	DefaultWait             = false
	DefaultTimeoutSeconds   = int(60)
	DefaultServerCores      = 2
	DefaultServerRAM        = 256
	DefaultServerCPUFamily  = "AMD_OPTERON"
	DefaultVolumeSize       = 10
	DefaultLanPublic        = false
	DefaultNicDhcp          = true
	DefaultNicLanId         = 1
	DefaultLoadBalancerDhcp = true
	Username                = "userdata.name"
	Password                = "userdata.password"
	Token                   = "userdata.token"
)

// Required Flags
const (
	RequiredFlag               = "(required)"
	RequiredFlagDatacenterId   = "The unique Data Center Id " + RequiredFlag
	RequiredFlagLanId          = "The unique LAN Id " + RequiredFlag
	RequiredFlagLoadBalancerId = "The unique Load Balancer Id " + RequiredFlag
	RequiredFlagNicId          = "The unique NIC Id " + RequiredFlag
	RequiredFlagRequestId      = "The unique Request Id " + RequiredFlag
	RequiredFlagServerId       = "The unique Server Id " + RequiredFlag
	RequiredFlagVolumeId       = "The unique Volume Id " + RequiredFlag
	RequiredFlagSnapshotId     = "The unique Snapshot Id " + RequiredFlag
	RequiredFlagImageId        = "The unique Image Id " + RequiredFlag
	RequiredFlagIpBlockId      = "The unique IpBlock Id " + RequiredFlag
	RequiredFlagFirewallRuleId = "The unique FirewallRule Id " + RequiredFlag
	RequiredFlagLabelKey       = "The unique Label Key " + RequiredFlag
	RequiredFlagLabelValue     = "The unique Label Value " + RequiredFlag
	RequiredFlagUserId         = "The unique User Id " + RequiredFlag
	RequiredFlagGroupId        = "The unique Group Id " + RequiredFlag
	RequiredFlagResourceId     = "The unique Resource Id " + RequiredFlag
	RequiredFlagS3KeyId        = "The unique User S3Key Id " + RequiredFlag
	RequiredFlagBackupUnitId   = "The unique BackupUnit Id " + RequiredFlag
	RequiredFlagPccId          = "The unique Private Cross-Connect Id " + RequiredFlag
	RequiredFlagK8sClusterId   = "The unique K8s Cluster Id " + RequiredFlag
	RequiredFlagK8sNodePoolId  = "The unique K8s Node Pool Id " + RequiredFlag
	RequiredFlagK8sNodeId      = "The unique K8s Node Id " + RequiredFlag
)
