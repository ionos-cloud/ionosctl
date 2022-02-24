# Changelog

## [6.1.2] (upcoming release)

- enhancement: added `BackupLocation` & `SynchronizationMode` as `--cols` options for `ionosctl dbaas postgres cluster` commands
- enhancement: renamed methods for `ionosctl token` commands according to the new updates
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
