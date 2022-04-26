package handlers

import (
	"context"

	"github.com/ViaQ/logerr/kverrors"
	"github.com/go-logr/logr"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/external/k8s"
	"github.com/grafana/loki/operator/internal/handlers/internal/rules"
	"github.com/grafana/loki/operator/internal/status"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ValidateLokiRule handles LokiRule validation status.
func ValidateLokiRule(
	ctx context.Context,
	log logr.Logger,
	req ctrl.Request,
	k k8s.Client,
) error {
	ll := log.WithValues("lokirule", req.NamespacedName, "event", "validate")

	var r lokiv1beta1.LokiRule
	if err := k.Get(ctx, req.NamespacedName, &r); err != nil {
		if apierrors.IsNotFound(err) {
			// maybe the user deleted it before we could react? Either way this isn't an issue
			ll.Error(err, "could not find the requested loki rule", "name", req.NamespacedName)
			return nil
		}
		return kverrors.Wrap(err, "failed to lookup lokirule", "name", req.NamespacedName)
	}

	err := rules.IsValid(&r.Spec)
	if err != nil {
		switch kverrors.Root(err) {
		case rules.ErrAmbiguousRuleType:
			return &status.ValidationError{
				Message: "Ambiguous configuration mix of alerting and recording rule fields",
				Reason:  lokiv1beta1.ReasonAmbiguousRuleConfig,
			}
		case rules.ErrParseAlertForPeriod:
			return &status.ValidationError{
				Message: "Invalid alerting for period",
				Reason:  lokiv1beta1.ReasonInvalidAlertingRuleConfig,
			}
		case rules.ErrParseEvaluationInterval, rules.ErrInvalidRecordMetricName:
			return &status.ValidationError{
				Message: "Invalid recording rule configuration",
				Reason:  lokiv1beta1.ReasonInvalidRecordingRuleConfig,
			}
		case rules.ErrParseLogQLExpression:
			return &status.ValidationError{
				Message: "Invalid rule expression syntax",
				Reason:  lokiv1beta1.ReasonInvalidRuleExpression,
			}
		case rules.ErrGroupNamesNotUnique:
			return &status.ValidationError{
				Message: "Group names not unique",
				Reason:  lokiv1beta1.ReasonNotUniqueRuleGroupName,
			}
		}
	}

	return nil
}
