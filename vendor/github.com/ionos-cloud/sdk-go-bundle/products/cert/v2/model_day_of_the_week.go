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

	"fmt"
)

// DayOfTheWeek The name of the week day.
type DayOfTheWeek string

// List of DayOfTheWeek
const (
	DAYOFTHEWEEK_SUNDAY    DayOfTheWeek = "Sunday"
	DAYOFTHEWEEK_MONDAY    DayOfTheWeek = "Monday"
	DAYOFTHEWEEK_TUESDAY   DayOfTheWeek = "Tuesday"
	DAYOFTHEWEEK_WEDNESDAY DayOfTheWeek = "Wednesday"
	DAYOFTHEWEEK_THURSDAY  DayOfTheWeek = "Thursday"
	DAYOFTHEWEEK_FRIDAY    DayOfTheWeek = "Friday"
	DAYOFTHEWEEK_SATURDAY  DayOfTheWeek = "Saturday"
)

func (v *DayOfTheWeek) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DayOfTheWeek(value)
	for _, existing := range []DayOfTheWeek{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"} {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DayOfTheWeek", value)
}

// Ptr returns reference to DayOfTheWeek value
func (v DayOfTheWeek) Ptr() *DayOfTheWeek {
	return &v
}

type NullableDayOfTheWeek struct {
	value *DayOfTheWeek
	isSet bool
}

func (v NullableDayOfTheWeek) Get() *DayOfTheWeek {
	return v.value
}

func (v *NullableDayOfTheWeek) Set(val *DayOfTheWeek) {
	v.value = val
	v.isSet = true
}

func (v NullableDayOfTheWeek) IsSet() bool {
	return v.isSet
}

func (v *NullableDayOfTheWeek) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDayOfTheWeek(val *DayOfTheWeek) *NullableDayOfTheWeek {
	return &NullableDayOfTheWeek{value: val, isSet: true}
}

func (v NullableDayOfTheWeek) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDayOfTheWeek) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
