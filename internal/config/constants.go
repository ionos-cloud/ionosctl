package config

const (
	RequestInfoMessage     = "Request ID: %v Execution Time: %v"
	RequestTimeMessage     = "Request Execution Time: %v"
	StatusDeletingAll      = "Status: Deleting %v with ID: %v..."
	StatusRemovingAll      = "Status: Removing %v with ID: %v..."
	DeleteAllAppendErr     = "error occurred removing %v with ID: %v. error: %w"
	WaitDeleteAllAppendErr = "error occurred waiting for request on removing %v with ID: %v. error: %w"

	ArgConfig              = "config"
	ArgConfigShort         = "c"
	ArgOutput              = "output"
	ArgOutputShort         = "o"
	ArgQuiet               = "quiet"
	ArgQuietShort          = "q"
	ArgVerbose             = "verbose"
	ArgVerboseShort        = "v"
	ArgServerUrl           = "api-url"
	ArgServerUrlShort      = "u"
	ArgForce               = "force"
	ArgForceShort          = "f"
	ArgWaitForRequest      = "wait-for-request"
	ArgWaitForRequestShort = "w"
	ArgWaitForState        = "wait-for-state"
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

	DefaultApiURL         = "https://api.ionos.com"
	DefaultConfigFileName = "/config.json"
	DefaultOutputFormat   = "text"
	DefaultWait           = false
	DefaultTimeoutSeconds = int(60)

	Username         = "userdata.name"
	Password         = "userdata.password"
	Token            = "userdata.token"
	ServerUrl        = "userdata.api-url"
	CLIHttpUserAgent = "cli-user-agent"
)
