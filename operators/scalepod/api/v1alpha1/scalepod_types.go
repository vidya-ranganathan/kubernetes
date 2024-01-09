/*


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

// scalepodSpec defines the desired state of scalepod
type scalepodSpec struct {
	// Replicas is the desired number of pods for the scalepod
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10

	// json flag is added for serialisation..
	Replicas int32 `json:"replicas,omitempty"`
}

// scalepodStatus defines the current status of scalepod
type scalepodStatus struct {
	PodNames          []string `json:"podNames"`
	AvailableReplicas int32    `json:"availableReplicas"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// scalepod is the Schema for the scalepods API
// +kubebuilder:printcolumn:JSONPath=".spec.replicas",name=Desired,type=string
// +kubebuilder:printcolumn:JSONPath=".status.availableReplicas",name=Available,type=string
type scalepod struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   scalepodSpec   `json:"spec,omitempty"`
	Status scalepodStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// scalepodList contains a list of scalepod
type scalepodList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []scalepod `json:"items"`
}

func init() {
	SchemeBuilder.Register(&scalepod{}, &scalepodList{})
}
