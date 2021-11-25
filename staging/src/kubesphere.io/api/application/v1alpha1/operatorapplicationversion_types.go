/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// OperatorApplicationVersionSpec defines the desired state of OperatorApplicationVersion
type OperatorApplicationVersionSpec struct {
	AppName         string `json:"name"`
	Screenshots     string `json:"screenshots,omitempty"`
	ScreenshotsZh   string `json:"screenshots_zh,omitempty"`
	ChangeLog       string `json:"changeLog"`
	ChangeLogZh     string `json:"changeLog_zh,omitempty"`
	OperatorVersion string `json:"operatorVersion"`
	AppVersion      string `json:"appVersion"`
	Owner           string `json:"owner,omitempty"`
}

// OperatorApplicationVersionStatus defines the observed state of OperatorApplicationVersion
type OperatorApplicationVersionStatus struct {
	State string `json:"state,omitempty"`
}

//+genclient
//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+genclient:nonNamespaced
//+k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatorApplicationVersion is the Schema for the operatorapplicationversions API
type OperatorApplicationVersion struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorApplicationVersionSpec   `json:"spec,omitempty"`
	Status OperatorApplicationVersionStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// OperatorApplicationVersionList contains a list of OperatorApplicationVersion
type OperatorApplicationVersionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OperatorApplicationVersion `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OperatorApplicationVersion{}, &OperatorApplicationVersionList{})
}

func (in *OperatorApplicationVersion) GetVersionName() string {
	return in.Spec.OperatorVersion
}
