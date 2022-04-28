package v1beta1

import (
	"github.com/grafana/loki/pkg/logql/syntax"
	"github.com/prometheus/common/model"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var recordingrulelog = logf.Log.WithName("recordingrule-resource")

// SetupWebhookWithManager registers the RecordingRuleWebhook to the controller-runtime manager
// or returns an error.
func (r *RecordingRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-loki-grafana-com-v1beta1-recordingrule,mutating=false,failurePolicy=fail,sideEffects=None,groups=loki.grafana.com,resources=recordingrules,verbs=create;update,versions=v1beta1,name=vrecordingrule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &RecordingRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *RecordingRule) ValidateCreate() error {
	recordingrulelog.Info("validate create", "name", r.Name)

	errs := r.validate()
	if len(errs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "loki.grafana.com", Kind: "RecordingRule"},
		r.Name,
		errs,
	)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *RecordingRule) ValidateUpdate(old runtime.Object) error {
	recordingrulelog.Info("validate update", "name", r.Name)

	errs := r.validate()
	if len(errs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "loki.grafana.com", Kind: "RecordingRule"},
		r.Name,
		errs,
	)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *RecordingRule) ValidateDelete() error {
	recordingrulelog.Info("validate delete", "name", r.Name)
	// Do nothing
	return nil
}

func (r *RecordingRule) validate() field.ErrorList {
	var allErrs field.ErrorList

	found := make([]string, 0)

	for i, g := range r.Spec.Groups {
		// Check for group name uniqueness
		for _, n := range found {
			if n == g.Name {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(i).Child("Name"),
					g.Name,
					ErrGroupNamesNotUnique.Error(),
				))
			}
		}

		found = append(found, g.Name)

		// Check if rule evaluation period is a valid PromQL duration
		_, err := model.ParseDuration(string(g.Interval))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(
				field.NewPath("Spec").Child("Groups").Index(i).Child("Interval"),
				g.Interval,
				ErrParseEvaluationInterval.Error(),
			))
		}

		for j, r := range g.Rules {
			// Check if recording rule name is a valid PromQL Label Name
			if r.Record != "" {
				if !model.IsValidMetricName(model.LabelValue(r.Record)) {
					allErrs = append(allErrs, field.Invalid(
						field.NewPath("Spec").Child("Groups").Index(i).Child("Rules").Index(j).Child("Record"),
						r.Record,
						ErrInvalidRecordMetricName.Error(),
					))
				}
			}

			// Check if the LogQL parser can parse the rule expression
			_, err := syntax.ParseExpr(r.Expr)
			if err != nil {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(i).Child("Rules").Index(j).Child("Expr"),
					r.Expr,
					ErrParseLogQLExpression.Error(),
				))
			}
		}
	}

	return allErrs
}
