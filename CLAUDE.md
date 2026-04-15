# CLAUDE.md

## Git Workflow (MANDATORY)

**Never work directly on `master` or `<base-branch>`.** Always create a feature branch and open a PR.

```bash
# Before starting any work:
git checkout master && git pull # git checkout <base-branch> && git pull
git checkout -b feat/<short-description>   # or fix:/, doc:/, test:/, refactor:/

# When done:
git push -u origin feat/<short-description>
gh pr create --title "feat: <description>" --body "..."
```

- **Never commit directly to `master` or `<base-branch>`** — all changes go through a PR reviewed and merged by a human.
- Branch naming: `feat/<name>`, `fix/<name>`, `doc/<name>`, `test/<name>`, `refactor/<name>`.
- One logical change per branch. Keep PRs focused.

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

`ionosctl` is a Go CLI tool for managing IONOS Cloud resources. It uses Cobra for the CLI framework, Viper for configuration, and multiple IONOS SDK packages for different services (Cloud API v6, DBaaS, DNS, CDN, VPN, Kafka, etc.).

## Common Commands

```bash
make build          # Build the binary (runs tools/build.sh build)
make install        # Install locally
make utest          # Run unit tests with coverage
make itest          # Run integration + unit tests (requires build tag: integration)
make test           # Full test suite: bats-core shell tests + go tests
make lint           # Run golangci-lint
make gofmt_check    # Check code formatting
make gofmt          # Format code
make mocks          # Regenerate mocks (uses golang/mock)
make vendor         # Update vendor dependencies
```
## Testing (MANDATORY)

**YOU ARE NOT ALLOWED TO RUN TESTS WITHOUT MY EXPLICIT APPROVAL.**

Run a single Go test:
```bash
go test ./commands/cloudapi-v6/... -run TestFunctionName -v
go test ./internal/printer/... -run TestFunctionName -v
```

Run tests with integration tag:
```bash
go test -tags integration ./...
```

## Architecture

### Layer Structure

```
main.go → commands/ → internal/core/ → services/ → IONOS SDKs
                    ↘ internal/printer/
                    ↘ internal/client/
```

1. **`commands/`** — Command definitions organized by service (`cloudapi-v6/`, `dns/`, `cdn/`, `dbaas/`, etc.)
2. **`internal/core/`** — Command framework wrapping Cobra (`Command`, `CommandBuilder`, `CommandConfig`)
3. **`services/cloudapi-v6/resources/`** — Service layer wrapping Cloud API v6 SDK calls
4. **`internal/printer/jsontabwriter/`** — Output formatting (text tables, JSON, api-json)
5. **`internal/client/`** — API client initialization, credentials, multi-SDK support
6. **`internal/constants/`** — Flag name constants and shared strings

### Command Structure Pattern

Every command follows a **Namespace.Resource.Verb** hierarchy and is built with `core.CommandBuilder`:

```go
core.NewCommand(ctx, parentCmd, core.CommandBuilder{
    Namespace:  "datacenter",
    Resource:   "datacenter",
    Verb:       "list",
    Aliases:    []string{"l", "ls"},
    ShortDesc:  "List Data Centers",
    PreCmdRun:  core.NoPreRun,       // validation (returns error to abort)
    CmdRun:     RunDataCenterList,   // main logic
    InitClient: true,                // creates API client before CmdRun
})
```

- **`PreCmdRun`** validates flags and preconditions; return an error to abort execution
- **`CmdRun`** receives `*core.CommandConfig` which provides access to services, viper config, and the cobra command

### Adding a New Command

1. Create `commands/{service}/{resource}.go`
2. Define a root function returning `*core.Command` (see `commands/cloudapi-v6/location.go` as a reference)
3. Add subcommands via `core.NewCommand` with a `CommandBuilder`
4. Implement `PreCmdRun` (validation) and `CmdRun` (logic) functions
5. Use `jsontabwriter.GenerateOutput()` for output; define column headers with `tabheaders`
6. Register in `commands/root.go` → `addCommands()` function

### Flag Naming Convention

Flag names are constants in `internal/constants/`. Access flag values in `CmdRun` via:
```go
viper.GetString(core.GetFlagName(c.NS, constants.FlagDatacenterId))
```

### Authentication & Configuration

Credentials are loaded in priority order:
1. Environment variables: `IONOS_USERNAME`, `IONOS_PASSWORD`, `IONOS_TOKEN`, `IONOS_API_URL`
2. Config file: `~/.config/ionosctl/config.json` (Linux)

### Output Formatting

The `--output` global flag controls format: `text` (default table), `json`, or `api-json` (raw API response).
The `--cols` flag controls which columns to display.
The `--filters` flag supports case-insensitive key filtering.

### Services (Cloud API v6)

Services are lazy-loaded in `services/cloudapi-v6/services.go`:
```go
type Services struct {
    Locations   func() resources.LocationsService
    DataCenters func() resources.DatacentersService
    // ...
}
```

Each resource service implements CRUD operations and is accessible via `c.CloudApiV6Services` in `CmdRun`.

### Mocks

Mocks for service interfaces are generated using `golang/mock`. Regenerate with `make mocks`. Mock files live alongside their interfaces, typically in `services/cloudapi-v6/resources/mocks/`.

### Prerequisites
- Remember to run make gofmt after implementing changes
- Remember to run make docs after implementing changes
- ionosctl uses bats for e2e tests, implement e2e tests for new features
- Remember to run gofmt on the files to correctly import tests
- Cyclomatic complexity for new functions should not exceed 15; refactor into smaller functions if necessary. Only for new code.
- Add bats test if adding a new feature
