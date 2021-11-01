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
	ArgPostgresVersion      = "pgsql-version"
	ArgPostgresVersionShort = "p"
	ArgReplicas             = "replicas"
	ArgReplicasShort        = "R"
	ArgCpuCoreCount         = "cpu-core-count"
	ArgRamSize              = "ram-size"
	ArgStorageSize          = "storage-size"
	ArgStorageType          = "storage-type"
	ArgVdcId                = "vdc-id"
	ArgVdcIdShort           = "V"
	ArgBackupId             = "backup-id"
	ArgBackupIdShort        = "b"
	ArgTime                 = "time"
	ArgIpAddress            = "ip"
	ArgLanId                = "lan-id"
	ArgLocation             = "location-id"
	ArgName                 = "name"
	ArgNameShort            = "n"
	ArgUsername             = "username"
	ArgPassword             = "password"
	ArgMaintenanceTime      = "maintenance-time"
	ArgMaintenanceTimeShort = "T"
	ArgMaintenanceDay       = "maintenance-day"
	ArgMaintenanceDayShort  = "d"
)

const (
	ClusterId             = "The unique ID of the Cluster"
	BackupId              = ""
	DefaultClusterTimeout = int(600)
)
