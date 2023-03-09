package constants

/*
 * Global level constants.
 */

// flags
const (
	FlagClusterId       = "cluster-id"
	FlagSnapshotId      = "snapshot-id"
	FlagIdShort         = "i"
	FlagName            = "name"
	FlagNameShort       = "n"
	FlagTemplateId      = "template-id"
	FlagMongoVersion    = "mongo-version"
	FlagInstances       = "instances"
	FlagMaintenanceTime = "maintenance-time"
	FlagMaintenanceDay  = "maintenance-day"
	FlagLocation        = "location"
	FlagLocationShort   = "l"
	FlagOffset          = "offset"
	FlagMaxResults      = "max-results"
	FlagMaxResultsShort = "M"
	FlagDatacenterId    = "datacenter-id"
	FlagCidr            = "cidr"
	FlagCidrShort       = "c"
	FlagLanId           = "lan-id"
)

// Flag descriptions. Prefixed with "Desc" for easy find and replace
const (
	DescMaxResults = "The maximum number of elements to return"
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
	DefaultConfigFileName = "/config.json"
	DefaultOutputFormat   = "text"
	DefaultWait           = false
	DefaultTimeoutSeconds = int(60)
	DefaultParentIndex    = int(1)
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
