package manifests_test

import (
	"strings"
	"testing"

	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/manifests"
	"github.com/stretchr/testify/require"
)

func TestNewRulerStatefulSet_HasTemplateConfigHashAnnotation(t *testing.T) {
	ss := manifests.NewRulerStatefulSet(manifests.Options{
		Name:               "abcd",
		Namespace:          "efgh",
		ConfigSHA1:         "deadbeef",
		AlertringRulesSHA1: "deadbeef",
		RecordingRulesSHA1: "deadbeef",
		Stack: lokiv1beta1.LokiStackSpec{
			StorageClassName: "standard",
			Template: &lokiv1beta1.LokiTemplateSpec{
				Ruler: &lokiv1beta1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	expected := "loki.grafana.com/config-hash"
	alertsExpected := "loki.grafana.com/alerting-rules-hash"
	recordingsExpected := "loki.grafana.com/recording-rules-hash"
	annotations := ss.Spec.Template.Annotations
	require.Contains(t, annotations, expected)
	require.Equal(t, annotations[expected], "deadbeef")
	require.Equal(t, annotations[alertsExpected], "deadbeef")
	require.Equal(t, annotations[recordingsExpected], "deadbeef")
}

func TestNewRulerStatefulSet_SelectorMatchesLabels(t *testing.T) {
	// You must set the .spec.selector field of a StatefulSet to match the labels of
	// its .spec.template.metadata.labels. Prior to Kubernetes 1.8, the
	// .spec.selector field was defaulted when omitted. In 1.8 and later versions,
	// failing to specify a matching Pod Selector will result in a validation error
	// during StatefulSet creation.
	// See https://kubernetes.io/docs/concepts/workloads/controllers/statefulset/#pod-selector
	sts := manifests.NewRulerStatefulSet(manifests.Options{
		Name:      "abcd",
		Namespace: "efgh",
		Stack: lokiv1beta1.LokiStackSpec{
			StorageClassName: "standard",
			Template: &lokiv1beta1.LokiTemplateSpec{
				Ruler: &lokiv1beta1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	l := sts.Spec.Template.GetObjectMeta().GetLabels()
	for key, value := range sts.Spec.Selector.MatchLabels {
		require.Contains(t, l, key)
		require.Equal(t, l[key], value)
	}
}

func TestNewRulerStatefulSet_MountsRulesInIndependentConfigMapVolumes(t *testing.T) {
	sts := manifests.NewRulerStatefulSet(manifests.Options{
		Name:      "abcd",
		Namespace: "efgh",
		Stack: lokiv1beta1.LokiStackSpec{
			StorageClassName: "standard",
			Template: &lokiv1beta1.LokiTemplateSpec{
				Ruler: &lokiv1beta1.LokiComponentSpec{
					Replicas: 1,
				},
			},
		},
	})

	vs := sts.Spec.Template.Spec.Volumes
	vms := sts.Spec.Template.Spec.Containers[0].VolumeMounts

	var (
		volumeNames      []string
		volumeMountNames []string
		volumeMountPaths = make(map[string]string)
	)
	for _, v := range vs {
		volumeNames = append(volumeNames, v.Name)
	}
	for _, v := range vms {
		volumeMountNames = append(volumeMountNames, v.Name)
		volumeMountPaths[v.Name] = v.MountPath
	}

	require.Contains(t, volumeNames, "alertingrules")
	require.Contains(t, volumeNames, "recordingrules")
	require.Contains(t, volumeMountNames, "alertingrules")
	require.Contains(t, volumeMountNames, "recordingrules")
	require.True(t, strings.HasPrefix(volumeMountPaths["alertingrules"], "/tmp/rules"))
	require.True(t, strings.HasSuffix(volumeMountPaths["alertingrules"], "/alerting-rules"))
	require.True(t, strings.HasPrefix(volumeMountPaths["recordingrules"], "/tmp/rules"))
	require.True(t, strings.HasSuffix(volumeMountPaths["recordingrules"], "/recording-rules"))
}
