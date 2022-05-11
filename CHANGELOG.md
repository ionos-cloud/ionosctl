# Changelog

## [6.1.6] (May 2022)

### Features
- added new values for `--storage-type` option for `ionosctl dbaas postgres cluster create` command: **SSD_PREMIUM**, **SSD_STANDARD**.
  - Note: Value **SSD** is deprecated. Use the equivalent **SSD_PREMIUM** instead.
- added option to do certificate pinning by using `IONOS_PINNED_CERT` environment variable for the `ionosctl dbaas postgres` commands.
  - Note: Set the `IONOS_PINNED_CERT` environment variable to be the public sha256 fingerprint of the certificate to be pinned.

### Dependency-update
- updated SDK-Go-DBaaS Postgres version from [v1.0.2](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.2) to [v1.0.3](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.3)

## [6.1.5] (April 2022)

### Fixes
- added `--labels`,`--annotations` options for `ionosctl k8s nodepool create` command to set one or multiple labels, annotations
- added `--labels`,`--annotations` options for `ionosctl k8s nodepool update` command to set one or multiple labels, annotations (fixes [164](https://github.com/ionos-cloud/ionosctl/issues/164))
- added `Annotations,Labels` values for `--cols` option for `ionosctl k8s nodepool` commands 

### Deprecated
- marked as deprecated the following options: `--label-key`,`--label-value`,`--annotation-key`,`--annotation-value` for `ionosctl k8s nodepool update` command (use `--labels`,`--annotations` options instead).

## [6.1.4] (April 2022)

- enhancement: added `--backup-location` option for `ionosctl dbaas postgres cluster create` command
- enhancement: added `--direction` option for `ionosctl dbaas postgres logs list` command
- enhancement: added `--since` and `--until` option for `ionosctl dbaas postgres logs list` command, to easily specify timeframe for getting logs
- enhancement: added `--public` option for `ionosctl k8s cluster create` command
- enhancement: added `--gateway-ip` option for `ionosctl k8s nodepool create` command
- enhancement: added `BootServerId` value for `--cols` option for `ionosctl volume` commands
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.2](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.2)
- dependency-update: added SDK-Go version [v6.0.2](https://github.com/ionos-cloud/sdk-go/releases/tag/v6.0.2)
- dependency-update: updated Go version from 1.16 to 1.17

## [6.1.3]

- enhancement: added `--no-headers` option for GET commands, for text output (PR #158)

## [6.1.2]

- enhancement: added `SynchronizationMode` as `--cols` option for dbaas postgres cluster commands
- enhancement: renamed methods for token commands according to the new updates from [v1.0.2](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.2)
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.1](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.1)
- dependency-update: added SDK-Go-Auth version [v1.0.3](https://github.com/ionos-cloud/sdk-go-auth/releases/tag/v1.0.3)

## [6.1.1]

- bug-fix: `ionosctl k8s cluster` command supports now `--cols` option

## [6.1.0]

- new service: **Database as a Service (DBaaS) - Postgres**
  - added the CLI commands for DBaaS Postgres under `dbaas postgres` namespace (PR #155):
    - `ionosctl dbaas postgres cluster` 
    - `ionosctl dbaas postgres logs`
    - `ionosctl dbaas postgres backup`
    - `ionosctl dbaas postgres version`
    - `ionosctl dbaas postgres api-version`
- dependency-update: added SDK-Go-DBaaS Postgres version [v1.0.0](https://github.com/ionos-cloud/sdk-go-dbaas-postgres/releases/tag/v1.0.0)

## [6.0.2]

- enhancement: added `--no-headers` option for list commands, for text output (PR #153)
- documentation: separate documentation per service (PR #148)

## [6.0.1]

- enhancement: `--all` option iterates through all resources even if it hits error
- enhancement: improved logs messages on `--all` option and request info and request info

## [6.0.0]

- feature: added `--password` on `ionosctl user update` command
- feature: updated code for `ionosctl k8s nodepool` commands in sync with the changes from SDK Go
- bug-fix: `ionosctl lan create` command supports now `--cols` option
- dependency-update: added SDK-Go version `v6.0.0-beta.9` to `v6.0.0`

## [6.0.0-beta.8]

- feature: added `token` commands, added support for Auth API, to generate, list, delete Tokens
- dependency-update: added SDK-Go-Auth version v1.0.1

## [6.0.0-beta.7]

- feature: added `--filters`, `--max-results`, `--order-by` options on all list commands
- feature: added `-all` option for remove and detach commands
- enhancement: added completion support for `--k8s-version` option
- dependency-update: SDK-Go version from v6.0.0-beta.6 to v6.0.0-beta.9

## [6.0.0-beta.6]

- bug-fix: fixed `login` command to support username and password or token authentication

## [6.0.0-beta.5]

- Added `--all` option on delete commands
- Updated SDK-Go version to v6.0.0-beta.6
- Added `--image-alias` option to volume commands
- Removed `--public` and `--gateway-ip` options from k8s cluster commands
- Renamed `--ssh-keys` to `--ssh-key-paths` on volume commands and support uploading SSH Keys from files
- Added BootVolume, `--volume-id` and BootCdrom, `--cdrom-id` to server update command
- Renamed `--target-ip` to `--destination-ip`, `--type` to `--direction` from firewall rule commands
- Updated documentation with usage of boolean flags

## [6.0.0-beta.4]

- Added usage message on required flags
- Improved pkg modularization
- Added request time on verbose print
- Fixed [#113](https://github.com/ionos-cloud/ionosctl/issues/113)

## [6.0.0-beta.3]

- Added K8s Cluster security improvements
- Renamed `--bucket-name` flag to `--s3bucket` flag
- Added `--verbose` flag
- Updated Cobra version to [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.0), improving completions with descriptions
- Updated Go version to 1.16
- Updated SDK-Go version to v6.0.0-beta.4

## [6.0.0-beta2]

- Added Template, FlowLog, NAT Gateway, Network Load Balancer commands
- Updated Server commands to support Server of type CUBE
- Updated Datacenter, Location, Group, Contract, Kubernetes Node Pool Lan properties
- Updated Image, Request commands to support fetching the latest N Images/Requests
