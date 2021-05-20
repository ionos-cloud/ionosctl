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
	ArgConfig              = "config"
	ArgOutput              = "output"
	ArgQuiet               = "quiet"
	ArgWaitForRequest      = "wait-for-request"
	ArgWaitForState        = "wait-for-state"
	ArgTimeout             = "timeout"
	ArgForce               = "force"
	ArgServerUrl           = "api-url"
	ArgCols                = "cols"
	ArgDataCenterId        = "datacenter-id"
	ArgName                = "name"
	ArgDescription         = "description"
	ArgLocation            = "location"
	ArgServerId            = "server-id"
	ArgAvailabilityZone    = "availability-zone"
	ArgCores               = "cores"
	ArgRamSize             = "ram-size"
	ArgCPUFamily           = "cpu-family"
	ArgVolumeId            = "volume-id"
	ArgSize                = "size"
	ArgBus                 = "bus"
	ArgLicenceType         = "licence-type"
	ArgType                = "type"
	ArgSshKeys             = "ssh-keys"
	ArgLocationId          = "location-id"
	ArgLanId               = "lan-id"
	ArgPublic              = "public"
	ArgNicId               = "nic-id"
	ArgIps                 = "ips"
	ArgDhcp                = "dhcp"
	ArgLoadBalancerId      = "loadbalancer-id"
	ArgIp                  = "ip"
	ArgRequestId           = "request-id"
	ArgCpuHotPlug          = "cpu-hot-plug"
	ArgCpuHotUnplug        = "cpu-hot-unplug"
	ArgRamHotPlug          = "ram-hot-plug"
	ArgRamHotUnplug        = "ram-hot-unplug"
	ArgNicHotPlug          = "nic-hot-plug"
	ArgNicHotUnplug        = "nic-hot-unplug"
	ArgDiscVirtioHotPlug   = "disc-virtio-hot-plug"
	ArgDiscVirtioHotUnplug = "disc-virtio-hot-unplug"
	ArgDiscScsiHotPlug     = "disc-scsi-hot-plug"
	ArgDiscScsiHotUnplug   = "disc-scsi-hot-unplug"
	ArgSecAuthProtection   = "sec-auth-protection"
	ArgSnapshotId          = "snapshot-id"
	ArgImageId             = "image-id"
	ArgImageAlias          = "image-alias"
	ArgPassword            = "password"
	ArgIpBlockId           = "ipblock-id"
	ArgFirewallRuleId      = "firewallrule-id"
	ArgProtocol            = "protocol"
	ArgSourceMac           = "source-mac"
	ArgSourceIp            = "source-ip"
	ArgTargetIp            = "target-ip"
	ArgIcmpCode            = "icmp-code"
	ArgIcmpType            = "icmp-type"
	ArgPortRangeStart      = "port-range-start"
	ArgPortRangeStop       = "port-range-end"
	ArgLabelUrn            = "label-urn"
	ArgLabelKey            = "label-key"
	ArgLabelValue          = "label-value"
	ArgResourceLimits      = "resource-limits"
	ArgUserId              = "user-id"
	ArgUserData            = "user-data"
	ArgFirstName           = "first-name"
	ArgLastName            = "last-name"
	ArgEmail               = "email"
	ArgAdministrator       = "administrator"
	ArgForceSecAuth        = "force-secure-auth"
	ArgGroupId             = "group-id"
	ArgCreateDc            = "create-dc"
	ArgCreateSnapshot      = "create-snapshot"
	ArgReserveIp           = "reserve-ip"
	ArgAccessLog           = "access-logs"
	ArgS3Privilege         = "s3privilege"
	ArgCreateBackUpUnit    = "create-backup"
	ArgCreatePcc           = "create-pcc"
	ArgCreateNic           = "create-nic"
	ArgCreateK8s           = "create-k8s"
	ArgResourceId          = "resource-id"
	ArgEditPrivilege       = "edit-privilege"
	ArgSharePrivilege      = "share-privilege"
	ArgS3KeyId             = "s3key-id"
	ArgS3KeyActive         = "s3key-active"
	ArgBackupUnitId        = "backupunit-id"
	ArgPccId               = "pcc-id"
	ArgK8sClusterId        = "cluster-id"
	ArgK8sVersion          = "k8s-version"
	ArgK8sNodePoolId       = "nodepool-id"
	ArgK8sNodeCount        = "node-count"
	ArgCoresCount          = "cores-count"
	ArgCpuFamily           = "cpu-family"
	ArgStorageType         = "storage-type"
	ArgStorageSize         = "storage-size"
	ArgK8sMinNodeCount     = "min-node-count"
	ArgK8sMaxNodeCount     = "max-node-count"
	ArgK8sNodeId           = "node-id"
	ArgK8sMaintenanceDay   = "maintenance-day"
	ArgK8sMaintenanceTime  = "maintenance-time"
	ArgK8sAnnotationKey    = "annotation-key"
	ArgK8sAnnotationValue  = "annotation-value"
	ArgPublicIps           = "public-ips"
	ArgGatewayIp           = "gateway-ip"
	ArgCdromId             = "cdrom-id"
)

// Default values
const (
	DefaultApiURL          = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName  = "/config.json"
	DefaultOutputFormat    = "text"
	DefaultWait            = false
	DefaultPublic          = false
	DefaultDhcp            = true
	DefaultTimeoutSeconds  = int(60)
	K8sTimeoutSeconds      = int(600)
	DefaultServerCores     = 2
	DefaultServerRAM       = 256
	DefaultVolumeSize      = 10
	DefaultNicLanId        = 1
	DefaultServerCPUFamily = "AMD_OPTERON"
	Username               = "userdata.name"
	Password               = "userdata.password"
	Token                  = "userdata.token"
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
