# Changelog

## \[6.0.0-beta.6\]

* bug-fix: fixed `login` command to support username and password or token authentication

## \[6.0.0-beta.5\]

* Added `--all` option on delete commands
* Updated SDK-Go version to v6.0.0-beta.6
* Added `--image-alias` option to volume commands
* Removed `--public` and `--gateway-ip` options from k8s cluster commands
* Renamed `--ssh-keys` to `--ssh-key-paths` on volume commands and support uploading SSH Keys from files
* Added BootVolume, `--volume-id` and BootCdrom, `--cdrom-id` to server update command
* Renamed `--target-ip` to `--destination-ip`, `--type` to `--direction` from firewall rule commands
* Updated documentation with usage of boolean flags

## \[6.0.0-beta.4\]

* Added usage message on required flags
* Improved pkg modularization
* Added request time on verbose print
* Fixed [#113](https://github.com/ionos-cloud/ionosctl/issues/113)

## \[6.0.0-beta.3\]

* Added K8s Cluster security improvements
* Renamed `--bucket-name` flag to `--s3bucket` flag
* Added `--verbose` flag
* Updated Cobra version to [v1.2.1](https://github.com/spf13/cobra/releases/tag/v1.2.0), improving completions with descriptions
* Updated Go version to 1.16
* Updated SDK-Go version to v6.0.0-beta.4

## \[6.0.0-beta2\]

* Added Template, FlowLog, NAT Gateway, Network Load Balancer commands
* Updated Server commands to support Server of type CUBE
* Updated Datacenter, Location, Group, Contract, Kubernetes Node Pool Lan properties
* Updated Image, Request commands to support fetching the latest N Images/Requests

