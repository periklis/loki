package rules

import lokiv1beta1 "github.com/grafana/loki/operator/api/v1beta1"

// Options is used to render the loki-rules-config.yaml file template
type Options struct {
	Groups []*lokiv1beta1.LokiRuleGroup
}
