package config

const (
	// Global Flags for Root Command
	ArgConfig      = "config"
	ArgOutput      = "output"
	ArgQuiet       = "quiet"
	ArgWait        = "wait"
	ArgTimeout     = "timeout"
	ArgIgnoreStdin = "ignore-stdin"
	ArgServerUrl   = "api-url"
	ArgCols        = "cols"
	// Data Center Flags
	ArgDataCenterId          = "datacenter-id"
	ArgDataCenterName        = "datacenter-name"
	ArgDataCenterDescription = "datacenter-description"
	ArgDataCenterRegion      = "datacenter-location"
	// Server Flags
	ArgServerId        = "server-id"
	ArgServerName      = "server-name"
	ArgServerZone      = "server-zone"
	ArgServerCores     = "server-cores"
	ArgServerRAM       = "server-ram"
	ArgServerCPUFamily = "server-cpufamily"
	// Volume Flags
	ArgVolumeId          = "volume-id"
	ArgVolumeSize        = "volume-size"
	ArgVolumeBus         = "volume-bus"
	ArgVolumeLicenceType = "volume-licencetype"
	ArgVolumeType        = "volume-type"
	ArgVolumeName        = "volume-name"
	ArgVolumeZone        = "volume-zone"
	ArgVolumeSshKey      = "volume-sshkeys"
	// Lan Flags
	ArgLanId     = "lan-id"
	ArgLanName   = "lan-name"
	ArgLanPublic = "lan-public"
	// Nic Flags
	ArgNicId   = "nic-id"
	ArgNicName = "nic-name"
	ArgNicIps  = "nic-ips"
	ArgNicDhcp = "nic-dhcp"
	// Load Balancer Flags
	ArgLoadbalancerId   = "loadbalancer-id"
	ArgLoadbalancerName = "loadbalancer-name"
	ArgLoadbalancerIp   = "loadbalancer-ip"
	ArgLoadbalancerDhcp = "loadbalancer-dhcp"
	// Request Flags
	ArgRequestId = "request-id"
	// Snapshot Flags
	ArgSnapshotName        = "snapshot-name"
	ArgSnapshotDescription = "snapshot-description"
	ArgSnapshotLicenceType = "snapshot-licencetype"
	ArgSnapshotId          = "snapshot-id"

	// Default values for Global Flags
	DefaultApiURL         = "https://api.ionos.com/cloudapi/v5"
	DefaultConfigFileName = "/config.json"
	DefaultOutputFormat   = "text"
	DefaultWait           = false
	DefaultTimeoutSeconds = 60
	// Default values for Server Flags
	DefaultServerCores     = 2
	DefaultServerRAM       = 256
	DefaultServerCPUFamily = "AMD_OPTERON"
	// Default values for Volume Flags
	DefaultVolumeSize = 10
	// Default values for Lan Flags
	DefaultLanPublic = false
	// Default values for Nic Flags
	DefaultNicDhcp  = true
	DefaultNicLanId = 1
	// Default values for Load Balancer Flags
	DefaultLoadBalancerDhcp = true

	Username = "userdata.name"
	Password = "userdata.password"
)
