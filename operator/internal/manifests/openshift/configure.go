package openshift

import (
	"fmt"
	"path"
	"regexp"
	"strings"

	"github.com/ViaQ/logerr/v2/kverrors"
	"github.com/imdario/mergo"

	monitoringv1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

const (
	// tenantApplication is the name of the tenant holding application logs.
	tenantApplication = "application"
	// tenantInfrastructure is the name of the tenant holding infrastructure logs.
	tenantInfrastructure = "infrastructure"
	// tenantAudit is the name of the tenant holding audit logs.
	tenantAudit = "audit"
	// alertManagerServerName is the certificate server name to verify when accessing
	// OpenShift Monitoring's AlertManager instance on via https://alertmanager-operated.openshift-monitoring.svc:9095
	alertManagerServerName = "alertmanager-main.openshift-monitoring.svc"
)

var (
	// defaultTenants represents the slice of all supported LokiStack on OpenShift.
	defaultTenants = []string{
		tenantApplication,
		tenantInfrastructure,
		tenantAudit,
	}

	logsEndpointRe = regexp.MustCompile(`.*logs..*.endpoint.*`)
)

// ConfigureGatewayDeployment merges an OpenPolicyAgent sidecar into the deployment spec.
// With this, the deployment will route authorization request to the OpenShift
// apiserver through the sidecar.
func ConfigureGatewayDeployment(
	d *appsv1.Deployment,
	stackName string,
	gwContainerName string,
	sercretVolumeName, tlsDir, certFile, keyFile string,
	caDir, caFile string,
	withTLS, withCertSigningService bool,
) error {
	var gwIndex int
	for i, c := range d.Spec.Template.Spec.Containers {
		if c.Name == gwContainerName {
			gwIndex = i
			break
		}
	}

	gwContainer := d.Spec.Template.Spec.Containers[gwIndex].DeepCopy()
	gwArgs := gwContainer.Args
	gwVolumes := d.Spec.Template.Spec.Volumes

	if withCertSigningService {
		for i, a := range gwArgs {
			if logsEndpointRe.MatchString(a) {
				gwContainer.Args[i] = strings.Replace(a, "http", "https", 1)
			}
		}

		gwArgs = append(gwArgs, fmt.Sprintf("--logs.tls.ca-file=%s/%s", caDir, caFile))

		caBundleVolumeName := serviceCABundleName(Options{
			BuildOpts: BuildOptions{
				LokiStackName: stackName,
			},
		})

		gwContainer.VolumeMounts = append(gwContainer.VolumeMounts, corev1.VolumeMount{
			Name:      caBundleVolumeName,
			ReadOnly:  true,
			MountPath: caDir,
		})

		gwVolumes = append(gwVolumes, corev1.Volume{
			Name: caBundleVolumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					DefaultMode: &defaultConfigMapMode,
					LocalObjectReference: corev1.LocalObjectReference{
						Name: caBundleVolumeName,
					},
				},
			},
		})
	}

	gwContainer.Args = gwArgs

	p := corev1.PodSpec{
		ServiceAccountName: d.GetName(),
		Containers: []corev1.Container{
			*gwContainer,
			newOPAOpenShiftContainer(sercretVolumeName, tlsDir, certFile, keyFile, withTLS),
		},
		Volumes: gwVolumes,
	}

	if err := mergo.Merge(&d.Spec.Template.Spec, p, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed to merge sidecar container spec ")
	}

	return nil
}

// ConfigureGatewayService merges the OpenPolicyAgent sidecar metrics port into
// the service spec. With this the metrics are exposed through the same service.
func ConfigureGatewayService(s *corev1.ServiceSpec) error {
	spec := corev1.ServiceSpec{
		Ports: []corev1.ServicePort{
			{
				Name: opaMetricsPortName,
				Port: GatewayOPAInternalPort,
			},
		},
	}

	if err := mergo.Merge(s, spec, mergo.WithAppendSlice); err != nil {
		return kverrors.Wrap(err, "failed to merge sidecar service ports")
	}

	return nil
}

// ConfigureGatewayServiceMonitor merges the OpenPolicyAgent sidecar endpoint into
// the service monitor. With this cluster-monitoring prometheus can scrape
// the sidecar metrics.
func ConfigureGatewayServiceMonitor(sm *monitoringv1.ServiceMonitor, withTLS bool) error {
	var opaEndpoint monitoringv1.Endpoint

	if withTLS {
		tlsConfig := sm.Spec.Endpoints[0].TLSConfig
		opaEndpoint = monitoringv1.Endpoint{
			Port:            opaMetricsPortName,
			Path:            "/metrics",
			Scheme:          "https",
			BearerTokenFile: bearerTokenFile,
			TLSConfig:       tlsConfig,
		}
	} else {
		opaEndpoint = monitoringv1.Endpoint{
			Port:   opaMetricsPortName,
			Path:   "/metrics",
			Scheme: "http",
		}
	}

	spec := monitoringv1.ServiceMonitorSpec{
		Endpoints: []monitoringv1.Endpoint{opaEndpoint},
	}

	if err := mergo.Merge(&sm.Spec, spec, mergo.WithAppendSlice); err != nil {
		return kverrors.Wrap(err, "failed to merge sidecar service monitor endpoints")
	}

	return nil
}

// ConfigureRulerStatefulSet merges the serviceaccount name into the ruler component pod spec.
func ConfigureRulerStatefulSet(
	sts *appsv1.StatefulSet,
	stackName string,
	containerName string,
	caDir, caFile string,
) error {
	var gwIndex int
	for i, c := range sts.Spec.Template.Spec.Containers {
		if c.Name == containerName {
			gwIndex = i
			break
		}
	}

	rulerContainer := sts.Spec.Template.Spec.Containers[gwIndex].DeepCopy()
	rulerVolumes := sts.Spec.Template.Spec.Volumes

	rulerContainer.Args = append(rulerContainer.Args,
		fmt.Sprintf("-ruler.alertmanager-client.tls-ca-path=%s", path.Join(caDir, caFile)),
		fmt.Sprintf("-ruler.alertmanager-client.tls-server-name=%s", alertManagerServerName),
		fmt.Sprintf("-ruler.alertmanager-client.credentials-file=%s", bearerTokenFile),
	)

	caBundleVolumeName := serviceCABundleName(Options{
		BuildOpts: BuildOptions{
			LokiStackName: stackName,
		},
	})

	rulerContainer.VolumeMounts = append(rulerContainer.VolumeMounts, corev1.VolumeMount{
		Name:      caBundleVolumeName,
		ReadOnly:  true,
		MountPath: caDir,
	})

	rulerVolumes = append(rulerVolumes, corev1.Volume{
		Name: caBundleVolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				DefaultMode: &defaultConfigMapMode,
				LocalObjectReference: corev1.LocalObjectReference{
					Name: caBundleVolumeName,
				},
			},
		},
	})

	p := corev1.PodSpec{
		ServiceAccountName: sts.GetName(),
		Containers: []corev1.Container{
			*rulerContainer,
		},
		Volumes: rulerVolumes,
	}

	if err := mergo.Merge(&sts.Spec.Template.Spec, p, mergo.WithOverride); err != nil {
		return kverrors.Wrap(err, "failed to merge ruler pod spec")
	}

	return nil
}
