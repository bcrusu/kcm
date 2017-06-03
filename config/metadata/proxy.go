package metadata

type proxyTemplateParams struct {
	ImageTag        string
	PodsNetworkCIDR string
}

const proxyTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-proxy
spec:
  hostNetwork: true
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/kube-proxy:{{ .ImageTag }}
    command:
    - kube-proxy
    - "--bind-address=0.0.0.0"
    - "--master=127.0.0.1:8080"
    - "--cluster-cidr={{ .PodsNetworkCIDR }}"
    livenessProbe:
      httpGet:
        scheme: HTTP
        host: 127.0.0.1
        port: 10249
        path: "/healthz"
      initialDelaySeconds: 15
      timeoutSeconds: 15
`
