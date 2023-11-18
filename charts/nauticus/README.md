# nauticus

![Version: 0.2.0](https://img.shields.io/badge/Version-0.2.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: v0.2.0](https://img.shields.io/badge/AppVersion-v0.2.0-informational?style=flat-square)

**Homepage:** <https://github.com/edixos/nauticus>

## Prerequisites

- Helm v3

- [CertManager](https://cert-manager.io) when admission webhooks is enabled

## Maintainers

| Name | Email | Url |
| ---- | ------ | --- |
| Ismail KABOUBI | <ikaboubi@gmail.com> | <https://smileisak.com> |

## Description

A Helm chart to deploy the Nauticus Controller Simplifying Kubernetes cluster management with fully-managed Spaces

## Source Code

* <https://github.com/edixos/nauticus>

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| affinity | object | `{}` | Affinity to apply to the controller pod |
| autoscaling.enabled | bool | `false` | Enable/Disable autoscaling |
| autoscaling.maxReplicas | int | `100` | Max replicas |
| autoscaling.minReplicas | int | `1` | Min replicas |
| autoscaling.targetCPUUtilizationPercentage | int | `80` | CPU Utilization Percentage |
| customAnnotations | object | `{}` | customAnnotations to add to resources |
| fullnameOverride | string | `""` | Override the default fullname |
| image.pullPolicy | string | `"IfNotPresent"` | Image Pull Policy |
| image.repository | string | `"ghcr.io/edixos/nauticus"` | Image repository operator |
| image.tag | string | chart version | Tag of the image of operator |
| imagePullSecrets | list | `[]` | Reference to one or more secrets to be used when pulling images ref: https://kubernetes.io/docs/tasks/configure-pod-container/pull-image-private-registry/ |
| kubeRbacProxy.image.repository | string | `"gcr.io/kubebuilder/kube-rbac-proxy"` | Image repository for rbac proxy |
| kubeRbacProxy.image.tag | string | `"v0.13.1"` | Tag of the image of rbac proxy |
| nameOverride | string | `""` | Override the default name |
| nodeSelector | string | `nil` | NodeSelector to apply to the controller pod |
| podAnnotations | object | `{}` | Annotations to add to the controller pod. |
| podSecurityContext | object | `{}` | podSecurityContext holds pod-level security attributes and common container settings. |
| replicaCount | int | `1` |  |
| resources | object | `{"limits":{"cpu":"500m","memory":"128Mi"},"requests":{"cpu":"10m","memory":"64Mi"}}` | Adds resources limits and request to controller pod. |
| securityContext | object | `{}` | securityContext holds container-level security attributes and common container settings. |
| serviceAccount.annotations | object | `{}` | Annotations to add to the service account |
| serviceAccount.create | bool | `true` | Specifies whether a service account should be created |
| serviceAccount.name | string | `""` | The name of the service account to use. If not set and create is true, a name is generated using the fullname template |
| serviceMonitor.annotations | object | `{}` | Assign additional Annotations |
| serviceMonitor.enabled | bool | `false` | Enable ServiceMonitor |
| serviceMonitor.endpoint.interval | string | `"15s"` | Set the scrape interval for the endpoint of the serviceMonitor |
| serviceMonitor.endpoint.metricRelabelings | list | `[]` | Set metricRelabelings for the endpoint of the serviceMonitor |
| serviceMonitor.endpoint.relabelings | list | `[]` | Set relabelings for the endpoint of the serviceMonitor |
| serviceMonitor.endpoint.scrapeTimeout | string | `""` | Set the scrape timeout for the endpoint of the serviceMonitor |
| serviceMonitor.labels | object | `{}` | Assign additional labels according to Prometheus' serviceMonitorSelector matching labels |
| serviceMonitor.matchLabels | object | `{}` | Change matching labels |
| serviceMonitor.namespace | string | `""` | Install the ServiceMonitor into a different Namespace, as the monitoring stack one (default: the release one) |
| serviceMonitor.targetLabels | list | `[]` | Set targetLabels for the serviceMonitor |
| tolerations | list | `[]` | Toleration to apply to the controller pod |

## Installing the Chart

### With Helm

To install the chart with the release name `my-release`:

```bash
helm repo add edixos https://edixos.github.io/charts
helm install edixos/nauticus
```

### With ArgoCD

#### Cluster k8s with access to public registry

Add new application as:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: nauticus
spec:
  project: infra

  source:
    repoURL: "https://edixos.github.io/charts"
    targetRevision: "0.2.0"
    chart: nauticus
    path: ''

  destination:
    server: https://kubernetes.default.svc
    namespace: "nauticus-system"

```