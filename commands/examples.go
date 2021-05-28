package commands

const (

	/*
		Login Examples
	*/
	loginExamples = `ionosctl login --user USERNAME --password PASSWORD
Status: Authentication successful!

ionosctl login 
Enter your username:
USERNAME
Enter your password:

Status: Authentication successful!`

	/*
		Location Examples
	*/
	listLocationExample = `ionosctl location list`
	getLocationExample  = `ionosctl location get --location-id LOCATION_ID`

	/*
		Data Center Examples
	*/
	listDatacenterExample = `ionosctl datacenter list

ionosctl datacenter list --cols "DatacenterId,Name,Location,Version"`
	getDatacenterExample    = `ionosctl datacenter get --datacenter-id DATACENTER_ID`
	createDatacenterExample = `ionosctl datacenter create --name NAME --location LOCATION_ID

ionosctl datacenter create --name NAME --location LOCATION_ID --wait-for-request`
	updateDatacenterExample = `ionosctl datacenter update --datacenter-id DATACENTER_ID --description DESCRIPTION --cols "DatacenterId,Description"`
	deleteDatacenterExample = `ionosctl datacenter delete --datacenter-id DATACENTER_ID

ionosctl datacenter delete --datacenter-id DATACENTER_ID --force --wait-for-request`

	/*
		Server Examples
	*/
	listServerExample   = `ionosctl server list --datacenter-id DATACENTER_ID`
	getServerExample    = `ionosctl server get --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	createServerExample = `ionosctl server create --datacenter-id DATACENTER_ID --name NAME --cores 2 --ram 512MB -w -W`
	updateServerExample = `ionosctl server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 4`
	deleteServerExample = `ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID

ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --force`
	startServerExample        = `ionosctl server start --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	stopServerExample         = `ionosctl server stop --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	resetServerExample        = `ionosctl server reset --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	attachVolumeServerExample = `ionosctl server volume attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`
	listVolumesServerExample  = `ionosctl server volume list --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getVolumeServerExample    = `ionosctl server volume get --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`
	detachVolumeServerExample = `ionosctl server volume detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`

	/*
		Volume Examples
	*/
	createVolumeExample = `ionosctl volume create --datacenter-id DATACENTER_ID --name NAME`
	updateVolumeExample = `ionosctl volume update --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --size 20`
	listVolumeExample   = `ionosctl volume list --datacenter-id DATACENTER_ID`
	getVolumeExample    = `ionosctl volume get --datacenter-id DATACENTER_ID --volume-id VOLUME_ID`
	deleteVolumeExample = `ionosctl volume delete --datacenter-id DATACENTER_ID --volume-id VOLUME_ID`

	/*
		Load Balancer Examples
	*/
	createLoadbalancerExample    = `ionosctl loadbalancer create --datacenter-id DATACENTER_ID --name NAME`
	updateLoadbalancerExample    = `ionosctl loadbalancer update --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --dhcp=false --wait-for-request`
	listLoadbalancerExample      = `ionosctl loadbalancer list --datacenter-id DATACENTER_ID`
	getLoadbalancerExample       = `ionosctl loadbalancer get --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID`
	deleteLoadbalancerExample    = `ionosctl loadbalancer delete --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --force --wait-for-request`
	attachNicLoadbalancerExample = `ionosctl loadbalancer nic attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --loadbalancer-id LOADBALANCER_ID`
	listNicsLoadbalancerExample  = `ionosctl loadbalancer nic list --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID`
	getNicLoadbalancerExample    = `ionosctl loadbalancer nic get --datacenter-id DATACENTER_ID --loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID`
	detachNicLoadbalancerExample = `ionosctl loadbalancer nic detach --datacenter-id DATACENTER_ID--loadbalancer-id LOADBALANCER_ID --nic-id NIC_ID`

	/*
		NIC Examples
	*/
	createNicExample = `ionosctl nic create --datacenter-id DATACENTER_ID --server-id SERVER_ID --name NAME`
	updateNicExample = `ionosctl nic update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --wait-for-request`
	listNicExample   = `ionosctl nic list --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getNicExample    = `ionosctl nic get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID`
	deleteNicExample = `ionosctl nic delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --force`

	/*
		Lan Examples
	*/
	createLanExample = `ionosctl lan create --datacenter-id DATACENTER_ID --name NAMEd`
	updateLanExample = `ionosctl lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --name NAME --public=true`
	listLanExample   = `ionosctl lan list --datacenter-id DATACENTER_ID`
	getLanExample    = `ionosctl lan get --datacenter-id DATACENTER_ID --lan-id LAN_ID`
	deleteLanExample = `ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID

ionosctl lan delete --datacenter-id DATACENTER_ID --lan-id LAN_ID --wait-for-request`

	/*
		IP Failover Examples
	*/
	addIpFailoverExample    = `ionosctl ipfailover add --datacenter-id DATACENTER_ID --server-id SERVER_ID --lan-id LAN_ID --nic-id NIC_ID --ip "x.x.x.x"`
	removeIpFailoverExample = `ionosctl ipfailover remove --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --lan-id LAN_ID --ip "x.x.x.x"`
	listIpFailoverExample   = `ionosctl ipfailover list --datacenter-id DATACENTER_ID --lan-id LAN_ID`

	/*
		Request Examples
	*/
	getRequestExample  = `ionosctl request get --request-id REQUEST_ID`
	waitRequestExample = `ionosctl request wait --request-id REQUEST_ID`

	/*
		Image Examples
	*/
	listImagesExample = `ionosctl image list

ionosctl image list --location us/las --type HDD --licence-type LINUX`
	getImageExample = `ionosctl image get --image-id IMAGE_ID`

	/*
		Snapshot Examples
	*/
	listSnapshotsExample   = `ionosctl snapshot list`
	getSnapshotExample     = `ionosctl snapshot get --snapshot-id SNAPSHOT_ID`
	createSnapshotExample  = `ionosctl snapshot create --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --name NAME`
	updateSnapshotExample  = `ionosctl snapshot update --snapshot-id SNAPSHOT_ID --name NAME`
	restoreSnapshotExample = `ionosctl snapshot restore --snapshot-id SNAPSHOT_ID --datacenter-id DATACENTER_ID --volume-id VOLUME_ID --wait-for-request`
	deleteSnapshotExample  = `ionosctl snapshot delete --snapshot-id SNAPSHOT_ID --wait-for-request`

	/*
		IpBlock Examples
	*/
	listIpBlockExample     = `ionosctl ipblock list`
	getIpBlockExample      = `ionosctl ipblock get --ipblock-id IPBLOCK_ID`
	createIpBlockExample   = `ionosctl ipblock create --ipblock-name NAME --ipblock-location LOCATION_ID --ipblock-size IPBLOCK_SIZE`
	updateIpBlockExample   = `ionosctl ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME`
	deleteIpBlockExample   = `ionosctl ipblock delete --ipblock-id IPBLOCK_ID --wait-for-request`
	listIpConsumersExample = `ionosctl ipconsumer list --ipblock-id IPBLOCK_ID`

	/*
		Firewall Rule Examples
	*/
	listFirewallRuleExample   = `ionosctl firewallrule list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID`
	getFirewallRuleExample    = `ionosctl firewallrule get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID`
	createFirewallRuleExample = `ionosctl firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol PROTOCOL --name NAME --port-range-start PORT_START --port-range-end PORT_END`
	updateFirewallRuleExample = `ionosctl firewallrule update --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID --name NAME --wait-for-request`
	deleteFirewallRuleExample = `ionosctl firewallrule delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID`

	/*
		Label Examples
	*/
	listLabelsExample = `ionosctl label list

ionosctl label list --resource-type datacenter --datacenter-id DATACENTER_ID`
	getLabelByUrnExample = `ionosctl label get-by-urn --label-urn "urn:label:server:SERVER_ID:test"`
	getLabelExample      = `ionosctl label get --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY`
	addLabelExample      = `ionosctl label add --resource-type server --datacenter-id DATACENTER_ID --server-id SERVER_ID  --label-key LABEL_KEY --label-value LABEL_VALUE

ionosctl label add --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY --label-value LABEL_VALUE`
	removeLabelExample = `ionosctl label remove --resource-type datacenter --datacenter-id DATACENTER_ID --label-key LABEL_KEY`

	/*
		Contract Resources Examples
	*/
	getContractExample = `ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]`

	/*
		User Examples
	*/
	listUserExample   = `ionosctl user list`
	getUserExample    = `ionosctl user get --user-id USER_ID`
	createUserExample = `ionosctl user create --first-name NAME --last-name NAME --email EMAIL --password PASSWORD`
	updateUserExample = `ionosctl user update --user-id USER_ID --admin=true`
	deleteUserExample = `ionosctl user delete --user-id USER_ID --force`

	/*
		Group Examples
	*/
	createGroupExample = `ionosctl group create --name NAME --wait-for-request`
	getGroupExample    = `ionosctl group get --group-id GROUP_ID`
	listGroupExample   = `ionosctl group list`
	updateGroupExample = `ionosctl group update --group-id GROUP_ID --reserve-ip`
	deleteGroupExample = `ionosctl group delete --group-id GROUP_ID`

	/*
		Group Users Examples
	*/
	listGroupUsersExample  = `ionosctl group user list --group-id GROUP_ID`
	removeGroupUserExample = `ionosctl group user remove --group-id GROUP_ID --user-id USER_ID`
	addGroupUserExample    = `ionosctl group user add --group-id GROUP_ID --user-id USER_ID`

	/*
		Group Resources Example
	*/
	listGroupResourcesExample = `ionosctl group resource list --group-id GROUP_ID`

	/*
		Resources Example
	*/
	listResourcesExample = `ionosctl resource list`
	getResourceExample   = `ionosctl resource get --resource-type ipblock`

	/*
		Share Example
	*/
	listSharesExample  = `ionosctl share list --group-id GROUP_ID`
	getShareExample    = `ionosctl share get --group-id GROUP_ID --resource-id RESOURCE_ID`
	createShareExample = `ionosctl share create --group-id GROUP_ID --resource-id RESOURCE_ID`
	updateShareExample = `ionosctl share update --group-id GROUP_ID --resource-id RESOURCE_ID --share-privilege`
	deleteShareExample = `ionosctl share delete --group-id GROUP_ID --resource-id RESOURCE_ID --wait-for-request`

	/*
		S3Keys Example
	*/
	listS3KeysExample  = `ionosctl user s3key list --user-id USER_ID`
	getS3KeyExample    = `ionosctl user s3key get --user-id USER_ID --s3key-id S3KEY_ID`
	createS3KeyExample = `ionosctl user s3key create --user-id USER_ID`
	updateS3KeyExample = `ionosctl user s3key update --user-id USER_ID --s3key-id S3KEY_ID --s3key-active=false`
	deleteS3KeyExample = `ionosctl user s3key delete --user-id USER_ID --s3key-id S3KEY_ID --force`

	/*
		BackupUnit Example
	*/
	listBackupUnitsExample  = `ionosctl backupunit list`
	getBackupUnitExample    = `ionosctl backupunit get --backupunit-id BACKUPUNIT_ID`
	getBackupUnitSSOExample = `ionosctl backupunit get-sso-url --backupunit-id BACKUPUNIT_ID`
	createBackupUnitExample = `ionosctl backupunit create --name NAME --email EMAIL --password PASSWORD`
	updateBackupUnitExample = `ionosctl backupunit update --backupunit-id BACKUPUNIT_ID --email EMAIL`
	deleteBackupUnitExample = `ionosctl backupunit delete --backupunit-id BACKUPUNIT_ID`

	/*
		Private Cross-Connect Example
	*/
	listPccsExample     = `ionosctl pcc list`
	getPccExample       = `ionosctl pcc get --pcc-id PCC_ID`
	listPccPeersExample = `ionosctl pcc peers list --pcc-id PCC_ID`
	createPccExample    = `ionosctl pcc create --name NAME --description DESCRIPTION --wait-for-request`
	updatePccExample    = `ionosctl pcc update --pcc-id PCC_ID --description DESCRIPTION`
	deletePccExample    = `ionosctl pcc delete --pcc-id PCC_ID --wait-for-request`

	/*
		K8s Example
	*/
	listK8sClustersExample  = `ionosctl k8s cluster list`
	getK8sClusterExample    = `ionosctl k8s cluster get --cluster-id CLUSTER_ID`
	createK8sClusterExample = `ionosctl k8s cluster create --name NAME`
	updateK8sClusterExample = `ionosctl k8s cluster update --cluster-id CLUSTER_ID --name NAME`
	deleteK8sClusterExample = `ionosctl k8s cluster delete --cluster-id CLUSTER_ID`

	listK8sNodePoolsExample  = `ionosctl k8s nodepool list --cluster-id CLUSTER_ID`
	getK8sNodePoolExample    = `ionosctl k8s nodepool get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`
	createK8sNodePoolExample = `ionosctl k8s nodepool create --datacenter-id DATACENTER_ID --cluster-id CLUSTER_ID --name NAME`
	updateK8sNodePoolExample = `ionosctl k8s nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count NODE_COUNT`
	deleteK8sNodePoolExample = `ionosctl k8s nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`

	deleteK8sNodeExample   = `ionosctl k8s node delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	recreateK8sNodeExample = `ionosctl k8s node recreate --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	getK8sNodeExample      = `ionosctl k8s node get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	listK8sNodesExample    = `ionosctl k8s node list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`

	getK8sKubeconfigExample = `ionosctl k8s kubeconfig get --cluster-id CLUSTER_ID`
	listK8sVersionsExample  = `ionosctl k8s version list`
	getK8sVersionExample    = `ionosctl k8s version get`

	/*
		Server Cdrom Example
	*/
	listCdromServerExample   = `ionosctl server cdrom list --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getCdromServerExample    = `ionosctl server cdrom get --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID`
	attachCdromServerExample = `ionosctl server cdrom attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request`
	detachCdromServerExample = `ionosctl server cdrom detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --cdrom-id CDROM_ID --wait-for-request --force`
)
