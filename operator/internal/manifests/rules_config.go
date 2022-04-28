package manifests

import (
	"crypto/sha1"
	"fmt"

	"github.com/grafana/loki/operator/internal/manifests/internal/rules"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// AlertingRulesConfigMap returns a ConfigMap resource that contains
// all loki alerting rules as YAML data.
func AlertingRulesConfigMap(opts Options) (*corev1.ConfigMap, string, error) {
	var (
		data = make(map[string]string)
		rStr []byte
	)

	for _, r := range opts.AlertingRules {
		rOpts := rules.Options{AlertingGroups: r.Spec.Groups}

		c, err := rules.Build(rOpts)
		if err != nil {
			return nil, "", err
		}

		key := fmt.Sprintf("%s-%s.yaml", r.Namespace, r.Name)
		data[key] = c
		rStr = append(rStr, []byte(c)...)
	}

	sha1C := dataSHA1(rStr)
	l := commonLabels(opts.Name)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      AlertringRulesConfigMapName(opts.Name),
			Namespace: opts.Namespace,
			Labels:    l,
		},
		Data: data,
	}, sha1C, nil
}

// RecordingRulesConfigMap returns a ConfigMap resource that contains
// all loki recording rules as YAML data.
func RecordingRulesConfigMap(opts Options) (*corev1.ConfigMap, string, error) {
	var (
		data = make(map[string]string)
		rStr []byte
	)

	for _, r := range opts.RecordingRules {
		rOpts := rules.Options{RecordingGroups: r.Spec.Groups}

		c, err := rules.Build(rOpts)
		if err != nil {
			return nil, "", err
		}

		key := fmt.Sprintf("%s-%s.yaml", r.Namespace, r.Name)
		data[key] = c
		rStr = append(rStr, []byte(c)...)
	}

	sha1C := dataSHA1(rStr)
	l := commonLabels(opts.Name)

	return &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      RecordingRulesConfigMapName(opts.Name),
			Namespace: opts.Namespace,
			Labels:    l,
		},
		Data: data,
	}, sha1C, nil
}

func dataSHA1(b []byte) string {
	s := sha1.New()
	_, err := s.Write(b)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", s.Sum(nil))
}
