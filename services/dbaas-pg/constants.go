package dbaas_pg

const (
	ArgIdShort              = "i"
	ArgClusterId            = "cluster-id"
	ArgStartTime            = "start-time"
	ArgStartTimeShort       = "s"
	ArgEndTime              = "end-time"
	ArgEndTimeShort         = "e"
	ArgLimit                = "limit"
	ArgLimitShort           = "l"
	ArgVersion              = "version"
	ArgVersionShort         = "V"
	ArgInstances            = "instances"
	ArgInstancesShort       = "I"
	ArgSyncMode             = "sync"
	ArgSyncModeShort        = "S"
	ArgCores                = "cores"
	ArgRam                  = "ram"
	ArgStorageSize          = "storage-size"
	ArgStorageType          = "storage-type"
	ArgDatacenterId         = "datacenter-id"
	ArgDatacenterIdShort    = "D"
	ArgBackupId             = "backup-id"
	ArgBackupIdShort        = "b"
	ArgTime                 = "time"
	ArgIpAddress            = "ip"
	ArgLanId                = "lan-id"
	ArgLocation             = "location-id"
	ArgName                 = "name"
	ArgNameShort            = "n"
	ArgUsername             = "username"
	ArgUsernameShort        = "U"
	ArgPassword             = "password"
	ArgPasswordShort        = "P"
	ArgMaintenanceTime      = "maintenance-time"
	ArgMaintenanceTimeShort = "T"
	ArgMaintenanceDay       = "maintenance-day"
	ArgMaintenanceDayShort  = "d"
	ArgSaveToFile           = "save-to-file"
	ArgSaveToFileShort      = "S"
)

const (
	ClusterId             = "The unique ID of the Cluster"
	BackupId              = "The unique ID of the Backup"
	DefaultClusterTimeout = int(1200)
)
