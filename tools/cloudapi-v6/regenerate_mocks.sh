#!/bin/bash

# Regenerate mocks

# For CloudApiV6 Resources
pushd "services/cloudapi-v6/resources" >/dev/null || exit

GO111MODULE=off go get -d github.com/golang/mock/mockgen

mockgen -source datacenter.go >mocks/DataCenterService.go
mockgen -source lan.go >mocks/LanService.go
mockgen -source loadbalancer.go >mocks/LoadBalancerService.go
mockgen -source location.go >mocks/LocationService.go
mockgen -source natgateway.go >mocks/NatGatewayService.go
mockgen -source networkloadbalancer.go >mocks/NetworkLoadBalancerService.go
mockgen -source nic.go >mocks/NicService.go
mockgen -source request.go >mocks/RequestService.go
mockgen -source server.go >mocks/ServerService.go
mockgen -source volume.go >mocks/VolumeService.go
mockgen -source image.go >mocks/ImageService.go
mockgen -source snapshot.go >mocks/SnapshotService.go
mockgen -source ipblock.go >mocks/IpBlockService.go
mockgen -source firewallrule.go >mocks/FirewallRuleService.go
mockgen -source flowlog.go >mocks/FlowLogService.go
mockgen -source label.go >mocks/LabelResourceService.go
mockgen -source contract.go >mocks/ContractService.go
mockgen -source user.go >mocks/UserService.go
mockgen -source group.go >mocks/UserGroupService.go
mockgen -source s3key.go >mocks/S3KeyService.go
mockgen -source backupunit.go >mocks/BackupUnitService.go
mockgen -source pcc.go >mocks/PccService.go
mockgen -source k8s.go >mocks/K8sService.go
mockgen -source template.go >mocks/TemplateService.go
mockgen -source applicationloadbalancer.go >mocks/ApplicationLoadBalancerService.go
mockgen -source targetgroup.go >mocks/TargetGroupService.go

pushd >/dev/null || exit
