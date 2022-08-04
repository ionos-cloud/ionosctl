package dataplatform

const (
	listClusterExample   = `ionosctl dataplatform cluster list`
	getClusterExample    = `ionosctl dataplatform cluster get -i CLUSTER_ID`
	createClusterExample = `ionosctl dataplatform cluster create --datacenter-id DATACENTER_ID --name NAME --version DATA_PLATFORM_VERSION

ionosctl dbaas postgres cluster create -D DATACENTER_ID -L LAN_ID -C CIDR -U DB_USERNAME -P DB_PASSWORD`
	updateClusterExample  = `ionosctl dataplatform cluster update -i CLUSTER_ID -n CLUSTER_NAME`
	restoreClusterExample = `ionosctl dataplatform cluster restore -i CLUSTER_ID --backup-id BACKUP_ID`
	deleteClusterExample  = `ionosctl dataplatform cluster delete -i CLUSTER_ID`
	listBackupExample     = `ionosctl dataplatform backup list`
	getBackupExample      = `ionosctl dataplatform backup get -i BACKUP_ID`
	listLogsExample       = `ionosctl dataplatform logs list --cluster-id CLUSTER_ID --since 5h --until 1h`
	listVersionExample    = `ionosctl dataplatform version list`
	getVersionExample     = `ionosctl dataplatform version get --cluster-id CLUSTER_ID`
	listAPIVersionExample = `ionosctl dataplatform api-version list`
	getAPIVersionExample  = `ionosctl dataplatform api-version get`
)
