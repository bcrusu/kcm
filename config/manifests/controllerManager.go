package manifests

type controllerManagerTemplateParams struct {
	ImageTag        string
	ClusterName     string
	PodsNetworkCIDR string
}

const controllerManagerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-controller-manager
  namespace: kube-system
spec:
  hostNetwork: true
  affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: node-role.kubernetes.io/master
            operator: Exists
  containers:
  - name: kube-controller-manager
    image: gcr.io/google_containers/kube-controller-manager:{{ .ImageTag }}
    command:
    - kube-controller-manager
    - "--address=127.0.0.1"
    - "--kubeconfig=/opt/kubernetes/kubeconfig/kube-controller-manager"
    - "--cluster-name={{ .ClusterName }}"
    - "--root-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--cluster-signing-cert-file=/opt/kubernetes/certs/ca.pem"
    - "--cluster-signing-key-file=/opt/kubernetes/certs/ca-key.pem"
    - "--use-service-account-credentials=true"
    - "--allocate-node-cidrs=true"
    - "--cluster-cidr={{ .PodsNetworkCIDR }}"
    - "--leader-elect=true"
    - "--controllers=*,serviceaccount-token,bootstrapsigner,tokencleaner"
    - "--service-account-private-key-file=/opt/kubernetes/certs/tls-server-key.pem"
    volumeMounts:
    - name: srvkube
      mountPath: "/srv/kubernetes"
      readOnly: true
    - name: etcssl
      mountPath: "/etc/ssl"
      readOnly: true
    - name: opt-kubernetes
      mountPath: "/opt/kubernetes"
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
  - name: opt-kubernetes
    hostPath:
      path: /opt/kubernetes  
`
