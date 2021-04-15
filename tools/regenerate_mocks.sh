#!/bin/bash

# regenerate generated mocks

pushd "pkg/resources" >/dev/null || exit

GO111MODULE=off go get -u github.com/golang/mock/mockgen

mockgen -source client.go >mocks/ClientService.go
mockgen -source datacenter.go >mocks/DataCenterService.go
mockgen -source lan.go >mocks/LanService.go
mockgen -source loadbalancer.go >mocks/LoadBalancerService.go
mockgen -source location.go >mocks/LocationService.go
mockgen -source nic.go >mocks/NicService.go
mockgen -source request.go >mocks/RequestService.go
mockgen -source server.go >mocks/ServerService.go
mockgen -source volume.go >mocks/VolumeService.go
mockgen -source image.go >mocks/ImageService.go
mockgen -source snapshot.go >mocks/SnapshotService.go
mockgen -source ipblock.go >mocks/IpBlockService.go
mockgen -source firewallrule.go >mocks/FirewallRuleService.go
mockgen -source label.go >mocks/LabelResourceService.go
mockgen -source contract.go >mocks/ContractsService.go

pushd >/dev/null || exit
