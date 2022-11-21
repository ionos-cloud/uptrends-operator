package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CRDResourceKind ...
	CRDResourceKind = "Uptrends"
	// AnnotationPrefix ...
	AnnotationPrefix = "uptrends.ionos-cloud.github.io/monitor."
)

func init() {
	SchemeBuilder.Register(&Uptrends{}, &UptrendsList{})
}

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "make generate" or "go generate ./..." to regenerate code after modifying this file

//+kubebuilder:object:root=true

// Uptrends is the Schema for the uptrends API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +operator-sdk:csv:customresourcedefinitions:resources={{Uptrends,v1alpha1,""}}
// +operator-sdk:csv:customresourcedefinitions:resources={{Ingress,v1,""}}
type Uptrends struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UptrendsSpec   `json:"spec,omitempty"`
	Status UptrendsStatus `json:"status,omitempty"`
}

// UptrendsSpec defines the desired state of Uptrends
// +k8s:openapi-gen=true
type UptrendsSpec struct {
	// Type of the Monitor.
	Type string `json:"type"`
	// Name of the Monitor.
	Name string `json:"name"`
	// Description of the Monitor.
	Description string `json:"description"`
	// Interval of the Monitor.
	Interval int `json:"interval"`
	// Url of the Monitor.
	Url string `json:"url"`
}

//+kubebuilder:object:root=true

// UptrendsList contains a list of Uptrends
type UptrendsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Uptrends `json:"items"`
}

type UptrendsPhase string

const (
	UptrendsPhaseNone     UptrendsPhase = ""
	UptrendsPhaseCreating UptrendsPhase = "Creating"
	UptrendsPhaseRunning  UptrendsPhase = "Running"
	UptrendsPhaseFailed   UptrendsPhase = "Failed"
)

// UptrendsStatus defines the observed state of Uptrends
// +k8s:openapi-gen=true
type UptrendsStatus struct {
	// Phase is the Uptrends running phase.
	Phase UptrendsPhase `json:"phase"`

	// ControlPaused indicates the operator pauses the control of
	// Uptrends.
	ControlPaused bool `json:"controlPaused,omitempty"`
}

// IsFailed ...
func (cs *UptrendsStatus) IsFailed() bool {
	if cs == nil {
		return false
	}

	return cs.Phase == UptrendsPhaseFailed
}

// SetPhase ...
func (cs *UptrendsStatus) SetPhase(p UptrendsPhase) {
	cs.Phase = p
}

// PauseControl ...
func (cs *UptrendsStatus) PauseControl() {
	cs.ControlPaused = true
}

// Control ...
func (cs *UptrendsStatus) Control() {
	cs.ControlPaused = false
}
