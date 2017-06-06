package manifests

type schedulerTemplateParams struct {
	ImageTag string
}

const schedulerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-scheduler
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
  - name: kube-scheduler
    image: gcr.io/google_containers/kube-scheduler:{{ .ImageTag }}
    command:
    - kube-scheduler
    - "--address=0.0.0.0"
    - "--kubeconfig=/opt/kubernetes/kubeconfig"
    volumeMounts:
    - name: opt-kubernetes
      mountPath: "/opt/kubernetes"
      readOnly: true
    livenessProbe:
      httpGet:
        scheme: HTTP
        host: 127.0.0.1
        port: 10251
        path: "/healthz"
      initialDelaySeconds: 15
      timeoutSeconds: 15
  volumes:
  - name: opt-kubernetes
    hostPath:
      path: /opt/kubernetes  
`
