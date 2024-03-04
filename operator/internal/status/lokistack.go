package status

import (
	"context"
	"fmt"

	"github.com/ViaQ/logerr/v2/kverrors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	ctrl "sigs.k8s.io/controller-runtime"

	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
	"github.com/grafana/loki/operator/internal/external/k8s"
)

const (
	messageReady   = "All components ready"
	messageFailed  = "One or more LokiStack components failed"
	messagePending = "One or more LokiStack components pending on dependencies"
	messageRunning = "All components are running, but some readiness checks are failing"
)

var (
	conditionFailed = metav1.Condition{
		Type:    string(lokiv1.ConditionFailed),
		Message: messageFailed,
		Reason:  string(lokiv1.ReasonFailedComponents),
	}
	conditionPending = metav1.Condition{
		Type:    string(lokiv1.ConditionPending),
		Message: messagePending,
		Reason:  string(lokiv1.ReasonPendingComponents),
	}
	conditionRunning = metav1.Condition{
		Type:    string(lokiv1.ConditionPending),
		Message: messageRunning,
		Reason:  string(lokiv1.ReasonPendingComponents),
	}
	conditionReady = metav1.Condition{
		Type:    string(lokiv1.ConditionReady),
		Message: messageReady,
		Reason:  string(lokiv1.ReasonReadyComponents),
	}
)

// DegradedError contains information about why the managed LokiStack has an invalid configuration.
type DegradedError struct {
	Message string
	Reason  lokiv1.LokiStackConditionReason
	Requeue bool
}

func (e *DegradedError) Error() string {
	return fmt.Sprintf("cluster degraded: %s", e.Message)
}

// SetDegradedCondition appends the condition Degraded to the lokistack status conditions.
func SetDegradedCondition(ctx context.Context, k k8s.Client, req ctrl.Request, msg string, reason lokiv1.LokiStackConditionReason) error {
	degraded := metav1.Condition{
		Type:    string(lokiv1.ConditionDegraded),
		Message: msg,
		Reason:  string(reason),
	}

	return updateCondition(ctx, k, req, degraded)
}

func generateCondition(cs *lokiv1.LokiStackComponentStatus) metav1.Condition {
	// Check for failed pods first
	failed := len(cs.Compactor[lokiv1.PodFailed]) +
		len(cs.Distributor[lokiv1.PodFailed]) +
		len(cs.Ingester[lokiv1.PodFailed]) +
		len(cs.Querier[lokiv1.PodFailed]) +
		len(cs.QueryFrontend[lokiv1.PodFailed]) +
		len(cs.Gateway[lokiv1.PodFailed]) +
		len(cs.IndexGateway[lokiv1.PodFailed]) +
		len(cs.Ruler[lokiv1.PodFailed])

	if failed != 0 {
		return conditionFailed
	}

	// Check for pending pods
	pending := len(cs.Compactor[lokiv1.PodPending]) +
		len(cs.Distributor[lokiv1.PodPending]) +
		len(cs.Ingester[lokiv1.PodPending]) +
		len(cs.Querier[lokiv1.PodPending]) +
		len(cs.QueryFrontend[lokiv1.PodPending]) +
		len(cs.Gateway[lokiv1.PodPending]) +
		len(cs.IndexGateway[lokiv1.PodPending]) +
		len(cs.Ruler[lokiv1.PodPending])

	if pending != 0 {
		return conditionPending
	}

	// Check if there are pods that are running but not ready
	running := len(cs.Compactor[lokiv1.PodRunning]) +
		len(cs.Distributor[lokiv1.PodRunning]) +
		len(cs.Ingester[lokiv1.PodRunning]) +
		len(cs.Querier[lokiv1.PodRunning]) +
		len(cs.QueryFrontend[lokiv1.PodRunning]) +
		len(cs.Gateway[lokiv1.PodRunning]) +
		len(cs.IndexGateway[lokiv1.PodRunning]) +
		len(cs.Ruler[lokiv1.PodRunning])

	if running > 0 {
		return conditionRunning
	}

	return conditionReady
}

func updateCondition(ctx context.Context, k k8s.Client, req ctrl.Request, condition metav1.Condition) error {
	var stack lokiv1.LokiStack
	if err := k.Get(ctx, req.NamespacedName, &stack); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return kverrors.Wrap(err, "failed to lookup LokiStack", "name", req.NamespacedName)
	}

	for _, c := range stack.Status.Conditions {
		if c.Type == condition.Type &&
			c.Reason == condition.Reason &&
			c.Message == condition.Message &&
			c.Status == metav1.ConditionTrue {
			// resource already has desired condition
			return nil
		}
	}

	condition.Status = metav1.ConditionTrue

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := k.Get(ctx, req.NamespacedName, &stack); err != nil {
			return err
		}

		now := metav1.Now()
		condition.LastTransitionTime = now

		index := -1
		for i := range stack.Status.Conditions {
			// Reset all other conditions first
			stack.Status.Conditions[i].Status = metav1.ConditionFalse
			stack.Status.Conditions[i].LastTransitionTime = now

			// Locate existing pending condition if any
			if stack.Status.Conditions[i].Type == condition.Type {
				index = i
			}
		}

		if index == -1 {
			stack.Status.Conditions = append(stack.Status.Conditions, condition)
		} else {
			stack.Status.Conditions[index] = condition
		}

		return k.Status().Update(ctx, &stack)
	})
}
