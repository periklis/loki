package rules_test

import (
	"testing"

	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/handlers/internal/rules"
	"github.com/stretchr/testify/require"
)

func TestIsValid(t *testing.T) {
	tt := []struct {
		desc    string
		spec    *lokiv1beta1.LokiRuleSpec
		err     error
		wantErr bool
	}{
		{
			desc: "valid spec",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Limit:    10,
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Alert: "first-alert",
								For:   lokiv1beta1.EvaluationDuration("10m"),
								Expr:  `sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)`,
								Annotations: map[string]string{
									"annot": "something",
								},
								Labels: map[string]string{
									"severity": "critical",
								},
							},
							{
								Alert: "second-alert",
								For:   lokiv1beta1.EvaluationDuration("10m"),
								Expr:  `sum(rate({app="foo", env="stage"} |= "error" [5m])) by (job)`,
								Annotations: map[string]string{
									"env": "something",
								},
								Labels: map[string]string{
									"severity": "warning",
								},
							},
						},
					},
					{
						Name:     "second",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Limit:    10,
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Record: "nginx:requests:rate1m",
								Expr:   `sum(rate({container="nginx"}[1m]))`,
							},
							{
								Record: "banana:requests:rate5m",
								Expr:   `sum(rate({container="banana"}[1m]))`,
							},
						},
					},
				},
			},
		},
		{
			desc: "not unique group names",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
					},
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
					},
				},
			},
			err:     rules.ErrGroupNamesNotUnique,
			wantErr: true,
		},
		{
			desc: "ambiguous rule type",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Alert:  "an-alert",
								Record: "a_record_name",
							},
						},
					},
				},
			},
			err:     rules.ErrAmbiguousRuleType,
			wantErr: true,
		},
		{
			desc: "parse eval interval err",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1mo"),
					},
				},
			},
			err:     rules.ErrParseEvaluationInterval,
			wantErr: true,
		},
		{
			desc: "parse for interval err",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Alert: "an-alert",
								For:   lokiv1beta1.EvaluationDuration("10years"),
							},
						},
					},
				},
			},
			err:     rules.ErrParseAlertForPeriod,
			wantErr: true,
		},
		{
			desc: "invalid record metric name",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Record: "invalid:metric:name",
							},
						},
					},
				},
			},
			err:     rules.ErrInvalidRecordMetricName,
			wantErr: true,
		},
		{
			desc: "parse LogQL expression err",
			spec: &lokiv1beta1.LokiRuleSpec{
				Groups: []*lokiv1beta1.LokiRuleGroup{
					{
						Name:     "first",
						Interval: lokiv1beta1.EvaluationDuration("1m"),
						Rules: []*lokiv1beta1.LokiRuleGroupSpec{
							{
								Expr: "this is not a valid expression",
							},
						},
					},
				},
			},
			err:     rules.ErrParseLogQLExpression,
			wantErr: true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			err := rules.IsValid(tc.spec)
			if tc.wantErr {
				require.ErrorAs(t, err, &tc.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
