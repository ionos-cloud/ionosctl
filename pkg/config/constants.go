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
	ArgConfig                     = "config"
	ArgOutput                     = "output"
	ArgQuiet                      = "quiet"
	ArgWaitForRequest             = "wait-for-request"
	ArgWaitForState               = "wait-for-state"
	ArgTimeout                    = "timeout"
	ArgForce                      = "force"
	ArgServerUrl                  = "api-url"
	ArgCols                       = "cols"
	ArgDataCenterId               = "datacenter-id"
	ArgDataCenterName             = "datacenter-name"
	ArgDataCenterDescription      = "datacenter-description"
	ArgDataCenterRegion           = "datacenter-location"
	ArgServerId                   = "server-id"
	ArgServerName                 = "server-name"
	ArgServerZone                 = "server-zone"
	ArgServerCores                = "server-cores"
	ArgServerRAM                  = "server-ram"
	ArgServerCPUFamily            = "server-cpu-family"
	ArgVolumeId                   = "volume-id"
	ArgVolumeSize                 = "volume-size"
	ArgVolumeBus                  = "volume-bus"
	ArgLicenceType                = "licence-type"
	ArgVolumeType                 = "volume-type"
	ArgVolumeName                 = "volume-name"
	ArgVolumeZone                 = "volume-zone"
	ArgSshKeys                    = "ssh-keys"
	ArgLocationId                 = "location-id"
	ArgLanId                      = "lan-id"
	ArgLanName                    = "lan-name"
	ArgLanPublic                  = "lan-public"
	ArgNicId                      = "nic-id"
	ArgNicName                    = "nic-name"
	ArgNicIps                     = "nic-ips"
	ArgNicDhcp                    = "nic-dhcp"
	ArgLoadBalancerId             = "loadbalancer-id"
	ArgLoadBalancerName           = "loadbalancer-name"
	ArgLoadBalancerIp             = "loadbalancer-ip"
	ArgLoadBalancerDhcp           = "loadbalancer-dhcp"
	ArgRequestId                  = "request-id"
	ArgSnapshotName               = "snapshot-name"
	ArgSnapshotDescription        = "snapshot-description"
	ArgCpuHotPlug                 = "cpu-hot-plug"
	ArgCpuHotUnplug               = "cpu-hot-unplug"
	ArgRamHotPlug                 = "ram-hot-plug"
	ArgRamHotUnplug               = "ram-hot-unplug"
	ArgNicHotPlug                 = "nic-hot-plug"
	ArgNicHotUnplug               = "nic-hot-unplug"
	ArgDiscVirtioHotPlug          = "disc-virtio-hot-plug"
	ArgDiscVirtioHotUnplug        = "disc-virtio-hot-unplug"
	ArgDiscScsiHotPlug            = "disc-scsi-hot-plug"
	ArgDiscScsiHotUnplug          = "disc-scsi-hot-unplug"
	ArgSecAuthProtection          = "sec-auth-protection"
	ArgSnapshotId                 = "snapshot-id"
	ArgImageId                    = "image-id"
	ArgImageAlias                 = "image-alias"
	ArgImagePassword              = "image-password"
	ArgImageLocation              = "image-location"
	ArgImageLicenceType           = "image-licence-type"
	ArgImageType                  = "image-type"
	ArgImageSize                  = "image-size"
	ArgIpBlockId                  = "ipblock-id"
	ArgIpBlockName                = "ipblock-name"
	ArgIpBlockLocation            = "ipblock-location"
	ArgIpBlockSize                = "ipblock-size"
	ArgFirewallRuleId             = "firewallrule-id"
	ArgFirewallRuleName           = "firewallrule-name"
	ArgFirewallRuleProtocol       = "firewallrule-protocol"
	ArgFirewallRuleSourceMac      = "firewallrule-source-mac"
	ArgFirewallRuleSourceIp       = "firewallrule-source-ip"
	ArgFirewallRuleTargetIp       = "firewallrule-target-ip"
	ArgFirewallRuleIcmpCode       = "firewallrule-icmp-code"
	ArgFirewallRuleIcmpType       = "firewallrule-icmp-type"
	ArgFirewallRulePortRangeStart = "firewallrule-port-range-start"
	ArgFirewallRulePortRangeStop  = "firewallrule-port-range-end"
	ArgLabelUrn                   = "label-urn"
	ArgLabelKey                   = "label-key"
	ArgLabelValue                 = "label-value"
	ArgResourceLimits             = "resource-limits"
	ArgUserId                     = "user-id"
	ArgUserData                   = "user-data"
	ArgUserFirstName              = "user-first-name"
	ArgUserLastName               = "user-last-name"
	ArgUserEmail                  = "user-email"
	ArgUserPassword               = "user-password"
	ArgUserAdministrator          = "user-administrator"
	ArgUserForceSecAuth           = "user-force-secure"
	ArgGroupId                    = "group-id"
	ArgGroupName                  = "group-name"
	ArgGroupCreateDc              = "group-create-dc"
	ArgGroupCreateSnapshot        = "group-create-snapshot"
	ArgGroupReserveIp             = "group-reserve-ip"
	ArgGroupAccessLog             = "group-access-logs"
	ArgGroupS3Privilege           = "group-s3privilege"
	ArgGroupCreateBackUpUnit      = "group-create-backup"
	ArgGroupCreatePcc             = "group-create-pcc"
	ArgGroupCreateNic             = "group-create-nic"
	ArgGroupCreateK8s             = "group-create-k8s"
	ArgResourceId                 = "resource-id"
	ArgResourceType               = "resource-type"
	ArgEditPrivilege              = "edit-privilege"
	ArgSharePrivilege             = "share-privilege"
	ArgS3KeyId                    = "s3key-id"
	ArgS3KeyActive                = "s3key-active"
	ArgBackupUnitId               = "backupunit-id"
	ArgBackupUnitName             = "backupunit-name"
	ArgBackupUnitPassword         = "backupunit-password"
	ArgBackupUnitEmail            = "backupunit-email"
	ArgPccId                      = "pcc-id"
	ArgPccName                    = "pcc-name"
	ArgPccDescription             = "pcc-description"
	ArgK8sClusterId               = "cluster-id"
	ArgK8sClusterName             = "cluster-name"
	ArgK8sVersion                 = "k8s-version"
	ArgK8sNodePoolId              = "nodepool-id"
	ArgK8sNodePoolName            = "nodepool-name"
	ArgK8sNodeCount               = "node-count"
	ArgCoresCount                 = "cores-count"
	ArgCpuFamily                  = "cpu-family"
	ArgRamSize                    = "ram-size"
	ArgK8sNodeZone                = "node-zone"
	ArgStorageType                = "storage-type"
	ArgStorageSize                = "storage-size"
	ArgK8sMinNodeCount            = "min-node-count"
	ArgK8sMaxNodeCount            = "max-node-count"
	ArgK8sNodeId                  = "node-id"
	ArgK8sMaintenanceDay          = "maintenance-day"
	ArgK8sMaintenanceTime         = "maintenance-time"
	ArgK8sAnnotationKey           = "annotation-key"
	ArgK8sAnnotationValue         = "annotation-value"
	ArgPublicIps                  = "public-ips"
	ArgPublic                     = "public"
	ArgGatewayIp                  = "gateway-ip"
	ArgCdromId                    = "cdrom-id"
)

// Default values
const (
	DefaultApiURL           = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName   = "/config.json"
	DefaultOutputFormat     = "text"
	DefaultWait             = false
	DefaultTimeoutSeconds   = int(60)
	K8sTimeoutSeconds       = int(600)
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
	RequiredFlagLocationId     = "The unique Location Id " + RequiredFlag
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
	RequiredFlagCdromId        = "The unique Cdrom Id " + RequiredFlag
)
