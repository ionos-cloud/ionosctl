# NicFirewallRule

## Properties

|Name | Type | Description | Notes|
|------------ | ------------- | ------------- | -------------|
|**Name** | Pointer to **NullableString** | The name of the firewall rule. | [optional] |
|**Protocol** | **string** | The protocol for the rule. The property cannot be modified after its creation (not allowed in update requests). | |
|**SourceMac** | Pointer to **NullableString** | Only traffic originating from the respective MAC address is permitted. Valid format: &#39;aa:bb:cc:dd:ee:ff&#39;. The value &#39;null&#39; allows traffic from any MAC address. | [optional] |
|**SourceIp** | Pointer to **NullableString** | Only traffic originating from the respective IPv4 address is permitted. The value &#39;null&#39; allows traffic from any IP address. | [optional] |
|**TargetIp** | Pointer to **NullableString** | If the target NIC has multiple IP addresses, only the traffic directed to the respective IP address of the NIC is allowed. The value &#39;null&#39; allows traffic to any target IP address. | [optional] |
|**IcmpCode** | Pointer to **NullableInt32** | Sets the allowed code (from 0 to 254) when ICMP protocol is selected. The value &#39;null&#39;&#39; allows all codes. | [optional] |
|**IcmpType** | Pointer to **NullableInt32** | Sets the allowed type (from 0 to 254) if the protocol ICMP is selected. The value &#39;null&#39; allows all types. | [optional] |
|**PortRangeStart** | Pointer to **NullableInt32** | Sets the initial range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value &#39;null&#39; for &#39;portRangeStart&#39; and &#39;portRangeEnd&#39; allows all ports. | [optional] |
|**PortRangeEnd** | Pointer to **NullableInt32** | Sets the end range of the allowed port (from 1 to 65535) if the protocol TCP or UDP is selected. The value &#39;null&#39; for &#39;portRangeStart&#39; and &#39;portRangeEnd&#39; allows all ports. | [optional] |
|**Type** | Pointer to **NullableString** | The firewall rule type. If not specified, the default value &#39;INGRESS&#39; is used. | [optional] |

## Methods

### NewNicFirewallRule

`func NewNicFirewallRule(protocol string, ) *NicFirewallRule`

NewNicFirewallRule instantiates a new NicFirewallRule object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewNicFirewallRuleWithDefaults

`func NewNicFirewallRuleWithDefaults() *NicFirewallRule`

NewNicFirewallRuleWithDefaults instantiates a new NicFirewallRule object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *NicFirewallRule) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *NicFirewallRule) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *NicFirewallRule) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *NicFirewallRule) HasName() bool`

HasName returns a boolean if a field has been set.

### SetNameNil

`func (o *NicFirewallRule) SetNameNil(b bool)`

 SetNameNil sets the value for Name to be an explicit nil

### UnsetName
`func (o *NicFirewallRule) UnsetName()`

UnsetName ensures that no value is present for Name, not even an explicit nil
### GetProtocol

`func (o *NicFirewallRule) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *NicFirewallRule) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *NicFirewallRule) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.


### GetSourceMac

`func (o *NicFirewallRule) GetSourceMac() string`

GetSourceMac returns the SourceMac field if non-nil, zero value otherwise.

### GetSourceMacOk

`func (o *NicFirewallRule) GetSourceMacOk() (*string, bool)`

GetSourceMacOk returns a tuple with the SourceMac field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceMac

`func (o *NicFirewallRule) SetSourceMac(v string)`

SetSourceMac sets SourceMac field to given value.

### HasSourceMac

`func (o *NicFirewallRule) HasSourceMac() bool`

HasSourceMac returns a boolean if a field has been set.

### SetSourceMacNil

`func (o *NicFirewallRule) SetSourceMacNil(b bool)`

 SetSourceMacNil sets the value for SourceMac to be an explicit nil

### UnsetSourceMac
`func (o *NicFirewallRule) UnsetSourceMac()`

UnsetSourceMac ensures that no value is present for SourceMac, not even an explicit nil
### GetSourceIp

`func (o *NicFirewallRule) GetSourceIp() string`

GetSourceIp returns the SourceIp field if non-nil, zero value otherwise.

### GetSourceIpOk

`func (o *NicFirewallRule) GetSourceIpOk() (*string, bool)`

GetSourceIpOk returns a tuple with the SourceIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceIp

`func (o *NicFirewallRule) SetSourceIp(v string)`

SetSourceIp sets SourceIp field to given value.

### HasSourceIp

`func (o *NicFirewallRule) HasSourceIp() bool`

HasSourceIp returns a boolean if a field has been set.

### SetSourceIpNil

`func (o *NicFirewallRule) SetSourceIpNil(b bool)`

 SetSourceIpNil sets the value for SourceIp to be an explicit nil

### UnsetSourceIp
`func (o *NicFirewallRule) UnsetSourceIp()`

UnsetSourceIp ensures that no value is present for SourceIp, not even an explicit nil
### GetTargetIp

`func (o *NicFirewallRule) GetTargetIp() string`

GetTargetIp returns the TargetIp field if non-nil, zero value otherwise.

### GetTargetIpOk

`func (o *NicFirewallRule) GetTargetIpOk() (*string, bool)`

GetTargetIpOk returns a tuple with the TargetIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetIp

`func (o *NicFirewallRule) SetTargetIp(v string)`

SetTargetIp sets TargetIp field to given value.

### HasTargetIp

`func (o *NicFirewallRule) HasTargetIp() bool`

HasTargetIp returns a boolean if a field has been set.

### SetTargetIpNil

`func (o *NicFirewallRule) SetTargetIpNil(b bool)`

 SetTargetIpNil sets the value for TargetIp to be an explicit nil

### UnsetTargetIp
`func (o *NicFirewallRule) UnsetTargetIp()`

UnsetTargetIp ensures that no value is present for TargetIp, not even an explicit nil
### GetIcmpCode

`func (o *NicFirewallRule) GetIcmpCode() int32`

GetIcmpCode returns the IcmpCode field if non-nil, zero value otherwise.

### GetIcmpCodeOk

`func (o *NicFirewallRule) GetIcmpCodeOk() (*int32, bool)`

GetIcmpCodeOk returns a tuple with the IcmpCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcmpCode

`func (o *NicFirewallRule) SetIcmpCode(v int32)`

SetIcmpCode sets IcmpCode field to given value.

### HasIcmpCode

`func (o *NicFirewallRule) HasIcmpCode() bool`

HasIcmpCode returns a boolean if a field has been set.

### SetIcmpCodeNil

`func (o *NicFirewallRule) SetIcmpCodeNil(b bool)`

 SetIcmpCodeNil sets the value for IcmpCode to be an explicit nil

### UnsetIcmpCode
`func (o *NicFirewallRule) UnsetIcmpCode()`

UnsetIcmpCode ensures that no value is present for IcmpCode, not even an explicit nil
### GetIcmpType

`func (o *NicFirewallRule) GetIcmpType() int32`

GetIcmpType returns the IcmpType field if non-nil, zero value otherwise.

### GetIcmpTypeOk

`func (o *NicFirewallRule) GetIcmpTypeOk() (*int32, bool)`

GetIcmpTypeOk returns a tuple with the IcmpType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIcmpType

`func (o *NicFirewallRule) SetIcmpType(v int32)`

SetIcmpType sets IcmpType field to given value.

### HasIcmpType

`func (o *NicFirewallRule) HasIcmpType() bool`

HasIcmpType returns a boolean if a field has been set.

### SetIcmpTypeNil

`func (o *NicFirewallRule) SetIcmpTypeNil(b bool)`

 SetIcmpTypeNil sets the value for IcmpType to be an explicit nil

### UnsetIcmpType
`func (o *NicFirewallRule) UnsetIcmpType()`

UnsetIcmpType ensures that no value is present for IcmpType, not even an explicit nil
### GetPortRangeStart

`func (o *NicFirewallRule) GetPortRangeStart() int32`

GetPortRangeStart returns the PortRangeStart field if non-nil, zero value otherwise.

### GetPortRangeStartOk

`func (o *NicFirewallRule) GetPortRangeStartOk() (*int32, bool)`

GetPortRangeStartOk returns a tuple with the PortRangeStart field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortRangeStart

`func (o *NicFirewallRule) SetPortRangeStart(v int32)`

SetPortRangeStart sets PortRangeStart field to given value.

### HasPortRangeStart

`func (o *NicFirewallRule) HasPortRangeStart() bool`

HasPortRangeStart returns a boolean if a field has been set.

### SetPortRangeStartNil

`func (o *NicFirewallRule) SetPortRangeStartNil(b bool)`

 SetPortRangeStartNil sets the value for PortRangeStart to be an explicit nil

### UnsetPortRangeStart
`func (o *NicFirewallRule) UnsetPortRangeStart()`

UnsetPortRangeStart ensures that no value is present for PortRangeStart, not even an explicit nil
### GetPortRangeEnd

`func (o *NicFirewallRule) GetPortRangeEnd() int32`

GetPortRangeEnd returns the PortRangeEnd field if non-nil, zero value otherwise.

### GetPortRangeEndOk

`func (o *NicFirewallRule) GetPortRangeEndOk() (*int32, bool)`

GetPortRangeEndOk returns a tuple with the PortRangeEnd field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortRangeEnd

`func (o *NicFirewallRule) SetPortRangeEnd(v int32)`

SetPortRangeEnd sets PortRangeEnd field to given value.

### HasPortRangeEnd

`func (o *NicFirewallRule) HasPortRangeEnd() bool`

HasPortRangeEnd returns a boolean if a field has been set.

### SetPortRangeEndNil

`func (o *NicFirewallRule) SetPortRangeEndNil(b bool)`

 SetPortRangeEndNil sets the value for PortRangeEnd to be an explicit nil

### UnsetPortRangeEnd
`func (o *NicFirewallRule) UnsetPortRangeEnd()`

UnsetPortRangeEnd ensures that no value is present for PortRangeEnd, not even an explicit nil
### GetType

`func (o *NicFirewallRule) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *NicFirewallRule) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *NicFirewallRule) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *NicFirewallRule) HasType() bool`

HasType returns a boolean if a field has been set.

### SetTypeNil

`func (o *NicFirewallRule) SetTypeNil(b bool)`

 SetTypeNil sets the value for Type to be an explicit nil

### UnsetType
`func (o *NicFirewallRule) UnsetType()`

UnsetType ensures that no value is present for Type, not even an explicit nil

