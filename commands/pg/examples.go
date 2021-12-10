package pg

const (
	listClusterExample   = `ionosctl pg cluster list`
	getClusterExample    = `ionosctl pg cluster get -i CLUSTER_ID`
	createClusterExample = `ionosctl pg cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --db-username DB_USERNAME --db-password DB_PASSWORD

ionosctl pg cluster create -D DATACENTER_ID -L LAN_ID -C CIDR -U DB_USERNAME -P DB_PASSWORD`
	updateClusterExample  = `ionosctl pg cluster update -i CLUSTER_ID -n CLUSTER_NAME`
	restoreClusterExample = `ionosctl pg cluster restore -i CLUSTER_ID --backup-id BACKUP_ID`
	deleteClusterExample  = `ionosctl pg cluster delete -i CLUSTER_ID`
	listBackupExample     = `ionosctl pg backup list`
	getBackupExample      = `ionosctl pg backup get -i BACKUP_ID`
	listLogsExample       = `ionosctl pg logs list --cluster-id CLUSTER_ID`
	listVersionExample    = `ionosctl pg version list`
	getVersionExample     = `ionosctl pg version get --cluster-id CLUSTER_ID`
	listAPIVersionExample = `ionosctl pg api-version list`
	getAPIVersionExample  = `ionosctl pg api-version get`
)
