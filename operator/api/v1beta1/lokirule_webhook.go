package v1beta1

import (
	"errors"

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

var (
	// ErrAmbiguousRuleType is the error type when loki groups have not distinct alerting or recording rule config.
	ErrAmbiguousRuleType = errors.New("Ambiguous rule configuration")
	// ErrGroupNamesNotUnique is the error type when loki groups have not unique names.
	ErrGroupNamesNotUnique = errors.New("Group names are not unique")
	// ErrInvalidRecordMetricName when any loki recording rule has a invalid PromQL metric name.
	ErrInvalidRecordMetricName = errors.New("Failed to parse record metric name")
	// ErrParseAlertForPeriod when any loki alerting rule for period is not a valid PromQL duration.
	ErrParseAlertForPeriod = errors.New("Failed to parse alert firing period")
	// ErrParseEvaluationInterval when any loki group evaluation internal is not a valid PromQL duration.
	ErrParseEvaluationInterval = errors.New("Failed to parse evaluation")
	// ErrParseLogQLExpression when any loki rule expression is not a valid LogQL expression.
	ErrParseLogQLExpression = errors.New("Failed to parse LogQL expression")
)

// log is for logging in this package.
var lokirulelog = logf.Log.WithName("lokirule-resource")

// SetupWebhookWithManager registers the LokiRuleWebhook to the controller-runtime manager
// or returns an error.
func (r *LokiRule) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-loki-grafana-com-v1beta1-lokirule,mutating=false,failurePolicy=fail,sideEffects=None,groups=loki.grafana.com,resources=lokirules,verbs=create;update,versions=v1beta1,name=vlokirule.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &LokiRule{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LokiRule) ValidateCreate() error {
	lokirulelog.Info("validate create", "name", r.Name)

	errs := r.validateLokiRule()
	if len(errs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
		r.Name,
		errs,
	)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LokiRule) ValidateUpdate(old runtime.Object) error {
	lokirulelog.Info("validate update", "name", r.Name)

	errs := r.validateLokiRule()
	if len(errs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
		r.Name,
		errs,
	)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LokiRule) ValidateDelete() error {
	lokirulelog.Info("validate delete", "name", r.Name)
	// Do nothing
	return nil
}

func (r *LokiRule) validateLokiRule() field.ErrorList {
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
			// Check if we have a mix of alerting and recording rule syntax
			if r.Alert != "" && r.Record != "" {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(i).Child("Rules").Index(j).Child("Alert"),
					r.Alert,
					ErrAmbiguousRuleType.Error(),
				))

				allErrs = append(allErrs, field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(i).Child("Rules").Index(j).Child("Record"),
					r.Record,
					ErrAmbiguousRuleType.Error(),
				))
			}

			// Check if alert for period is a valid PromQL duration
			if r.Alert != "" {
				_, err := model.ParseDuration(string(r.For))
				if err != nil {
					allErrs = append(allErrs, field.Invalid(
						field.NewPath("Spec").Child("Groups").Index(i).Child("Rules").Index(j).Child("For"),
						r.For,
						ErrParseAlertForPeriod.Error(),
					))
				}
			}

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
