package manifests

type controllerManagerTemplateParams struct {
	ImageTag            string
	ClusterName         string
	PodsNetworkCIDR     string
	ServicesNetworkCIDR string
}

const controllerManagerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-controller-manager
spec:
  hostNetwork: true
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/node
  containers:
  - name: kube-controller-manager
    image: gcr.io/google_containers/kube-controller-manager:{{ .ImageTag }}
    command:
    - kube-controller-manager
    - "--address=0.0.0.0"
    - "--kubeconfig=/opt/kubernetes/kubeconfig"
    - "--cluster-name={{ .ClusterName }}"
    - "--root-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--service-account-private-key-file=/opt/kubernetes/certs/tls-key.pem"
    - "--service-cluster-ip-range={{ .ServicesNetworkCIDR }}"
    - "--cluster-cidr={{ .PodsNetworkCIDR }}"
    volumeMounts:
    - name: srvkube
      mountPath: "/srv/kubernetes"
      readOnly: true
    - name: etcssl
      mountPath: "/etc/ssl"
      readOnly: true
    livenessProbe:
      httpGet:
        scheme: HTTP
        host: 127.0.0.1
        port: 10252
        path: "/healthz"
      initialDelaySeconds: 15
      timeoutSeconds: 15
  volumes:
  - name: srvkube
    hostPath:
      path: "/srv/kubernetes"
  - name: etcssl
    hostPath:
      path: "/etc/ssl"
`
