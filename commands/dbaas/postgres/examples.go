package postgres

const (
	listClusterExample   = `ionosctl dbaas postgres cluster list`
	getClusterExample    = `ionosctl dbaas postgres cluster get -i CLUSTER_ID`
	createClusterExample = `ionosctl dbaas postgres cluster create --datacenter-id DATACENTER_ID --lan-id LAN_ID --cidr CIDR --db-username DB_USERNAME --db-password DB_PASSWORD

ionosctl dbaas postgres cluster create -D DATACENTER_ID -L LAN_ID -C CIDR -U DB_USERNAME -P DB_PASSWORD`
	updateClusterExample  = `ionosctl dbaas postgres cluster update -i CLUSTER_ID -n CLUSTER_NAME`
	restoreClusterExample = `ionosctl dbaas postgres cluster restore -i CLUSTER_ID --backup-id BACKUP_ID`
	deleteClusterExample  = `ionosctl dbaas postgres cluster delete -i CLUSTER_ID`
	listBackupExample     = `ionosctl dbaas postgres backup list`
	getBackupExample      = `ionosctl dbaas postgres backup get -i BACKUP_ID`
	listLogsExample       = `ionosctl dbaas postgres logs list --cluster-id CLUSTER_ID`
	listVersionExample    = `ionosctl dbaas postgres version list`
	getVersionExample     = `ionosctl dbaas postgres version get --cluster-id CLUSTER_ID`
	listAPIVersionExample = `ionosctl dbaas postgres api-version list`
	getAPIVersionExample  = `ionosctl dbaas postgres api-version get`
)
