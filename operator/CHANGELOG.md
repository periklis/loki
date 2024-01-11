## Main

## Release 5.6.16

- [11624](https://github.com/grafana/loki/pull/11624) **xperimental**: React to changes in ConfigMap used for storage CA

## Release 5.6.15

- [11448](https://github.com/grafana/loki/pull/11448) **periklis**: Update Loki operand to v2.9.3
- [11357](https://github.com/grafana/loki/pull/11357) **periklis**: Fix storing authentication credentials in the Loki ConfigMap

## Release 5.6.14

- [11393](https://github.com/grafana/loki/pull/11393) **periklis**: Add infra annotations for OpenShift based deployments
- [11288](https://github.com/grafana/loki/pull/11288) **periklis**: Fix custom CA for object-store in ruler component

## Release 5.6.13

No changes.

## Release 5.6.12

- [10924](https://github.com/grafana/loki/pull/10924) **periklis**: Update Loki operand to v2.9.2
- [10874](https://github.com/grafana/loki/pull/10874) **periklis**: Bump deps to address CVE-2023-39325 and CVE-2023-44487
- [10562](https://github.com/grafana/loki/pull/10562) **periklis**: Add memberlist IPv6 support
- [10600](https://github.com/grafana/loki/pull/10600) **periklis**: Update Loki operand to v2.9.1

## Release 5.6.11

No changes.

## Release 5.6.10

No changes.

## Release 5.6.9

- [10019](https://github.com/grafana/loki/pull/10019) **periklis**: Update Loki operand to v2.8.3

## Release 5.6.8

- [9830](https://github.com/grafana/loki/pull/9830) **periklis**: Expose limits config setting cardinality_limit
- [9630](https://github.com/grafana/loki/pull/9630) **jpinsonneau**: Expose per_stream_rate_limit & burst

## Release 5.6.7

- [9623](https://github.com/grafana/loki/pull/9623) **periklis**: Fix timeout config constructor when only tenants limits
- [9511](https://github.com/grafana/loki/pull/9511) **xperimental**: Do not update status after setting degraded condition
- [9405](https://github.com/grafana/loki/pull/9405) **periklis**: Add support for configuring HTTP server timeouts
- [9346](https://github.com/grafana/loki/pull/9346) **periklis**: Enable Route by default on OpenShift clusters

## Release 5.6.6

- [9036](https://github.com/grafana/loki/pull/9036) **periklis**: Update Loki operand to v2.8.0

## Release 5.6.5

- [8978](https://github.com/grafana/loki/pull/8978) **aminesnow**: Add watch for the object storage secret
- [8958](https://github.com/grafana/loki/pull/8958) **periklis**: Align common instance addr with memberlist advertise addr

## Release 5.6.4

- [8672](https://github.com/grafana/loki/pull/8672) **periklis**: Add support for memberlist bind network configuration
- [8771](https://github.com/grafana/loki/pull/8771) **periklis**: Update Loki operand to v2.7.4
- [8707](https://github.com/grafana/loki/pull/8707) **aminesnow**: Fix gateway's nodeSelector and toleration
- [8577](https://github.com/grafana/loki/pull/8577) **Red-GV**: Store gateway tenant information in secret instead of configmap
- [8397](https://github.com/grafana/loki/pull/8397) **periklis**: Update Loki operand to v2.7.3
- [8336](https://github.com/grafana/loki/pull/8336) **periklis**: Update Loki operand to v2.7.2
- [8265](https://github.com/grafana/loki/pull/8265) **Red-GV**: Use gRPC compactor service instead of http for retention
- [8087](https://github.com/grafana/loki/pull/8087) **xperimental**: Fix status not updating when state of pods changes
- [8173](https://github.com/grafana/loki/pull/8173) **periklis**: Remove custom webhook cert mounts for OLM-based deployment (OpenShift)
- [8068](https://github.com/grafana/loki/pull/8068) **periklis**: Use lokistack-gateway replicas from size table
- [7910](https://github.com/grafana/loki/pull/7910) **periklis**: Update Loki operand to v2.7.1
- [7815](https://github.com/grafana/loki/pull/7815) **periklis**: Apply delete client changes for compat with release-2.7.x
- [7809](https://github.com/grafana/loki/pull/7809) **xperimental**: Fix histogram-based alerting rules
- [7808](https://github.com/grafana/loki/pull/7808) **xperimental**: Replace fifocache usage by embedded_cache
- [7753](https://github.com/grafana/loki/pull/7753) **periklis**: Check for mandatory CA configmap name in ObjectStorageTLS spec
- [7744](https://github.com/grafana/loki/pull/7744) **periklis**: Fix object storage TLS spec CAKey descriptor
- [7716](https://github.com/grafana/loki/pull/7716) **aminesnow**: Migrate API docs generation tool
- [7710](https://github.com/grafana/loki/pull/7710) **periklis**: Fix LokiStackController watches for cluster-scoped resources
- [7682](https://github.com/grafana/loki/pull/7682) **periklis**: Refactor cluster proxy to use configv1.Proxy on OpenShift
- [7711](https://github.com/grafana/loki/pull/7711) **Red-GV**: Remove default value from replicationFactor field
- [7617](https://github.com/grafana/loki/pull/7617) **Red-GV**: Modify ingestionRate for respective shirt size
- [7592](https://github.com/grafana/loki/pull/7592) **aminesnow**: Update API docs generation using gen-crd-api-reference-docs
- [7448](https://github.com/grafana/loki/pull/7448) **periklis**: Add TLS support for compactor delete client
- [7596](https://github.com/grafana/loki/pull/7596) **periklis**: Fix fresh-installs with built-in cert management enabled
- [7064](https://github.com/grafana/loki/pull/7064) **periklis**: Add support for built-in cert management
- [7471](https://github.com/grafana/loki/pull/7471) **aminesnow**: Expose and migrate query_timeout in limits config
- [7437](https://github.com/grafana/loki/pull/7437) **aminesnow**: Fix Custom TLS profile setting for LokiStack on OpenShift
- [7415](https://github.com/grafana/loki/pull/7415) **aminesnow**: Add alert relabel config
- [7418](https://github.com/grafana/loki/pull/7418) **Red-GV**: Update golang to v1.19 and k8s dependencies to v0.25.2
- [7322](https://github.com/grafana/loki/pull/7322) **Red-GV**: Configuring server and client HTTP and GRPC TLS options
- [7272](https://github.com/grafana/loki/pull/7272) **aminesnow**: Use cluster monitoring alertmanager by default on openshift clusters
- [7295](https://github.com/grafana/loki/pull/7295) **xperimental**: Add extended-validation for rules on OpenShift
- [6951](https://github.com/grafana/loki/pull/6951) **Red-GV**: Adding operational Lokistack alerts
- [7254](https://github.com/grafana/loki/pull/7254) **periklis**: Expose Loki Ruler API via the lokistack-gateway
- [7214](https://github.com/grafana/loki/pull/7214) **periklis**: Fix ruler GRPC tls client configuration
- [7201](https://github.com/grafana/loki/pull/7201) **xperimental**: Write configuration for per-tenant retention
- [7037](https://github.com/grafana/loki/pull/7037) **xperimental**: Skip enforcing matcher for certain tenants on OpenShift
- [7106](https://github.com/grafana/loki/pull/7106) **xperimental**: Manage global stream-based retention
- [7092](https://github.com/grafana/loki/pull/7092) **aminesnow**: Configure kube-rbac-proxy sidecar to use Intermediate TLS security profile in OCP
- [6870](https://github.com/grafana/loki/pull/6870) **aminesnow**: Configure gateway to honor the global tlsSecurityProfile on Openshift
- [6999](https://github.com/grafana/loki/pull/6999) **Red-GV**: Adding LokiStack Gateway alerts
- [7000](https://github.com/grafana/loki/pull/7000) **xperimental**: Configure default node affinity for all pods
- [6923](https://github.com/grafana/loki/pull/6923) **xperimental**: Reconcile owner reference for existing objects
- [6907](https://github.com/grafana/loki/pull/6907) **Red-GV**: Adding valid subscription annotation to operator metadata
- [6479](https://github.com/grafana/loki/pull/6749) **periklis**: Update Loki operand to v2.6.1
- [6748](https://github.com/grafana/loki/pull/6748) **periklis**: Update go4.org/unsafe/assume-no-moving-gc to latest
- [6741](https://github.com/grafana/loki/pull/6741) **aminesnow**: Golang version to 1.18 and k8s client to 1.24
- [6669](https://github.com/grafana/loki/pull/6669) **xperimental**: Set minimum TLS version to 1.2 to support FIPS
- [6663](https://github.com/grafana/loki/pull/6663) **aminesnow**: Generalize live tail fix to all clusters using TLS
- [6443](https://github.com/grafana/loki/pull/6443) **aminesnow**: Fix live tail of logs not working on OpenShift-based clusters
- [6646](https://github.com/grafana/loki/pull/6646) **periklis**: Update Loki operand to v2.6.0
- [6594](https://github.com/grafana/loki/pull/6594) **xperimental**: Disable client certificate authentication on gateway
- [6551](https://github.com/grafana/loki/pull/6561) **periklis**: Add operator docs for object storage
- [6549](https://github.com/grafana/loki/pull/6549) **periklis**: Refactor feature gates to use custom resource definition
- [6514](https://github.com/grafana/loki/pull/6514) **Red-GV** Update all pods and containers to be compliant with restricted Pod Security Standard
- [6531](https://github.com/grafana/loki/pull/6531) **periklis**: Use default interface_names for lokistack clusters (IPv6 Support)
- [6411](https://github.com/grafana/loki/pull/6478) **aminesnow**: Support TLS enabled lokistack-gateway for vanilla kubernetes deployments
- [6504](https://github.com/grafana/loki/pull/6504) **periklis**: Disable usage report on OpenShift
- [6474](https://github.com/grafana/loki/pull/6474) **periklis**: Bump loki.grafana.com/LokiStack from v1beta to v1
- [6411](https://github.com/grafana/loki/pull/6411) **Red-GV**: Extend schema validation in LokiStack webhook
- [6334](https://github.com/grafana/loki/pull/6433) **periklis**: Move operator cli flags to component config
- [6224](https://github.com/grafana/loki/pull/6224) **periklis**: Add support for GRPC over TLS for Loki components
- [5952](https://github.com/grafana/loki/pull/5952) **Red-GV**: Add api to change storage schema version
- [6363](https://github.com/grafana/loki/pull/6363) **periklis**: Allow optional installation of webhooks (Kind)
- [6362](https://github.com/grafana/loki/pull/6362) **periklis**: Allow reduced tenant OIDC authentication requirements
- [6288](https://github.com/grafana/loki/pull/6288) **aminesnow**: Expose only an HTTPS gateway when in openshift mode
- [6195](https://github.com/grafana/loki/pull/6195) **periklis**: Add ruler config support
- [6198](https://github.com/grafana/loki/pull/6198) **periklis**: Add support for custom S3 CA
- [6199](https://github.com/grafana/loki/pull/6199) **Red-GV**: Update GCP secret volume path
- [6125](https://github.com/grafana/loki/pull/6125) **sasagarw**: Add method to get authenticated from GCP
- [5986](https://github.com/grafana/loki/pull/5986) **periklis**: Add support for Loki Rules reconciliation
- [5987](https://github.com/grafana/loki/pull/5987) **Red-GV**: Update logerr to v2.0.0
- [5907](https://github.com/grafana/loki/pull/5907) **xperimental**: Do not include non-static labels in pod selectors
- [5893](https://github.com/grafana/loki/pull/5893) **periklis**: Align PVC storage size requests for all lokistack t-shirt sizes
- [5884](https://github.com/grafana/loki/pull/5884) **periklis**: Update Loki operand to v2.5.0
- [5748](https://github.com/grafana/loki/pull/5748) **Red-GV**: Update Prometheus go client to 12.1
- [5739](https://github.com/grafana/loki/pull/5739) **sasagarw**: Change UUIDs to tenant name in doc
- [5729](https://github.com/grafana/loki/pull/5729) **periklis**: Add missing label matcher for openshift logging tenant mode (OpenShift)
- [5691](https://github.com/grafana/loki/pull/5691) **sasagarw**: Fix immediate reset of degraded condition
- [5704](https://github.com/grafana/loki/pull/5704) **xperimental**: Update operator-sdk to 1.18.1
- [5693](https://github.com/grafana/loki/pull/5693) **periklis**: Replace frontend_worker parallelism with match_max_concurrent
- [5699](https://github.com/grafana/loki/pull/5699) **Red-GV**: Configure boltdb_shipper and schema to use Azure, GCS, and Swift storage
- [5701](https://github.com/grafana/loki/pull/5701) **sasagarw**: Make ReplicationFactor optional in LokiStack API
- [5695](https://github.com/grafana/loki/pull/5695) **xperimental**: Update Go to 1.17
- [5615](https://github.com/grafana/loki/pull/5615) **sasagarw**: Document how to connect to LokiStack gateway component
- [5655](https://github.com/grafana/loki/pull/5655) **xperimental**: Update Loki operand to 2.4.2
- [5579](https://github.com/grafana/loki/pull/5579) **Red-GV**: Add playbook for responding to operator alerts
- [5640](https://github.com/grafana/loki/pull/5640) **sasagarw**: Update CSV to point to candidate channel and use openshift-operators-redhat ns (OpenShift)
- [5551](https://github.com/grafana/loki/pull/5551) **sasagarw**: Document how to connect to distributor component
- [5624](https://github.com/grafana/loki/pull/5624) **periklis**: Use tenant name as id for mode openshift-logging (OpenShift)
- [5621](https://github.com/grafana/loki/pull/5621) **periklis**: Use recommended labels for LokiStack components
- [5607](https://github.com/grafana/loki/pull/5607) **periklis**: Use lokistack name as prefix for owned resources
- [5588](https://github.com/grafana/loki/pull/5588) **periklis**: Add RBAC for Prometheus service discovery to Loki component metrics (OpenShift)
- [5576](https://github.com/grafana/loki/pull/5576) **xperimental**: Change endpoints for generated liveness and readiness probes
- [5560](https://github.com/grafana/loki/pull/5560) **periklis**: Fix service monitor's server name for operator metrics
- [5345](https://github.com/grafana/loki/pull/5345) **ronensc**: Add flag to create Prometheus rules
- [5432](https://github.com/grafana/loki/pull/5432) **Red-GV**: Provide storage configuration for Azure, GCS, and Swift through common_config
- [4975](https://github.com/grafana/loki/pull/4975) **periklis**: Provide saner default for loki-operator managed chunk_target_size
