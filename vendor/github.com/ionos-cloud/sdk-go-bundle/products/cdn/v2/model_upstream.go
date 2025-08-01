/*
 * IONOS Cloud - CDN Distribution API
 *
 * This API manages CDN distributions.
 *
 * API version: 1.2.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cdn

import (
	"encoding/json"
)

// checks if the Upstream type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Upstream{}

// Upstream struct for Upstream
type Upstream struct {
	// The upstream host that handles the requests if not already cached. This host will be protected by the WAF if the option is enabled.
	Host string `json:"host"`
	// Enable or disable caching. If enabled, the CDN will cache the responses from the upstream host. Subsequent requests for the same resource will be served from the cache.
	Caching bool `json:"caching"`
	// Enable or disable WAF to protect the upstream host.
	Waf             bool                     `json:"waf"`
	GeoRestrictions *UpstreamGeoRestrictions `json:"geoRestrictions,omitempty"`
	// Rate limit class that will be applied to limit the number of incoming requests per IP.
	RateLimitClass string `json:"rateLimitClass"`
	// The SNI (Server Name Indication) mode of the upstream host. It supports two modes: - distribution: for outgoing connections to the upstream host, the CDN requires the upstream host to present a valid certificate that matches the configured domain of the CDN distribution. - origin: for outgoing connections to the upstream host, the CDN requires the upstream host to present a valid certificate that matches the configured upstream/origin hostname.
	SniMode string `json:"sniMode"`
}

// NewUpstream instantiates a new Upstream object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewUpstream(host string, caching bool, waf bool, rateLimitClass string, sniMode string) *Upstream {
	this := Upstream{}

	this.Host = host
	this.Caching = caching
	this.Waf = waf
	this.RateLimitClass = rateLimitClass
	this.SniMode = sniMode

	return &this
}

// NewUpstreamWithDefaults instantiates a new Upstream object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewUpstreamWithDefaults() *Upstream {
	this := Upstream{}
	return &this
}

// GetHost returns the Host field value
func (o *Upstream) GetHost() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Host
}

// GetHostOk returns a tuple with the Host field value
// and a boolean to check if the value has been set.
func (o *Upstream) GetHostOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Host, true
}

// SetHost sets field value
func (o *Upstream) SetHost(v string) {
	o.Host = v
}

// GetCaching returns the Caching field value
func (o *Upstream) GetCaching() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Caching
}

// GetCachingOk returns a tuple with the Caching field value
// and a boolean to check if the value has been set.
func (o *Upstream) GetCachingOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Caching, true
}

// SetCaching sets field value
func (o *Upstream) SetCaching(v bool) {
	o.Caching = v
}

// GetWaf returns the Waf field value
func (o *Upstream) GetWaf() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Waf
}

// GetWafOk returns a tuple with the Waf field value
// and a boolean to check if the value has been set.
func (o *Upstream) GetWafOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Waf, true
}

// SetWaf sets field value
func (o *Upstream) SetWaf(v bool) {
	o.Waf = v
}

// GetGeoRestrictions returns the GeoRestrictions field value if set, zero value otherwise.
func (o *Upstream) GetGeoRestrictions() UpstreamGeoRestrictions {
	if o == nil || IsNil(o.GeoRestrictions) {
		var ret UpstreamGeoRestrictions
		return ret
	}
	return *o.GeoRestrictions
}

// GetGeoRestrictionsOk returns a tuple with the GeoRestrictions field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Upstream) GetGeoRestrictionsOk() (*UpstreamGeoRestrictions, bool) {
	if o == nil || IsNil(o.GeoRestrictions) {
		return nil, false
	}
	return o.GeoRestrictions, true
}

// HasGeoRestrictions returns a boolean if a field has been set.
func (o *Upstream) HasGeoRestrictions() bool {
	if o != nil && !IsNil(o.GeoRestrictions) {
		return true
	}

	return false
}

// SetGeoRestrictions gets a reference to the given UpstreamGeoRestrictions and assigns it to the GeoRestrictions field.
func (o *Upstream) SetGeoRestrictions(v UpstreamGeoRestrictions) {
	o.GeoRestrictions = &v
}

// GetRateLimitClass returns the RateLimitClass field value
func (o *Upstream) GetRateLimitClass() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RateLimitClass
}

// GetRateLimitClassOk returns a tuple with the RateLimitClass field value
// and a boolean to check if the value has been set.
func (o *Upstream) GetRateLimitClassOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RateLimitClass, true
}

// SetRateLimitClass sets field value
func (o *Upstream) SetRateLimitClass(v string) {
	o.RateLimitClass = v
}

// GetSniMode returns the SniMode field value
func (o *Upstream) GetSniMode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SniMode
}

// GetSniModeOk returns a tuple with the SniMode field value
// and a boolean to check if the value has been set.
func (o *Upstream) GetSniModeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SniMode, true
}

// SetSniMode sets field value
func (o *Upstream) SetSniMode(v string) {
	o.SniMode = v
}

func (o Upstream) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Upstream) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["host"] = o.Host
	toSerialize["caching"] = o.Caching
	toSerialize["waf"] = o.Waf
	if !IsNil(o.GeoRestrictions) {
		toSerialize["geoRestrictions"] = o.GeoRestrictions
	}
	toSerialize["rateLimitClass"] = o.RateLimitClass
	toSerialize["sniMode"] = o.SniMode
	return toSerialize, nil
}

type NullableUpstream struct {
	value *Upstream
	isSet bool
}

func (v NullableUpstream) Get() *Upstream {
	return v.value
}

func (v *NullableUpstream) Set(val *Upstream) {
	v.value = val
	v.isSet = true
}

func (v NullableUpstream) IsSet() bool {
	return v.isSet
}

func (v *NullableUpstream) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableUpstream(val *Upstream) *NullableUpstream {
	return &NullableUpstream{value: val, isSet: true}
}

func (v NullableUpstream) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableUpstream) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
