package manifests_test

import (
	"testing"

	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/manifests"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestLokiRulesConfigMap_ReturnsSHA1OfAllDataContents(t *testing.T) {
	_, sha1C, err := manifests.LokiRulesConfigMap(testOptions())
	require.NoError(t, err)
	require.NotEmpty(t, sha1C)
}

func TestLokiRulesConfigMap_ReturnsDataEntriesPerRule(t *testing.T) {
	cm, _, err := manifests.LokiRulesConfigMap(testOptions())
	require.NoError(t, err)
	require.NotNil(t, cm)
	require.Len(t, cm.Data, 4)
	require.Contains(t, cm.Data, "dev-alerting-rules.yaml")
	require.Contains(t, cm.Data, "prod-alerting-rules.yaml")
	require.Contains(t, cm.Data, "dev-recording-rules.yaml")
	require.Contains(t, cm.Data, "prod-recording-rules.yaml")
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
