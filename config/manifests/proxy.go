package manifests

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
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/kube-proxy:{{ .ImageTag }}
    command:
    - kube-proxy
    - "--bind-address=0.0.0.0"
    - "--kubeconfig=/opt/kubernetes/kubeconfig"
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
