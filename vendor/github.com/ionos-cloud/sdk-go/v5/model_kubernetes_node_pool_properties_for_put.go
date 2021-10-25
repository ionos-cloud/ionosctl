/*
 * CLOUD API
 *
 * An enterprise-grade Infrastructure is provided as a Service (IaaS) solution that can be managed through a browser-based \"Data Center Designer\" (DCD) tool or via an easy to use API.   The API allows you to perform a variety of management tasks such as spinning up additional servers, adding volumes, adjusting networking, and so forth. It is designed to allow users to leverage the same power and flexibility found within the DCD visual tool. Both tools are consistent with their concepts and lend well to making the experience smooth and intuitive.
 *
 * API version: 5.0
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package ionoscloud

import (
	"encoding/json"
)

// KubernetesNodePoolPropertiesForPut struct for KubernetesNodePoolPropertiesForPut
type KubernetesNodePoolPropertiesForPut struct {
	// A Kubernetes Node Pool Name. Valid Kubernetes Node Pool name must be 63 characters or less and must be empty or begin and end with an alphanumeric character ([a-z0-9A-Z]) with dashes (-), underscores (_), dots (.), and alphanumerics between.
	Name *string `json:"name"`
	// Number of nodes part of the Node Pool
	NodeCount *int32 `json:"nodeCount"`
	// The kubernetes version in which a nodepool is running. This imposes restrictions on what kubernetes versions can be run in a cluster's nodepools. Additionally, not all kubernetes versions are viable upgrade targets for all prior versions.
	K8sVersion        *string                      `json:"k8sVersion,omitempty"`
	MaintenanceWindow *KubernetesMaintenanceWindow `json:"maintenanceWindow,omitempty"`
	AutoScaling       *KubernetesAutoScaling       `json:"autoScaling,omitempty"`
	// array of additional LANs attached to worker nodes
	Lans *[]KubernetesNodePoolLan `json:"lans,omitempty"`
	// map of labels attached to node pool
	Labels *map[string]string `json:"labels,omitempty"`
	// map of annotations attached to node pool
	Annotations *map[string]string `json:"annotations,omitempty"`
	// Optional array of reserved public IP addresses to be used by the nodes. IPs must be from same location as the data center used for the node pool. The array must contain one extra IP than maximum number of nodes could be. (nodeCount+1 if fixed node amount or maxNodeCount+1 if auto scaling is used) The extra provided IP Will be used during rebuilding of nodes.
	PublicIps *[]string `json:"publicIps,omitempty"`
}

// GetName returns the Name field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetName() *string {
	if o == nil {
		return nil
	}

	return o.Name

}

// GetNameOk returns a tuple with the Name field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Name, true
}

// SetName sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetName(v string) {

	o.Name = &v

}

// HasName returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// GetNodeCount returns the NodeCount field value
// If the value is explicit nil, the zero value for int32 will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetNodeCount() *int32 {
	if o == nil {
		return nil
	}

	return o.NodeCount

}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}

	return o.NodeCount, true
}

// SetNodeCount sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetNodeCount(v int32) {

	o.NodeCount = &v

}

// HasNodeCount returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasNodeCount() bool {
	if o != nil && o.NodeCount != nil {
		return true
	}

	return false
}

// GetK8sVersion returns the K8sVersion field value
// If the value is explicit nil, the zero value for string will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetK8sVersion() *string {
	if o == nil {
		return nil
	}

	return o.K8sVersion

}

// GetK8sVersionOk returns a tuple with the K8sVersion field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetK8sVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}

	return o.K8sVersion, true
}

// SetK8sVersion sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetK8sVersion(v string) {

	o.K8sVersion = &v

}

// HasK8sVersion returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasK8sVersion() bool {
	if o != nil && o.K8sVersion != nil {
		return true
	}

	return false
}

// GetMaintenanceWindow returns the MaintenanceWindow field value
// If the value is explicit nil, the zero value for KubernetesMaintenanceWindow will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetMaintenanceWindow() *KubernetesMaintenanceWindow {
	if o == nil {
		return nil
	}

	return o.MaintenanceWindow

}

// GetMaintenanceWindowOk returns a tuple with the MaintenanceWindow field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetMaintenanceWindowOk() (*KubernetesMaintenanceWindow, bool) {
	if o == nil {
		return nil, false
	}

	return o.MaintenanceWindow, true
}

// SetMaintenanceWindow sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetMaintenanceWindow(v KubernetesMaintenanceWindow) {

	o.MaintenanceWindow = &v

}

// HasMaintenanceWindow returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasMaintenanceWindow() bool {
	if o != nil && o.MaintenanceWindow != nil {
		return true
	}

	return false
}

// GetAutoScaling returns the AutoScaling field value
// If the value is explicit nil, the zero value for KubernetesAutoScaling will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetAutoScaling() *KubernetesAutoScaling {
	if o == nil {
		return nil
	}

	return o.AutoScaling

}

// GetAutoScalingOk returns a tuple with the AutoScaling field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetAutoScalingOk() (*KubernetesAutoScaling, bool) {
	if o == nil {
		return nil, false
	}

	return o.AutoScaling, true
}

// SetAutoScaling sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetAutoScaling(v KubernetesAutoScaling) {

	o.AutoScaling = &v

}

// HasAutoScaling returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasAutoScaling() bool {
	if o != nil && o.AutoScaling != nil {
		return true
	}

	return false
}

// GetLans returns the Lans field value
// If the value is explicit nil, the zero value for []KubernetesNodePoolLan will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetLans() *[]KubernetesNodePoolLan {
	if o == nil {
		return nil
	}

	return o.Lans

}

// GetLansOk returns a tuple with the Lans field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetLansOk() (*[]KubernetesNodePoolLan, bool) {
	if o == nil {
		return nil, false
	}

	return o.Lans, true
}

// SetLans sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetLans(v []KubernetesNodePoolLan) {

	o.Lans = &v

}

// HasLans returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasLans() bool {
	if o != nil && o.Lans != nil {
		return true
	}

	return false
}

// GetLabels returns the Labels field value
// If the value is explicit nil, the zero value for map[string]string will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetLabels() *map[string]string {
	if o == nil {
		return nil
	}

	return o.Labels

}

// GetLabelsOk returns a tuple with the Labels field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetLabelsOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Labels, true
}

// SetLabels sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetLabels(v map[string]string) {

	o.Labels = &v

}

// HasLabels returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasLabels() bool {
	if o != nil && o.Labels != nil {
		return true
	}

	return false
}

// GetAnnotations returns the Annotations field value
// If the value is explicit nil, the zero value for map[string]string will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetAnnotations() *map[string]string {
	if o == nil {
		return nil
	}

	return o.Annotations

}

// GetAnnotationsOk returns a tuple with the Annotations field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetAnnotationsOk() (*map[string]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.Annotations, true
}

// SetAnnotations sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetAnnotations(v map[string]string) {

	o.Annotations = &v

}

// HasAnnotations returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasAnnotations() bool {
	if o != nil && o.Annotations != nil {
		return true
	}

	return false
}

// GetPublicIps returns the PublicIps field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetPublicIps() *[]string {
	if o == nil {
		return nil
	}

	return o.PublicIps

}

// GetPublicIpsOk returns a tuple with the PublicIps field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *KubernetesNodePoolPropertiesForPut) GetPublicIpsOk() (*[]string, bool) {
	if o == nil {
		return nil, false
	}

	return o.PublicIps, true
}

// SetPublicIps sets field value
func (o *KubernetesNodePoolPropertiesForPut) SetPublicIps(v []string) {

	o.PublicIps = &v

}

// HasPublicIps returns a boolean if a field has been set.
func (o *KubernetesNodePoolPropertiesForPut) HasPublicIps() bool {
	if o != nil && o.PublicIps != nil {
		return true
	}

	return false
}

func (o KubernetesNodePoolPropertiesForPut) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}

	if o.Name != nil {
		toSerialize["name"] = o.Name
	}

	if o.NodeCount != nil {
		toSerialize["nodeCount"] = o.NodeCount
	}

	if o.K8sVersion != nil {
		toSerialize["k8sVersion"] = o.K8sVersion
	}

	if o.MaintenanceWindow != nil {
		toSerialize["maintenanceWindow"] = o.MaintenanceWindow
	}

	if o.AutoScaling != nil {
		toSerialize["autoScaling"] = o.AutoScaling
	}

	if o.Lans != nil {
		toSerialize["lans"] = o.Lans
	}

	if o.Labels != nil {
		toSerialize["labels"] = o.Labels
	}

	if o.Annotations != nil {
		toSerialize["annotations"] = o.Annotations
	}

	if o.PublicIps != nil {
		toSerialize["publicIps"] = o.PublicIps
	}
	return json.Marshal(toSerialize)
}

type NullableKubernetesNodePoolPropertiesForPut struct {
	value *KubernetesNodePoolPropertiesForPut
	isSet bool
}

func (v NullableKubernetesNodePoolPropertiesForPut) Get() *KubernetesNodePoolPropertiesForPut {
	return v.value
}

func (v *NullableKubernetesNodePoolPropertiesForPut) Set(val *KubernetesNodePoolPropertiesForPut) {
	v.value = val
	v.isSet = true
}

func (v NullableKubernetesNodePoolPropertiesForPut) IsSet() bool {
	return v.isSet
}

func (v *NullableKubernetesNodePoolPropertiesForPut) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableKubernetesNodePoolPropertiesForPut(val *KubernetesNodePoolPropertiesForPut) *NullableKubernetesNodePoolPropertiesForPut {
	return &NullableKubernetesNodePoolPropertiesForPut{value: val, isSet: true}
}

func (v NullableKubernetesNodePoolPropertiesForPut) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableKubernetesNodePoolPropertiesForPut) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
