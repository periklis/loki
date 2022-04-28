package rules

import (
	"context"

	"github.com/ViaQ/logerr/v2/kverrors"
	lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"
	"github.com/grafana/loki/operator/internal/external/k8s"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ListAlertingRules returns a k8s resource list of AlertingRules for the given spec or an error. Three cases apply:
// - Return only matching rules in the stack namespace if no namespace selector given.
// - Return only matching rules in the stack namespace and in namespaces matching the namespace selector.
// - Return no rules if rules selector does not apply at all.
func ListAlertingRules(ctx context.Context, k k8s.Client, stackNs string, rs *lokiv1beta1.RulesSpec) (lokiv1beta1.AlertingRuleList, error) {
	nsl, err := selectRulesNamespaces(ctx, k, stackNs, rs)
	if err != nil {
		return lokiv1beta1.AlertingRuleList{}, err
	}

	rules, err := selectAlertingRules(ctx, k, rs)
	if err != nil {
		return lokiv1beta1.AlertingRuleList{}, err
	}

	var lrl lokiv1beta1.AlertingRuleList
	for _, rule := range rules.Items {
		for _, ns := range nsl.Items {
			if rule.Namespace == ns.Name {
				lrl.Items = append(lrl.Items, rule)
				break
			}
		}
	}

	return lrl, nil
}

// ListRecordingRules returns a k8s resource list of AlertingRules for the given spec or an error. Three cases apply:
// - Return only matching rules in the stack namespace if no namespace selector given.
// - Return only matching rules in the stack namespace and in namespaces matching the namespace selector.
// - Return no rules if rules selector does not apply at all.
func ListRecordingRules(ctx context.Context, k k8s.Client, stackNs string, rs *lokiv1beta1.RulesSpec) (lokiv1beta1.RecordingRuleList, error) {
	nsl, err := selectRulesNamespaces(ctx, k, stackNs, rs)
	if err != nil {
		return lokiv1beta1.RecordingRuleList{}, err
	}

	rules, err := selectRecordingRules(ctx, k, rs)
	if err != nil {
		return lokiv1beta1.RecordingRuleList{}, err
	}

	var lrl lokiv1beta1.RecordingRuleList
	for _, rule := range rules.Items {
		for _, ns := range nsl.Items {
			if rule.Namespace == ns.Name {
				lrl.Items = append(lrl.Items, rule)
				break
			}
		}
	}

	return lrl, nil
}

func selectRulesNamespaces(ctx context.Context, k k8s.Client, stackNs string, rs *lokiv1beta1.RulesSpec) (corev1.NamespaceList, error) {
	var stackNamespace corev1.Namespace
	key := client.ObjectKey{Name: stackNs}

	err := k.Get(ctx, key, &stackNamespace)
	if err != nil {
		return corev1.NamespaceList{}, kverrors.Wrap(err, "failed to get LokiStack namespace", "namespace", stackNs)
	}

	nsList := corev1.NamespaceList{Items: []corev1.Namespace{stackNamespace}}

	nsSelector, err := metav1.LabelSelectorAsSelector(rs.NamespaceSelector)
	if err != nil {
		return nsList, kverrors.Wrap(err, "failed to create LokiRule namespace selector", "namespaceSelector", rs.NamespaceSelector)
	}

	var nsl v1.NamespaceList
	err = k.List(ctx, &nsl, &client.MatchingLabelsSelector{Selector: nsSelector})
	if err != nil {
		return nsList, kverrors.Wrap(err, "failed to list namespaces for selector", "namespaceSelector", rs.NamespaceSelector)
	}

	for _, ns := range nsl.Items {
		if ns.Name == stackNs {
			continue
		}

		nsList.Items = append(nsList.Items, ns)
	}

	return nsList, nil
}

func selectAlertingRules(ctx context.Context, k k8s.Client, rs *lokiv1beta1.RulesSpec) (lokiv1beta1.AlertingRuleList, error) {
	rulesSelector, err := metav1.LabelSelectorAsSelector(rs.Selector)
	if err != nil {
		return lokiv1beta1.AlertingRuleList{}, kverrors.Wrap(err, "failed to create AlertingRules selector", "selector", rs.Selector)
	}

	var rl lokiv1beta1.AlertingRuleList
	err = k.List(ctx, &rl, &client.MatchingLabelsSelector{Selector: rulesSelector})
	if err != nil {
		return lokiv1beta1.AlertingRuleList{}, kverrors.Wrap(err, "failed to list AlertingRules for selector", "selector", rs.Selector)
	}

	return rl, nil
}

func selectRecordingRules(ctx context.Context, k k8s.Client, rs *lokiv1beta1.RulesSpec) (lokiv1beta1.RecordingRuleList, error) {
	rulesSelector, err := metav1.LabelSelectorAsSelector(rs.Selector)
	if err != nil {
		return lokiv1beta1.RecordingRuleList{}, kverrors.Wrap(err, "failed to create RecordingRules selector", "selector", rs.Selector)
	}

	var rl lokiv1beta1.RecordingRuleList
	err = k.List(ctx, &rl, &client.MatchingLabelsSelector{Selector: rulesSelector})
	if err != nil {
		return lokiv1beta1.RecordingRuleList{}, kverrors.Wrap(err, "failed to list RecordingRules for selector", "selector", rs.Selector)
	}

	return rl, nil
}
