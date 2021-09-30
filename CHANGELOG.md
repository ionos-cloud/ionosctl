# Changelog

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
- Update Go version to 1.16
- Update SDK-Go version to v6.0.0-beta.4

## [6.0.0-beta2]

- Added Template, FlowLog, NAT Gateway, Network Load Balancer commands
- Updated Server commands to support Server of type CUBE
- Updated Datacenter, Location, Group, Contract, Kubernetes Node Pool Lan properties
- Updated Image, Request commands to support fetching the latest N Images/Requests
