package config

const (
	Username         = "userdata.name"
	Password         = "userdata.password"
	Token            = "userdata.token"
	ServerUrl        = "userdata.api-url"
	CLIHttpUserAgent = "cli-user-agent"
)

// Some legacy messages, which might need looking into. Too hard to move to pkg/constants
const (
	RequestInfoMessage     = "Request ID: %v Execution Time: %v"
	RequestTimeMessage     = "Request Execution Time: %v"
	StatusDeletingAll      = "Status: Deleting %v with ID: %v..."
	StatusRemovingAll      = "Status: Removing %v with ID: %v..."
	DeleteAllAppendErr     = "error occurred removing %v with ID: %v. error: %w"
	WaitDeleteAllAppendErr = "error occurred waiting on removing %v with ID: %v. error: %w"
)
