package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EvaluationDuration defines the type for Prometheus durations.
//
// +kubebuilder:validation:Pattern:="((([0-9]+)y)?(([0-9]+)w)?(([0-9]+)d)?(([0-9]+)h)?(([0-9]+)m)?(([0-9]+)s)?(([0-9]+)ms)?|0)"
type EvaluationDuration string

// LokiRuleSpec defines the desired state of LokiRule
type LokiRuleSpec struct {
	// List of groups for alerting and/or recording rules.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Groups"
	Groups []*LokiRuleGroup `json:"groups"`
}

// LokiRuleGroup defines a group of Loki alerting and/or recording rules.
type LokiRuleGroup struct {
	// Name defines a name of the present recoding/alerting rule. Must be unique
	// within all loki rules.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Name"
	Name string `json:"name"`

	// Interval defines the time interval between evaluation of the given
	// recoding rule.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Evaluation Interval"
	Interval EvaluationDuration `json:"interval"`

	// Limit defines the number of alerts an alerting rule and series a recording
	// rule can produce. 0 is no limit.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Limit of firing alerts"
	Limit int32 `json:"limit,omitempty"`

	// Rules defines a list of alerting and/or recording rules
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Rules"
	Rules []*LokiRuleGroupSpec `json:"rules"`
}

// LokiRuleGroupSpec defines the spec for a Loki alerting or recording rule.
type LokiRuleGroupSpec struct {
	// The name of the alert. Must be a valid label value.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Alert name"
	Alert string `json:"alert,omitempty"`

	// The name of the time series to output to. Must be a valid metric name.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Recording Metric Name"
	Record string `json:"record,omitempty"`

	// The LogQL expression to evaluate. Every evaluation cycle this is
	// evaluated at the current time, and all resultant time series become
	// pending/firing alerts.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="LogQL Expression"
	Expr string `json:"expr"`

	// Alerts are considered firing once they have been returned for this long.
	// Alerts which have not yet fired for long enough are considered pending.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Alert Firing Threshold"
	For EvaluationDuration `json:"for,omitempty"`

	// Annotations to add to each alert.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Alert Annotations"
	Annotations map[string]string `json:"annotations,omitempty"`

	// Labels to add to each alert.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Alert Labels"
	Labels map[string]string `json:"labels,omitempty"`
}

// LokiRuleStatus defines the observed state of LokiRule
type LokiRuleStatus struct {
	// Conditions of the LokiRule generation health.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=status,xDescriptors="urn:alm:descriptor:io.kubernetes.conditions"
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// LokiRule is the Schema for the lokirules API
// +operator-sdk:csv:customresourcedefinitions:displayName="LokiRule",resources={{LokiStack,v1beta1}}
type LokiRule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   LokiRuleSpec   `json:"spec,omitempty"`
	Status LokiRuleStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// LokiRuleList contains a list of LokiRule
type LokiRuleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []LokiRule `json:"items"`
}

func init() {
	SchemeBuilder.Register(&LokiRule{}, &LokiRuleList{})
}
