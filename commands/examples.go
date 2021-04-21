package commands

const (
	/*
		Login Examples
	*/
	loginExamples = `ionosctl login --user USERNAME --password PASSWORD
Status: Authentication successful!

ionosctl login 
Enter your username:
USERNAME
Enter your password:

Status: Authentication successful!`
	/*
		Location Examples
	*/
	listLocationExample = `ionosctl location list 
LocationId   Name        Features
de/fra       frankfurt   [SSD]
us/las       lasvegas    [SSD]
us/ewr       newark      [SSD]
de/txl       berlin      [SSD]
gb/lhr       london      [SSD]`
	/*
		Data Center Examples
	*/
	listDatacenterExample = `ionosctl datacenter list 
DatacenterId                           Name             Location
ff279ffd-ac61-4e5d-ba5e-058296c77774   demoDatacenter   us/las

ionosctl datacenter list --cols "DatacenterId,Name,Location,Version"
DatacenterId                           Name             Location   Version
ff279ffd-ac61-4e5d-ba5e-058296c77774   demoDatacenter   us/las     1`
	getDatacenterExample = `ionosctl datacenter get --datacenter-id ff279ffd-ac61-4e5d-ba5e-058296c77774
DatacenterId                           Name             Location
ff279ffd-ac61-4e5d-ba5e-058296c77774   demoDatacenter   us/las`
	createDatacenterExample = `ionosctl datacenter create --datacenter-name demoDatacenter --datacenter-location us/las
DatacenterId                           Name             Location
f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d   demoDatacenter   us/las
RequestId: 98ab8148-96c4-4091-90e8-9ee2b8a172f4
Status: Command datacenter create has been successfully executed

ionosctl datacenter create --datacenter-name demoDatacenter --datacenter-location gb/lhr --wait 
Waiting for request: 2401b498-8afb-4728-a22a-d2b26f5e31c3
DatacenterId                           Name             Location
8e543958-04f5-4872-bbf3-b28d46393ac7   demoDatacenter   gb/lhr
RequestId: 2401b498-8afb-4728-a22a-d2b26f5e31c3
Status: Command datacenter create and request have been successfully executed`
	updateDatacenterExample = `ionosctl datacenter update --datacenter-id 8e543958-04f5-4872-bbf3-b28d46393ac7 --datacenter-description demoDescription --cols "DatacenterId,Description"
DatacenterId                           Description
8e543958-04f5-4872-bbf3-b28d46393ac7   demoDescription
RequestId: 46af6915-9003-4f11-a1fe-bab1eac9bccc
Status: Command datacenter update has been successfully executed`
	deleteDatacenterExample = `ionosctl datacenter delete --datacenter-id 8e543958-04f5-4872-bbf3-b28d46393ac7
Warning: Are you sure you want to delete data center (y/N) ? y
RequestId: 12547a71-9768-483b-8a8e-e03e58df6dc3
Status: Command datacenter delete has been successfully executed

ionosctl datacenter delete --datacenter-id ff279ffd-ac61-4e5d-ba5e-058296c77774 --ignore-stdin --wait 
Waiting for request: a2f71ef3-f81c-4b15-8f8f-5dfd1bdb3c26
RequestId: a2f71ef3-f81c-4b15-8f8f-5dfd1bdb3c26
Status: Command datacenter delete and request have been successfully executed`

	/*
		Server Examples
	*/
	listServerExample = `ionosctl server list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d 
ServerId                               Name         AvailabilityZone   State       Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               AVAILABLE   4       256MB   AMD_OPTERON`
	getServerExample = `ionosctl server get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff 
ServerId                               Name         AvailabilityZone   State       Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               AVAILABLE   4       256MB   AMD_OPTERON`
	createServerExample = `ionosctl server create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-name demoServer
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    2       256MB   AMD_OPTERON
RequestId: 07fd3682-8642-4a5e-a57a-056e909a2af8
Status: Command server create has been successfully executed

ionosctl server create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-name demoServer --wait 
Waiting for request: e9d12f57-3513-4ae3-ab39-179aacb8c072
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
35201d04-0ea2-43e7-abc4-56f92737bb9d   demoServer                      BUSY    2       256MB   AMD_OPTERON
RequestId: e9d12f57-3513-4ae3-ab39-179aacb8c072
Status: Command server create and request have been successfully executed`
	updateServerExample = `ionosctl server update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff --server-cores 4
ServerId                               Name         AvailabilityZone   State   Cores   Ram     CpuFamily
f45f435e-8d6c-4170-ab90-858b59dab9ff   demoServer   AUTO               BUSY    4       256MB   AMD_OPTERON
RequestId: 571a1bbb-26b3-449d-9885-a20e50dc3b95
Status: Command server update has been successfully executed`
	deleteServerExample = `ionosctl server delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id f45f435e-8d6c-4170-ab90-858b59dab9ff 
Warning: Are you sure you want to delete server (y/N) ? Y
RequestId: 1f00c6d9-072a-4dd0-8c09-c46f2f20a230
Status: Command server delete has been successfully executed

ionosctl server delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 35201d04-0ea2-43e7-abc4-56f92737bb9d --ignore-stdin 
RequestId: f596caba-78b7-4c99-8c9d-56198d3754b6
Status: Command server delete has been successfully executed`
	startServerExample = `ionosctl server start --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
Warning: Are you sure you want to start server (y/N) ? y
RequestId: 9f03a764-5f6c-4740-87e2-d9e9589265dc
Status: Command server start has been successfully executed`
	stopServerExample = `ionosctl server stop --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
Warning: Are you sure you want to stop server (y/N) ? y
RequestId: 8c06523d-8838-4409-aee3-68c042f5a256
Status: Command server stop has been successfully executed`
	resetServerExample = `ionosctl server reset --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa
Warning: Are you sure you want to reboot server (y/N) ? y
RequestId: e6720605-2fa4-46d9-be74-42b733eb1128
Status: Command server reset has been successfully executed`
	attachVolumeServerExample = `ionosctl server attach-volume --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --volume-id 1ceb4b02-ed41-4651-a90b-9a30bc284e74 
VolumeId                               Name   Size   Type   LicenceType   State   Image
1ceb4b02-ed41-4651-a90b-9a30bc284e74   test   10GB   HDD    LINUX         BUSY    
RequestId: b1d06478-b3d5-4069-bb75-86b56422698d
Status: Command server attach-volume has been successfully executed`
	listVolumesServerExample = `ionosctl server list-volumes --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 
VolumeId                               Name   Size   Type   LicenceType   State       Image
1ceb4b02-ed41-4651-a90b-9a30bc284e74   test   10GB   HDD    LINUX         AVAILABLE`
	getVolumeServerExample = `ionosctl server get-volume --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --volume-id 1ceb4b02-ed41-4651-a90b-9a30bc284e74 
VolumeId                               Name   Size   Type   LicenceType   State       Image
1ceb4b02-ed41-4651-a90b-9a30bc284e74   test   10GB   HDD    LINUX         AVAILABLE`
	detachVolumeServerExample = `ionosctl server detach-volume --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --volume-id 1ceb4b02-ed41-4651-a90b-9a30bc284e74 
Warning: Are you sure you want to detach volume from server (y/N) ? 
y
RequestId: 0fd9d6eb-25a1-496c-b0c9-bbe18a989f18
Status: Command server detach-volume has been successfully executed`

	/*
		Volume Examples
	*/
	createVolumeExample = `ionosctl volume create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-name demoVolume
VolumeId                               Name         Size   Type   LicenceType   State   Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   10GB   HDD    LINUX         BUSY    
RequestId: a2da3bb7-3851-4e80-a5e9-6e98a66cebab
Status: Command volume create has been successfully executed`
	updateVolumeExample = `ionosctl volume update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 --volume-size 20
VolumeId                               Name         Size   Type   LicenceType   State   Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         BUSY    
RequestId: ad4080a9-a51f-4d81-ae40-660cbfe009f4
Status: Command volume update has been successfully executed`
	listVolumeExample = `ionosctl volume list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d 
VolumeId                               Name         Size   Type   LicenceType   State       Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         AVAILABLE`
	getVolumeExample = `ionosctl volume get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 
VolumeId                               Name         Size   Type   LicenceType   State       Image
ce510144-9bc6-4115-bd3d-b9cd232dd422   demoVolume   20GB   HDD    LINUX         AVAILABLE`
	deleteVolumeExample = `ionosctl volume delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --volume-id ce510144-9bc6-4115-bd3d-b9cd232dd422 
Warning: Are you sure you want to delete volume (y/N) ? y
RequestId: 6958b90b-54fa-4967-8be2-e32412559f9c
Status: Command volume delete has been successfully executed`

	/*
		Load Balancer Examples
	*/
	createLoadbalancerExample = `ionosctl loadbalancer create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-name demoLoadBalancer
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   true
RequestId: 74441964-1134-4009-8b81-d7189170885e
Status: Command loadbalancer create has been successfully executed`
	updateLoadbalancerExample = `ionosctl loadbalancer update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c --loadbalancer-dhcp=false --wait
Waiting for request: 0a9279d8-9757-41e0-b64f-b4cd2baf4717
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false
RequestId: 0a9279d8-9757-41e0-b64f-b4cd2baf4717
Status: Command loadbalancer update and request have been successfully executed`
	listLoadbalancerExample = `ionosctl loadbalancer list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d 
LoadbalancerId                         Name               Dhcp
f16dfcc1-9181-400b-a08d-7fe15ca0e9af   demoLoadbalancer   true
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false`
	getLoadbalancerExample = `ionosctl loadbalancer get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c 
LoadbalancerId                         Name               Dhcp
3f9f14a9-5fa8-4786-ba86-a91f9daded2c   demoLoadBalancer   false`
	deleteLoadbalancerExample = `ionosctl loadbalancer delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --loadbalancer-id 3f9f14a9-5fa8-4786-ba86-a91f9daded2c --ignore-stdin --wait 
Waiting for request: 29c4e7bb-8ce8-4153-8b42-3734d8ede034
RequestId: 29c4e7bb-8ce8-4153-8b42-3734d8ede034
Status: Command loadbalancer delete and request have been successfully executed`
	attachNicLoadbalancerExample = `ionosctl loadbalancer attach-nic --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --server-id 2bf04e0d-86e4-4f13-b405-442363b25e28 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   1       []
RequestId: 01b8468f-b489-40af-a4fd-3606d06da8d7
Status: Command loadbalancer attach-nic has been successfully executed`
	listNicsLoadbalancerExample = `ionosctl loadbalancer list-nics --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   2       []`
	getNicLoadbalancerExample = `ionosctl loadbalancer get-nic --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 
NicId                                  Name   Dhcp   LanId   Ips
6e8faa79-1e7e-4e99-be76-f3b3179ed3c3   test   true   2       []`
	detachNicLoadbalancerExample = `ionosctl loadbalancer detach-nic --datacenter-id 154360e9-3930-46f1-a29e-a7704ea7abc2 --loadbalancer-id 4450e35a-e89d-4769-af60-4957c3deaf33 --nic-id 6e8faa79-1e7e-4e99-be76-f3b3179ed3c3 
Warning: Are you sure you want to detach nic from loadbalancer (y/N) ? 
y
RequestId: a2a136cd-7bce-40fe-ae53-ad0d7b322387
Status: Command loadbalancer detach-nic has been successfully executed`

	/*
		NIC Examples
	*/
	createNicExample = `ionosctl nic create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-name demoNic
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   1       []
RequestId: 67bdb2fb-b1ee-419a-9bcf-f8ea4b800653
Status: Command nic create has been successfully executed`
	updateNicExample = `ionosctl nic update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 --lan-id 2 --wait 
Waiting for request: b0361cf3-06b2-4cca-ae13-4035ace9f265
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []
RequestId: b0361cf3-06b2-4cca-ae13-4035ace9f265
Status: Command nic update and request have been successfully executed`
	listNicExample = `ionosctl nic list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa 
NicId                                  Name      Dhcp   LanId   Ips
c7903181-daa1-4e16-a65a-e9b495c1b324   demoNIC   true   1       []
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []`
	getNicExample = `ionosctl nic get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 
NicId                                  Name      Dhcp   LanId   Ips
2978400e-da90-405f-905e-8200d4f48158   demoNic   true   2       []`
	deleteNicExample = `ionosctl nic delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --server-id 25baee29-d79a-4b5e-aae6-080feea977aa --nic-id 2978400e-da90-405f-905e-8200d4f48158 --ignore-stdin 
RequestId: 14a4bf17-48aa-4f87-b0dc-9c769a4cbdcb
Status: Command nic delete has been successfully executed`

	/*
		Lan Examples
	*/
	createLanExample = `ionosctl lan create --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-name demoLan
LanId   Name      Public
4       demoLan   false
RequestId: da824a69-a12a-4153-b302-a797b3581c2b
Status: Command lan create has been successfully executed`
	updateLanExample = `ionosctl lan update --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --lan-name demoLAN --lan-public=true
LanId   Name      Public
3       demoLAN   true
RequestId: 0a174dca-62b1-4360-aef8-89fd31c196f2
Status: Command lan update has been successfully executed`
	listLanExample = `ionosctl lan list --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d 
LanId   Name                                                Public
4       demoLan                                             false
3       demoLAN                                             true
2       Switch of LB f16dfcc1-9181-400b-a08d-7fe15ca0e9af   false
1       Switch of LB 3f9f14a9-5fa8-4786-ba86-a91f9daded2c   false`
	getLanExample = `ionosctl lan get --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 4
LanId   Name      Public
4       demoLan   false`
	deleteLanExample = `ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 4
Warning: Are you sure you want to delete lan (y/N) ? y
RequestId: bd5ffcf4-1b05-4cb2-917b-a0140d5f7a2b
Status: Command lan delete has been successfully executed

ionosctl lan delete --datacenter-id f28c0edd-d5ef-48f2-b8a3-aa8f6b55da3d --lan-id 3 --wait 
Warning: Are you sure you want to delete lan (y/N) ? y
Waiting for request: e65fc2fe-8005-48a5-9d06-f1a4f8bc9ef1
RequestId: e65fc2fe-8005-48a5-9d06-f1a4f8bc9ef1
Status: Command lan delete and request have been successfully executed`

	/*
		Request Examples
	*/
	getRequestExample = `ionosctl request get --request-id 20333e60-d65c-4a95-846b-08c48b871186 
RequestId                              Status   Message
20333e60-d65c-4a95-846b-08c48b871186   DONE     Request has been successfully executed`
	waitRequestExample = `ionosctl request wait --request-id 20333e60-d65c-4a95-846b-08c48b871186 
RequestId                              Status   Message
20333e60-d65c-4a95-846b-08c48b871186   DONE     Request has been successfully executed`

	/*
		Image Examples
	*/
	listImagesExample = `ionosctl image list --image-location us/las --image-type HDD
ImageId                                Name                                 Location   Size   LicenceType   ImageType
8991cf6c-8706-11eb-a1d6-72dfddd36b99   windows-2012-r2-server-2021-03       us/las     14     WINDOWS       HDD
7ab978cb-870a-11eb-a1d6-72dfddd36b99   windows-2016-server-2021-03          us/las     14     WINDOWS2016   HDD
aca2279d-870b-11eb-a1d6-72dfddd36b99   windows-2019-server-2021-03          us/las     15     WINDOWS2016   HDD
33915e02-9291-11eb-b68e-9ad3ea4b1420   CentOS-7-server-2021-04-01           us/las     4      LINUX         HDD
bc81846b-929c-11eb-b68e-9ad3ea4b1420   Debian-testing-server-2021-04-01     us/las     2      LINUX         HDD
8eca30f5-92a0-11eb-b68e-9ad3ea4b1420   Ubuntu-16.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
f1316493-92a4-11eb-b68e-9ad3ea4b1420   Ubuntu-18.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
a86dd807-9297-11eb-b68e-9ad3ea4b1420   Debian-10-server-2021-04-01          us/las     2      LINUX         HDD
4bf9c128-929a-11eb-b68e-9ad3ea4b1420   Debian-9-server-2021-04-01           us/las     2      LINUX         HDD
378ecd34-9295-11eb-b68e-9ad3ea4b1420   CentOS-8-server-2021-04-01           us/las     4      LINUX         HDD
02e06d34-92aa-11eb-b68e-9ad3ea4b1420   Ubuntu-20.04-LTS-server-2021-04-01   us/las     3      LINUX         HDD
81adeb01-3379-11eb-a681-1e659523cb7b   CentOS-6-server-2020-12-01           us/las     2      LINUX         HDD
6f9bae91-3386-11eb-a681-1e659523cb7b   Debian-8-server-2020-12-01           us/las     2      LINUX         HDD
8fc5f591-338e-11eb-a681-1e659523cb7b   Ubuntu-19.10-server-2020-12-01       us/las     3      LINUX         HDD`
	getImageExample = `ionosctl image get --image-id 8fc5f591-338e-11eb-a681-1e659523cb7b 
ImageId                                Name                             Location   Size   LicenceType   ImageType
8fc5f591-338e-11eb-a681-1e659523cb7b   Ubuntu-19.10-server-2020-12-01   us/las     3      LINUX         HDD`

	/*
		Snapshot Examples
	*/
	listSnapshotsExample = `ionosctl snapshot list 
SnapshotId                             Name           LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   testSnapshot   LINUX         10
8e0bc509-87ee-47f4-a382-302e4f7e103d   image          LINUX         10`
	getSnapshotExample = `ionosctl snapshot get --snapshot-id dc688daf-8e54-4db8-ac4a-487ad5a34e9c 
SnapshotId                             Name           LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   testSNapshot   LINUX         10`
	createSnapshotExample = `ionosctl snapshot create --datacenter-id 451cc0c1-883a-44aa-9ae4-336c0c3eaa5d --volume-id 4acddd40-959f-4517-b628-dc24e37df942 --snapshot-name testSnapshot
SnapshotId                             Name           LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   testSnapshot   LINUX         0
RequestId: fed5555a-ac00-41c8-abbe-cc53c8179716
Status: Command snapshot create has been successfully executed`
	updateSnapshotExample = `ionosctl snapshot update --snapshot-id dc688daf-8e54-4db8-ac4a-487ad5a34e9c --snapshot-name test
SnapshotId                             Name   LicenceType   Size
dc688daf-8e54-4db8-ac4a-487ad5a34e9c   test   LINUX         10
RequestId: 3540e9be-ed35-41c0-83d9-923882bfa9bd
Status: Command snapshot update has been successfully executed`
	restoreSnapshotExample = `ionosctl snapshot restore --snapshot-id dc688daf-8e54-4db8-ac4a-487ad5a34e9c --datacenter-id 451cc0c1-883a-44aa-9ae4-336c0c3eaa5d --volume-id 4acddd40-959f-4517-b628-dc24e37df942 --wait 
Warning: Are you sure you want to restore snapshot (y/N) ? 
y
RequestId: 21ca5546-9314-4cd5-8832-6029638b1237
Status: Command snapshot restore and request have been successfully executed`
	deleteSnapshotExample = `ionosctl snapshot delete --snapshot-id 8e0bc509-87ee-47f4-a382-302e4f7e103d --wait 
Warning: Are you sure you want to delete snapshot (y/N) ? 
y
RequestId: 6e029eb6-47e6-4dcd-a333-d620b49c01e5
Status: Command snapshot delete and request have been successfully executed`

	/*
		IpBlock Examples
	*/
	listIpBlockExample = `ionosctl ipblock list 
IpBlockId                              Name   Location   Size   Ips                 State
bf932826-d71b-4759-a7d0-0028261c1e8d   demo   us/las     1      [x.x.x.x]           AVAILABLE
3bb77993-dd2a-4845-8115-5001ae87d5e4   test   us/las     2      [x.x.x.x x.x.x.x]   AVAILABLE`
	getIpBlockExample = `ionosctl ipblock get --ipblock-id 3bb77993-dd2a-4845-8115-5001ae87d5e4 
IpBlockId                              Name   Location   Size   Ips                 State
3bb77993-dd2a-4845-8115-5001ae87d5e4   test   us/las     2      [x.x.x.x x.x.x.x]   AVAILABLE`
	createIpBlockExample = `ionosctl ipblock create --ipblock-name testing --ipblock-location us/las --ipblock-size 1
IpBlockId                              Name      Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   testing   us/las     1      [x.x.x.x]   BUSY
RequestId: a99bd16c-bf7b-4966-8a30-437b5182226b
Status: Command ipblock create has been successfully executed`
	updateIpBlockExample = `ionosctl ipblock update --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --ipblock-name demo
IpBlockId                              Name   Location   Size   Ips         State
bf932826-d71b-4759-a7d0-0028261c1e8d   demo   us/las     1      [x.x.x.x]   BUSY
RequestId: 5864afe5-4df5-4843-b548-4489857dc3c5
Status: Command ipblock update has been successfully executed`
	deleteIpBlockExample = `ionosctl ipblock delete --ipblock-id bf932826-d71b-4759-a7d0-0028261c1e8d --wait 
Warning: Are you sure you want to delete ipblock (y/N) ? 
y
Waiting for request: 6b1aa258-799f-4712-9f90-ba4494d84026
RequestId: 6b1aa258-799f-4712-9f90-ba4494d84026
Status: Command ipblock delete and request have been successfully executed`

	/*
		Firewall Rule Examples
	*/
	listFirewallRuleExample = `ionosctl firewallrule list --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 
FirewallRuleId                         Name        Protocol   PortRangeStart   PortRangeStop   State
f537ff0e-8b2c-4ce6-8a92-297a5ad08ca1   test        TCP        80               80              AVAILABLE`
	getFirewallRuleExample = `ionosctl firewallrule get --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id f537ff0e-8b2c-4ce6-8a92-297a5ad08ca1 
FirewallRuleId                         Name        Protocol   PortRangeStart   PortRangeEnd   State
f537ff0e-8b2c-4ce6-8a92-297a5ad08ca1   test        TCP        80               80             AVAILABLE`
	createFirewallRuleExample = `ionosctl firewallrule create --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-protocol TCP --firewallrule-name demo --firewallrule-port-range-start 2476 --firewallrule-port-range-end 2476
FirewallRuleId                         Name   Protocol   PortRangeStart   PortRangeEnd   State
4221e2c8-0316-447c-aeed-69ac92e585be   demo   TCP        2476             2476           BUSY
RequestId: 09a47137-e377-4a79-b2b9-16744e298ad5
Status: Command firewallrule create has been successfully executed`
	updateFirewallRuleExample = `ionosctl firewallrule update --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id 4221e2c8-0316-447c-aeed-69ac92e585be --firewallrule-name new-test --wait 
Waiting for request: 2e3d6e81-2830-4d68-82ff-daee6f115864
FirewallRuleId                         Name       Protocol   PortRangeStart   PortRangeEnd   State
4221e2c8-0316-447c-aeed-69ac92e585be   new-test   TCP        2476             2476           BUSY
RequestId: 2e3d6e81-2830-4d68-82ff-daee6f115864
Status: Command firewallrule update and request have been successfully executed`
	deleteFirewallRuleExample = `ionosctl firewallrule delete --datacenter-id f2d82ba9-7dc4-4945-89b6-3d194f6be29b --server-id d776e064-a3f9-4fbd-8729-93818b7459bb --nic-id 029c05a4-f5f7-4398-9469-2eb3d6db3460 --firewallrule-id e7c4e91a-d3e3-42db-bfb1-2d5e9ebc952b 
Warning: Are you sure you want to delete firewall rule (y/N) ? 
y
RequestId: 481b6e7c-0c31-4395-81e4-36fad877b77b
Status: Command firewallrule delete has been successfully executed`

	/*
		Label Examples
	*/
	listLabelsExample = `ionosctl label list 
Key       Value            ResourceType   ResourceId
test      testserver       server         27dde318-f0d4-4f97-a04d-9dafe4a89637
test      testdatacenter   datacenter     ed612a0a-9506-4b56-8d1b-ce2b04090f19
test      testsnapshot     snapshot       df7f4ad9-b942-4e79-939d-d1c10fb6fbff`
	getLabelExample = `ionosctl label get --label-urn "urn:label:server:27dde318-f0d4-4f97-a04d-9dafe4a89637:test"
Key    Value        ResourceType   ResourceId
test   testserver   server         27dde318-f0d4-4f97-a04d-9dafe4a89637`

	/*
		Label Resources Examples
	*/
	listDataCenterLabelsExample = `ionosctl datacenter list-labels --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 
Key    Value
test   testdatacenter`
	getDataCenterLabelExample = `ionosctl datacenter get-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --label-key test
Key    Value
test   testdatacenter`
	addDataCenterLabelExample = `ionosctl datacenter add-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --label-key test --label-value testdatacenter
Key    Value
test   testdatacenter`
	removeDataCenterLabelExample = `ionosctl datacenter remove-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --label-key test`
	listServerLabelsExample      = `ionosctl server list-labels --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 
Key    Value
test   test`
	getServerLabelExample = `ionosctl server get-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 --label-key test
Key    Value
test   test`
	addServerLabelExample = `ionosctl server add-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 --label-key test --label-value test
Key    Value
test   test`
	removeServerLabelExample = `ionosctl server remove-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --server-id 27dde318-f0d4-4f97-a04d-9dafe4a89637 --label-key test`
	listVolumeLabelsExample  = `ionosctl volume list-labels --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --volume-id 5d23eee2-45e5-44fe-96fe-e15aba2c48f5 
Key    Value
test   testvolume`
	getVolumeLabelExample = `ionosctl volume get-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --volume-id 5d23eee2-45e5-44fe-96fe-e15aba2c48f5 --label-key test
Key    Value
test   testvolume`
	addVolumeLabelExample = `ionosctl volume add-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --volume-id 5d23eee2-45e5-44fe-96fe-e15aba2c48f5 --label-key test --label-value testvolume
Key    Value
test   testvolume`
	removeVolumeLabelExample = `ionosctl volume remove-label --datacenter-id ed612a0a-9506-4b56-8d1b-ce2b04090f19 --volume-id 5d23eee2-45e5-44fe-96fe-e15aba2c48f5 --label-key test`
	listIpBlockLabelsExample = `ionosctl ipblock list-labels --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 
Key    Value
test   testipblock`
	getIpBlockLabelExample = `ionosctl ipblock get-label --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 --label-key test
Key    Value
test   testipblock`
	addIpBlockLabelExample = `ionosctl ipblock add-label --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 --label-key test --label-value testipblock
Key    Value
test   testipblock`
	removeIpBlockLabelExample = `ionosctl ipblock remove-label --ipblock-id 379a995b-f285-493e-a56a-f32e1cb6dd06 --label-key test`
	listSnapshotLabelsExample = `ionosctl snapshot list-labels --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff
Key    Value
test   testsnapshot`
	getSnapshotLabelExample = ` ionosctl snapshot get-label --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff --label-key test
Key    Value
test   testsnapshot`
	addSnapshotLabelExample = `ionosctl snapshot add-label --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff --label-key test --label-value testsnapshot
Key    Value
test   testsnapshot`
	removeSnapshotLabelExample = `ionosctl snapshot remove-label --snapshot-id df7f4ad9-b942-4e79-939d-d1c10fb6fbff --label-key test`

	/*
		Contract Resources Examples
	*/
	getContractExample = `ionosctl contract get --resource-limits [ CORES|RAM|HDD|SSD|IPS|K8S ]`

	/*
		User Examples
	*/
	listUserExample = `ionosctl user list 
UserId                                 Firstname   Lastname   Email                      Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com     false           false          false           a74101e7c1948450432d5b6512f2712c   true
53d68de9-931a-4b61-b532-82f7b27afef3   test1       test1      testrandom13@ionos.com     false           false          false           8b9dd6f39e613adb7a837127edb67d38   true`
	getUserExample = `ionosctl user get --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com   false           false          false           a74101e7c1948450432d5b6512f2712c   true`
	createUserExample = `ionosctl user create --user-first-name test1 --user-last-name test1 --user-email testrandom16@gmail.com --user-password test123
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId   Active
99499053-059e-4ee6-b56f-66b0df93262d   test1       test1      testrandom16@ionos.com   false           false          false                               true
RequestId: ca349e08-5820-41ba-8252-ee4c8dd2ccdb
Status: Command user create has been successfully executed`
	updateUserExample = `ionosctl user update --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 --user-administrator=true
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
2470f439-1d73-42f8-90a9-f78cf2776c74   test1       test1      testrandom12@ionos.com   true            false          false           a74101e7c1948450432d5b6512f2712c   true
RequestId: 439f79fc-5bfc-43da-92f3-0d804ebb28ac
Status: Command user update has been successfully executed`
	deleteUserExample = `ionosctl user delete --user-id 2470f439-1d73-42f8-90a9-f78cf2776c74 --ignore-stdin 
RequestId: a2f6e7fa-6030-4267-950e-1e0886316475
Status: Command user delete has been successfully executed`

	/*
		Group Examples
	*/
	createGroupExample = `ionosctl group create --group-name test --wait 
Waiting for request: eae6bb8b-3736-4cf0-bc71-72a95d1b2a63
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false`
	getGroupExample = `ionosctl group get --group-id 1d500d7a-43af-488a-a656-79e902433767 
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false`
	listGroupExample = `ionosctl group list
GroupId                                Name   CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
1d500d7a-43af-488a-a656-79e902433767   test   false              false            false       false               false       false         false              false                  false`
	updateGroupExample = `ionosctl group update --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --group-reserve-ip 
GroupId                                Name         CreateDataCenter   CreateSnapshot   ReserveIp   AccessActivityLog   CreatePcc   S3Privilege   CreateBackupUnit   CreateInternetAccess   CreateK8s
e99f4cdb-746d-4c3c-b38c-b749ca23f917   testUpdate   true               true             true        false               false       false         false              false                  true
RequestId: 2bfe43a4-ea09-48fc-bb53-136c7f7d061f
Status: Command group update has been successfully executed`
	deleteGroupExample = `ionosctl group delete --group-id 1d500d7a-43af-488a-a656-79e902433767 
Warning: Are you sure you want to delete group (y/N) ? 
y
RequestId: e20d2851-0d20-453d-b752-ed1c34a83625
Status: Command group delete has been successfully executed`
	/*
		Group Users Examples
	*/
	listGroupUsersExample = `ionosctl group list-users --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
53d68de9-931a-4b61-b532-82f7b27afef3   test1       test1      testrandom13@ionos.com   false           false          false           8b9dd6f39e613adb7a837127edb67d38   true
e3b5fefb-27b0-4eea-a7c4-c57934ad23cb   test1       test1      testrandom14@ionos.com   false           false          false           25e754c1f9f0169213ec4ad5e5e02dcd   true`
	removeGroupUserExample = `ionosctl group remove-user --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --user-id e3b5fefb-27b0-4eea-a7c4-c57934ad23cb --wait 
Warning: Are you sure you want to remove user from group (y/N) ? 
y
Waiting for request: 353eff98-120f-4d91-82f5-a8aff1ddb277
RequestId: 353eff98-120f-4d91-82f5-a8aff1ddb277
Status: Command group remove-user and request have been successfully executed`
	addGroupUserExample = `ionosctl group add-user --group-id e99f4cdb-746d-4c3c-b38c-b749ca23f917 --user-id 53d68de9-931a-4b61-b532-82f7b27afef3
UserId                                 Firstname   Lastname   Email                    Administrator   ForceSecAuth   SecAuthActive   S3CanonicalUserId                  Active
53d68de9-931a-4b61-b532-82f7b27afef3   test1       test1      testrandom13@ionos.com   false           false          false           8b9dd6f39e613adb7a837127edb67d38   true`
	/*
		Group Shares Example
	*/
	listGroupSharesExample = `ionosctl group list-shares --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           false`
	getGroupShareExample = `ionosctl group get-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           true`
	addGroupShareExample = `ionosctl group add-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           false
RequestId: ffb8e7ba-4a49-4ea5-a97e-e3a61e55c277
Status: Command group add-share has been successfully executed`
	updateGroupShareExample = `ionosctl group update-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e --share-privilege 
ShareId                                EditPrivilege   SharePrivilege
cefc2175-001f-4b94-8693-6263d731fe8e   false           true
RequestId: 0dfccab0-c148-40c8-9794-067d23f79f0e
Status: Command group update-share has been successfully executed`
	removeGroupShareExample = `ionosctl group remove-share --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f --resource-id cefc2175-001f-4b94-8693-6263d731fe8e --wait 
Warning: Are you sure you want to remove share from group (y/N) ? 
y
Waiting for request: 9ff7e57f-b568-4257-b27f-13a4cf11a7fc
RequestId: 9ff7e57f-b568-4257-b27f-13a4cf11a7fc
Status: Command group remove-share and request have been successfully executed`
	/*
		Group Resources Example
	*/
	listGroupResourcesExample = `ionosctl group list-resources --group-id 83ad9b77-7598-44d7-a817-d3f12f92387f 
ResourceId                             Name   SecAuthProtection   Type
cefc2175-001f-4b94-8693-6263d731fe8e          false               datacenter`
	/*
		Resources Example
	*/
	listResourcesExample = `ionosctl resource list 
ResourceId                             Name                            SecAuthProtection   Type
cefc2175-001f-4b94-8693-6263d731fe8e                                   false               datacenter
d8922413-05f1-48bb-90ed-c2d407e05b1d   IP_BLOCK_2021-04-20T11:02:52Z   false               ipblock`
	getResourceExample = `ionosctl resource get --resource-type ipblock
ResourceId                             Name                            SecAuthProtection   Type
d8922413-05f1-48bb-90ed-c2d407e05b1d   IP_BLOCK_2021-04-20T11:02:52Z   false               ipblock`

	/*
		S3Keys Example
	*/
	listS3KeysExample = `ionosctl s3key list --user-id 013188d4-af9a-4207-b495-de36cb2dc344 
S3KeyId                Active
00a29d110b48daa3a18b   false`
	getS3KeyExample = `ionosctl s3key get --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 00a29d110b48daa3a18b 
S3KeyId                Active
00a29d110b48daa3a18b   false`
	createS3KeyExample = `ionosctl s3key create --user-id 013188d4-af9a-4207-b495-de36cb2dc344 
S3KeyId                Active
75f4319cbf3f6d538da7   true
RequestId: 869fc059-165d-480b-a913-a410d38d20e0
Status: Command s3key create has been successfully executed`
	updateS3KeyExample = `ionosctl s3key update --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 75f4319cbf3f6d538da7 --s3key-active=false
S3KeyId                Active
75f4319cbf3f6d538da7   false
RequestId: 4cda4b65-f58b-492a-bf45-6f1d8fb42928
Status: Command s3key update has been successfully executed`
	deleteS3KeyExample = `ionosctl s3key delete --user-id 013188d4-af9a-4207-b495-de36cb2dc344 --s3key-id 75f4319cbf3f6d538da7 --wait 
Warning: Are you sure you want to delete S3Key (y/N) ? 
y
RequestId: 1529f8b7-08bb-4321-a996-08865660dee8
Status: Command s3key delete and request have been successfully executed`

	/*
		Private Cross-Connect Example
	*/
	listPccsExample = `ionosctl pcc list 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test
4b9c6a43-a338-11eb-b70c-7ade62b52cc0   test   test`
	getPccExample = `ionosctl pcc get --pcc-id e2337b40-52d9-48d2-bcbc-41c5abc29d11 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test`
	getPccPeersExample = `ionosctl pcc get-peers --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 
LanId   LanName     DatacenterId                           DatacenterName   Location
1       testlan2    1ef56b51-98be-487e-925a-c9f3dfa4a076   test2            us/las
1       testlan1    95b7f7f0-a6f3-4fc9-8d06-018d2c1efc89   test1            us/las`
	createPccExample = `ionosctl pcc create --pcc-name test --pcc-description "test test" --wait 
PccId                                  Name   Description
e2337b40-52d9-48d2-bcbc-41c5abc29d11   test   test test
RequestId: 64720266-c6e8-4e78-8e31-6754f006dcb1
Status: Command pcc create and request have been successfully executed`
	updatePccExample = `ionosctl pcc update --pcc-id 4b9c6a43-a338-11eb-b70c-7ade62b52cc0 --pcc-description test
PccId                                  Name   Description
4b9c6a43-a338-11eb-b70c-7ade62b52cc0   test   test
RequestId: 81525f2d-cc91-4c55-84b8-07fac9a47e35
Status: Command pcc update has been successfully executed`
	deletePccExample = `ionosctl pcc delete --pcc-id e2337b40-52d9-48d2-bcbc-41c5abc29d11 --wait 
Warning: Are you sure you want to delete private cross-connect (y/N) ? 
y
RequestId: 7fa56e7f-1d63-4c5f-a7ea-eec6a015282a
Status: Command pcc delete and request have been successfully executed`
)
