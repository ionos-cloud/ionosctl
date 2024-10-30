package constants

/*
 * Global level constants.
 */

// flags
const (
	FlagDatacenterId    = "datacenter-id"
	FlagSnapshotId      = "snapshot-id"
	FlagIdShort         = "i"
	FlagName            = "name"
	FlagNameShort       = "n"
	FlagTemplate        = "template"
	FlagTemplateId      = "template-id"
	FlagInstances       = "instances"
	FlagShards          = "shards"
	FlagMaintenanceTime = "maintenance-time"
	FlagMaintenanceDay  = "maintenance-day"
	FlagLocation        = "location"
	FlagLocationShort   = "l"
	FlagOffset          = "offset"
	// DescAuthenticationOrder explains auth order. Embed this in any auth-related commands
	DescAuthenticationOrder = `AUTHENTICATION ORDER
ionosctl uses a layered approach for authentication, prioritizing sources in this order:
  1. Global flags
  2. Environment variables
  3. Config file entries
Within each layer, a token takes precedence over a username and password combination. For instance, if a token and a username/password pair are both defined in environment variables, ionosctl will prioritize the token. However, higher layers can override the use of a token from a lower layer. For example, username and password environment variables will supersede a token found in the config file.`
	FlagMaxResults      = "max-results"
	FlagMaxResultsShort = "M"
	FlagCidr            = "cidr"
	FlagIp              = "ip"
	FlagIps             = "ips"
	FlagLanId           = "lan-id"
	FlagEdition         = "edition"

	FlagCores                 = "cores"
	FlagRam                   = "ram"
	FlagAvailabilityZone      = "availability-zone"
	FlagAvailabilityZoneShort = "z"
	FlagCpuFamily             = "cpu-family"
	FlagStorageType           = "storage-type"
	FlagStorageSize           = "storage-size"

	FlagClusterId        = "cluster-id"
	FlagNodepoolId       = "nodepool-id"
	FlagBackupId         = "backup-id"
	FlagNodeCount        = "node-count"
	FlagNodeSubnet       = "node-subnet"
	FlagLabels           = "labels"
	FlagLabelsShort      = "L"
	FlagAnnotations      = "annotations"
	FlagAnnotationsShort = "A"
	FlagVersion          = "version"

	FlagZone          = "zone"
	FlagZoneShort     = "z"
	FlagRecord        = "record"
	FlagRecordShort   = "r"
	FlagState         = "state"
	FlagDescription   = "description"
	FlagEnabled       = "enabled"
	FlagContent       = "content"
	FlagTtl           = "ttl"
	FlagPriority      = "priority"
	FlagType          = "type"
	FlagPrimaryIPs    = "primary-ips"
	FlagZoneFile      = "zone-file"
	FlagSecondaryZone = "secondary-zone"

	FlagCloudInit                       = "cloud-init"
	FlagLoggingPipelineId               = "pipeline-id"
	FlagLoggingPipelineLogTag           = "log-tag"
	FlagLoggingPipelineLogSource        = "log-source"
	FlagLoggingPipelineLogProtocol      = "log-protocol"
	FlagLoggingPipelineLogLabels        = "log-labels"
	FlagLoggingPipelineLogType          = "log-type"
	FlagLoggingPipelineLogRetentionTime = "log-retention-time"

	FlagCDNDistributionFilterDomain        = "domain"
	FlagCDNDistributionFilterState         = "state"
	FlagCDNDistributionID                  = "distribution-id"
	FlagCDNDistributionDomain              = "domain"
	FlagCDNDistributionCertificateID       = "certificate-id"
	FlagCDNDistributionRoutingRules        = "routing-rules"
	FlagCDNDistributionRoutingRulesExample = "routing-rules-example"

	FlagGroupId  = "group-id"
	FlagServerId = "server-id"
	FlagActionId = "action-id"

	FlagRegistryId       = "registry-id"
	FlagRegistryIdShort  = "r"
	FlagArtifactId       = "artifact-id"
	FlagVulnerabilityId  = "vulnerability-id"
	FlagRegistryVulnScan = "vulnerability-scanning"

	FlagDatabase = "database"
	FlagOwner    = "owner"
)

// Flag descriptions. Prefixed with "Desc" for easy find and replace
const (
	DescMaxResults = "The maximum number of elements to return"
	DescZone       = "The name or ID of the DNS zone"
)

// legacy flags. TODO: Arg should be renamed to Flag.
const (
	ArgOutput       = "output"
	ArgOutputShort  = "o"
	ArgQuiet        = "quiet"
	ArgQuietShort   = "q"
	ArgVerbose      = "verbose"
	ArgVerboseShort = "v"
	ArgDepth        = "depth"
	ArgDepthShort   = "D"

	ArgAll                 = "all"
	ArgAllShort            = "a"
	ArgForce               = "force"
	ArgForceShort          = "f"
	ArgWaitForRequest      = "wait-for-request"
	ArgWaitForRequestShort = "w"
	ArgWaitForState        = "wait-for-state"
	ArgWaitForDelete       = "wait-for-deletion"
	ArgWaitForStateShort   = "W"
	ArgTimeout             = "timeout"
	ArgTimeoutShort        = "t"
	ArgCols                = "cols"
	ArgUpdates             = "updates"
	ArgUser                = "user"
	ArgPassword            = "password"
	ArgPasswordShort       = "p"
	ArgNoHeaders           = "no-headers"
)

// Defaults
const (
	DefaultApiURL               = "https://api.ionos.com"
	DefaultDnsApiURL            = "dns.de-fra.ionos.com"
	DefaultLoggingServiceApiURL = "logging.de-txl.ionos.com"
	DefaultConfigFileName       = "/config.json"
	DefaultOutputFormat         = "text"
	DefaultWait                 = false
	DefaultTimeoutSeconds       = int(60)
	DefaultParentIndex          = int(1)
)

// enum values. TODO: ideally i'd like these handled by the SDK
var (
	EnumLicenceType      = []string{"LINUX", "RHEL", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}
	EnumLogProtocols     = []string{"http", "tcp"}
	EnumLogSources       = []string{"docker", "systemd", "generic", "kubernetes"}
	EnumLogRetentionTime = []string{"7", "14", "30"}
)

// Some legacy messages, which might need looking into
const (
	MessageRequestInfo = "Request ID: %v Execution Time: %v"
	MessageRequestTime = "Request Execution Time: %v"
	MessageDeletingAll = "Status: Deleting %v with ID: %v..."
	MessageRemovingAll = "Status: Removing %v with ID: %v..." // TODO: cleanup constant. reduce duplication
)

const (
	ErrDeleteAll     = "error occurred removing %v with ID: %v. error: %w"
	ErrWaitDeleteAll = "error occurred waiting on removing %v with ID: %v. error: %w" // TODO: cleanup constant. reduce duplication
)

// Config
const (
	FlagJsonProperties        = "json-properties"
	FlagJsonPropertiesExample = "json-properties-example"

	ArgConfig         = "config"
	ArgConfigShort    = "c"
	ArgServerUrl      = "api-url"
	ArgServerUrlShort = "u"
	ArgToken          = "token"
	ArgTokenShort     = "t"

	EnvUsername  = "IONOS_USERNAME"
	EnvPassword  = "IONOS_PASSWORD"
	EnvToken     = "IONOS_TOKEN"
	EnvServerUrl = "IONOS_API_URL"

	CfgToken     = "userdata.token"
	CfgServerUrl = "userdata.api-url"
	CfgUsername  = "userdata.name"
	CfgPassword  = "userdata.password"

	CLIHttpUserAgent = "cli-user-agent"

	FlagProvenance      = "provenance"
	FlagProvenanceShort = "p"
	FlagSkipVerify      = "skip-verify"
)

// Manpages
const (
	FlagTargetDir       = "target-dir"
	FlagSkipCompression = "skip-compression"
)

// Resource info
const (
	DatacenterId              = "Datacenter ID: %v"
	ApplicationLoadBalancerId = "Application Load Balancer ID: %v"
	TargetGroupId             = "Target Group ID: %v"
	ClusterId                 = "Cluster ID: %v"
	ForwardingRuleId          = "Forwarding Rule ID: %v"
)
