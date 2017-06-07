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
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: node-role.kubernetes.io/node
            operator: Exists
  containers:
  - name: kube-proxy
    image: gcr.io/google_containers/kube-proxy:{{ .ImageTag }}
    command:
    - kube-proxy
    - "--bind-address=127.0.0.1"
    - "--kubeconfig=/opt/kubernetes/kubeconfig-kube-proxy"
    - "--cluster-cidr={{ .PodsNetworkCIDR }}"
    volumeMounts:
    - name: opt-kubernetes
      mountPath: "/opt/kubernetes"
      readOnly: true
    livenessProbe:
      httpGet:
        scheme: HTTP
        host: 127.0.0.1
        port: 10249
        path: "/healthz"
      initialDelaySeconds: 15
      timeoutSeconds: 15
  volumes:
  - name: opt-kubernetes
    hostPath:
      path: /opt/kubernetes
`
