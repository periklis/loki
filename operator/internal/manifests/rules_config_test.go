package manifests_test

import (
	"testing"

	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/manifests"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRulesConfigMap_ReturnsDataEntriesPerRule(t *testing.T) {
	cm, _, err := manifests.RulesConfigMap(testOptions())
	require.NoError(t, err)
	require.NotNil(t, cm)
	require.Len(t, cm.Data, 4)
	require.Contains(t, cm.Data, "dev-alerting-rules-alerts.yaml")
	require.Contains(t, cm.Data, "dev-recording-rules-recs.yaml")
	require.Contains(t, cm.Data, "prod-alerting-rules-alerts.yaml")
	require.Contains(t, cm.Data, "prod-recording-rules-recs.yaml")
}

func TestRulesConfigMap_ReturnsTenantMapPerRule(t *testing.T) {
	cm, tenants, err := manifests.RulesConfigMap(testOptions())
	require.NoError(t, err)
	require.NotNil(t, cm)
	require.Len(t, cm.Data, 4)
	require.Contains(t, tenants, "tenant-a")
	require.Contains(t, tenants, "tenant-b")
	require.Contains(t, tenants["tenant-a"], "dev-alerting-rules-alerts.yaml")
	require.Contains(t, tenants["tenant-a"], "prod-alerting-rules-alerts.yaml")
	require.Contains(t, tenants["tenant-b"], "dev-recording-rules-recs.yaml")
	require.Contains(t, tenants["tenant-b"], "prod-recording-rules-recs.yaml")
}

func testOptions() manifests.Options {
	return manifests.Options{
		AlertingRules: []lokiv1beta1.AlertingRule{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "alerting-rules",
					Namespace: "dev",
				},
				Spec: lokiv1beta1.AlertingRuleSpec{
					TenantID: "tenant-a",
					Groups: []*lokiv1beta1.AlertingRuleGroup{
						{
							Name: "rule-a",
						},
						{
							Name: "rule-b",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "alerting-rules",
					Namespace: "prod",
				},
				Spec: lokiv1beta1.AlertingRuleSpec{
					TenantID: "tenant-a",
					Groups: []*lokiv1beta1.AlertingRuleGroup{
						{
							Name: "rule-c",
						},
						{
							Name: "rule-d",
						},
					},
				},
			},
		},
		RecordingRules: []lokiv1beta1.RecordingRule{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "recording-rules",
					Namespace: "dev",
				},
				Spec: lokiv1beta1.RecordingRuleSpec{
					TenantID: "tenant-b",
					Groups: []*lokiv1beta1.RecordingRuleGroup{
						{
							Name: "rule-a",
						},
						{
							Name: "rule-b",
						},
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "recording-rules",
					Namespace: "prod",
				},
				Spec: lokiv1beta1.RecordingRuleSpec{
					TenantID: "tenant-b",
					Groups: []*lokiv1beta1.RecordingRuleGroup{
						{
							Name: "rule-c",
						},
						{
							Name: "rule-d",
						},
					},
				},
			},
		},
	}
}
