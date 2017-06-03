package metadata

type apiServerTemplateParams struct {
	ImageTag            string
	ServicesNetworkCIDR string
}

const apiServerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-apiserver
spec:
  hostNetwork: true
  containers:
  - name: kube-apiserver
    image: gcr.io/google_containers/kube-apiserver:{{ .ImageTag }}
    command:
    - kube-apiserver
    - "--apiserver-count=1"
    - "--allow-privileged=true"
    - "--etcd-servers=http://127.0.0.1:2379"
    - "--bind-address=0.0.0.0"
    - "--anonymous-auth=false"
    - "--tls-cert-file=/opt/kubernetes/certs/tls.pem"
    - "--tls-private-key-file=/opt/kubernetes/certs/tls-key.pem"
    - "--tls-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--client-ca-file=/opt/kubernetes/certs/ca.pem"
    - "--kubelet-certificate-authority=/opt/kubernetes/certs/ca.pem"
    - "--kubelet-client-certificate=/opt/kubernetes/certs/tls.pem"
    - "--kubelet-client-key=/opt/kubernetes/certs/tls-key.pem"
    - "--service-cluster-ip-range={{ .ServicesNetworkCIDR }}"
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
`
