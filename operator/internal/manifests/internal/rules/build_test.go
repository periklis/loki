package rules_test

import (
	"testing"

	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/manifests/internal/rules"
	"github.com/stretchr/testify/require"
)

func TestBuild_RulesConfig_RenderAlertingRules(t *testing.T) {
	expCfg := `
groups:
  - name: an-alert
    limit: 2
    rules:
      - expr: |
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05
        alert: HighPercentageErrors
        for: 10m
        annotations:
          - playbook: http://link/to/playbook
          - summary: High Percentage Latency
        labels:
          - environment: production
          - severity: page
      - expr: |
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05
        alert: LowPercentageErrors
        for: 10m
        annotations:
          - playbook: http://link/to/playbook
          - summary: Low Percentage Latency
        labels:
          - environment: production
          - severity: low
`

	opts := rules.Options{
		AlertingGroups: []*lokiv1beta1.AlertingRuleGroup{
			{
				Name:  "an-alert",
				Limit: 2,
				Rules: []*lokiv1beta1.AlertingRuleGroupSpec{
					{
						Alert: "HighPercentageErrors",
						Expr: `
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05`,
						For: lokiv1beta1.PrometheusDuration("10m"),
						Labels: map[string]string{
							"severity":    "page",
							"environment": "production",
						},
						Annotations: map[string]string{
							"summary":  "High Percentage Latency",
							"playbook": "http://link/to/playbook",
						},
					},
					{
						Alert: "LowPercentageErrors",
						Expr: `
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05`,
						For: lokiv1beta1.PrometheusDuration("10m"),
						Labels: map[string]string{
							"severity":    "low",
							"environment": "production",
						},
						Annotations: map[string]string{
							"summary":  "Low Percentage Latency",
							"playbook": "http://link/to/playbook",
						},
					},
				},
			},
		},
	}

	cfg, err := rules.Build(opts)
	require.NoError(t, err)
	require.YAMLEq(t, expCfg, cfg)
}

func TestBuild_RulesConfig_RenderRecordingRules(t *testing.T) {
	expCfg := `
groups:
  - name: a-recording
    interval: 2d
    rules:
      - expr: |
          sum(
            rate({container="nginx"}[1m])
          )
        record: nginx:requests:rate1m
      - expr: |
          sum(
            rate({container="banana"}[5m])
          )
        record: banana:requests:rate5m
`

	opts := rules.Options{
		RecordingGroups: []*lokiv1beta1.RecordingRuleGroup{
			{
				Name:     "a-recording",
				Interval: lokiv1beta1.PrometheusDuration("2d"),
				Rules: []*lokiv1beta1.RecordingRuleGroupSpec{
					{
						Expr: `
          sum(
            rate({container="nginx"}[1m])
          )`,
						Record: "nginx:requests:rate1m",
					},
					{
						Expr: `
          sum(
            rate({container="banana"}[5m])
          )`,
						Record: "banana:requests:rate5m",
					},
				},
			},
		},
	}

	cfg, err := rules.Build(opts)
	require.NoError(t, err)
	require.YAMLEq(t, expCfg, cfg)
}

func TestBuild_RulesConfig_MultipleGroups(t *testing.T) {
	expCfg := `
groups:
  - name: an-alert
    limit: 2
    rules:
      - expr: |
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05
        alert: HighPercentageErrors
        for: 10m
        annotations:
          - playbook: http://link/to/playbook
          - summary: High Percentage Latency
        labels:
          - environment: production
          - severity: page
      - expr: |
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05
        alert: LowPercentageErrors
        for: 10m
        annotations:
          - playbook: http://link/to/playbook
          - summary: Low Percentage Latency
        labels:
          - environment: production
          - severity: low
  - name: a-recording
    interval: 2d
    rules:
      - expr: |
          sum(
            rate({container="nginx"}[1m])
          )
        record: nginx:requests:rate1m
      - expr: |
          sum(
            rate({container="banana"}[5m])
          )
        record: banana:requests:rate5m
`

	opts := rules.Options{
		AlertingGroups: []*lokiv1beta1.AlertingRuleGroup{
			{
				Name:  "an-alert",
				Limit: 2,
				Rules: []*lokiv1beta1.AlertingRuleGroupSpec{
					{
						Alert: "HighPercentageErrors",
						Expr: `
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05`,
						For: lokiv1beta1.PrometheusDuration("10m"),
						Labels: map[string]string{
							"severity":    "page",
							"environment": "production",
						},
						Annotations: map[string]string{
							"summary":  "High Percentage Latency",
							"playbook": "http://link/to/playbook",
						},
					},
					{
						Alert: "LowPercentageErrors",
						Expr: `
          sum(rate({app="foo", env="production"} |= "error" [5m])) by (job)
            /
          sum(rate({app="foo", env="production"}[5m])) by (job)
            > 0.05`,
						For: lokiv1beta1.PrometheusDuration("10m"),
						Labels: map[string]string{
							"severity":    "low",
							"environment": "production",
						},
						Annotations: map[string]string{
							"summary":  "Low Percentage Latency",
							"playbook": "http://link/to/playbook",
						},
					},
				},
			},
		},
		RecordingGroups: []*lokiv1beta1.RecordingRuleGroup{
			{
				Name:     "a-recording",
				Interval: lokiv1beta1.PrometheusDuration("2d"),
				Rules: []*lokiv1beta1.RecordingRuleGroupSpec{
					{
						Expr: `
          sum(
            rate({container="nginx"}[1m])
          )`,
						Record: "nginx:requests:rate1m",
					},
					{
						Expr: `
          sum(
            rate({container="banana"}[5m])
          )`,
						Record: "banana:requests:rate5m",
					},
				},
			},
		},
	}

	cfg, err := rules.Build(opts)
	require.NoError(t, err)
	require.YAMLEq(t, expCfg, cfg)
}
