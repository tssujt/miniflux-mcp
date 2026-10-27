# miniflux-mcp chart values

Configuration for the miniflux-mcp Helm chart.

### Type: `object`

| Property | Type | Required | Possible values | Default | Description |
| -------- | ---- | -------- | --------------- | ------- | ----------- |
| replicaCount | `integer` |  | integer | `1` | Number of application replicas. |
| image | `object` |  | object |  | Container image configuration. |
| image.repository | `string` |  | string | `"jwonder/miniflux-mcp"` | Container image repository. |
| image.tag | `string` |  | string | `"latest"` | Container image tag. |
| image.pullPolicy | `string` |  | `Always` `IfNotPresent` `Never` | `"IfNotPresent"` | Container image pull policy. |
| imagePullSecrets | `array` |  | object | `[]` | References to Secrets used to pull the container image. |
| nameOverride | `string` |  | string | `""` | Override the chart name. |
| fullnameOverride | `string` |  | string | `""` | Override the fully qualified application name. |
| service | `object` |  | object |  | Kubernetes Service configuration. |
| service.type | `string` |  | string | `"ClusterIP"` | Kubernetes Service type. |
| service.port | `integer` |  | integer | `8080` | Service port. |
| miniflux | `object` |  | object |  | Miniflux connection and authentication settings. |
| miniflux.url | `string` |  | string | `"http://miniflux:8080"` | URL of the Miniflux server. |
| miniflux.apiKey | `string` |  | string | `""` | Miniflux API key. Prefer existingSecret for production use. |
| miniflux.username | `string` |  | string | `""` | Miniflux username used when no API key is configured. |
| miniflux.password | `string` |  | string | `""` | Miniflux password used when no API key is configured. |
| mcp | `object` |  | object |  | MCP HTTP server configuration. |
| mcp.port | `integer` |  | integer | `8080` | Container port used by the MCP HTTP server. |
| mcp.path | `string` |  | string | `"/mcp"` | HTTP path of the MCP endpoint. |
| mcp.authToken | `string` |  | string | `""` | Bearer token for the MCP endpoint. Prefer existingSecret for production use. |
| existingSecret | `string` |  | string | `""` | Name of an existing Secret containing MCP and Miniflux credentials. |
| secretKeys | `object` |  | object |  | Key names used in the generated or existing Secret. |
| secretKeys.mcpAuthToken | `string` |  | string | `"MCP_AUTH_TOKEN"` | Secret key containing the MCP bearer token. |
| secretKeys.minifluxApiKey | `string` |  | string | `"MINIFLUX_API_KEY"` | Secret key containing the Miniflux API key. |
| secretKeys.minifluxUsername | `string` |  | string | `"MINIFLUX_USERNAME"` | Secret key containing the Miniflux username. |
| secretKeys.minifluxPassword | `string` |  | string | `"MINIFLUX_PASSWORD"` | Secret key containing the Miniflux password. |
| ingress | `object` |  | object |  | Kubernetes Ingress configuration. |
| ingress.enabled | `boolean` |  | boolean | `true` | Create an Ingress resource. |
| ingress.className | `string` |  | string | `"nginx"` | Ingress class name. |
| ingress.annotations | `object` |  | object | `{}` | Annotations added to the Ingress. |
| ingress.hosts | `array` |  | object | `[{"host": "mcp.example.com", "paths": [{"path": "/mcp", "pathType": "Prefix"}]}]` | Ingress hosts and paths. |
| ingress.hosts[].host | `string` |  | string |  | Hostname matched by the Ingress rule. |
| ingress.hosts[].paths | `array` |  | object |  | Paths exposed for the host. |
| ingress.hosts[].paths[].path | `string` |  | string |  | URL path matched by the Ingress rule. |
| ingress.hosts[].paths[].pathType | `string` |  | `Exact` `Prefix` `ImplementationSpecific` |  | Kubernetes Ingress path matching mode. |
| ingress.tls | `array` |  | object | `[]` | Ingress TLS configuration. |
| ingress.tls[].secretName | `string` |  | string |  | Secret containing the TLS certificate. |
| ingress.tls[].hosts | `array` |  | string |  | Hosts covered by the TLS certificate. |
| resources | `object` |  | object | `{"requests": {"cpu": "10m", "memory": "32Mi"}, "limits": {"cpu": "250m", "memory": "128Mi"}}` | Container resource requests and limits. |
| podSecurityContext | `object` |  | object | `{"runAsNonRoot": true, "runAsUser": 65532, "runAsGroup": 65532, "fsGroup": 65532}` | Pod-level security context. |
| securityContext | `object` |  | object | `{"allowPrivilegeEscalation": false, "readOnlyRootFilesystem": true, "capabilities": {"drop": ["ALL"]}}` | Container-level security context. |
| nodeSelector | `object` |  | object | `{}` | Node labels used for pod assignment. |
| tolerations | `array` |  | object | `[]` | Pod tolerations. |
| affinity | `object` |  | object | `{}` | Pod affinity rules. |
