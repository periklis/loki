package controllers

import (
	"context"
	"errors"
	"time"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/handlers"
	"github.com/grafana/loki/operator/internal/status"
)

// LokiRuleReconciler reconciles a LokiRule object
type LokiRuleReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=loki.grafana.com,resources=lokirules,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=loki.grafana.com,resources=lokirules/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=loki.grafana.com,resources=lokirules/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the LokiRule object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *LokiRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	err := handlers.ValidateLokiRule(ctx, r.Log, req, r.Client)

	var invalid *status.ValidationError
	if errors.As(err, &invalid) {
		err = status.SetInvalidCondition(ctx, r.Client, req, invalid.Reason)
		if err != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, err
		}

		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, nil
	}

	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	err = status.SetValidCondition(ctx, r.Client, req)
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	var stacks lokiv1beta1.LokiStackList
	err = r.Client.List(ctx, &stacks, client.MatchingLabelsSelector{Selector: labels.Everything()})
	if err != nil {
		return ctrl.Result{
			Requeue:      true,
			RequeueAfter: time.Second,
		}, err
	}

	for _, s := range stacks.Items {
		ss := s.DeepCopy()
		if ss.Annotations == nil {
			ss.Annotations = make(map[string]string)
		}

		ss.Annotations["loki.grafana.com/rulesDiscoveredAt"] = time.Now().Format(time.RFC3339)

		if err := r.Client.Update(ctx, ss); err != nil {
			return ctrl.Result{
				Requeue:      true,
				RequeueAfter: time.Second,
			}, err
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *LokiRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1beta1.LokiRule{}, createOrUpdateOnlyPred).
		Complete(r)
}
