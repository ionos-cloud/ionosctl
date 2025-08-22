package dbaas_postgres

const (
	ArgStorageSize       = "storage-size"
	ArgStorageType       = "storage-type"
	ArgDatacenterId      = "datacenter-id"
	ArgDatacenterIdShort = "D"
	ArgRecoveryTime      = "recovery-time"
	ArgRecoveryTimeShort = "R"
	ArgCidr              = "cidr"
	ArgCidrShort         = "C"
	ArgLanId             = "lan-id"
	ArgLanIdShort        = "L"
	ArgLocation          = "location-id"
	ArgName              = "name"
	ArgNameShort         = "n"
	ArgDbUsername        = "db-username"
	ArgDbUsernameShort   = "U"
	ArgDbPassword        = "db-password"
	ArgDbPasswordShort   = "P"
	ArgRemoveConnection  = "remove-connection"
)

const (
	ClusterId             = "The unique ID of the Cluster"
	BackupId              = "The unique ID of the Backup"
	DefaultClusterTimeout = int(1200)
)
