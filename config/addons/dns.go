package addons

type dnsTemplateParams struct {
	ClusterDomain string
	ServiceIP     string
	ImageTag      string
}

const dnsTemplate = `kind: ServiceAccount
apiVersion: v1
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-app: kube-dns
    k8s-addon: kube-dns.addons.k8s.io
---
kind: Service
apiVersion: v1
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-addon: kube-dns.addons.k8s.io
    k8s-app: kube-dns
    kubernetes.io/cluster-service: "true"
    kubernetes.io/name: "KubeDNS"
spec:
  selector:
    k8s-app: kube-dns
  clusterIP: {{ .ServiceIP }}
  ports:
  - name: dns
    port: 53
    protocol: UDP
  - name: dns-tcp
    port: 53
protocol: TCP
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    addonmanager.kubernetes.io/mode: EnsureExists
---
kind: Deployment
apiVersion: extensions/v1beta1
metadata:
  name: kube-dns
  namespace: kube-system
  labels:
    k8s-addon: kube-dns.addons.k8s.io
    k8s-app: kube-dns
    kubernetes.io/cluster-service: "true"
spec:
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 10%
      maxUnavailable: 0
  selector:
    matchLabels:
      k8s-app: kube-dns
  template:
    metadata:
      labels:
        k8s-app: kube-dns
      annotations:
        scheduler.alpha.kubernetes.io/critical-pod: ''
    spec:
      dnsPolicy: Default  # Don't use cluster DNS.
      serviceAccountName: kube-dns
      volumes:
      - name: kube-dns-config
        configMap:
          name: kube-dns
          optional: true
      nodeSelector:
        beta.kubernetes.io/arch: amd64

      containers:
      - name: kubedns
        image: gcr.io/google_containers/k8s-dns-kube-dns-amd64:{{ .ImageTag }}
        resources:
          limits:
            memory: 170Mi
          requests:
            cpu: 100m
            memory: 70Mi
        livenessProbe:
          httpGet:
            path: /healthcheck/kubedns
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        readinessProbe:
          httpGet:
            path: /readiness
            port: 8081
            scheme: HTTP
          initialDelaySeconds: 3
          timeoutSeconds: 5
        args:
        - --domain={{ .ClusterDomain }}.
        - --dns-port=10053
        - --config-dir=/kube-dns-config
        - --v=2
        env:
        - name: PROMETHEUS_PORT
          value: "10055"
        ports:
        - containerPort: 10053
          name: dns-local
          protocol: UDP
        - containerPort: 10053
          name: dns-tcp-local
          protocol: TCP
        - containerPort: 10055
          name: metrics
          protocol: TCP
        volumeMounts:
        - name: kube-dns-config
          mountPath: /kube-dns-config

      - name: dnsmasq
        image: gcr.io/google_containers/k8s-dns-dnsmasq-nanny-amd64:{{ .ImageTag }}
        livenessProbe:
          httpGet:
            path: /healthcheck/dnsmasq
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        args:
        - -v=2
        - -logtostderr
        - -configDir=/etc/k8s/dns/dnsmasq-nanny
        - -restartDnsmasq=true
        - --
        - -k
        - --cache-size=1000
        - --log-facility=-
        - --server=/{{ .ClusterDomain }}/127.0.0.1#10053
        - --server=/in-addr.arpa/127.0.0.1#10053
        - --server=/in6.arpa/127.0.0.1#10053
        ports:
        - containerPort: 53
          name: dns
          protocol: UDP
        - containerPort: 53
          name: dns-tcp
          protocol: TCP
        resources:
          requests:
            cpu: 150m
            memory: 20Mi
        volumeMounts:
        - name: kube-dns-config
          mountPath: /etc/k8s/dns/dnsmasq-nanny

      - name: sidecar
        image: gcr.io/google_containers/k8s-dns-sidecar-amd64:{{ .ImageTag }}
        livenessProbe:
          httpGet:
            path: /metrics
            port: 10054
            scheme: HTTP
          initialDelaySeconds: 60
          timeoutSeconds: 5
          successThreshold: 1
          failureThreshold: 5
        args:
        - --v=2
        - --logtostderr
        - --probe=kubedns,127.0.0.1:10053,kubernetes.default.svc.{{ .ClusterDomain }},5,A
        - --probe=dnsmasq,127.0.0.1:53,kubernetes.default.svc.{{ .ClusterDomain }},5,A
        ports:
        - containerPort: 10054
          name: metrics
          protocol: TCP
        resources:
          requests:
            memory: 20Mi
            cpu: 10m
`
