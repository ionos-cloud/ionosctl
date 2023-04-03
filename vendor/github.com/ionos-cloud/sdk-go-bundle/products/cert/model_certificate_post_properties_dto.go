/*
 * Certificate Manager Service API
 *
 * Using the Certificate Manager Service, you can conveniently provision and manage SSL certificates with IONOS services and your internal connected resources. For the [Application Load Balancer](https://api.ionos.com/docs/cloud/v6/#Application-Load-Balancers-get-datacenters-datacenterId-applicationloadbalancers), you usually need a certificate to encrypt your HTTPS traffic.  The service provides the basic functions of uploading and deleting your certificates for this purpose.
 *
 * API version: 1.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cert

import (
	"encoding/json"
)

// CertificatePostPropertiesDto struct for CertificatePostPropertiesDto
type CertificatePostPropertiesDto struct {
	// The certificate name.
	Name *string `json:"name"`
	// The certificate body.
	Certificate *string `json:"certificate"`
	// The certificate chain.
	CertificateChain *string `json:"certificateChain"`
	// The RSA private key is used for authentication and symmetric key exchange when establishing an SSL session. It is a part of the public key infrastructure generally used with SSL certificates.
	PrivateKey *string `json:"privateKey"`
}

// NewCertificatePostPropertiesDto instantiates a new CertificatePostPropertiesDto object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCertificatePostPropertiesDto(name string, certificate string, certificateChain string, privateKey string) *CertificatePostPropertiesDto {
	this := CertificatePostPropertiesDto{}

	this.Name = &name
	this.Certificate = &certificate
	this.CertificateChain = &certificateChain
	this.PrivateKey = &privateKey

	return &this
}

// NewCertificatePostPropertiesDtoWithDefaults instantiates a new CertificatePostPropertiesDto object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCertificatePostPropertiesDtoWithDefaults() *CertificatePostPropertiesDto {
	this := CertificatePostPropertiesDto{}
	return &this
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificatePostPropertiesDto) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePostPropertiesDto) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *CertificatePostPropertiesDto) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *CertificatePostPropertiesDto) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetCertificate returns the Certificate field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificatePostPropertiesDto) GetCertificate() *string {
	if o == nil {
		return nil
	}

	return o.Certificate

}

// GetCertificateOk returns a tuple with the Certificate field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePostPropertiesDto) GetCertificateOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Certificate, true
}

// SetCertificate sets field value
func (o *CertificatePostPropertiesDto) SetCertificate(v string) {

	o.Certificate = &v

}

// HasCertificate returns a boolean if a field has been set.
func (o *CertificatePostPropertiesDto) HasCertificate() bool {
	if o != nil && o.Certificate != nil {
		return true
	}

	return false
}

// GetCertificateChain returns the CertificateChain field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificatePostPropertiesDto) GetCertificateChain() *string {
	if o == nil {
		return nil
	}

	return o.CertificateChain

}

// GetCertificateChainOk returns a tuple with the CertificateChain field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePostPropertiesDto) GetCertificateChainOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.CertificateChain, true
}

// SetCertificateChain sets field value
func (o *CertificatePostPropertiesDto) SetCertificateChain(v string) {

	o.CertificateChain = &v

}

// HasCertificateChain returns a boolean if a field has been set.
func (o *CertificatePostPropertiesDto) HasCertificateChain() bool {
	if o != nil && o.CertificateChain != nil {
		return true
	}

	return false
}

// GetPrivateKey returns the PrivateKey field value
// If the value is explicit nil, the zero value for string will be returned
func (o *CertificatePostPropertiesDto) GetPrivateKey() *string {
	if o == nil {
		return nil
	}

	return o.PrivateKey

}

// GetPrivateKeyOk returns a tuple with the PrivateKey field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *CertificatePostPropertiesDto) GetPrivateKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.PrivateKey, true
}

// SetPrivateKey sets field value
func (o *CertificatePostPropertiesDto) SetPrivateKey(v string) {

	o.PrivateKey = &v

}

// HasPrivateKey returns a boolean if a field has been set.
func (o *CertificatePostPropertiesDto) HasPrivateKey() bool {
	if o != nil && o.PrivateKey != nil {
		return true
	}

	return false
}

func (o CertificatePostPropertiesDto) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.Certificate != nil {
		toSerialize["certificate"] = o.Certificate
	}

	if o.CertificateChain != nil {
		toSerialize["certificateChain"] = o.CertificateChain
	}

	if o.PrivateKey != nil {
		toSerialize["privateKey"] = o.PrivateKey
	}

	return json.Marshal(toSerialize)
}

type NullableCertificatePostPropertiesDto struct {
	value *CertificatePostPropertiesDto
	isSet bool
}

func (v NullableCertificatePostPropertiesDto) Get() *CertificatePostPropertiesDto {
	return v.value
}

func (v *NullableCertificatePostPropertiesDto) Set(val *CertificatePostPropertiesDto) {
	v.value = val
	v.isSet = true
}

func (v NullableCertificatePostPropertiesDto) IsSet() bool {
	return v.isSet
}

func (v *NullableCertificatePostPropertiesDto) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCertificatePostPropertiesDto(val *CertificatePostPropertiesDto) *NullableCertificatePostPropertiesDto {
	return &NullableCertificatePostPropertiesDto{value: val, isSet: true}
}

func (v NullableCertificatePostPropertiesDto) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCertificatePostPropertiesDto) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
