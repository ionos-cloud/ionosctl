package config

// Flags
const (
	ArgConfig                      = "config"
	ArgOutput                      = "output"
	ArgQuiet                       = "quiet"
	ArgWait                        = "wait"
	ArgTimeout                     = "timeout"
	ArgIgnoreStdin                 = "ignore-stdin"
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
	ArgSnapshotSize                = "snapshot-size"
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
)

// Default values
const (
	DefaultApiURL           = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName   = "/config.json"
	DefaultOutputFormat     = "text"
	DefaultWait             = false
	DefaultTimeoutSeconds   = 60
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
	RequiredFlag               = "[Required flag]"
	RequiredFlagDatacenterId   = "The unique Data Center Id " + RequiredFlag
	RequiredFlagLanId          = "The unique LAN Id " + RequiredFlag
	RequiredFlagLoadBalancerId = "The unique Load Balancer Id " + RequiredFlag
	RequiredFlagNicId          = "The unique NIC Id " + RequiredFlag
	RequiredFlagRequestId      = "The unique Request Id " + RequiredFlag
	RequiredFlagServerId       = "The unique Server Id " + RequiredFlag
	RequiredFlagVolumeId       = "The unique Volume Id " + RequiredFlag
	RequiredFlagSnapshotId     = "The unique Snapshot Id " + RequiredFlag
	RequiredFlagImageId        = "The unique Image Id " + RequiredFlag
	RequiredFlagIpBlockId      = "The unique IPBlock Id " + RequiredFlag
)
