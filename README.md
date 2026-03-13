[![CI](https://github.com/ionos-cloud/ionosctl/workflows/CI/badge.svg)](https://github.com/ionos-cloud/ionosctl/actions)
[![Release](https://img.shields.io/github/v/release/ionos-cloud/ionosctl.svg)](https://github.com/ionos-cloud/ionosctl/releases/latest)
[![Go](https://img.shields.io/github/go-mod/go-version/ionos-cloud/ionosctl.svg)](https://github.com/ionos-cloud/ionosctl)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=cli-ionosctl&metric=alert_status)](https://sonarcloud.io/dashboard?id=cli-ionosctl)
[![Gitter](https://img.shields.io/gitter/room/ionos-cloud/sdk-general)](https://gitter.im/ionos-cloud/sdk-general)

![IONOS Cloud](.github/IONOS.CLOUD.BLU.svg?raw=true "IONOS Cloud")

# ionosctl

The command-line interface for [IONOS Cloud](https://www.ionos.com/enterprise-cloud/signup). Create and manage cloud resources -- virtual machines, networks, storage, databases, Kubernetes clusters, DNS, and more -- directly from your terminal.

[Ionosctl usage overview](https://github.com/user-attachments/assets/78b2b920-70bb-4df1-8144-32b860b7ff70)

> **Prerequisites:** You need an [IONOS Cloud account](https://www.ionos.com/enterprise-cloud/signup) to use ionosctl.

## Table of Contents

- [Quick Start](#quick-start)
- [Installation](#installation)
- [Authentication](#authentication)
- [Usage](#usage)
  - [Interactive Shell](#interactive-shell)
  - [Output Formats](#output-formats)
  - [Filtering & Querying](#filtering--querying)
  - [Waiting for Resources](#waiting-for-resources)
  - [Scripting & Automation](#scripting--automation)
  - [Getting Help](#getting-help)
  - [Global Flags](#global-flags)
- [Shell Auto-Completion](#shell-auto-completion)
- [Configuration](#configuration)
- [Environment Variables](#environment-variables)
- [Advanced Configuration](#advanced-configuration)
- [Documentation](#documentation)
- [Uninstalling](#uninstalling)
- [Contributing](#contributing)

## Quick Start

```bash
# 1. Install (macOS example -- see Installation for all platforms)
brew tap ionos-cloud/homebrew-ionos-cloud && brew install ionosctl

# 2. Authenticate
ionosctl login

# 3. Verify your identity
ionosctl cfg whoami

# 4. Run your first command
ionosctl datacenter list

# 5. Explore available commands
ionosctl --help
```

## Installation

### macOS (Homebrew)

```bash
brew tap ionos-cloud/homebrew-ionos-cloud
brew install ionosctl
```

### Linux (Snap)

```bash
snap install ionosctl
```

### Windows (Scoop)

```bash
scoop bucket add ionos-cloud https://github.com/ionos-cloud/scoop-bucket.git
scoop install ionos-cloud/ionosctl
```

### Binary Download

Download the archive for your OS and architecture from the [Releases](https://github.com/ionos-cloud/ionosctl/releases/latest) page, or use:

```bash
# Download and extract (replace <version> with the full semantic version)
curl -sL https://github.com/ionos-cloud/ionosctl/releases/download/v<version>/ionosctl-<version>-linux-amd64.tar.gz | tar -xzv

# Move to a directory in your PATH
sudo mv ionosctl /usr/local/bin/

# Verify
ionosctl version
```

### Building from Source

Requires [Go](https://go.dev/dl/) (see `go.mod` for minimum version).

```bash
git clone https://github.com/ionos-cloud/ionosctl.git
cd ionosctl
make build    # or: make install
```

Dependencies are managed with [Go Modules](https://github.com/golang/go/wiki/Modules) and vendored.

> **Note:** The development version may contain unreleased changes. For production use, prefer [official releases](https://github.com/ionos-cloud/ionosctl/releases).

## Authentication

ionosctl supports multiple authentication methods. Environment variables take priority over the config file.

### Token Authentication

Token-based authentication is the recommended approach, especially for accounts with two-factor authentication (2FA) enabled:

```bash
# Via environment variable
export IONOS_TOKEN="your-bearer-token"

# Or via login command (persists to config file)
ionosctl login --token "$IONOS_TOKEN"
```

You can create tokens via the [DCD](https://dcd.ionos.com/) or the CLI (`ionosctl token create`).

### Username & Password

For accounts without 2FA, username/password authentication is also supported:

```bash
# Interactive login (prompts for credentials)
ionosctl login

# Or pass credentials directly
ionosctl login --user <username> --password <password>

# Or via environment variables
export IONOS_USERNAME="your-email@example.com"
export IONOS_PASSWORD="your-password"
```

### Verifying Your Identity

Use `whoami` to check who you're logged in as, and `--provenance` to debug the authentication source:

```bash
ionosctl cfg whoami
ionosctl cfg whoami --provenance
```

## Usage

### Command Structure

```
ionosctl [service] [resource] [command] [flags]
```

Examples:

```bash
# List all data centers
ionosctl datacenter list

# Create a server in a data center
ionosctl server create --datacenter-id <dc-id> --name "web-server" --cores 2 --ram 4096

# Get a specific resource as JSON
ionosctl server get --datacenter-id <dc-id> --server-id <server-id> --output json

# Delete with auto-confirm (useful for scripts)
ionosctl server delete --datacenter-id <dc-id> --server-id <server-id> --force

# List Kubernetes clusters
ionosctl compute k8s cluster list

# Create a DNS zone
ionosctl dns zone create --name example.com

# List PostgreSQL clusters
ionosctl dbaas postgres cluster list
```

### Interactive Shell

ionosctl includes a built-in interactive shell with auto-completion, command history, and inline suggestions -- ideal for exploration and ad-hoc management:

```bash
ionosctl shell
```

Inside the shell, you get:
- **Real-time auto-completion** for commands, subcommands, and flags
- **Keyboard shortcuts** (Ctrl+A/E for line start/end, Ctrl+P/N for history, Ctrl+W to cut word, etc.)
- **Persistent flag values** between commands with `--persist-flag-values`

> **Note:** The interactive shell is a BETA feature. Destructive commands (e.g., `delete`) require the `--force` flag instead of interactive confirmation.

### Output Formats

Control how results are displayed with `--output` (`-o`):

| Format | Flag | Description |
|--------|------|-------------|
| Table (default) | `--output text` | Human-readable tabular output |
| JSON | `--output json` | Parsed JSON, suitable for `jq` piping |
| API JSON | `--output api-json` | Raw API response JSON |

```bash
# Default table output
ionosctl datacenter list
# DatacenterId                          Name         Location
# 12345678-abcd-1234-abcd-123456789012  production   de/fra

# JSON output
ionosctl datacenter list --output json
# [{"id": "12345678-...", "properties": {"name": "production", ...}}]

# Select specific columns
ionosctl datacenter list --cols "DatacenterId,Name,Location"

# Hide column headers (useful for scripting)
ionosctl datacenter list --no-headers

# Quiet mode -- suppress all output except errors
ionosctl server delete --datacenter-id <id> --server-id <id> --force --quiet
```

### Filtering & Querying

```bash
# Server-side filtering on list commands
ionosctl datacenter list --filters "name=production"
ionosctl server list --datacenter-id <id> --filters "vmState=RUNNING,cores=4"

# Order results
ionosctl datacenter list --order-by name

# Pagination
ionosctl server list --datacenter-id <id> --limit 10 --offset 20

# JMESPath query for advanced output filtering
ionosctl datacenter list --output json --query "[?properties.location=='de/fra']"

# Control API response depth
ionosctl datacenter get --datacenter-id <id> --depth 3
```

### Waiting for Resources

Many create/update/delete operations return immediately while the resource is being provisioned. ionosctl commands that modify resources typically support a `--wait-for-request` (`-w`) flag to block until the operation completes:

```bash
# Wait for server to be fully provisioned
ionosctl server create --datacenter-id <id> --name "my-server" --cores 2 --ram 4096 --wait-for-request

# Wait for deletion to complete
ionosctl server delete --datacenter-id <id> --server-id <id> --force --wait-for-request
```

### Scripting & Automation

ionosctl is designed to work well in scripts and CI/CD pipelines:

```bash
# Use JSON output + jq for programmatic access
DC_ID=$(ionosctl datacenter list --output json | jq -r '.[0].id')

# Combine --force and --wait-for-request for unattended operations
ionosctl server delete --datacenter-id "$DC_ID" --server-id "$SRV_ID" --force --wait-for-request

# Use --quiet to suppress output in scripts (only errors go to stderr)
ionosctl volume create --datacenter-id "$DC_ID" --name "data" --size 50 --quiet --wait-for-request

# Delete all servers in a datacenter
ionosctl server delete --datacenter-id "$DC_ID" --all --force

# Use --no-headers and --cols for clean parseable output
ionosctl server list --datacenter-id "$DC_ID" --cols ServerId --no-headers
```

### Getting Help

```bash
# Top-level help
ionosctl --help

# Help for a specific command
ionosctl server create --help

# Help for a service group
ionosctl dbaas --help

# Nested help
ionosctl compute k8s cluster --help
```

### Global Flags

These flags are available on all (or most) commands:

| Flag | Short | Description |
|------|-------|-------------|
| `--output` | `-o` | Output format: `text`, `json`, `api-json` |
| `--quiet` | `-q` | Suppress all output except errors |
| `--force` | `-f` | Skip confirmation prompts (for destructive commands) |
| `--all` | `-a` | Target all resources (for delete/remove commands) |
| `--wait-for-request` | `-w` | Block until the API operation completes |
| `--verbose` | `-v` | Increase verbosity (`-v`, `-vv`, `-vvv`) |
| `--no-headers` | | Hide table column headers |
| `--cols` | | Select specific output columns |
| `--filters` | `-F` | Server-side filtering (`KEY=VALUE,...`) |
| `--order-by` | | Sort results by property |
| `--limit` | | Max items per request (default: 50) |
| `--offset` | | Skip N items for pagination |
| `--query` | | JMESPath query to filter output |
| `--depth` | `-D` | API response detail level (default: 1) |
| `--api-url` | `-u` | Override API endpoint |
| `--config` | | Path to config file |

## Shell Auto-Completion

ionosctl supports auto-completion for **Bash**, **Zsh**, **Fish**, and **PowerShell**. Completions include commands, subcommands, flags, and even flag values (like available data center IDs).

### Bash

```bash
# Current session
source <(ionosctl completion bash)

# Permanent (add to ~/.bashrc)
echo 'source <(ionosctl completion bash)' >> ~/.bashrc
```

### Zsh

```bash
# Setup (add to ~/.zshrc BEFORE compinit)
mkdir -p ~/.config/ionosctl/completion/zsh
ionosctl completion zsh > ~/.config/ionosctl/completion/zsh/_ionosctl
```

Add to `~/.zshrc`:

```zsh
fpath+=(~/.config/ionosctl/completion/zsh)
autoload -Uz compinit; compinit
```

### Fish

```bash
ionosctl completion fish > ~/.config/fish/completions/ionosctl.fish
```

### PowerShell

```powershell
# Current session
ionosctl completion powershell | Out-String | Invoke-Expression

# Permanent (add to your PowerShell profile)
ionosctl completion powershell > ionosctl.ps1
# Then source it from your $PROFILE
```

> **Tip:** Use `--no-descriptions` to disable completion descriptions if your shell feels too cluttered.

## Configuration

Credentials and settings are stored in a JSON config file at:

| OS | Path |
|----|------|
| macOS | `~/Library/Application Support/ionosctl/config.json` |
| Linux | `$XDG_CONFIG_HOME/ionosctl/config.json` |
| Windows | `%APPDATA%\ionosctl\config.json` |

Manage your configuration with:

```bash
ionosctl cfg login         # Store credentials
ionosctl cfg logout        # Remove stored credentials
ionosctl cfg whoami        # Show current identity
ionosctl cfg location list # List available IONOS Cloud locations
```

You can also point to a custom config file:

```bash
ionosctl --config /path/to/config.json datacenter list
```

## Environment Variables

| Variable | Description |
|----------|-------------|
| `IONOS_USERNAME` | Username for basic authentication |
| `IONOS_PASSWORD` | Password for basic authentication |
| `IONOS_TOKEN` | Bearer token (takes precedence over username/password) |
| `IONOS_API_URL` | Override the default API endpoint (`api.ionos.com`) |
| `IONOS_LOG_LEVEL` | SDK log level: `Off`, `Debug`, `Trace` |
| `IONOS_PINNED_CERT` | SHA-256 fingerprint for certificate pinning |

> **Warning:** Set `IONOS_LOG_LEVEL=Trace` only for debugging. It logs full request/response payloads including sensitive data, and can significantly impact performance.

## Advanced Configuration

### API Endpoint Override

Override the default `https://api.ionos.com` endpoint for testing or private deployments:

```bash
export IONOS_API_URL="https://custom-api.example.com"
# or per-command:
ionosctl --api-url https://custom-api.example.com datacenter list
```

### Certificate Pinning

Bypass standard certificate validation by pinning a specific SHA-256 fingerprint (useful for self-signed certificates):

```bash
export IONOS_PINNED_CERT="<sha256-public-fingerprint>"
```

You can obtain the SHA-256 fingerprint from your browser's certificate inspector.

### Man Pages (Linux)

Generate man pages for offline reference:

```bash
ionosctl man --target-dir /tmp/ionosctl-man
# Then copy to your man path and run: sudo mandb
```

### Checking for Updates

```bash
ionosctl version --updates
```

## Documentation

| Resource | Link |
|----------|------|
| Full CLI Reference | [docs.ionos.com/cli-reference](https://docs.ionos.com/cli-reference) |
| IONOS Cloud User Guide | [docs.ionos.com/cloud](https://docs.ionos.com/cloud) |
| API Reference | [api.ionos.com/docs](https://api.ionos.com/docs/) |
| Cloud Console (DCD) | [dcd.ionos.com](https://dcd.ionos.com/) |
| Changelog | [CHANGELOG.md](CHANGELOG.md) |

## Uninstalling

```bash
# Homebrew (macOS)
brew uninstall ionosctl

# Snap (Linux)
snap remove ionosctl

# Scoop (Windows)
scoop uninstall ionosctl

# Manual install
sudo rm /usr/local/bin/ionosctl

# Local build
make clean
```

## Contributing

Bugs & feature requests: [Open an issue](https://github.com/ionos-cloud/ionosctl/issues/new/choose)

Pull requests are welcome! Fork the repository, make your changes, and submit a PR. We'll review it and work together to get it released.

### Running Tests

```bash
make test
```
