package manifests

import (
	"fmt"

	"github.com/grafana/loki/operator/internal/manifests/internal/rules"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RulesConfigMap returns a ConfigMap resource that contains
// all loki alerting and recording rules as YAML data.
func RulesConfigMap(opts Options) (*corev1.ConfigMap, map[string][]string, error) {
	var (
		data    = make(map[string]string)
		tenants = make(map[string][]string)
	)

	for _, r := range opts.AlertingRules {
		rOpts := rules.Options{AlertingGroups: r.Spec.Groups}

		c, err := rules.Build(rOpts)
		if err != nil {
			return nil, nil, err
		}

		key := fmt.Sprintf("%s-%s-alerts.yaml", r.Namespace, r.Name)
		tenants[r.Spec.TenantID] = append(tenants[r.Spec.TenantID], key)
		data[key] = c
	}

	for _, r := range opts.RecordingRules {
		rOpts := rules.Options{RecordingGroups: r.Spec.Groups}

		c, err := rules.Build(rOpts)
		if err != nil {
			return nil, nil, err
		}

		key := fmt.Sprintf("%s-%s-recs.yaml", r.Namespace, r.Name)
		tenants[r.Spec.TenantID] = append(tenants[r.Spec.TenantID], key)
		data[key] = c
	}

	l := commonLabels(opts.Name)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      RulesConfigMapName(opts.Name),
			Namespace: opts.Namespace,
			Labels:    l,
		},
		Data: data,
	}, tenants, nil
}
