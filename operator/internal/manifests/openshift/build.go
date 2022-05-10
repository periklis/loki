package openshift

import "sigs.k8s.io/controller-runtime/pkg/client"

// Build returns a list of auxiliary openshift/k8s objects
// for lokistack gateway deployments on OpenShift.
func Build(opts Options) []client.Object {
	objs := []client.Object{
		BuildRoute(opts),
		BuildGatewayServiceAccount(opts),
		BuildGatewayClusterRole(opts),
		BuildGatewayClusterRoleBinding(opts),
		BuildServiceCAConfigMap(opts),
	}

	if opts.BuildOpts.EnableServiceMonitors {
		objs = append(
			objs,
			BuildMonitoringRole(opts),
			BuildMonitoringRoleBinding(opts),
		)
	}

	if opts.BuildOpts.EnableRulerAlertManager {
		objs = append(objs,
			BuildRulerServiceAccount(opts),
			BuildRulerClusterRole(opts),
			BuildRulerClusterRoleBinding(opts),
		)
	}

	return objs
}
