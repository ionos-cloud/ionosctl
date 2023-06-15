package commands

/*
 * TODO: This is hard to maintain, they should be aware of the flags they're using instead of being hardcoded
 * TODO: This is incompatible with cobra command structure: Having subfolders for commands means losing access to these unexported fields
 */

const (
	/*
		Location Examples
	*/
	listLocationExample    = `ionosctl location list`
	getLocationExample     = `ionosctl location get --location-id LOCATION_ID`
	listLocationCpuExample = `ionosctl location cpu list --location-id LOCATION_ID`

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
	createServerExample = `ionosctl server create --datacenter-id DATACENTER_ID --cores 2 --ram 512MB -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --licence-type LICENCE_TYPE -w -W

ionosctl server create --datacenter-id DATACENTER_ID --type CUBE --template-id TEMPLATE_ID --image-id IMAGE_ID --password IMAGE_PASSWORD -w -W`
	updateServerExample = `ionosctl server update --datacenter-id DATACENTER_ID --server-id SERVER_ID --cores 4`
	deleteServerExample = `ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID

ionosctl server delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --force`
	startServerExample        = `ionosctl server start --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	stopServerExample         = `ionosctl server stop --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	rebootServerExample       = `ionosctl server reboot --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	attachVolumeServerExample = `ionosctl server volume attach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`
	listVolumesServerExample  = `ionosctl server volume list --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getVolumeServerExample    = `ionosctl server volume get --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`
	detachVolumeServerExample = `ionosctl server volume detach --datacenter-id DATACENTER_ID --server-id SERVER_ID --volume-id VOLUME_ID`
	suspendServerExample      = `ionosctl server suspend --datacenter-id DATACENTER_ID -i SERVER_ID`
	resumeServerExample       = `ionosctl server resume --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getTokenServerExample     = `ionosctl server token get --datacenter-id DATACENTER_ID --server-id SERVER_ID`
	getConsoleServerExample   = `ionosctl server console get --datacenter-id DATACENTER_ID --server-id SERVER_ID`

	/*
		Volume Examples
	*/
	createVolumeExample = `ionosctl volume create --datacenter-id DATACENTER_ID --name NAME

ionosctl volume create --datacenter-id DATACENTER_ID --name NAME --image-alias IMAGE_ALIAS --ssh-keys-path "SSH_KEY_PATH1,SSH_KEY_PATH2"`
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
	createLanExample = `ionosctl lan create --datacenter-id DATACENTER_ID --name NAME --public=true`
	updateLanExample = `ionosctl lan update --datacenter-id DATACENTER_ID --lan-id LAN_ID --name NAME --public=false`
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
	listRequestExample = `ionosctl request list --latest N`
	getRequestExample  = `ionosctl request get --request-id REQUEST_ID`
	waitRequestExample = `ionosctl request wait --request-id REQUEST_ID`

	/*
		Image Examples
	*/

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
	createIpBlockExample   = `ionosctl ipblock create --name NAME --location LOCATION_ID --size IPBLOCK_SIZE`
	updateIpBlockExample   = `ionosctl ipblock update --ipblock-id IPBLOCK_ID --ipblock-name NAME`
	deleteIpBlockExample   = `ionosctl ipblock delete --ipblock-id IPBLOCK_ID --wait-for-request`
	listIpConsumersExample = `ionosctl ipconsumer list --ipblock-id IPBLOCK_ID`

	/*
		Firewall Rule Examples
	*/
	listFirewallRuleExample   = `ionosctl firewallrule list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID`
	getFirewallRuleExample    = `ionosctl firewallrule get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --firewallrule-id FIREWALLRULE_ID`
	createFirewallRuleExample = `ionosctl firewallrule create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --protocol PROTOCOL --direction DIRECTION --destination-ip DESTINATION_IP`
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

	listK8sNodePoolLanExample   = `ionosctl k8s nodepool lan list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`
	addK8sNodePoolLanExample    = `ionosctl k8s nodepool lan add --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID`
	removeK8sNodePoolLanExample = `ionosctl k8s nodepool lan remove --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --lan-id LAN_ID`

	deleteK8sNodeExample   = `ionosctl k8s node delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	recreateK8sNodeExample = `ionosctl k8s node recreate --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	getK8sNodeExample      = `ionosctl k8s node get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-id NODE_ID`
	listK8sNodesExample    = `ionosctl k8s node list --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`

	getK8sKubeconfigExample = `ionosctl k8s kubeconfig get --cluster-id CLUSTER_ID`
	listK8sVersionsExample  = `ionosctl k8s version list`
	getK8sVersionExample    = `ionosctl k8s version get`

	/*
		Template Example
	*/
	listTemplateExample = `ionosctl template list`
	getTemplateExample  = `ionosctl template get -i TEMPLATE_ID`

	/*
		FlowLog Example
	*/
	listFlowLogExample   = `ionosctl flowlog list --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID`
	getFlowLogExample    = `ionosctl flowlog get --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID`
	createFlowLogExample = `ionosctl flowlog create --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --name NAME --action ACTION --direction DIRECTION --s3bucket BUCKET_NAME`
	deleteFlowLogExample = `ionosctl flowlog delete --datacenter-id DATACENTER_ID --server-id SERVER_ID --nic-id NIC_ID --flowlog-id FLOWLOG_ID -f -w`

	/*
		NatGateway Example
	*/
	listNatGatewayExample   = `ionosctl natgateway list --datacenter-id DATACENTER_ID`
	getNatGatewayExample    = `ionosctl natgateway get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID`
	createNatGatewayExample = `ionosctl natgateway create --datacenter-id DATACENTER_ID --name NAME --ips IP_1,IP_2`
	updateNatGatewayExample = `ionosctl natgateway update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME`
	deleteNatGatewayExample = `ionosctl natgateway delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID`

	/*
		NatGateway Lan Example
	*/
	listNatGatewayLanExample = `ionosctl natgateway lan list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID`
	addNatGatewayLanExample  = `ionosctl natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID

ionosctl natgateway lan add --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID --ips IP_1,IP_2`
	removeNatGatewayLanExample = `ionosctl natgateway lan remove --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --lan-id LAN_ID`

	/*
		NatGateway Rule Example
	*/
	listNatGatewayRuleExample   = `ionosctl natgateway rule list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID`
	getNatGatewayRuleExample    = `ionosctl natgateway rule get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID`
	createNatGatewayRuleExample = `ionosctl natgateway rule create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET`
	updateNatGatewayRuleExample = `ionosctl natgateway rule update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name NAME`
	deleteNatGatewayRuleExample = `ionosctl natgateway rule delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID`

	/*
		NatGateway FlowLog Example
	*/
	listNatGatewayFlowLogExample   = `ionosctl natgateway flowlog list --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID`
	getNatGatewayFlowLogExample    = `ionosctl natgateway flowlog get --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID`
	createNatGatewayFlowLogExample = `ionosctl natgateway flowlog create --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --name NAME --ip IP_1 --source-subnet SOURCE_SUBNET --target-subnet TARGET_SUBNET`
	updateNatGatewayFlowLogExample = `ionosctl natgateway flowlog update --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID --name NAME`
	deleteNatGatewayFlowLogExample = `ionosctl natgateway flowlog delete --datacenter-id DATACENTER_ID --natgateway-id NATGATEWAY_ID --rule-id RULE_ID`

	/*
		Network Load Balancer Example
	*/
	listNetworkLoadBalancerExample   = `ionosctl networkloadbalancer list --datacenter-id DATACENTER_ID`
	getNetworkLoadBalancerExample    = `ionosctl networkloadbalancer get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID`
	createNetworkLoadBalancerExample = `ionosctl networkloadbalancer create --datacenter-id DATACENTER_ID`
	updateNetworkLoadBalancerExample = `ionosctl networkloadbalancer update --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID --name NAME`
	deleteNetworkLoadBalancerExample = `ionosctl networkloadbalancer delete --datacenter-id DATACENTER_ID -i NETWORKLOADBALANCER_ID`

	/*
		Network Load Balancer FlowLog Example
	*/
	listNetworkLoadBalancerFlowLogExample   = `ionosctl networkloadbalancer flowlog list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID`
	getNetworkLoadBalancerFlowLogExample    = `ionosctl networkloadbalancer flowlog get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID`
	createNetworkLoadBalancerFlowLogExample = `ionosctl networkloadbalancer flowlog create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --s3bucket BUCKET_NAME`
	updateNetworkLoadBalancerFlowLogExample = `ionosctl networkloadbalancer flowlog update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID --name NAME`
	deleteNetworkLoadBalancerFlowLogExample = `ionosctl networkloadbalancer flowlog delete --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FLOWLOG_ID`

	/*
		Network Load Balancer ForwardingRule Example
	*/
	listNetworkLoadBalancerForwardingRuleExample   = `ionosctl networkloadbalancer rule list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID`
	getNetworkLoadBalancerForwardingRuleExample    = `ionosctl nlb rule get --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID`
	createNetworkLoadBalancerForwardingRuleExample = `ionosctl networkloadbalancer rule create --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --listener-ip LISTENER_IP --listener-port LISTENER_PORT`
	updateNetworkLoadBalancerForwardingRuleExample = `ionosctl nlb rule update --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME`
	deleteNetworkLoadBalancerForwardingRuleExample = `ionosctl networkloadbalancer rule delete --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID -i FORWARDINGRULE_ID`

	/*
		Network Load Balancer ForwardingRule Target Example
	*/
	listNetworkLoadBalancerRuleTargetExample   = `ionosctl networkloadbalancer rule target list --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID`
	addNetworkLoadBalancerRuleTargetExample    = `ionosctl networkloadbalancer rule target add --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip TARGET_IP --port TARGET_PORT`
	removeNetworkLoadBalancerRuleTargetExample = `ionosctl nlb rule target remove --datacenter-id DATACENTER_ID --networkloadbalancer-id NETWORKLOADBALANCER_ID --rule-id FORWARDINGRULE_ID --ip TARGET_IP --port TARGET_PORT`

	/*
		Target Group Example
	*/
	listTargetGroupExample   = `ionosctl targetgroup list`
	getTargetGroupExample    = `ionosctl targetgroup get -i TARGET_GROUP_ID`
	createTargetGroupExample = `ionosctl targetgroup create --name TARGET_GROUP_NAME`
	updateTargetGroupExample = `ionosctl targetgroup update --targetgroup-id TARGET_GROUP_ID --name TARGET_GROUP_NEW_NAME -w`
	deleteTargetGroupExample = `ionosctl targetgroup delete --targetgroup-id TARGET_GROUP_ID --force`

	listTargetGroupTargetExample   = `ionosctl targetgroup target list --targetgroup-id TARGET_GROUP_ID`
	addTargetGroupTargetExample    = `ionosctl targetgroup target add --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT`
	removeTargetGroupTargetExample = `ionosctl targetgroup target remove --targetgroup-id TARGET_GROUP_ID --ip TARGET_IP --port TARGET_PORT`

	/*
		Application Load Balancer Example
	*/
	listApplicationLoadBalancerExample   = `ionosctl applicationloadbalancer list --datacenter-id DATACENTER_ID`
	getApplicationLoadBalancerExample    = `ionosctl applicationloadbalancer get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID`
	createApplicationLoadBalancerExample = `ionosctl applicationloadbalancer create --datacenter-id DATACENTER_ID`
	updateApplicationLoadBalancerExample = `ionosctl applicationloadbalancer update --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID --name NAME`
	deleteApplicationLoadBalancerExample = `ionosctl applicationloadbalancer delete --datacenter-id DATACENTER_ID -i APPLICATIONLOADBALANCER_ID`

	/*
		Application Load Balancer ForwardingRule Example
	*/
	listApplicationLoadBalancerForwardingRuleExample   = `ionosctl applicationloadbalancer rule list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID`
	getApplicationLoadBalancerForwardingRuleExample    = `ionosctl alb rule get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID`
	createApplicationLoadBalancerForwardingRuleExample = `ionosctl applicationloadbalancer rule create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --listener-ip LISTENER_IP --listener-port LISTENER_PORT`
	updateApplicationLoadBalancerForwardingRuleExample = `ionosctl alb rule update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID --name NAME`
	deleteApplicationLoadBalancerForwardingRuleExample = `ionosctl applicationloadbalancer rule delete --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FORWARDINGRULE_ID`

	/*
		Application Load Balancer ForwardingRule HttpRule Example
	*/
	listApplicationLoadBalancerForwardingRuleHttpExample   = `ionosctl alb rule http list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID`
	addApplicationLoadBalancerForwardingRuleHttpExample    = `ionosctl alb rule http add --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME --type TYPE`
	removeApplicationLoadBalancerForwardingRuleHttpExample = `ionosctl alb rule httprule remove --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --rule-id FORWARDINGRULE_ID -n NAME`

	/*
		Application Load Balancer FlowLog Example
	*/
	listApplicationLoadBalancerFlowLogExample   = `ionosctl applicationloadbalancer flowlog list --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID`
	getApplicationLoadBalancerFlowLogExample    = `ionosctl applicationloadbalancer flowlog get --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID`
	createApplicationLoadBalancerFlowLogExample = `ionosctl applicationloadbalancer flowlog create --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID --action ACTION --name NAME --direction DIRECTION --s3bucket BUCKET_NAME`
	updateApplicationLoadBalancerFlowLogExample = `ionosctl applicationloadbalancer flowlog update --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID --name NAME`
	deleteApplicationLoadBalancerFlowLogExample = `ionosctl applicationloadbalancer flowlog delete --datacenter-id DATACENTER_ID --applicationloadbalancer-id APPLICATIONLOADBALANCER_ID -i FLOWLOG_ID`
)
