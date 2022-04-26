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
	require.Len(t, cm.Data, 2)
	require.Contains(t, cm.Data, "dev-rules.yaml")
	require.Contains(t, cm.Data, "prod-rules.yaml")
}

func testOptions() manifests.Options {
	return manifests.Options{
		Rules: []lokiv1beta1.LokiRule{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "rules",
					Namespace: "dev",
				},
				Spec: lokiv1beta1.LokiRuleSpec{
					Groups: []*lokiv1beta1.LokiRuleGroup{
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
					Name:      "rules",
					Namespace: "prod",
				},
				Spec: lokiv1beta1.LokiRuleSpec{
					Groups: []*lokiv1beta1.LokiRuleGroup{
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
