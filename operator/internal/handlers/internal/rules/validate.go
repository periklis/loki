package rules

import (
	"errors"

	"github.com/ViaQ/logerr/kverrors"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/pkg/logql/syntax"
	"github.com/prometheus/common/model"
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

// IsValid checks if the given Loki rule spec is valid and if not returns and error if:
// - The group names are not unique.
// - The evaluation interval is not a valid PromQL duration.
// - Any rule is of ambiguous type, i.e. Alert and Record fields set.
// - Any rule's alert for period is not a valid PromQL duration.
// - Any rule's record field valud is not a valid PromQL metric name.
// - Any rule's expression is not a valid LogQL expression.
func IsValid(s *lokiv1beta1.LokiRuleSpec) error {
	found := make([]string, 0)

	for _, g := range s.Groups {
		// Check for group name uniqueness
		for _, n := range found {
			if n == g.Name {
				return kverrors.Wrap(ErrGroupNamesNotUnique, "group names are not unique across all defined rules", "name", g.Name)
			}
		}

		found = append(found, g.Name)

		// Check if rule evaluation period is a valid PromQL duration
		_, err := model.ParseDuration(string(g.Interval))
		if err != nil {
			return kverrors.Wrap(ErrParseEvaluationInterval, "invalid evaluation interval format", "interval", g.Interval)
		}

		for _, r := range g.Rules {
			// Check if we have a mix of alerting and recording rule syntax
			if r.Alert != "" && r.Record != "" {
				return kverrors.Wrap(ErrAmbiguousRuleType, "invalid rule syntax, cannot define Alert and Record at the same time")
			}

			// Check if alert for period is a valid PromQL duration
			if r.Alert != "" {
				_, err := model.ParseDuration(string(r.For))
				if err != nil {
					return kverrors.Wrap(ErrParseAlertForPeriod, "invalid alert for period as duration", "alert", r.Alert, "for", r.For)
				}
			}

			// Check if recording rule name is a valid PromQL Label Name
			if r.Record != "" {
				if !model.IsValidMetricName(model.LabelValue(r.Record)) {
					return kverrors.Wrap(ErrInvalidRecordMetricName, "not a valid label name", "record", r.Record)
				}
			}

			// Check if the LogQL parser can parse the rule expression
			_, err := syntax.ParseExpr(r.Expr)
			if err != nil {
				return kverrors.Wrap(ErrParseLogQLExpression, "unable parse LogQL expression for rule")
			}
		}
	}

	return nil
}
