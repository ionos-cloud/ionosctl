package postgres

const (
	listClusterExample   = `ionosctl dbaas postgres cluster list`
	getClusterExample    = `ionosctl dbaas postgres cluster get -i CLUSTER_ID`
	createClusterExample = `# Basic: smallest cluster with defaults (PostgreSQL 15, 1 instance, 2 cores, 4GB RAM, 20GB HDD, ASYNCHRONOUS). Location is inherited from the datacenter
ionosctl dbaas postgres cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --db-username dbadmin --db-password 'S3cr3tPassw0rd'

# Advanced: named 3-instance HA cluster, sized, on SSD Premium, with a maintenance window and a synchronous replication mode
ionosctl dbaas postgres cluster create --name prod-orders --version 16 --instances 3 --cores 4 --ram 8GB --storage-size 100GB --storage-type SSD_PREMIUM --sync STRICTLY_SYNCHRONOUS --location de/fra --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --backup-location eu-central-3 --maintenance-day Sunday --maintenance-time 03:00:00 --db-username dbadmin --db-password 'S3cr3tPassw0rd'

# Clone: create a new cluster from an existing backup, replayed to a point in time (PITR)
ionosctl dbaas postgres cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr 192.168.1.100/24 --db-username dbadmin --db-password 'S3cr3tPassw0rd' --backup-id BACKUP_ID --recovery-time 2024-01-15T10:00:00Z`
	updateClusterExample = `# Rename a cluster
ionosctl dbaas postgres cluster update -i CLUSTER_ID --name new-name

# Scale up compute and storage, and set a maintenance window (day and time must be given together)
ionosctl dbaas postgres cluster update -i CLUSTER_ID --cores 4 --ram 16GB --storage-size 200GB --maintenance-day Saturday --maintenance-time 04:00:00`
	restoreClusterExample = `# Restore in place to the end of a backup
ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID

# Restore in place to a specific point in time within the backup's recovery window (PITR)
ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID --recovery-time 2024-01-15T10:00:00Z`
	deleteClusterExample = `ionosctl dbaas postgres cluster delete -i CLUSTER_ID`
	listBackupExample    = `ionosctl dbaas postgres backup list`
	getBackupExample     = `ionosctl dbaas postgres backup get -i BACKUP_ID`
	listLogsExample      = `# Last 5 hours up to 1 hour ago, oldest first
ionosctl dbaas postgres logs list --cluster-id CLUSTER_ID --since 5h --until 1h --direction FORWARD

# Absolute time range with an explicit line limit
ionosctl dbaas postgres logs list --cluster-id CLUSTER_ID --start-time 2021-10-05T11:30:17Z --end-time 2021-10-05T12:30:17Z --limit 500`
	listVersionExample    = `ionosctl dbaas postgres version list`
	getVersionExample     = `ionosctl dbaas postgres version get --cluster-id CLUSTER_ID`
	listAPIVersionExample = `ionosctl dbaas postgres api-version list`
	getAPIVersionExample  = `ionosctl dbaas postgres api-version get`
)
