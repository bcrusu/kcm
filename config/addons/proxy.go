package addons

type proxyTemplateParams struct {
	ImageTag        string
	PodsNetworkCIDR string
}

const proxyTemplate = `kind: ServiceAccount
apiVersion: v1
metadata:
  name: kube-proxy
  namespace: kube-system
---
kind: DaemonSet
apiVersion: extensions/v1beta1
metadata:
  name: kube-proxy
  namespace: kube-system
spec:
  template:
    metadata:
      labels:
        k8s-app: kube-proxy
    spec:
      hostNetwork: true
      nodeSelector:
        beta.kubernetes.io/arch: amd64
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      containers:
      - name: kube-proxy
        image: gcr.io/google_containers/kube-proxy:{{ .ImageTag }}
        securityContext:
          privileged: true
        command:
        - kube-proxy
        - "--bind-address=127.0.0.1"
        - "--kubeconfig=/opt/kubernetes/kubeconfig/kube-proxy"
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
      serviceAccountName: kube-proxy
      volumes:
      - name: opt-kubernetes
        hostPath:
          path: /opt/kubernetes
`
