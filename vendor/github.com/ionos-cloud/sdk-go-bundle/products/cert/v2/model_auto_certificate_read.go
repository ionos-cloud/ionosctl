/*
 * Certificate Manager Service API
 *
 * Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates  with IONOS services and your internal connected resources.   For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic. The service provides the basic functions of uploading and deleting your certificates for this purpose.
 *
 * API version: 2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cert

import (
	"encoding/json"
)

// checks if the AutoCertificateRead type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AutoCertificateRead{}

// AutoCertificateRead struct for AutoCertificateRead
type AutoCertificateRead struct {
	// The ID (UUID) of the AutoCertificate.
	Id string `json:"id"`
	// The type of the resource.
	Type string `json:"type"`
	// The URL of the AutoCertificate.
	Href       string                                 `json:"href"`
	Metadata   MetadataWithAutoCertificateInformation `json:"metadata"`
	Properties AutoCertificate                        `json:"properties"`
}

// NewAutoCertificateRead instantiates a new AutoCertificateRead object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAutoCertificateRead(id string, type_ string, href string, metadata MetadataWithAutoCertificateInformation, properties AutoCertificate) *AutoCertificateRead {
	this := AutoCertificateRead{}

	this.Id = id
	this.Type = type_
	this.Href = href
	this.Metadata = metadata
	this.Properties = properties

	return &this
}

// NewAutoCertificateReadWithDefaults instantiates a new AutoCertificateRead object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAutoCertificateReadWithDefaults() *AutoCertificateRead {
	this := AutoCertificateRead{}
	return &this
}

// GetId returns the Id field value
func (o *AutoCertificateRead) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *AutoCertificateRead) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *AutoCertificateRead) SetId(v string) {
	o.Id = v
}

// GetType returns the Type field value
func (o *AutoCertificateRead) GetType() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Type
}

// GetTypeOk returns a tuple with the Type field value
// and a boolean to check if the value has been set.
func (o *AutoCertificateRead) GetTypeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Type, true
}

// SetType sets field value
func (o *AutoCertificateRead) SetType(v string) {
	o.Type = v
}

// GetHref returns the Href field value
func (o *AutoCertificateRead) GetHref() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Href
}

// GetHrefOk returns a tuple with the Href field value
// and a boolean to check if the value has been set.
func (o *AutoCertificateRead) GetHrefOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Href, true
}

// SetHref sets field value
func (o *AutoCertificateRead) SetHref(v string) {
	o.Href = v
}

// GetMetadata returns the Metadata field value
func (o *AutoCertificateRead) GetMetadata() MetadataWithAutoCertificateInformation {
	if o == nil {
		var ret MetadataWithAutoCertificateInformation
		return ret
	}

	return o.Metadata
}

// GetMetadataOk returns a tuple with the Metadata field value
// and a boolean to check if the value has been set.
func (o *AutoCertificateRead) GetMetadataOk() (*MetadataWithAutoCertificateInformation, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Metadata, true
}

// SetMetadata sets field value
func (o *AutoCertificateRead) SetMetadata(v MetadataWithAutoCertificateInformation) {
	o.Metadata = v
}

// GetProperties returns the Properties field value
func (o *AutoCertificateRead) GetProperties() AutoCertificate {
	if o == nil {
		var ret AutoCertificate
		return ret
	}

	return o.Properties
}

// GetPropertiesOk returns a tuple with the Properties field value
// and a boolean to check if the value has been set.
func (o *AutoCertificateRead) GetPropertiesOk() (*AutoCertificate, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Properties, true
}

// SetProperties sets field value
func (o *AutoCertificateRead) SetProperties(v AutoCertificate) {
	o.Properties = v
}

func (o AutoCertificateRead) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["id"] = o.Id
	toSerialize["type"] = o.Type
	toSerialize["href"] = o.Href
	toSerialize["metadata"] = o.Metadata
	toSerialize["properties"] = o.Properties
	return toSerialize, nil
}

type NullableAutoCertificateRead struct {
	value *AutoCertificateRead
	isSet bool
}

func (v NullableAutoCertificateRead) Get() *AutoCertificateRead {
	return v.value
}

func (v *NullableAutoCertificateRead) Set(val *AutoCertificateRead) {
	v.value = val
	v.isSet = true
}

func (v NullableAutoCertificateRead) IsSet() bool {
	return v.isSet
}

func (v *NullableAutoCertificateRead) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAutoCertificateRead(val *AutoCertificateRead) *NullableAutoCertificateRead {
	return &NullableAutoCertificateRead{value: val, isSet: true}
}

func (v NullableAutoCertificateRead) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAutoCertificateRead) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
