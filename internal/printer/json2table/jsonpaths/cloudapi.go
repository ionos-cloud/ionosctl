package jsonpaths

// CloudAPI json paths
var (
	ApplicationLoadBalancer = map[string]string{
		"ApplicationLoadBalancerId": "id",
		"Name":                      "properties.name",
		"ListenerLan":               "properties.listenerLan",
		"Ips":                       "properties.ips",
		"TargetLan":                 "properties.targetLan",
		"PrivateIps":                "properties.lbPrivateIps",
		"State":                     "metadata.state",
	}

	ApplicationLoadBalancerForwardingRule = map[string]string{
		"ForwardingRuleId":   "id",
		"Name":               "properties.name",
		"Protocol":           "properties.protocol",
		"ListenerIp":         "properties.listenerIp",
		"ListenerPort":       "properties.listenerPort",
		"ClientTimeout":      "properties.clientTimeout",
		"ServerCertificates": "properties.serverCertificates",
		"State":              "metadata.state",
	}

	ApplicationLoadBalancerHTTPRule = map[string]string{
		"Name":            "name",
		"Type":            "type",
		"TargetGroupId":   "targetGroup",
		"DropQuery":       "dropQuery",
		"Location":        "location",
		"StatusCode":      "statusCode",
		"ResponseMessage": "responseMessage",
		"ContentType":     "contentType",
		"Condition":       "conditions",
	}

	BackupUnit = map[string]string{
		"BackupUnitId": "id",
		"Name":         "properties.name",
		"Email":        "properties.email",
		"State":        "metadata.state",
	}

	BackupUnitSSOUrl = map[string]string{
		"BackupUnitSsoUrl": "ssoUrl",
	}

	Console = map[string]string{
		"RemoteConsoleUrl": "url",
	}

	Contract = map[string]string{
		"ContractNumber":         "properties.contractNumber",
		"Owner":                  "properties.owner",
		"Status":                 "properties.status",
		"RegistrationDomain":     "properties.regDomain",
		"CoresPerServer":         "properties.resourceLimits.coresPerServer",
		"CoresPerContract":       "properties.resourceLimits.coresPerContract",
		"CoresProvisioned":       "properties.resourceLimits.coresProvisioned",
		"RamPerServer":           "properties.resourceLimits.ramPerServer",
		"RamPerContract":         "properties.resourceLimits.ramPerContract",
		"RamProvisioned":         "properties.resourceLimits.ramProvisioned",
		"HddLimitPerVolume":      "properties.resourceLimits.hddLimitPerVolume",
		"HddLimitPerContract":    "properties.resourceLimits.hddLimitPerContract",
		"HddVolumeProvisioned":   "properties.resourceLimits.hddVolumeProvisioned",
		"SsdLimitPerVolume":      "properties.resourceLimits.ssdLimitPerVolume",
		"SsdLimitPerContract":    "properties.resourceLimits.ssdLimitPerContract",
		"SsdVolumeProvisioned":   "properties.resourceLimits.ssdVolumeProvisioned",
		"DasVolumeProvisioned":   "properties.resourceLimits.dasVolumeProvisioned",
		"ReservableIps":          "properties.resourceLimits.reservableIps",
		"ReservedIpsOnContract":  "properties.resourceLimits.reservedIpsOnContract",
		"ReservedIpsInUse":       "properties.resourceLimits.reserverIpsInUse",
		"K8sClusterLimitTotal":   "k8sClusterLimitTotal",
		"K8sClustersProvisioned": "k8sClustersProvisioned",
		"NlbLimitTotal":          "properties.resourceLimits.nlbLimitTotal",
		"NlbProvisioned":         "properties.resourceLimits.nlbProvisioned",
		"NatGatewayLimitTotal":   "properties.resourceLimits.natGatewayLimitTotal",
		"NatGatewayProvisioned":  "properties.resourceLimits.natGatewayProvisioned",
	}

	Cpu = map[string]string{
		"CpuFamily": "cpuFamily",
		"MaxCores":  "maxCores",
		"MaxRam":    "maxRam",
		"Vendor":    "vendor",
	}

	Datacenter = map[string]string{
		"DatacenterId":      "id",
		"Name":              "properties.name",
		"Location":          "properties.location",
		"Description":       "properties.description",
		"Version":           "properties.version",
		"State":             "metadata.state",
		"Features":          "properties.features",
		"CpuFamily":         "properties.cpuArchitecture.*.cpuFamily",
		"SecAuthProtection": "properties.secAuthProtection",
		"IPv6CidrBlock":     "properties.ipv6CidrBlock",
	}

	FirewallRule = map[string]string{
		"FirewallRuleId": "id",
		"Name":           "properties.name",
		"Protocol":       "properties.protocol",
		"SourceMac":      "properties.sourceMac",
		"SourceIP":       "properties.sourceIp",
		"DestinationIP":  "properties.destinationIp",
		"PortRangeStart": "properties.portRangeStart",
		"PortRangeEnd":   "properties.portRangeEnd",
		"IcmpCode":       "properties.icmpCode",
		"IcmpType":       "properties.icmpType",
		"Direction":      "properties.type",
		"IPVersion":      "properties.ipVersion",
		"State":          "metadata.state",
	}

	Flowlog = map[string]string{
		"FlowLogId": "id",
		"Name":      "properties.name",
		"Action":    "properties.action",
		"Direction": "properties.direction",
		"Bucket":    "properties.bucket",
		"State":     "metadata.state",
	}

	Group = map[string]string{
		"GroupId":                     "id",
		"Name":                        "properties.name",
		"CreateDataCenter":            "properties.createDataCenter",
		"CreateSnapshot":              "properties.createSnapshot",
		"ReserveIp":                   "properties.reserveIp",
		"AccessActivityLog":           "properties.accessActivityLog",
		"CreatePcc":                   "properties.createPcc",
		"S3Privilege":                 "properties.s3Privilege",
		"CreateBackupUnit":            "properties.createBackupUnit",
		"CreateInternetAccess":        "properties.createInternetAccess",
		"CreateK8s":                   "properties.createK8sCluster",
		"CreateFlowLog":               "properties.createFlowLog",
		"AccessAndManageMonitoring":   "properties.accessAndManageMonitoring",
		"AccessAndManageCertificates": "properties.accessAndManageCertificates",
		"AccessAndManageDns":          "properties.accessAndManageDns",
		"ManageDBaaS":                 "properties.manageDBaaS",
		"ManageRegistry":              "properties.manageRegistry",
		"ManageDataplatform":          "properties.manageDataplatform",
	}

	Image = map[string]string{
		"ImageId":         "id",
		"Name":            "properties.name",
		"Description":     "properties.description",
		"Location":        "properties.location",
		"Size":            "properties.size",
		"LicenceType":     "properties.licenceType",
		"ImageType":       "properties.imageType",
		"Public":          "properties.public",
		"ImageAliases":    "properties.imageAliases",
		"CloudInit":       "properties.cloudInit",
		"CreatedBy":       "metadata.createdBy",
		"CreatedByUserId": "metadata.createdByUserId",
		"CreatedDate":     "metadata.createdDate",
	}

	IpBlock = map[string]string{
		"IpBlockId": "id",
		"Name":      "properties.name",
		"Location":  "properties.location",
		"Size":      "properties.size",
		"Ips":       "properties.ips",
		"State":     "metadata.state",
	}

	IpConsumer = map[string]string{
		"Ip":             "ip",
		"Mac":            "mac",
		"NicId":          "nicId",
		"ServerId":       "serverId",
		"ServerName":     "serverName",
		"DatacenterId":   "datacenterId",
		"DatacenterName": "datacenterName",
		"K8sNodePoolId":  "k8sNodePoolUuid",
		"K8sClusterId":   "k8sClusterUuid",
	}

	IpFailover = map[string]string{
		"NicId": "nicUuid",
		"Ip":    "ip",
	}

	K8sCluster = map[string]string{
		"ClusterId":                "id",
		"Name":                     "properties.name",
		"K8sVersion":               "properties.k8sVersion",
		"AvailableUpgradeVersions": "properties.availableUpgradeVersions",
		"ViableNodePoolVersions":   "properties.viableNodePoolVersions",
		"State":                    "metadata.state",
		"S3Bucket":                 "properties.s3Buckets",
		"ApiSubnetAllowList":       "properties.apiSubnetAllowList",
		"Public":                   "properties.public",
		"Location":                 "properties.location",
		"NatGatewayIp":             "properties.natGatewayIp",
		"NodeSubnet":               "properties.nodeSubnet",
	}

	K8sNode = map[string]string{
		"NodeId":     "id",
		"Name":       "properties.name",
		"K8sVersion": "properties.k8sVersion",
		"PublicIP":   "properties.publicIP",
		"PrivateIP":  "properties.privateIP",
		"State":      "metadata.state",
	}

	K8sNodepool = map[string]string{
		"NodePoolId":               "id",
		"Name":                     "properties.name",
		"K8sVersion":               "properties.k8sVersion",
		"DatacenterId":             "properties.datacenterId",
		"NodeCount":                "properties.nodeCount",
		"CpuFamily":                "properties.cpuFamily",
		"StorageType":              "properties.storageType",
		"State":                    "metadata.state",
		"LanIds":                   "properties.lan.*.id",
		"CoresCount":               "properties.coresCount",
		"RamSize":                  "properties.ramSize",
		"AvailabilityZone":         "properties.availabilityZone",
		"StorageSize":              "properties.storageSize",
		"AutoScaling":              "properties.autoScaling",
		"PublicIps":                "properties.publicIps",
		"AvailableUpgradeVersions": "properties.availableUpgradeVersions",
		"Annotations":              "properties.annotations",
		"Labels":                   "properties.labels",
	}

	K8sNodePoolLan = map[string]string{
		"LanId":           "id",
		"Dhcp":            "dhcp",
		"RoutesNetwork":   "routes.*.network",
		"RoutesGatewayIp": "routes.*.gatewayIp",
	}

	Lan = map[string]string{
		"LanId":         "id",
		"Name":          "properties.name",
		"Public":        "properties.public",
		"PccId":         "properties.pcc",
		"State":         "metadata.state",
		"IPv6CidrBlock": "properties.ipv6CidrBlock",
	}

	LoadBalancer = map[string]string{
		"LoadBalancerId": "id",
		"Name":           "properties.name",
		"Dhcp":           "properties.dhcp",
		"Ip":             "properties.ip",
		"State":          "metadata.state",
	}

	Location = map[string]string{
		"LocationId":   "id",
		"Name":         "properties.name",
		"Features":     "properties.features",
		"CpuFamily":    "properties.cpuArchitecture.*.cpuFamily",
		"ImageAliases": "properties.imageAliases",
	}

	NatGateway = map[string]string{
		"NatGatewayId": "id",
		"Name":         "properties.name",
		"PublicIps":    "properties.publicIps",
		"State":        "metadata.state",
	}

	NatGatewayLan = map[string]string{
		"NatGatewayLanId": "id",
		"GatewayIps":      "gatewayIps",
	}

	NatGatewayRule = map[string]string{
		"NatGatewayRuleId":     "id",
		"Name":                 "properties.name",
		"Type":                 "properties.type",
		"Protocol":             "properties.protocol",
		"SourceSubnet":         "properties.sourceSubnet",
		"PublicIp":             "properties.publicIp",
		"TargetSubnet":         "properties.targetSubnet",
		"TargetPortRangeStart": "properties.targetPortRange.start",
		"TargetPortRangeEnd":   "properties.targetPortRange.end",
		"State":                "metadata.state",
	}

	NetworkLoadBalancer = map[string]string{
		"NetworkLoadBalancerId": "id",
		"Name":                  "properties.name",
		"ListenerLan":           "properties.listenerLan",
		"Ips":                   "properties.ips",
		"TargetLan":             "properties.targetLan",
		"LbPrivateIps":          "properties.lbPrivateIps",
		"State":                 "metadata.state",
	}

	NetworkLoadBalancerRule = map[string]string{
		"ForwardingRuleId": "id",
		"Name":             "properties.name",
		"Algorithm":        "properties.algorithm",
		"Protocol":         "properties.protocol",
		"ListenerIp":       "properties.listenerIp",
		"ListenerPort":     "properties.listenerPort",
		"ClientTimeout":    "properties.healthCheck.clientTimeout",
		"ConnectTimeout":   "properties.healthCheck.connectTimeout",
		"TargetTimeout":    "properties.healthCheck.targetTimeout",
		"Retries":          "properties.healthCheck.retries",
		"State":            "metadata.state",
	}

	NetworkLoadBalancerRuleTarget = map[string]string{
		"TargetIp":      "ip",
		"TargetPort":    "port",
		"Weight":        "weight",
		"CheckInterval": "healthCheck.checkInterval",
		"Check":         "healthCheck.check",
		"Maintenance":   "healthCheck.maintenance",
	}

	Nic = map[string]string{
		"NicId":          "id",
		"Name":           "properties.name",
		"Dhcp":           "properties.dhcp",
		"LanId":          "properties.lan",
		"Ips":            "properties.ips",
		"FirewallActive": "properties.firewallActive",
		"FirewallType":   "properties.firewallType",
		"Mac":            "properties.mac",
		"State":          "metadata.state",
		"DeviceNumber":   "properties.deviceNumber",
		"PciSlot":        "properties.pciSlot",
		"IPv6Ips":        "properties.ipv6Ips",
		"IPv6CidrBlock":  "properties.ipv6CidrBlock",
		"DHCPv6":         "properties.dhcpv6",
	}

	Request = map[string]string{
		"RequestId":   "id",
		"Status":      "metadata.requestStatus.metadata.status",
		"Message":     "metadata.requestStatus.metadata.message",
		"Method":      "properties.method",
		"Url":         "properties.url",
		"Body":        "properties.body",
		"CreatedBy":   "metadata.createdBy",
		"CreatedDate": "metadata.createdDate",
	}

	Resource = map[string]string{
		"ResourceId":        "id",
		"Name":              "properties.name",
		"SecAuthProtection": "properties.secAuthProtection",
		"Type":              "type",
		"State":             "metadata.state",
	}

	S3Key = map[string]string{
		"S3KeyId":   "id",
		"Active":    "properties.active",
		"SecretKey": "properties.secretKey",
	}

	Server = map[string]string{
		"ServerId":         "id",
		"Name":             "properties.name",
		"AvailabilityZone": "properties.availabilityZone",
		"State":            "metadata.state",
		"Cores":            "properties.cores",
		"Ram":              "properties.ram",
		"CpuFamily":        "properties.cpuFamily",
		"VmState":          "properties.vmState",
		"BootVolumeId":     "properties.bootVolume.id",
		"BootCdromId":      "properties.bootCdrom.id",
		"TemplateId":       "properties.templateUuid",
		"Type":             "properties.type",
	}

	Share = map[string]string{
		"ShareId":        "id",
		"EditPrivilege":  "properties.editPrivilege",
		"SharePrivilege": "properties.sharePrivilege",
		"Type":           "type",
	}

	Snapshot = map[string]string{
		"SnapshotId":  "id",
		"Name":        "properties.name",
		"LicenceType": "properties.licenseType",
		"Size":        "properties.size",
		"State":       "metadata.state",
	}

	TargetGroup = map[string]string{
		"TargetGroupId": "id",
		"Name":          "properties.name",
		"Algorithm":     "properties.algorithm",
		"Protocol":      "properties.protocol",
		"CheckTimeout":  "properties.healthCheck.timeout",
		"CheckInterval": "properties.healthCheck.interval",
		"Retries":       "properties.healthCheck.retries",
		"Path":          "properties.httpHealthCheck.path",
		"Method":        "properties.httpHealthCheck.method",
		"MatchType":     "properties.httpHealthCheck.matchType",
		"Response":      "properties.httpHealthCheck.response",
		"Regex":         "properties.httpHealthCheck.regex",
		"Negate":        "properties.httpHealthCheck.negate",
		"State":         "metadata.state",
	}

	TargetGroupTarget = map[string]string{
		"TargetIp":           "ip",
		"TargetPort":         "port",
		"Weight":             "weight",
		"HealthCheckEnabled": "healthCheckEnabled",
		"MaintenanceEnabled": "maintenanceEnabled",
	}

	Template = map[string]string{
		"TemplateId":  "id",
		"Name":        "properties.name",
		"Cores":       "properties.cores",
		"Ram":         "properties.ram",
		"StorageSize": "properties.storageSize",
	}

	Token = map[string]string{
		"Token": "token",
	}

	User = map[string]string{
		"UserId":            "id",
		"Firstname":         "properties.firstName",
		"Lastname":          "properties.lastName",
		"Email":             "properties.email",
		"Administrator":     "properties.administrator",
		"ForceSecAuth":      "properties.forceSecAuth",
		"SecAuthActive":     "properties.secAuthActive",
		"S3CanonicalUserId": "properties.s3CanonicalUserId",
		"Active":            "properties.active",
	}

	Volume = map[string]string{
		"VolumeId":         "id",
		"Name":             "properties.name",
		"Size":             "properties.size",
		"Type":             "properties.type",
		"LicenceType":      "properties.licenceType",
		"Bus":              "properties.bus",
		"AvailabilityZone": "properties.availabilityZone",
		"State":            "metadata.state",
		"Image":            "properties.image",
		"DeviceNumber":     "properties.deviceNumber",
		"BackupunitId":     "properties.backupunitId",
		"UserData":         "properties.userData",
		"BootServerId":     "properties.bootServer",
	}

	Label = map[string]string{
		"URN":          "id",
		"Key":          "properties.key",
		"Value":        "properties.value",
		"ResourceType": "properties.resourceType",
		"ResourceId":   "properties.resourceId",
	}
)
