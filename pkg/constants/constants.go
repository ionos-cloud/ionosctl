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
	FlagMaxResults      = "max-results"
	FlagMaxResultsShort = "M"
	FlagCidr            = "cidr"
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
	FlagNodeCount        = "node-count"
	FlagLabels           = "labels"
	FlagLabelsShort      = "L"
	FlagAnnotations      = "annotations"
	FlagAnnotationsShort = "A"
	FlagVersion          = "version"

	FlagZone        = "zone"
	FlagZoneShort   = "z"
	FlagRecord      = "record"
	FlagRecordShort = "r"
	FlagState       = "state"
	FlagDescription = "description"
	FlagEnabled     = "enabled"
	FlagContent     = "content"
	FlagTtl         = "ttl"
	FlagPriority    = "priority"
	FlagType        = "type"

	FlagCloudInit = "cloud-init"
)

// Flag descriptions. Prefixed with "Desc" for easy find and replace
const (
	DescMaxResults = "The maximum number of elements to return"
	DescZone       = "The name or ID of the DNS zone"
)

// legacy flags. TODO: Arg should be renamed to Flag.
const (
	ArgConfig              = "config"
	ArgConfigShort         = "c"
	ArgOutput              = "output"
	ArgOutputShort         = "o"
	ArgQuiet               = "quiet"
	ArgQuietShort          = "q"
	ArgVerbose             = "verbose"
	ArgVerboseShort        = "v"
	ArgDepth               = "depth"
	ArgDepthShort          = "D"
	ArgServerUrl           = "api-url"
	EnvServerUrl           = "IONOS_API_URL"
	ArgServerUrlShort      = "u"
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
	ArgToken               = "token"
	ArgTokenShort          = "t"
	ArgUser                = "user"
	ArgPassword            = "password"
	ArgPasswordShort       = "p"
	ArgNoHeaders           = "no-headers"
)

// Defaults
const (
	DefaultApiURL         = "https://api.ionos.com"
	DefaultDnsApiURL      = "dns.de-fra.ionos.com"
	DefaultConfigFileName = "/config.json"
	DefaultOutputFormat   = "text"
	DefaultWait           = false
	DefaultTimeoutSeconds = int(60)
	DefaultParentIndex    = int(1)
)

// enum values. TODO: ideally i'd like these handled by the SDK
var (
	EnumLicenceType = []string{"LINUX", "RHEL", "WINDOWS", "WINDOWS2016", "UNKNOWN", "OTHER"}
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

const (
	Username         = "userdata.name"
	Password         = "userdata.password"
	Token            = "userdata.token"
	ServerUrl        = "userdata.api-url"
	CLIHttpUserAgent = "cli-user-agent"
)

// Resource info
const (
	DatacenterId              = "Datacenter ID: %v"
	ApplicationLoadBalancerId = "Application Load Balancer ID: %v"
	TargetGroupId             = "Target Group ID: %v"
	ClusterId                 = "Cluster ID: %v"
	ForwardingRuleId          = "Forwarding Rule ID: %v"
)
