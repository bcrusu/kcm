package manifests

type apiServerTemplateParams struct {
	ImageTag            string
	ServicesNetworkCIDR string
}

const apiServerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-apiserver
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
  - name: kube-apiserver
    image: gcr.io/google_containers/kube-apiserver:{{ .ImageTag }}
    command:
    - kube-apiserver
    - "--apiserver-count=1"
    - "--allow-privileged=true"
    - "--etcd-servers=http://127.0.0.1:2379"
    - "--bind-address=0.0.0.0"
    - "--secure-port=6443"
    - "--anonymous-auth=false"
    - "--tls-cert-file=/opt/kubernetes/certs/tls-server.pem"
    - "--tls-private-key-file=/opt/kubernetes/certs/tls-server-key.pem"
    - "--tls-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--client-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--kubelet-certificate-authority=/opt/kubernetes/certs/ca.pem"
    - "--kubelet-client-certificate=/opt/kubernetes/certs/tls-client.pem"
    - "--kubelet-client-key=/opt/kubernetes/certs/tls-client-key.pem"
    - "--service-cluster-ip-range={{ .ServicesNetworkCIDR }}"
    - "--storage-backend=etcd2"
    - "--storage-media-type=application/json"
    - "--service-account-key-file=/opt/kubernetes/certs/tls-server-key.pem"
    - "--admission-control=Initializers,NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,ResourceQuota"
    ports:
    - name: https
      hostPort: 443
      containerPort: 443
    - name: local
      hostPort: 8080
      containerPort: 8080
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
        port: 8080
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
