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
	FlagAll                 = "all"
	FlagAllShort            = "a"
	FlagCols                = "cols"
	FlagUserData            = "user-data"
	FlagFirstName           = "first-name"
	FlagLastName            = "last-name"
	FlagEmail               = "email"
	FlagEmailShort          = "e"
	FlagPassword            = "password"
	FlagPasswordShort       = "p"
	FlagAdmin               = "admin"
	FlagName                = "name"
	FlagVolumeName          = "volume-name"
	FlagNameShort           = "n"
	FlagApiSubnets          = "api-subnets"
	FlagDescription         = "description"
	FlagDescriptionShort    = "d"
	FlagLocation            = "location"
	FlagLocationShort       = "l"
	FlagDirection           = "direction"
	FlagDirectionShort      = "d"
	FlagAction              = "action"
	FlagActionShort         = "a"
	ArgS3Bucket             = "s3bucket"
	ArgS3BucketShort        = "b"
	FlagSize                = "size"
	FlagSizeShort           = "s"
	FlagBus                 = "bus"
	FlagLicenceType         = "licence-type"
	FlagSshKeyPaths         = "ssh-key-paths"
	FlagSshKeyPathsShort    = "k"
	FlagPublic              = "public"
	FlagPublicShort         = "p"
	FlagIps                 = "ips"
	FlagIp                  = "ip"
	FlagNatGatewayIp        = "nat-gateway-ip"
	FlagDhcp                = "dhcp"
	FlagNetwork             = "network"
	FlagListenerLan         = "listener-lan"
	FlagListenerIp          = "listener-ip"
	FlagListenerPort        = "listener-port"
	FlagAlgorithm           = "algorithm"
	FlagTargetLan           = "target-lan"
	FlagRetries             = "retries"
	FlagClientTimeout       = "client-timeout"
	FlagConnectionTimeout   = "connection-timeout"
	FlagTargetTimeout       = "target-timeout"
	FlagCheck               = "check"
	FlagCheckInterval       = "check-interval"
	FlagHealthCheckEnabled  = "health-check-enabled"
	FlagMaintenanceEnabled  = "maintenance-enabled"
	FlagMaintenance         = "maintenance"
	FlagMaintenanceShort    = "m"
	FlagFirewallActive      = "firewall-active"
	FlagFirewallType        = "firewall-type"
	FlagCpuHotPlug          = "cpu-hot-plug"
	FlagCpuHotUnplug        = "cpu-hot-unplug"
	FlagRamHotPlug          = "ram-hot-plug"
	FlagRamHotUnplug        = "ram-hot-unplug"
	FlagNicHotPlug          = "nic-hot-plug"
	FlagNicHotUnplug        = "nic-hot-unplug"
	FlagDiscVirtioHotPlug   = "disc-virtio-hot-plug"
	FlagDiscVirtioHotUnplug = "disc-virtio-hot-unplug"
	FlagDiscScsiHotPlug     = "disc-scsi-hot-plug"
	FlagDiscScsiHotUnplug   = "disc-scsi-hot-unplug"
	FlagExposeSerial        = "expose-serial"
	FlagRequireLegacyBios   = "require-legacy-bios"
	FlagApplicationType     = "application-type"
	FlagSecAuthProtection   = "sec-auth-protection"
	FlagImageAlias          = "image-alias"
	FlagImageAliasShort     = "a"
	FlagProtocol            = "protocol"
	FlagProtocolShort       = "p"
	FlagSourceSubnet        = "source-subnet"
	FlagTargetSubnet        = "target-subnet"
	FlagSourceMac           = "source-mac"
	FlagSourceIp            = "source-ip"
	FlagDestinationIp       = "destination-ip"
	FlagDestinationIpShort  = "D"
	FlagTargetIp            = "target-ip"
	FlagTargetPort          = "target-port"
	FlagPort                = "port"
	FlagPortShort           = "P"
	FlagWeight              = "weight"
	FlagWeightShort         = "W"
	FlagIcmpCode            = "icmp-code"
	FlagIcmpType            = "icmp-type"
	FlagPortRangeStart      = "port-range-start"
	FlagPortRangeEnd        = "port-range-end"
	FlagLabelUrn            = "label-urn"
	FlagLabelKey            = "label-key"
	FlagLabelValue          = "label-value"
	FlagResourceLimits      = "resource-limits"
	FlagResourceType        = "resource-type"
	FlagForceSecAuth        = "force-secure-auth"
	FlagCreateDc            = "create-dc"
	FlagCreateSnapshot      = "create-snapshot"
	FlagReserveIp           = "reserve-ip"
	FlagAccessLog           = "access-logs"
	ArgS3Privilege          = "s3privilege"
	FlagCreateBackUpUnit    = "create-backup"
	FlagCreatePcc           = "create-pcc"
	FlagCreateNic           = "create-nic"
	ArgCreateK8s            = "create-k8s"
	FlagCreateFlowLog       = "create-flowlog"
	FlagAccessMonitoring    = "access-monitoring"
	FlagAccessCerts         = "access-certs"
	FlagAccessDNS           = "access-dns"
	FlagManageDbaas         = "manage-dbaas"
	FlagManageDataplatform  = "manage-dataplatform"
	FlagManageRegistry      = "manage-registry"
	FlagEditPrivilege       = "edit-privilege"
	FlagSharePrivilege      = "share-privilege"
	ArgS3KeyActive          = "s3key-active"
	ArgK8sVersion           = "k8s-version"
	ArgK8sMinNodeCount      = "min-node-count"
	ArgK8sMaxNodeCount      = "max-node-count"
	ArgK8sMaintenanceDay    = "maintenance-day"
	ArgK8sMaintenanceTime   = "maintenance-time"
	ArgK8sAnnotationKey     = "annotation-key"
	ArgK8sAnnotationValue   = "annotation-value"
	FlagPublicIps           = "public-ips"
	FlagPrivateIps          = "private-ips"
	FlagGatewayIp           = "gateway-ip"
	FlagLatest              = "latest"
	FlagMethod              = "method"
	FlagFilters             = "filters"
	FlagFiltersShort        = "F"
	FlagOrderBy             = "order-by"

	FlagDepth               = "depth"
	FlagDepthShort          = "D"
	FlagCheckTimeout        = "check-timeout"
	FlagPath                = "path"
	FlagMatchType           = "match-type"
	FlagResponse            = "response"
	FlagMessage             = "message"
	FlagMessageShort        = "m"
	FlagRegex               = "regex"
	FlagNegate              = "negate"
	FlagServerCertificates  = "server-certificates"
	FlagQuery               = "query"
	FlagQueryShort          = "Q"
	FlagStatusCode          = "status-code"
	FlagContentType         = "content-type"
	FlagCondition           = "condition"
	FlagConditionShort      = "C"
	FlagConditionType       = "condition-type"
	FlagConditionTypeShort  = "T"
	FlagConditionKey        = "condition-key"
	FlagConditionKeyShort   = "K"
	FlagConditionValue      = "condition-value"
	FlagConditionValueShort = "V"
)

// IDs Flags
const (
	FlagIdShort                   = "i"
	FlagDataCenterId              = "datacenter-id"
	FlagServerId                  = "server-id"
	FlagNatGatewayId              = "natgateway-id"
	FlagApplicationLoadBalancerId = "applicationloadbalancer-id"
	FlagNetworkLoadBalancerId     = "networkloadbalancer-id"
	FlagNicId                     = "nic-id"
	FlagLanId                     = "lan-id"
	FlagLanIds                    = "lan-ids"
	FlagLocationId                = "location-id"
	FlagVolumeId                  = "volume-id"
	FlagLoadBalancerId            = "loadbalancer-id"
	FlagRequestId                 = "request-id"
	FlagSnapshotId                = "snapshot-id"
	FlagImageId                   = "image-id"
	FlagIpBlockId                 = "ipblock-id"
	FlagFirewallRuleId            = "firewallrule-id"
	FlagFlowLogId                 = "flowlog-id"
	FlagUserId                    = "user-id"
	FlagGroupId                   = "group-id"
	FlagResourceId                = "resource-id"
	FlagRuleId                    = "rule-id"
	ArgS3KeyId                    = "s3key-id"
	FlagBackupUnitId              = "backupunit-id"
	FlagPccId                     = "pcc-id"
	ArgK8sNodeId                  = "node-id"
	FlagCdromId                   = "cdrom-id"
	FlagTargetGroupId             = "targetgroup-id"
	FlagTemplateId                = "template-id"
	FlagIPv6CidrBlock             = "ipv6-cidr"
	FlagDHCPv6                    = "dhcpv6"
	FlagIPv6IPs                   = "ipv6-ips"
	FlagIPVersion                 = "ip-version"
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
	FlagDepthDescription   = "Controls the detail depth of the response objects. Max depth is 10."
	FlagOrderByDescription = "Limits results to those containing a matching value for a specific property"
	FlagFiltersDescription = "Limits results to those containing a matching value for a specific property. " +
		"Use the following format to set filters: --filters KEY1=VALUE1,KEY2=VALUE2"
	FlagNoHeadersDescription = "When using text output, don't print headers"
	FlagListAllDescription   = "List all resources without the need of specifying parent ID name."
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
