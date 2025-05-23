/*
 * CLOUD API
 *
 *  IONOS Enterprise-grade Infrastructure as a Service (IaaS) solutions can be managed through the Cloud API, in addition or as an alternative to the \"Data Center Designer\" (DCD) browser-based tool.    Both methods employ consistent concepts and features, deliver similar power and flexibility, and can be used to perform a multitude of management tasks, including adding servers, volumes, configuring networks, and so on.
 *
 * API version: 6.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"
)

var (
	// used to set a nullable field to nil. This is a sentinel address that will be checked in the MarshalJson function.
	// if set to this address, a nil value will be marshalled
	Nilstring string = "<<ExplicitNil>>"
	Nilint32  int32  = -334455
	Nilbool   bool   = false
)

// ToPtr - returns a pointer to the given value.
func ToPtr[T any](v T) *T {
	return &v
}

// ToValue - returns the value of the pointer passed in
func ToValue[T any](ptr *T) T {
	return *ptr
}

// ToValueDefault - returns the value of the pointer passed in, or the default type value if the pointer is nil
func ToValueDefault[T any](ptr *T) T {
	var defaultVal T
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

func SliceToValueDefault[T any](ptrSlice *[]T) []T {
	return append([]T{}, *ptrSlice...)
}

// PtrBool - returns a pointer to given boolean value.
func PtrBool(v bool) *bool { return &v }

// PtrInt - returns a pointer to given integer value.
func PtrInt(v int) *int { return &v }

// PtrInt32 - returns a pointer to given integer value.
func PtrInt32(v int32) *int32 { return &v }

// PtrInt64 - returns a pointer to given integer value.
func PtrInt64(v int64) *int64 { return &v }

// PtrFloat32 - returns a pointer to given float value.
func PtrFloat32(v float32) *float32 { return &v }

// PtrFloat64 - returns a pointer to given float value.
func PtrFloat64(v float64) *float64 { return &v }

// PtrString - returns a pointer to given string value.
func PtrString(v string) *string { return &v }

// PtrTime - returns a pointer to given Time value.
func PtrTime(v time.Time) *time.Time { return &v }

// ToBool - returns the value of the bool pointer passed in
func ToBool(ptr *bool) bool {
	return *ptr
}

// ToBoolDefault - returns the value of the bool pointer passed in, or false if the pointer is nil
func ToBoolDefault(ptr *bool) bool {
	var defaultVal bool
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToBoolSlice - returns a bool slice of the pointer passed in
func ToBoolSlice(ptrSlice *[]bool) []bool {
	valSlice := make([]bool, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToByte - returns the value of the byte pointer passed in
func ToByte(ptr *byte) byte {
	return *ptr
}

// ToByteDefault - returns the value of the byte pointer passed in, or 0 if the pointer is nil
func ToByteDefault(ptr *byte) byte {
	var defaultVal byte
	if ptr == nil {
		return defaultVal
	}

	return *ptr
}

// ToByteSlice - returns a byte slice of the pointer passed in
func ToByteSlice(ptrSlice *[]byte) []byte {
	valSlice := make([]byte, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToString - returns the value of the string pointer passed in
func ToString(ptr *string) string {
	return *ptr
}

// ToStringDefault - returns the value of the string pointer passed in, or "" if the pointer is nil
func ToStringDefault(ptr *string) string {
	var defaultVal string
	if ptr == nil {
		return defaultVal
	}

	return *ptr
}

// ToStringSlice - returns a string slice of the pointer passed in
func ToStringSlice(ptrSlice *[]string) []string {
	valSlice := make([]string, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt - returns the value of the int pointer passed in
func ToInt(ptr *int) int {
	return *ptr
}

// ToIntDefault - returns the value of the int pointer passed in, or 0 if the pointer is nil
func ToIntDefault(ptr *int) int {
	var defaultVal int
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToIntSlice - returns a int slice of the pointer passed in
func ToIntSlice(ptrSlice *[]int) []int {
	valSlice := make([]int, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt8 - returns the value of the int8 pointer passed in
func ToInt8(ptr *int8) int8 {
	return *ptr
}

// ToInt8Default - returns the value of the int8 pointer passed in, or 0 if the pointer is nil
func ToInt8Default(ptr *int8) int8 {
	var defaultVal int8
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt8Slice - returns a int8 slice of the pointer passed in
func ToInt8Slice(ptrSlice *[]int8) []int8 {
	valSlice := make([]int8, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt16 - returns the value of the int16 pointer passed in
func ToInt16(ptr *int16) int16 {
	return *ptr
}

// ToInt16Default - returns the value of the int16 pointer passed in, or 0 if the pointer is nil
func ToInt16Default(ptr *int16) int16 {
	var defaultVal int16
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt16Slice - returns a int16 slice of the pointer passed in
func ToInt16Slice(ptrSlice *[]int16) []int16 {
	valSlice := make([]int16, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt32 - returns the value of the int32 pointer passed in
func ToInt32(ptr *int32) int32 {
	return *ptr
}

// ToInt32Default - returns the value of the int32 pointer passed in, or 0 if the pointer is nil
func ToInt32Default(ptr *int32) int32 {
	var defaultVal int32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt32Slice - returns a int32 slice of the pointer passed in
func ToInt32Slice(ptrSlice *[]int32) []int32 {
	valSlice := make([]int32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToInt64 - returns the value of the int64 pointer passed in
func ToInt64(ptr *int64) int64 {
	return *ptr
}

// ToInt64Default - returns the value of the int64 pointer passed in, or 0 if the pointer is nil
func ToInt64Default(ptr *int64) int64 {
	var defaultVal int64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToInt64Slice - returns a int64 slice of the pointer passed in
func ToInt64Slice(ptrSlice *[]int64) []int64 {
	valSlice := make([]int64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint - returns the value of the uint pointer passed in
func ToUint(ptr *uint) uint {
	return *ptr
}

// ToUintDefault - returns the value of the uint pointer passed in, or 0 if the pointer is nil
func ToUintDefault(ptr *uint) uint {
	var defaultVal uint
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUintSlice - returns a uint slice of the pointer passed in
func ToUintSlice(ptrSlice *[]uint) []uint {
	valSlice := make([]uint, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint8 -returns the value of the uint8 pointer passed in
func ToUint8(ptr *uint8) uint8 {
	return *ptr
}

// ToUint8Default - returns the value of the uint8 pointer passed in, or 0 if the pointer is nil
func ToUint8Default(ptr *uint8) uint8 {
	var defaultVal uint8
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint8Slice - returns a uint8 slice of the pointer passed in
func ToUint8Slice(ptrSlice *[]uint8) []uint8 {
	valSlice := make([]uint8, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint16 - returns the value of the uint16 pointer passed in
func ToUint16(ptr *uint16) uint16 {
	return *ptr
}

// ToUint16Default - returns the value of the uint16 pointer passed in, or 0 if the pointer is nil
func ToUint16Default(ptr *uint16) uint16 {
	var defaultVal uint16
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint16Slice - returns a uint16 slice of the pointer passed in
func ToUint16Slice(ptrSlice *[]uint16) []uint16 {
	valSlice := make([]uint16, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint32 - returns the value of the uint32 pointer passed in
func ToUint32(ptr *uint32) uint32 {
	return *ptr
}

// ToUint32Default - returns the value of the uint32 pointer passed in, or 0 if the pointer is nil
func ToUint32Default(ptr *uint32) uint32 {
	var defaultVal uint32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint32Slice - returns a uint32 slice of the pointer passed in
func ToUint32Slice(ptrSlice *[]uint32) []uint32 {
	valSlice := make([]uint32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToUint64 - returns the value of the uint64 pointer passed in
func ToUint64(ptr *uint64) uint64 {
	return *ptr
}

// ToUint64Default - returns the value of the uint64 pointer passed in, or 0 if the pointer is nil
func ToUint64Default(ptr *uint64) uint64 {
	var defaultVal uint64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToUint64Slice - returns a uint63 slice of the pointer passed in
func ToUint64Slice(ptrSlice *[]uint64) []uint64 {
	valSlice := make([]uint64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToFloat32 - returns the value of the float32 pointer passed in
func ToFloat32(ptr *float32) float32 {
	return *ptr
}

// ToFloat32Default - returns the value of the float32 pointer passed in, or 0 if the pointer is nil
func ToFloat32Default(ptr *float32) float32 {
	var defaultVal float32
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToFloat32Slice - returns a float32 slice of the pointer passed in
func ToFloat32Slice(ptrSlice *[]float32) []float32 {
	valSlice := make([]float32, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToFloat64 - returns the value of the float64 pointer passed in
func ToFloat64(ptr *float64) float64 {
	return *ptr
}

// ToFloat64Default - returns the value of the float64 pointer passed in, or 0 if the pointer is nil
func ToFloat64Default(ptr *float64) float64 {
	var defaultVal float64
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToFloat64Slice - returns a float64 slice of the pointer passed in
func ToFloat64Slice(ptrSlice *[]float64) []float64 {
	valSlice := make([]float64, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

// ToTime - returns the value of the Time pointer passed in
func ToTime(ptr *time.Time) time.Time {
	return *ptr
}

// ToTimeDefault - returns the value of the Time pointer passed in, or 0001-01-01 00:00:00 +0000 UTC if the pointer is nil
func ToTimeDefault(ptr *time.Time) time.Time {
	var defaultVal time.Time
	if ptr == nil {
		return defaultVal
	}
	return *ptr
}

// ToTimeSlice - returns a Time slice of the pointer passed in
func ToTimeSlice(ptrSlice *[]time.Time) []time.Time {
	valSlice := make([]time.Time, len(*ptrSlice))
	for i, v := range *ptrSlice {
		valSlice[i] = v
	}

	return valSlice
}

type NullableBool struct {
	value *bool
	isSet bool
}

func (v NullableBool) Get() *bool {
	return v.value
}

func (v *NullableBool) Set(val *bool) {
	v.value = val
	v.isSet = true
}

func (v NullableBool) IsSet() bool {
	return v.isSet
}

func (v *NullableBool) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBool(val *bool) *NullableBool {
	return &NullableBool{value: val, isSet: true}
}

func (v NullableBool) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBool) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt struct {
	value *int
	isSet bool
}

func (v NullableInt) Get() *int {
	return v.value
}

func (v *NullableInt) Set(val *int) {
	v.value = val
	v.isSet = true
}

func (v NullableInt) IsSet() bool {
	return v.isSet
}

func (v *NullableInt) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt(val *int) *NullableInt {
	return &NullableInt{value: val, isSet: true}
}

func (v NullableInt) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt32 struct {
	value *int32
	isSet bool
}

func (v NullableInt32) Get() *int32 {
	return v.value
}

func (v *NullableInt32) Set(val *int32) {
	v.value = val
	v.isSet = true
}

func (v NullableInt32) IsSet() bool {
	return v.isSet
}

func (v *NullableInt32) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt32(val *int32) *NullableInt32 {
	return &NullableInt32{value: val, isSet: true}
}

func (v NullableInt32) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt32) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableInt64 struct {
	value *int64
	isSet bool
}

func (v NullableInt64) Get() *int64 {
	return v.value
}

func (v *NullableInt64) Set(val *int64) {
	v.value = val
	v.isSet = true
}

func (v NullableInt64) IsSet() bool {
	return v.isSet
}

func (v *NullableInt64) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInt64(val *int64) *NullableInt64 {
	return &NullableInt64{value: val, isSet: true}
}

func (v NullableInt64) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInt64) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableFloat32 struct {
	value *float32
	isSet bool
}

func (v NullableFloat32) Get() *float32 {
	return v.value
}

func (v *NullableFloat32) Set(val *float32) {
	v.value = val
	v.isSet = true
}

func (v NullableFloat32) IsSet() bool {
	return v.isSet
}

func (v *NullableFloat32) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFloat32(val *float32) *NullableFloat32 {
	return &NullableFloat32{value: val, isSet: true}
}

func (v NullableFloat32) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFloat32) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableFloat64 struct {
	value *float64
	isSet bool
}

func (v NullableFloat64) Get() *float64 {
	return v.value
}

func (v *NullableFloat64) Set(val *float64) {
	v.value = val
	v.isSet = true
}

func (v NullableFloat64) IsSet() bool {
	return v.isSet
}

func (v *NullableFloat64) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableFloat64(val *float64) *NullableFloat64 {
	return &NullableFloat64{value: val, isSet: true}
}

func (v NullableFloat64) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableFloat64) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableString struct {
	value *string
	isSet bool
}

func (v NullableString) Get() *string {
	return v.value
}

func (v *NullableString) Set(val *string) {
	v.value = val
	v.isSet = true
}

func (v NullableString) IsSet() bool {
	return v.isSet
}

func (v *NullableString) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableString(val *string) *NullableString {
	return &NullableString{value: val, isSet: true}
}

func (v NullableString) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableString) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type NullableTime struct {
	value *time.Time
	isSet bool
}

func (v NullableTime) Get() *time.Time {
	return v.value
}

func (v *NullableTime) Set(val *time.Time) {
	v.value = val
	v.isSet = true
}

func (v NullableTime) IsSet() bool {
	return v.isSet
}

func (v *NullableTime) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTime(val *time.Time) *NullableTime {
	return &NullableTime{value: val, isSet: true}
}

func (v NullableTime) MarshalJSON() ([]byte, error) {
	return v.value.MarshalJSON()
}

func (v *NullableTime) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}

type IonosTime struct {
	time.Time
}

func (t *IonosTime) UnmarshalJSON(data []byte) error {
	str := string(data)
	if strlen(str) == 0 {
		t = nil
		return nil
	}
	if str[0] == '"' {
		str = str[1:]
	}
	if str[len(str)-1] == '"' {
		str = str[:len(str)-1]
	}
	if !strings.Contains(str, "Z") {
		/* forcefully adding timezone suffix to be able to parse the
		 * string using RFC3339 */
		str += "Z"
	}
	tt, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return err
	}
	*t = IonosTime{tt}
	return nil
}

// IsNil checks if an input is nil
func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	case reflect.Array:
		return reflect.ValueOf(i).IsZero()
	}
	return false
}

// EnsureURLFormat checks that the URL has the correct format (no trailing slash,
// has http/https scheme prefix) and updates it if necessary
func EnsureURLFormat(url string) string {
	length := len(url)

	if length <= 1 {
		return url
	}

	if url[length-1] == '/' {
		url = url[:length-1]
	}

	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = fmt.Sprintf("https://%s", url)
	}

	return url
}
