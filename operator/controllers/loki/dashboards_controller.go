package controllers

import (
	"context"

	"github.com/ViaQ/logerr/v2/kverrors"
	"github.com/go-logr/logr"
	lokiv1 "github.com/grafana/loki/operator/apis/loki/v1"
	"github.com/grafana/loki/operator/internal/handlers"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

var createOrDeletesPred = builder.WithPredicates(predicate.Funcs{
	UpdateFunc:  func(e event.UpdateEvent) bool { return false },
	CreateFunc:  func(e event.CreateEvent) bool { return true },
	DeleteFunc:  func(e event.DeleteEvent) bool { return true },
	GenericFunc: func(e event.GenericEvent) bool { return true },
})

// DashboardsReconciler cleans up all remaining
// dashboards when all LokiStack objects are removed.
type DashboardsReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

// Reconcile creates all LokiStack dashboard ConfigMap and PrometheusRule objects on OpenShift clusters when
// the at least one LokiStack custom resource exists or removes all when none.
func (r *DashboardsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var stacks lokiv1.LokiStackList
	if err := r.List(ctx, &stacks, client.MatchingLabelsSelector{Selector: labels.Everything()}); err != nil {
		return ctrl.Result{}, kverrors.Wrap(err, "failed to list any lokistack instances", "req")
	}

	if len(stacks.Items) > 0 {
		// Creates all LokiStack dashboard resources on OpenShift clusters when
		// the first LokiStack custom resource is created.
		if err := handlers.CreateDashboards(ctx, r.Log, req, r.Client, r.Scheme); err != nil {
			return ctrl.Result{}, kverrors.Wrap(err, "failed to create dashboard resources", "req", req)
		}

		return ctrl.Result{}, nil
	}

	// Removes all LokiStack dashboard resources on OpenShift clusters when
	// the last LokiStack custom resource is deleted.
	if err := handlers.DeleteDashboards(ctx, r.Client, req); err != nil {
		return ctrl.Result{}, kverrors.Wrap(err, "failed to delete dashboard resources")
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager to only call this controller on create/delete/generic events.
func (r *DashboardsReconciler) SetupWithManager(mgr manager.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lokiv1.LokiStack{}, createOrDeletesPred).
		Complete(r)
}