package config

const (
	RequestInfoMessage     = "Request ID: %v Execution Time: %v"
	RequestTimeMessage     = "Request Execution Time: %v"
	StatusDeletingAll      = "Status: Deleting %v with ID: %v..."
	StatusRemovingAll      = "Status: Removing %v with ID: %v..."
	DeleteAllAppendErr     = "error occurred removing %v with ID: %v. error: %w"
	WaitDeleteAllAppendErr = "error occurred waiting on removing %v with ID: %v. error: %w"

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

	DefaultApiURL         = "https://api.ionos.com"
	DefaultConfigFileName = "/config.json"
	DefaultOutputFormat   = "text"
	DefaultWait           = false
	DefaultTimeoutSeconds = int(60)
	DefaultParentIndex    = int(1)
	DefaultListDepth      = int(1)
	DefaultGetDepth       = int(0)
	DefaultCreateDepth    = int(0)
	DefaultUpdateDepth    = int(0)
	DefaultDeleteDepth    = int(0)
	DefaultMiscDepth      = int(0) // Attach, Detach (and similar); Server start/stop/suspend/etc.;

	Username         = "userdata.name"
	Password         = "userdata.password"
	Token            = "userdata.token"
	ServerUrl        = "userdata.api-url"
	CLIHttpUserAgent = "cli-user-agent"
)
