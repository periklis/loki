package openshift

import (
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
)

func TestBuildDashboards_ReturnsDashboardConfigMaps(t *testing.T) {
	opts := Options{
		BuildOpts: BuildOptions{
			LokiStackName:      "test",
			LokiStackNamespace: "test-ns",
		},
	}

	objs, err := BuildDashboards(opts)
	require.NoError(t, err)

	for _, d := range objs {
		switch d.(type) {
		case *corev1.ConfigMap:
			require.Equal(t, d.GetNamespace(), managedConfigNamespace)
			require.Contains(t, d.GetLabels(), labelConsoleDashboard)
		}
	}
}

func TestBuildDashboards_ReturnsPrometheusRules(t *testing.T) {
	opts := Options{
		BuildOpts: BuildOptions{
			LokiStackName:      "test",
			LokiStackNamespace: "test-ns",
		},
	}

	objs, err := BuildDashboards(opts)
	require.NoError(t, err)

	rules := objs[len(objs)-1]
	require.Equal(t, rules.GetName(), dashboardPrometheusRulesName(opts))
	require.Equal(t, rules.GetNamespace(), opts.BuildOpts.LokiStackNamespace)
}