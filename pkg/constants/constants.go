package constants

/*
 * Global level constants.
 */

// DBaaS Mongo flags
const (
	FlagClusterId = "cluster-id"
	FlagIdP       = "i"
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
	RequestInfoMessage     = "Request ID: %v Execution Time: %v"
	RequestTimeMessage     = "Request Execution Time: %v"
	StatusDeletingAll      = "Status: Deleting %v with ID: %v..."
	StatusRemovingAll      = "Status: Removing %v with ID: %v..."
	DeleteAllAppendErr     = "error occurred removing %v with ID: %v. error: %w"
	WaitDeleteAllAppendErr = "error occurred waiting on removing %v with ID: %v. error: %w"
)
