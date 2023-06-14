package manifests

import (
	"fmt"
	"path"

	"github.com/grafana/loki/operator/internal/manifests/internal/config"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/pointer"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func BuildNATS(opts Options) ([]client.Object, error) {
	sts := NewNATSStatefulSet(opts)

	if opts.Gates.HTTPEncryption {
		if err := configureNATSHTTPServicePKI(sts, opts); err != nil {
			return nil, err
		}
	}

	if opts.Gates.GRPCEncryption {
		if err := configureNATSGRPCServicePKI(sts, opts); err != nil {
			return nil, err
		}
	}

	if opts.Gates.HTTPEncryption || opts.Gates.GRPCEncryption {
		caBundleName := signingCABundleName(opts.Name)
		if err := configureServiceCA(&sts.Spec.Template.Spec, caBundleName); err != nil {
			return nil, err
		}
	}

	if opts.Gates.RestrictedPodSecurityStandard {
		if err := configurePodSpecForRestrictedStandard(&sts.Spec.Template.Spec); err != nil {
			return nil, err
		}
	}

	return []client.Object{
		sts,
		NewNATSGRPCService(opts),
		NewNATSHTTPService(opts),
		NewNATSPodDisruptionBudget(opts),
	}, nil
}

func NewNATSStatefulSet(opts Options) *appsv1.StatefulSet {
	l := ComponentLabels(LabelNATSComponent, opts.Name)
	a := commonAnnotations(opts.ConfigSHA1, opts.CertRotationRequiredAt)
	podSpec := corev1.PodSpec{
		Affinity: configureAffinity(LabelNATSComponent, opts.Name, opts.Gates.DefaultNodeAffinity, opts.Stack.Template.Nats),
		Volumes: []corev1.Volume{
			{
				Name: configVolumeName,
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						DefaultMode: &defaultConfigMapMode,
						LocalObjectReference: corev1.LocalObjectReference{
							Name: lokiConfigMapName(opts.Name),
						},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Image: opts.Image,
				Name:  "loki-nats",
				Resources: corev1.ResourceRequirements{
					Limits:   opts.ResourceRequirements.Nats.Limits,
					Requests: opts.ResourceRequirements.Nats.Requests,
				},
				Args: []string{
					"-target=nats",
					fmt.Sprintf("-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiConfigFileName)),
					fmt.Sprintf("-runtime-config.file=%s", path.Join(config.LokiConfigMountDir, config.LokiRuntimeConfigFileName)),
					"-config.expand-env=true",
					fmt.Sprintf("-nats.data-path=%s", dataDirectory),
				},
				ReadinessProbe: lokiReadinessProbe(),
				LivenessProbe:  lokiLivenessProbe(),
				Ports: []corev1.ContainerPort{
					{
						Name:          lokiHTTPPortName,
						ContainerPort: httpPort,
						Protocol:      protocolTCP,
					},
					{
						Name:          lokiGRPCPortName,
						ContainerPort: grpcPort,
						Protocol:      protocolTCP,
					},
					{
						Name:          lokiNATSClusterPortName,
						ContainerPort: natsClusterPort,
						Protocol:      protocolTCP,
					},
					{
						Name:          lokiNATSClientPortName,
						ContainerPort: natsClientPort,
						Protocol:      protocolTCP,
					},
				},
				VolumeMounts: []corev1.VolumeMount{
					{
						Name:      configVolumeName,
						ReadOnly:  false,
						MountPath: config.LokiConfigMountDir,
					},
					{
						Name:      storageVolumeName,
						ReadOnly:  false,
						MountPath: dataDirectory,
					},
				},
				TerminationMessagePath:   "/dev/termination-log",
				TerminationMessagePolicy: "File",
				ImagePullPolicy:          "IfNotPresent",
			},
		},
	}

	if opts.Stack.Template != nil && opts.Stack.Template.Nats != nil {
		podSpec.Tolerations = opts.Stack.Template.Nats.Tolerations
		podSpec.NodeSelector = opts.Stack.Template.Nats.NodeSelector
	}

	if opts.Stack.Replication != nil {
		podSpec.TopologySpreadConstraints = topologySpreadConstraints(*opts.Stack.Replication, LabelNATSComponent, opts.Name)
	}

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: appsv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   NATSName(opts.Name),
			Labels: l,
		},
		Spec: appsv1.StatefulSetSpec{
			PodManagementPolicy:  appsv1.OrderedReadyPodManagement,
			RevisionHistoryLimit: pointer.Int32(10),
			Replicas:             pointer.Int32(opts.Stack.Template.Nats.Replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: labels.Merge(l, GossipLabels()),
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:        fmt.Sprintf("loki-nats-%s", opts.Name),
					Labels:      labels.Merge(l, GossipLabels()),
					Annotations: a,
				},
				Spec: podSpec,
			},
			VolumeClaimTemplates: []corev1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Labels: l,
						Name:   storageVolumeName,
					},
					Spec: corev1.PersistentVolumeClaimSpec{
						AccessModes: []corev1.PersistentVolumeAccessMode{
							// TODO: should we verify that this is possible with the given storage class first?
							corev1.ReadWriteOnce,
						},
						Resources: corev1.ResourceRequirements{
							Requests: map[corev1.ResourceName]resource.Quantity{
								corev1.ResourceStorage: opts.ResourceRequirements.Nats.PVCSize,
							},
						},
						StorageClassName: pointer.String(opts.Stack.StorageClassName),
						VolumeMode:       &volumeFileSystemMode,
					},
				},
			},
		},
	}
}

// NewNATSGRPCService creates a k8s service for the NATS GRPC endpoint
func NewNATSGRPCService(opts Options) *corev1.Service {
	serviceName := serviceNameNATSGRPC(opts.Name)
	labels := ComponentLabels(LabelNATSComponent, opts.Name)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			ClusterIP: "None",
			Ports: []corev1.ServicePort{
				{
					Name:       lokiGRPCPortName,
					Port:       grpcPort,
					Protocol:   protocolTCP,
					TargetPort: intstr.IntOrString{IntVal: grpcPort},
				},
			},
			Selector: labels,
		},
	}
}

// NewNATSHTTPService creates a k8s service for the NATS HTTP endpoint
func NewNATSHTTPService(opts Options) *corev1.Service {
	serviceName := serviceNameNATSHTTP(opts.Name)
	labels := ComponentLabels(LabelNATSComponent, opts.Name)

	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:   serviceName,
			Labels: labels,
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       lokiHTTPPortName,
					Port:       httpPort,
					Protocol:   protocolTCP,
					TargetPort: intstr.IntOrString{IntVal: httpPort},
				},
			},
			Selector: labels,
		},
	}
}

func configureNATSHTTPServicePKI(statefulSet *appsv1.StatefulSet, opts Options) error {
	serviceName := serviceNameNATSHTTP(opts.Name)
	return configureHTTPServicePKI(&statefulSet.Spec.Template.Spec, serviceName)
}

func configureNATSGRPCServicePKI(sts *appsv1.StatefulSet, opts Options) error {
	serviceName := serviceNameNATSGRPC(opts.Name)
	return configureGRPCServicePKI(&sts.Spec.Template.Spec, serviceName)
}

// NewNATSPodDisruptionBudget returns a PodDisruptionBudget for the LokiStack NATS pods.
func NewNATSPodDisruptionBudget(opts Options) *policyv1.PodDisruptionBudget {
	l := ComponentLabels(LabelNATSComponent, opts.Name)

	ma := intstr.FromInt(1)

	return &policyv1.PodDisruptionBudget{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PodDisruptionBudget",
			APIVersion: policyv1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Labels:    l,
			Name:      NATSName(opts.Name),
			Namespace: opts.Namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: l,
			},
			MinAvailable: &ma,
		},
	}
}
