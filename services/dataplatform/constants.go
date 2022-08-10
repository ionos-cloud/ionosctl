package dataplatform

const (
	ArgIdShort               = "i"
	ArgClusterId             = "cluster-id"
	ArgNodePoolId            = "nodepool-id"
	ArgCpuFamily             = "cpu-family"
	ArgAvailabilityZone      = "availability-zone"
	ArgAvailabilityZoneShort = "z"
	ArgLabels                = "labels"
	ArgLabelsShort           = "L"
	ArgAnnotations           = "annotations"
	ArgAnnotationsShort      = "A"
	ArgVersion               = "version"
	ArgVersionShort          = "V"
	ArgCores                 = "cores"
	ArgRam                   = "ram"
	ArgStorageSize           = "storage-size"
	ArgStorageType           = "storage-type"
	ArgDatacenterId          = "datacenter-id"
	ArgDatacenterIdShort     = "D"
	ArgName                  = "name"
	ArgNameShort             = "n"
	ArgNodeCount             = "node-count"
	ArgMaintenanceTime       = "maintenance-time"
	ArgMaintenanceTimeShort  = "T"
	ArgMaintenanceDay        = "maintenance-day"
	ArgMaintenanceDayShort   = "d"
)

const (
	ClusterId             = "The unique ID of the Cluster"
	NodePoolId            = "The unique ID of the Node Pool"
	BackupId              = "The unique ID of the Backup"
	DefaultClusterTimeout = int(1200)
)

// Default values
const (
	DefaultOutputFormat    = "text"
	DefaultWait            = false
	DefaultPublic          = false
	DefaultDhcp            = true
	DefaultFirewallActive  = false
	DefaultTimeoutSeconds  = int(60)
	NlbTimeoutSeconds      = int(300)
	LbTimeoutSeconds       = int(300)
	AlbTimeoutSeconds      = int(10000)
	TimeoutSeconds         = int(600)
	DefaultServerCores     = 2
	DefaultVolumeSize      = 10
	DefaultNicLanId        = 1
	DefaultServerCPUFamily = "AUTO"
)
