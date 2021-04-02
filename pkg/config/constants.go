package config

// Flags
const (
	ArgConfig                = "config"
	ArgOutput                = "output"
	ArgQuiet                 = "quiet"
	ArgWait                  = "wait"
	ArgTimeout               = "timeout"
	ArgIgnoreStdin           = "ignore-stdin"
	ArgServerUrl             = "api-url"
	ArgCols                  = "cols"
	ArgDataCenterId          = "datacenter-id"
	ArgDataCenterName        = "datacenter-name"
	ArgDataCenterDescription = "datacenter-description"
	ArgDataCenterRegion      = "datacenter-location"
	ArgServerId              = "server-id"
	ArgServerName            = "server-name"
	ArgServerZone            = "server-zone"
	ArgServerCores           = "server-cores"
	ArgServerRAM             = "server-ram"
	ArgServerCPUFamily       = "server-cpu-family"
	ArgVolumeId              = "volume-id"
	ArgVolumeSize            = "volume-size"
	ArgVolumeBus             = "volume-bus"
	ArgVolumeLicenceType     = "volume-licence-type"
	ArgVolumeType            = "volume-type"
	ArgVolumeName            = "volume-name"
	ArgVolumeZone            = "volume-zone"
	ArgVolumeSshKey          = "volume-ssh-keys"
	ArgLanId                 = "lan-id"
	ArgLanName               = "lan-name"
	ArgLanPublic             = "lan-public"
	ArgNicId                 = "nic-id"
	ArgNicName               = "nic-name"
	ArgNicIps                = "nic-ips"
	ArgNicDhcp               = "nic-dhcp"
	ArgLoadBalancerId        = "loadbalancer-id"
	ArgLoadBalancerName      = "loadbalancer-name"
	ArgLoadBalancerIp        = "loadbalancer-ip"
	ArgLoadBalancerDhcp      = "loadbalancer-dhcp"
	ArgRequestId             = "request-id"
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
)
