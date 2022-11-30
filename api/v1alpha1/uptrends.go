package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// CRDResourceKind ...
	CRDResourceKind = "Uptrends"
	// AnnotationPrefix ...
	AnnotationPrefix = "uptrends.ionos-cloud.github.io/monitor."
	// FinalizerName ...
	FinalizerName = "uptrends.ionos-cloud.github.io/finalizer"
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
	// MonitorGroup associates a monitor group.
	Group MonitorGroup `json:"group,omitempty"`
	// Checkpoints are the checkpoints to use for monitoring.
	Checkpoints MonitorCheckpoints `json:"checkpoints,omitempty"`
}

// MonitorCheckpoints defines the set of point of presence to check from.
type MonitorCheckpoints struct {
	// Regions is the set of entire regions to use.
	Regions []int32 `json:"regions,omitempty"`
	// Checkpoints are single point of presence to use.
	Checkpoints []int32 `json:"checkpoints,omitempty"`
	// ExcludeCheckpoints is a list of point of presence to execlude to use.
	ExcludeCheckpoints []int32 `json:"exclude,omitempty"`
}

// MonitorGroup defines a monitor group.
type MonitorGroup struct {
	// GUID is the id of the monitor group.
	GUID string `json:"guid"`
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

	// MonitorGuid is the ID of the Uptrends Monitor.
	MonitorGuid string `json:"monitorGuid,omitempty"`
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
