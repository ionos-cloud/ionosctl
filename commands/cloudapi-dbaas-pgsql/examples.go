package cloudapi_dbaas_pgsql

const (
	listClusterExample    = `ionosctl dbaas-pgsql cluster list`
	getClusterExample     = `ionosctl dbaas-pgsql cluster get -i CLUSTER_ID`
	createClusterExample  = `ionosctl dbaas-pgsql cluster create -p POSTGRES_VERSION --vdc-id VDC_ID --lan-id LAN_ID --ip IP_ADDRESS`
	updateClusterExample  = `ionosctl dbaas-pgsql cluster update -i CLUSTER_ID -n CLUSTER_NAME`
	restoreClusterExample = `ionosctl dbaas-pgsql cluster restore -i CLUSTER_ID --backup-id BACKUP_ID`
	deleteClusterExample  = `ionosctl dbaas-pgsql cluster delete -i CLUSTER_ID`
	listBackupExample     = `ionosctl dbaas-pgsql backup list`
	getBackupExample      = `ionosctl dbaas-pgsql backup get -i BACKUP_ID`
	getLogsExample        = `ionosctl dbaas-pgsql logs get --cluster-id CLUSTER_ID`
	getQuotaExample       = `ionosctl dbaas-pgsql quota list`
	listVersionExample    = `ionosctl dbaas-pgsql version list`
	getVersionExample     = `ionosctl dbaas-pgsql version get --cluster-id CLUSTER_ID`
	listAPIVersionExample = `ionosctl dbaas-pgsql api-version list`
	getAPIVersionExample  = `ionosctl dbaas-pgsql api-version get`
)
