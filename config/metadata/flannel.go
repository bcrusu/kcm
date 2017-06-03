package metadata

type flannelTemplateParams struct {
	ImageTag        string
	PodsNetworkCIDR string
}

const flannelTemplate = `kind: ServiceAccount
apiVersion: v1
metadata:
  name: flannel
  namespace: kube-system
  labels:
    role.kubernetes.io/networking: "1"
---
kind: ConfigMap
apiVersion: v1
metadata:
  name: kube-flannel-cfg
  namespace: kube-system
  labels:
    k8s-app: flannel
    role.kubernetes.io/networking: "1"
data:
  cni-conf.json: |
    {
      "name": "cbr0",
      "type": "flannel",
      "delegate": {
        "forceAddress": true,
        "isDefaultGateway": true,
        "hairpinMode": true
      }
    }
  net-conf.json: |
    {
      "Network": "{{ .PodsNetworkCIDR }}",
      "Backend": {
        "Type": "udp"
      }
    }
---
kind: DaemonSet
apiVersion: extensions/v1beta1
metadata:
  name: kube-flannel-ds
  namespace: kube-system
  labels:
    k8s-app: flannel
    role.kubernetes.io/networking: "1"
spec:
  template:
    metadata:
      labels:
        tier: node
        app: flannel
        role.kubernetes.io/networking: "1"
    spec:
      hostNetwork: true
      nodeSelector:
        beta.kubernetes.io/arch: amd64
      serviceAccountName: flannel
      tolerations:
      - effect: NoSchedule
        key: node-role.kubernetes.io/master
      containers:
      - name: kube-flannel
        image: quay.io/coreos/flannel:{{ .ImageTag }}
        command: [ "/opt/bin/flanneld", "--ip-masq", "--kube-subnet-mgr" ]
        securityContext:
          privileged: true
        env:
        - name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        - name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        resources:
          limits:
            cpu: 100m
            memory: 100Mi
          requests:
            cpu: 100m
            memory: 100Mi
        volumeMounts:
        - name: run
          mountPath: /run
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      - name: install-cni
        image: quay.io/coreos/flannel:{{ .ImageTag }}
        command: [ "/bin/sh", "-c", "set -e -x; cp -f /etc/kube-flannel/cni-conf.json /etc/cni/net.d/10-flannel.conf; while true; do sleep 3600; done" ]
        resources:
          limits:
            cpu: 10m
            memory: 25Mi
          requests:
            cpu: 10m
            memory: 25Mi
        volumeMounts:
        - name: cni
          mountPath: /etc/cni/net.d
        - name: flannel-cfg
          mountPath: /etc/kube-flannel/
      volumes:
        - name: run
          hostPath:
            path: /run
        - name: cni
          hostPath:
            path: /etc/cni/net.d
        - name: flannel-cfg
          configMap:
            name: kube-flannel-cfg
`
