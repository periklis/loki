package openshift

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BuildRulerClusterRole returns a k8s Role object for the ruler component
// to grant access to:
// nolint:misspell
// - monitoring.coreos.com/non-existant
// - monitoring.coreos.com/alertmanagers
func BuildRulerClusterRole(opts Options) *rbacv1.ClusterRole {
	return &rbacv1.ClusterRole{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRole",
			APIVersion: rbacv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   rulerRbacName(opts),
			Labels: opts.BuildOpts.Labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{
					"monitoring.coreos.com", //nolint:misspell
				},
				ResourceNames: []string{
					"non-existant",
				},
				Resources: []string{
					"alertmanagers",
				},
				Verbs: []string{
					"patch",
				},
			},
		},
	}
}

// BuildRulerClusterRoleBinding returns a k8s RoleBinding object for the ruler component
// to grant access to:
// nolint:misspell
// - monitoring.coreos.com/non-existant
// - monitoring.coreos.com/alertmanagers
func BuildRulerClusterRoleBinding(opts Options) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ClusterRoleBinding",
			APIVersion: rbacv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   rulerRbacName(opts),
			Labels: opts.BuildOpts.Labels,
		},
		RoleRef: rbacv1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     rulerRbacName(opts),
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      rulerServiceAccountName(opts),
				Namespace: opts.BuildOpts.LokiStackNamespace,
			},
		},
	}
}
