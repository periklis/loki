package status

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/kverrors"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/external/k8s"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ValidationError contains information about why the managed LokiRule has an invalid configuration.
type ValidationError struct {
	Message string
	Reason  lokiv1beta1.LokiRuleConditionReason
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid rule: %s", e.Message)
}

// SetValidCondition updates or appends the condition Valid to the LokiRule status conditions.
// In addition it resets all other Status conditions to false.
func SetValidCondition(ctx context.Context, k k8s.Client, req ctrl.Request) error {
	var r lokiv1beta1.LokiRule
	if err := k.Get(ctx, req.NamespacedName, &r); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return kverrors.Wrap(err, "failed to lookup lokirule", "name", req.NamespacedName)
	}

	for _, cond := range r.Status.Conditions {
		if cond.Type == string(lokiv1beta1.ConditionValid) && cond.Status == metav1.ConditionTrue {
			return nil
		}
	}

	valid := metav1.Condition{
		Type:               string(lokiv1beta1.ConditionValid),
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Message:            "All rules valid",
		Reason:             string(lokiv1beta1.ReasonAllRulesValid),
	}

	index := -1
	for i := range r.Status.Conditions {
		// Reset all other conditions first
		r.Status.Conditions[i].Status = metav1.ConditionFalse
		r.Status.Conditions[i].LastTransitionTime = metav1.Now()

		// Locate existing ready condition if any
		if r.Status.Conditions[i].Type == string(lokiv1beta1.ConditionValid) {
			index = i
		}
	}

	if index == -1 {
		r.Status.Conditions = append(r.Status.Conditions, valid)
	} else {
		r.Status.Conditions[index] = valid
	}

	return k.Status().Update(ctx, &r, &client.UpdateOptions{})
}

// SetInvalidCondition updates or appends the condition Invalid to the LokiRule status conditions.
// In addition it resets all other Status conditions to false.
func SetInvalidCondition(ctx context.Context, k k8s.Client, req ctrl.Request, reason lokiv1beta1.LokiRuleConditionReason) error {
	var r lokiv1beta1.LokiRule
	if err := k.Get(ctx, req.NamespacedName, &r); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return kverrors.Wrap(err, "failed to lookup lokirule", "name", req.NamespacedName)
	}

	reasonStr := string(reason)
	for _, cond := range r.Status.Conditions {
		if cond.Type == string(lokiv1beta1.ConditionInvalid) && cond.Reason == reasonStr && cond.Status == metav1.ConditionTrue {
			return nil
		}
	}

	invalid := metav1.Condition{
		Type:               string(lokiv1beta1.ConditionInvalid),
		Status:             metav1.ConditionTrue,
		LastTransitionTime: metav1.Now(),
		Message:            "Invalid Loki rules",
		Reason:             reasonStr,
	}

	index := -1
	for i := range r.Status.Conditions {
		// Reset all other conditions first
		r.Status.Conditions[i].Status = metav1.ConditionFalse
		r.Status.Conditions[i].LastTransitionTime = metav1.Now()

		// Locate existing ready condition if any
		if r.Status.Conditions[i].Type == string(lokiv1beta1.ConditionInvalid) {
			index = i
		}
	}

	if index == -1 {
		r.Status.Conditions = append(r.Status.Conditions, invalid)
	} else {
		r.Status.Conditions[index] = invalid
	}

	return k.Status().Update(ctx, &r, &client.UpdateOptions{})
}
