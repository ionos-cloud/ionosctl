package cloudapi_v6

import (
	"fmt"
	"github.com/ionos-cloud/ionosctl/v6/services/cloudapi-v6/resources"
)

// Resources
const (
	DatacenterResource = "datacenter"
	ServerResource     = "server"
	VolumeResource     = "volume"
	IpBlockResource    = "ipblock"
	SnapshotResource   = "snapshot"
	ImageResource      = "image"
)

// Flags
const (
	ArgAll                 = "all"
	ArgAllShort            = "a"
	ArgCols                = "cols"
	ArgUserData            = "user-data"
	ArgFirstName           = "first-name"
	ArgLastName            = "last-name"
	ArgEmail               = "email"
	ArgEmailShort          = "e"
	ArgPassword            = "password"
	ArgPasswordShort       = "p"
	ArgAdmin               = "admin"
	ArgName                = "name"
	ArgVolumeName          = "volume-name"
	ArgNameShort           = "n"
	ArgApiSubnets          = "api-subnets"
	ArgDescription         = "description"
	ArgDescriptionShort    = "d"
	ArgLocation            = "location"
	ArgLocationShort       = "l"
	ArgDirection           = "direction"
	ArgDirectionShort      = "d"
	ArgAction              = "action"
	ArgActionShort         = "a"
	ArgS3Bucket            = "s3bucket"
	ArgS3BucketShort       = "b"
	ArgSize                = "size"
	ArgSizeShort           = "s"
	ArgBus                 = "bus"
	ArgLicenceType         = "licence-type"
	ArgSshKeyPaths         = "ssh-key-paths"
	ArgSshKeyPathsShort    = "k"
	ArgPublic              = "public"
	ArgPublicShort         = "p"
	ArgIps                 = "ips"
	ArgIp                  = "ip"
	ArgNatGatewayIp        = "nat-gateway-ip"
	ArgDhcp                = "dhcp"
	ArgNetwork             = "network"
	ArgListenerLan         = "listener-lan"
	ArgListenerIp          = "listener-ip"
	ArgListenerPort        = "listener-port"
	ArgAlgorithm           = "algorithm"
	ArgTargetLan           = "target-lan"
	ArgRetries             = "retries"
	ArgClientTimeout       = "client-timeout"
	ArgConnectionTimeout   = "connection-timeout"
	ArgTargetTimeout       = "target-timeout"
	ArgCheck               = "check"
	ArgCheckInterval       = "check-interval"
	ArgHealthCheckEnabled  = "health-check-enabled"
	ArgMaintenanceEnabled  = "maintenance-enabled"
	ArgMaintenance         = "maintenance"
	ArgMaintenanceShort    = "m"
	ArgFirewallActive      = "firewall-active"
	ArgFirewallType        = "firewall-type"
	ArgCpuHotPlug          = "cpu-hot-plug"
	ArgCpuHotUnplug        = "cpu-hot-unplug"
	ArgRamHotPlug          = "ram-hot-plug"
	ArgRamHotUnplug        = "ram-hot-unplug"
	ArgNicHotPlug          = "nic-hot-plug"
	ArgNicHotUnplug        = "nic-hot-unplug"
	ArgDiscVirtioHotPlug   = "disc-virtio-hot-plug"
	ArgDiscVirtioHotUnplug = "disc-virtio-hot-unplug"
	ArgDiscScsiHotPlug     = "disc-scsi-hot-plug"
	ArgDiscScsiHotUnplug   = "disc-scsi-hot-unplug"
	ArgExposeSerial        = "expose-serial"
	ArgRequireLegacyBios   = "require-legacy-bios"
	ArgApplicationType     = "application-type"
	ArgSecAuthProtection   = "sec-auth-protection"
	ArgImageAlias          = "image-alias"
	ArgImageAliasShort     = "a"
	ArgProtocol            = "protocol"
	ArgProtocolShort       = "p"
	ArgSourceSubnet        = "source-subnet"
	ArgTargetSubnet        = "target-subnet"
	ArgSourceMac           = "source-mac"
	ArgSourceIp            = "source-ip"
	ArgDestinationIp       = "destination-ip"
	ArgDestinationIpShort  = "D"
	ArgTargetIp            = "target-ip"
	ArgTargetPort          = "target-port"
	ArgPort                = "port"
	ArgPortShort           = "P"
	ArgWeight              = "weight"
	ArgWeightShort         = "W"
	ArgIcmpCode            = "icmp-code"
	ArgIcmpType            = "icmp-type"
	ArgPortRangeStart      = "port-range-start"
	ArgPortRangeEnd        = "port-range-end"
	ArgLabelUrn            = "label-urn"
	ArgLabelKey            = "label-key"
	ArgLabelValue          = "label-value"
	ArgResourceLimits      = "resource-limits"
	ArgResourceType        = "resource-type"
	ArgForceSecAuth        = "force-secure-auth"
	ArgCreateDc            = "create-dc"
	ArgCreateSnapshot      = "create-snapshot"
	ArgReserveIp           = "reserve-ip"
	ArgAccessLog           = "access-logs"
	ArgS3Privilege         = "s3privilege"
	ArgCreateBackUpUnit    = "create-backup"
	ArgCreatePcc           = "create-pcc"
	ArgCreateNic           = "create-nic"
	ArgCreateK8s           = "create-k8s"
	ArgCreateFlowLog       = "create-flowlog"
	ArgAccessMonitoring    = "access-monitoring"
	ArgAccessCerts         = "access-certs"
	ArgAccessDNS           = "access-dns"
	ArgManageDbaas         = "manage-dbaas"
	ArgManageDataplatform  = "manage-dataplatform"
	ArgManageRegistry      = "manage-registry"
	ArgEditPrivilege       = "edit-privilege"
	ArgSharePrivilege      = "share-privilege"
	ArgS3KeyActive         = "s3key-active"
	ArgK8sVersion          = "k8s-version"
	ArgK8sMinNodeCount     = "min-node-count"
	ArgK8sMaxNodeCount     = "max-node-count"
	ArgK8sMaintenanceDay   = "maintenance-day"
	ArgK8sMaintenanceTime  = "maintenance-time"
	ArgK8sAnnotationKey    = "annotation-key"
	ArgK8sAnnotationValue  = "annotation-value"
	ArgPublicIps           = "public-ips"
	ArgPrivateIps          = "private-ips"
	ArgGatewayIp           = "gateway-ip"
	ArgLatest              = "latest"
	ArgMethod              = "method"
	ArgFilters             = "filters"
	ArgFiltersShort        = "F"
	ArgOrderBy             = "order-by"

	ArgDepth               = "depth"
	ArgDepthShort          = "D"
	ArgCheckTimeout        = "check-timeout"
	ArgPath                = "path"
	ArgMatchType           = "match-type"
	ArgResponse            = "response"
	ArgMessage             = "message"
	ArgMessageShort        = "m"
	ArgRegex               = "regex"
	ArgNegate              = "negate"
	ArgServerCertificates  = "server-certificates"
	ArgQuery               = "query"
	ArgQueryShort          = "Q"
	ArgStatusCode          = "status-code"
	ArgContentType         = "content-type"
	ArgCondition           = "condition"
	ArgConditionShort      = "C"
	ArgConditionType       = "condition-type"
	ArgConditionTypeShort  = "T"
	ArgConditionKey        = "condition-key"
	ArgConditionKeyShort   = "K"
	ArgConditionValue      = "condition-value"
	ArgConditionValueShort = "V"
)

// IDs Flags
const (
	ArgIdShort                   = "i"
	ArgDataCenterId              = "datacenter-id"
	ArgServerId                  = "server-id"
	ArgNatGatewayId              = "natgateway-id"
	ArgApplicationLoadBalancerId = "applicationloadbalancer-id"
	ArgNetworkLoadBalancerId     = "networkloadbalancer-id"
	ArgNicId                     = "nic-id"
	ArgLanId                     = "lan-id"
	ArgLanIds                    = "lan-ids"
	ArgLocationId                = "location-id"
	ArgVolumeId                  = "volume-id"
	ArgLoadBalancerId            = "loadbalancer-id"
	ArgRequestId                 = "request-id"
	ArgSnapshotId                = "snapshot-id"
	ArgImageId                   = "image-id"
	ArgIpBlockId                 = "ipblock-id"
	ArgFirewallRuleId            = "firewallrule-id"
	ArgFlowLogId                 = "flowlog-id"
	ArgUserId                    = "user-id"
	ArgGroupId                   = "group-id"
	ArgResourceId                = "resource-id"
	ArgRuleId                    = "rule-id"
	ArgS3KeyId                   = "s3key-id"
	ArgBackupUnitId              = "backupunit-id"
	ArgPccId                     = "pcc-id"
	ArgK8sNodeId                 = "node-id"
	ArgCdromId                   = "cdrom-id"
	ArgTargetGroupId             = "targetgroup-id"
	ArgTemplateId                = "template-id"
	FlagIPv6CidrBlock            = "ipv6-cidr"
	FlagDHCPv6                   = "dhcpv6"
	FlagIPv6IPs                  = "ipv6-ips"
	FlagIPVersion                = "ip-version"
)

// Descriptions for Flags Resources
const (
	DatacenterId              = "The unique Data Center Id"
	LanId                     = "The unique LAN Id"
	LoadBalancerId            = "The unique Load Balancer Id"
	NicId                     = "The unique NIC Id"
	RequestId                 = "The unique Request Id"
	ServerId                  = "The unique Server Id"
	VolumeId                  = "The unique Volume Id"
	SnapshotId                = "The unique Snapshot Id"
	ImageId                   = "The unique Image Id"
	IpBlockId                 = "The unique IpBlock Id"
	FirewallRuleId            = "The unique FirewallRule Id"
	LocationId                = "The unique Location Id"
	LabelKey                  = "The unique Label Key"
	LabelValue                = "The unique Label Value"
	UserId                    = "The unique User Id"
	GroupId                   = "The unique Group Id"
	ResourceId                = "The unique Resource Id"
	S3KeyId                   = "The unique User S3Key Id"
	BackupUnitId              = "The unique BackupUnit Id"
	PccId                     = "The unique Cross-Connect Id"
	K8sClusterId              = "The unique K8s Cluster Id"
	K8sNodePoolId             = "The unique K8s Node Pool Id"
	K8sNodeId                 = "The unique K8s Node Id"
	CdromId                   = "The unique Cdrom Id"
	TargetGroupId             = "The unique Target Group Id"
	TemplateId                = "The unique Template Id"
	FlowLogId                 = "The unique FlowLog Id"
	NatGatewayId              = "The unique NatGateway Id"
	RuleId                    = "The unique Rule Id"
	ApplicationLoadBalancerId = "The unique ApplicationLoadBalancer Id"
	NetworkLoadBalancerId     = "The unique NetworkLoadBalancer Id"
	ForwardingRuleId          = "The unique ForwardingRule Id"
)

// Descriptions
const (
	ArgDepthDescription   = "Controls the detail depth of the response objects. Max depth is 10."
	ArgOrderByDescription = "Limits results to those containing a matching value for a specific property"
	ArgFiltersDescription = "Limits results to those containing a matching value for a specific property. " +
		"Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2"
	ArgNoHeadersDescription = "When using text output, don't print headers"
	ArgListAllDescription   = "List all resources without the need of specifying parent ID name."
)

// Default values
const (
	DefaultOutputFormat    = "text"
	DefaultWait            = false
	DefaultPublic          = false
	DefaultDhcp            = true
	DefaultFirewallActive  = false
	DefaultTimeoutSeconds  = int(60)
	NlbTimeoutSeconds      = int(300)
	LbTimeoutSeconds       = int(300)
	AlbTimeoutSeconds      = int(10000)
	K8sTimeoutSeconds      = int(600)
	DefaultServerCores     = 2
	DefaultVolumeSize      = 10
	DefaultNicLanId        = 1
	DefaultMaxResults      = int32(0)
	DefaultServerCPUFamily = "AUTO"
	DefaultListDepth       = int32(1)
	DefaultGetDepth        = int32(0)
	DefaultCreateDepth     = int32(0)
	DefaultUpdateDepth     = int32(0)
	DefaultDeleteDepth     = int32(0)
	DefaultMiscDepth       = int32(0) // Attach, Detach (and similar); Server start/stop/suspend/etc.;
)

// Utils
var (
	// Parent resource depth for ListAll, DetachAll, DeleteAll, etc.
	ParentResourceListDepth       = int32(1)
	ParentResourceQueryParams     = resources.QueryParams{Depth: &ParentResourceListDepth}
	ParentResourceListQueryParams = resources.ListQueryParams{QueryParams: ParentResourceQueryParams}

	defaultIPv6CidrBlockDescription = `The /%d IPv6 Cidr as defined in RFC 4291. It needs to be within the %s ` +
		`IPv6 Cidr Block range.`

	FlagIPv6CidrBlockDescriptionForLAN = fmt.Sprintf(
		defaultIPv6CidrBlockDescription+` It can also be set to "AUTO" or "DISABLE".`, 64, "Datacenter",
	)
	FlagIPv6CidrBlockDescriptionForNIC = fmt.Sprintf(defaultIPv6CidrBlockDescription, 80, "LAN")
)
