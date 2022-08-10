package dataplatform

const (
	listClusterExample    = `ionosctl dataplatform cluster list`
	getClusterExample     = `ionosctl dataplatform cluster get -i CLUSTER_ID`
	createClusterExample  = `ionosctl dataplatform cluster create --datacenter-id DATACENTER_ID --name NAME --version DATA_PLATFORM_VERSION`
	updateClusterExample  = `ionosctl dataplatform cluster update -i CLUSTER_ID -n CLUSTER_NAME`
	deleteClusterExample  = `ionosctl dataplatform cluster delete -i CLUSTER_ID`
	listNodePoolsExample  = `ionosctl dataplatform nodepool list --cluster-id CLUSTER_ID`
	getNodePoolExample    = `ionosctl dataplatform nodepool get --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`
	createNodePoolExample = `ionosctl dataplatform nodepool create --datacenter-id DATACENTER_ID --cluster-id CLUSTER_ID --name NAME --node-count NODE_COUNT`
	updateNodePoolExample = `ionosctl dataplatform nodepool update --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID --node-count NODE_COUNT`
	deleteNodePoolExample = `ionosctl dataplatform nodepool delete --cluster-id CLUSTER_ID --nodepool-id NODEPOOL_ID`
	listVersionsExample   = `ionosctl dataplatform versions list --cluster-id CLUSTER_ID`
	getKubeConfigExample  = `ionosctl dataplatform kubeconfig get --cluster-id CLUSTER_ID`
)
