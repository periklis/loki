package v1beta1_test

import (
	"testing"

	"github.com/grafana/loki/operator/api/v1beta1"
	"github.com/stretchr/testify/require"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var tt = []struct {
	desc    string
	spec    v1beta1.LokiRuleSpec
	err     *apierrors.StatusError
	wantErr bool
}{
	{
		desc: "valid spec",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
					Limit:    10,
					Rules: []*v1beta1.LokiRuleGroupSpec{
						{
							Alert: "first-alert",
							For:   v1beta1.EvaluationDuration("10m"),
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
							For:   v1beta1.EvaluationDuration("10m"),
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
					Interval: v1beta1.EvaluationDuration("1m"),
					Limit:    10,
					Rules: []*v1beta1.LokiRuleGroupSpec{
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
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
				},
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(1).Child("Name"),
					"first",
					v1beta1.ErrGroupNamesNotUnique.Error(),
				),
			},
		),
		wantErr: true,
	},
	{
		desc: "ambiguous rule type",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
					Rules: []*v1beta1.LokiRuleGroupSpec{
						{
							Alert:  "an-alert",
							Record: "a_record_name",
							For:    v1beta1.EvaluationDuration("1m"),
							Expr:   `sum(rate({label="value"}[1m]))`,
						},
					},
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Rules").Index(0).Child("Alert"),
					"an-alert",
					v1beta1.ErrAmbiguousRuleType.Error(),
				),
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Rules").Index(0).Child("Record"),
					"a_record_name",
					v1beta1.ErrAmbiguousRuleType.Error(),
				),
			},
		),
		wantErr: true,
	},
	{
		desc: "parse eval interval err",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1mo"),
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Interval"),
					"1mo",
					v1beta1.ErrParseEvaluationInterval.Error(),
				),
			},
		),
		wantErr: true,
	},
	{
		desc: "parse for interval err",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
					Rules: []*v1beta1.LokiRuleGroupSpec{
						{
							Alert: "an-alert",
							For:   v1beta1.EvaluationDuration("10years"),
							Expr:  `sum(rate({label="value"}[1m]))`,
						},
					},
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Rules").Index(0).Child("For"),
					"10years",
					v1beta1.ErrParseAlertForPeriod.Error(),
				),
			},
		),
		wantErr: true,
	},
	{
		desc: "invalid record metric name",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
					Rules: []*v1beta1.LokiRuleGroupSpec{
						{
							Record: "invalid&metric:name",
							Expr:   `sum(rate({label="value"}[1m]))`,
						},
					},
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Rules").Index(0).Child("Record"),
					"invalid&metric:name",
					v1beta1.ErrInvalidRecordMetricName.Error(),
				),
			},
		),
		wantErr: true,
	},
	{
		desc: "parse LogQL expression err",
		spec: v1beta1.LokiRuleSpec{
			Groups: []*v1beta1.LokiRuleGroup{
				{
					Name:     "first",
					Interval: v1beta1.EvaluationDuration("1m"),
					Rules: []*v1beta1.LokiRuleGroupSpec{
						{
							Expr: "this is not a valid expression",
						},
					},
				},
			},
		},
		err: apierrors.NewInvalid(
			schema.GroupKind{Group: "loki.grafana.com", Kind: "LokiRule"},
			"testing-rule",
			field.ErrorList{
				field.Invalid(
					field.NewPath("Spec").Child("Groups").Index(0).Child("Rules").Index(0).Child("Expr"),
					"this is not a valid expression",
					v1beta1.ErrParseLogQLExpression.Error(),
				),
			},
		),
		wantErr: true,
	},
}

func TestLokiRuleValidationWebhook_ValidateCreate(t *testing.T) {
	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			l := v1beta1.LokiRule{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testing-rule",
				},
				Spec: tc.spec,
			}

			err := l.ValidateCreate()
			if tc.wantErr {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLokiRuleValidationWebhook_ValidateUpdate(t *testing.T) {
	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			l := v1beta1.LokiRule{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testing-rule",
				},
				Spec: tc.spec,
			}

			err := l.ValidateUpdate(&v1beta1.LokiRule{})
			if tc.wantErr {
				require.Equal(t, tc.err, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestLokiRuleValidationWebhook_ValidateDelete_DoNothing(t *testing.T) {
	for _, tc := range tt {
		tc := tc
		t.Run(tc.desc, func(t *testing.T) {
			t.Parallel()
			l := v1beta1.LokiRule{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testing-rule",
				},
				Spec: tc.spec,
			}

			err := l.ValidateDelete()
			require.NoError(t, err)
		})
	}
}
