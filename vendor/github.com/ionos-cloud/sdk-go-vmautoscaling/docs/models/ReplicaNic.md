# ReplicaNic

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Lan** | **int32** | The LAN ID of this replica NIC. | |
|**Name** | **string** | The replica NIC name. | |
|**Dhcp** | Pointer to **NullableBool** | DHCP for this replica NIC. This is an optional attribute with the default value &#39;TRUE&#39; if not specified in the request payload or as null. | [optional] |
|**FirewallActive** | Pointer to **NullableBool** | Activate or deactivate the firewall. By default, an active firewall without any defined rules will block all incoming network traffic except for the firewall rules that explicitly allows certain protocols, IP addresses and ports. | [optional] |
|**FirewallType** | Pointer to **NullableString** | The type of firewall rules that will be allowed on the NIC. If not specified, the default INGRESS value is used. | [optional] |
|**FlowLogs** | Pointer to [**[]NicFlowLog**](NicFlowLog.md) | List of all flow logs for the specified NIC. | [optional] |
|**FirewallRules** | Pointer to [**[]NicFirewallRule**](NicFirewallRule.md) | List of all firewall rules for the specified NIC. | [optional] |
|**TargetGroup** | Pointer to [**TargetGroup**](TargetGroup.md) |  | [optional] |

## Methods

### NewReplicaNic

`func NewReplicaNic(lan int32, name string, ) *ReplicaNic`

NewReplicaNic instantiates a new ReplicaNic object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewReplicaNicWithDefaults

`func NewReplicaNicWithDefaults() *ReplicaNic`

NewReplicaNicWithDefaults instantiates a new ReplicaNic object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLan

`func (o *ReplicaNic) GetLan() int32`

GetLan returns the Lan field if non-nil, zero value otherwise.

### GetLanOk

`func (o *ReplicaNic) GetLanOk() (*int32, bool)`

GetLanOk returns a tuple with the Lan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLan

`func (o *ReplicaNic) SetLan(v int32)`

SetLan sets Lan field to given value.


### GetName

`func (o *ReplicaNic) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ReplicaNic) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ReplicaNic) SetName(v string)`

SetName sets Name field to given value.


### GetDhcp

`func (o *ReplicaNic) GetDhcp() bool`

GetDhcp returns the Dhcp field if non-nil, zero value otherwise.

### GetDhcpOk

`func (o *ReplicaNic) GetDhcpOk() (*bool, bool)`

GetDhcpOk returns a tuple with the Dhcp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcp

`func (o *ReplicaNic) SetDhcp(v bool)`

SetDhcp sets Dhcp field to given value.

### HasDhcp

`func (o *ReplicaNic) HasDhcp() bool`

HasDhcp returns a boolean if a field has been set.

### SetDhcpNil

`func (o *ReplicaNic) SetDhcpNil(b bool)`

 SetDhcpNil sets the value for Dhcp to be an explicit nil

### UnsetDhcp
`func (o *ReplicaNic) UnsetDhcp()`

UnsetDhcp ensures that no value is present for Dhcp, not even an explicit nil
### GetFirewallActive

`func (o *ReplicaNic) GetFirewallActive() bool`

GetFirewallActive returns the FirewallActive field if non-nil, zero value otherwise.

### GetFirewallActiveOk

`func (o *ReplicaNic) GetFirewallActiveOk() (*bool, bool)`

GetFirewallActiveOk returns a tuple with the FirewallActive field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirewallActive

`func (o *ReplicaNic) SetFirewallActive(v bool)`

SetFirewallActive sets FirewallActive field to given value.

### HasFirewallActive

`func (o *ReplicaNic) HasFirewallActive() bool`

HasFirewallActive returns a boolean if a field has been set.

### SetFirewallActiveNil

`func (o *ReplicaNic) SetFirewallActiveNil(b bool)`

 SetFirewallActiveNil sets the value for FirewallActive to be an explicit nil

### UnsetFirewallActive
`func (o *ReplicaNic) UnsetFirewallActive()`

UnsetFirewallActive ensures that no value is present for FirewallActive, not even an explicit nil
### GetFirewallType

`func (o *ReplicaNic) GetFirewallType() string`

GetFirewallType returns the FirewallType field if non-nil, zero value otherwise.

### GetFirewallTypeOk

`func (o *ReplicaNic) GetFirewallTypeOk() (*string, bool)`

GetFirewallTypeOk returns a tuple with the FirewallType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirewallType

`func (o *ReplicaNic) SetFirewallType(v string)`

SetFirewallType sets FirewallType field to given value.

### HasFirewallType

`func (o *ReplicaNic) HasFirewallType() bool`

HasFirewallType returns a boolean if a field has been set.

### SetFirewallTypeNil

`func (o *ReplicaNic) SetFirewallTypeNil(b bool)`

 SetFirewallTypeNil sets the value for FirewallType to be an explicit nil

### UnsetFirewallType
`func (o *ReplicaNic) UnsetFirewallType()`

UnsetFirewallType ensures that no value is present for FirewallType, not even an explicit nil
### GetFlowLogs

`func (o *ReplicaNic) GetFlowLogs() []NicFlowLog`

GetFlowLogs returns the FlowLogs field if non-nil, zero value otherwise.

### GetFlowLogsOk

`func (o *ReplicaNic) GetFlowLogsOk() (*[]NicFlowLog, bool)`

GetFlowLogsOk returns a tuple with the FlowLogs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlowLogs

`func (o *ReplicaNic) SetFlowLogs(v []NicFlowLog)`

SetFlowLogs sets FlowLogs field to given value.

### HasFlowLogs

`func (o *ReplicaNic) HasFlowLogs() bool`

HasFlowLogs returns a boolean if a field has been set.

### SetFlowLogsNil

`func (o *ReplicaNic) SetFlowLogsNil(b bool)`

 SetFlowLogsNil sets the value for FlowLogs to be an explicit nil

### UnsetFlowLogs
`func (o *ReplicaNic) UnsetFlowLogs()`

UnsetFlowLogs ensures that no value is present for FlowLogs, not even an explicit nil
### GetFirewallRules

`func (o *ReplicaNic) GetFirewallRules() []NicFirewallRule`

GetFirewallRules returns the FirewallRules field if non-nil, zero value otherwise.

### GetFirewallRulesOk

`func (o *ReplicaNic) GetFirewallRulesOk() (*[]NicFirewallRule, bool)`

GetFirewallRulesOk returns a tuple with the FirewallRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFirewallRules

`func (o *ReplicaNic) SetFirewallRules(v []NicFirewallRule)`

SetFirewallRules sets FirewallRules field to given value.

### HasFirewallRules

`func (o *ReplicaNic) HasFirewallRules() bool`

HasFirewallRules returns a boolean if a field has been set.

### SetFirewallRulesNil

`func (o *ReplicaNic) SetFirewallRulesNil(b bool)`

 SetFirewallRulesNil sets the value for FirewallRules to be an explicit nil

### UnsetFirewallRules
`func (o *ReplicaNic) UnsetFirewallRules()`

UnsetFirewallRules ensures that no value is present for FirewallRules, not even an explicit nil
### GetTargetGroup

`func (o *ReplicaNic) GetTargetGroup() TargetGroup`

GetTargetGroup returns the TargetGroup field if non-nil, zero value otherwise.

### GetTargetGroupOk

`func (o *ReplicaNic) GetTargetGroupOk() (*TargetGroup, bool)`

GetTargetGroupOk returns a tuple with the TargetGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetGroup

`func (o *ReplicaNic) SetTargetGroup(v TargetGroup)`

SetTargetGroup sets TargetGroup field to given value.

### HasTargetGroup

`func (o *ReplicaNic) HasTargetGroup() bool`

HasTargetGroup returns a boolean if a field has been set.


