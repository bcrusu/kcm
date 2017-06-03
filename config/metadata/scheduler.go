package metadata

type schedulerTemplateParams struct {
	ImageTag string
}

const schedulerTemplate = `kind: Pod
apiVersion: v1
metadata:
  name: kube-scheduler
spec:
  hostNetwork: true
  containers:
  - name: kube-scheduler
    image: gcr.io/google_containers/kube-scheduler:{{ .ImageTag }}
    command:
    - kube-scheduler
    - "--address=0.0.0.0"
    - "--master=127.0.0.1:8080"
    livenessProbe:
      httpGet:
        scheme: HTTP
        host: 127.0.0.1
        port: 10251
        path: "/healthz"
      initialDelaySeconds: 15
      timeoutSeconds: 15
`
