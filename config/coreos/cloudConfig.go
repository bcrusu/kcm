package coreos

import "github.com/bcrusu/kcm/util"

type CloudConfigParams struct {
	Hostname          string
	IsMaster          bool
	SSHPublicKey      string
	NonMasqueradeCIDR string
	Network           util.NetworkInfo
	ClusterDomain     string
}

const cloudConfigTemplate = `#cloud-config

hostname: {{ .Hostname }}

ssh_authorized_keys:
  - {{ .SSHPublicKey }}

write_files:
  - path: /etc/systemd/journald.conf
    permissions: 0644
    content: |
      [Journal]
      SystemMaxUse=50M
      RuntimeMaxUse=50M

coreos:
{{ if .IsMaster }}
  etcd2:
    advertise-client-urls: http://0.0.0.0:2379
    listen-client-urls: http://0.0.0.0:2379
    listen-peer-urls: http://0.0.0.0:2380
    initial-cluster-state: new
    initial-cluster: {{ .Hostname }}=http://0.0.0.0:2380
    initial-advertise-peer-urls: http://0.0.0.0:2380
{{ end }}

  units:
{{ if .IsMaster }}
    - name: etcd2.service
      command: start
      drop-ins:
        - name: 10-override-name.conf
          content: |
            [Service]
            Environment=ETCD_NAME=%H

# only one master is supported atm. - it gets the 2nd IP in the network (1st IP is assigned to the bridge/gateway)
    - name: static.network
      command: start
      content: |
        [Match]
        Name=eth0
        [Network]
        Address={{ .Network.MasterAddress }}
        DNS={{ .Network.BridgeIP }}
        Gateway={{ .Network.BridgeIP }}
{{ end }}

    - name: docker.service
      command: start
      drop-ins:
        - name: 50-opts.conf
          content: |
            [Service]
            Environment='DOCKER_OPTS=--iptables=false'
    - name: docker-tcp.socket
      command: start
      enable: yes
      content: |
        [Unit]
        Description=Docker Socket for the API
        [Socket]
        ListenStream=2375
        BindIPv6Only=both
        Service=docker.service
        [Install]
        WantedBy=sockets.target

    - name: opt-kubernetes.mount
      command: start
      content: |
        [Unit]
        ConditionVirtualization=|vm
        [Mount]
        What=k8sConfig
        Where=/opt/kubernetes
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p
    - name: opt-kubernetes-bin.mount
      command: start
      content: |
        [Unit]
        ConditionVirtualization=|vm
        After=opt-kubernetes.mount
        Requires=opt-kubernetes.mount
        [Mount]
        What=k8sBin
        Where=/opt/kubernetes/bin
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p
    - name: opt-kubernetes-manifests.mount
      command: start
      content: |
        [Unit]
        ConditionVirtualization=|vm
        After=opt-kubernetes.mount
        Requires=opt-kubernetes.mount
        [Mount]
        What=k8sConfigManifests
        Where=/opt/kubernetes/manifests
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p

    - name: kubelet.service
      command: start
      content: |
        [Unit]
        After=opt-kubernetes-bin.mount opt-kubernetes-manifests.mount docker.service load-k8s-images.service
        ConditionFileIsExecutable=/opt/kubernetes/bin/kubelet
        Description=Kubernetes Kubelet Server
        Documentation=https://github.com/kubernetes/kubernetes
        Requires=opt-kubernetes-bin.mount opt-kubernetes-manifests.mount docker.service load-k8s-images.service

        [Service]
        Restart=always
        RestartSec=2
        StartLimitInterval=0
        KillMode=process
        ExecStart=/opt/kubernetes/bin/kubelet \
        --address=0.0.0.0 \
        --hostname-override={{ .Hostname }} \
        --cluster-domain={{ .ClusterDomain }} \
        --kubeconfig=/opt/kubernetes/kubeconfig \
        --require-kubeconfig=true \
        --register-node=true \
        --node-labels='kubernetes.io/role={{ Role }},node-role.kubernetes.io/{{ Role }}=' \
        --network-plugin=cni \
        --non-masquerade-cidr={{ .NonMasqueradeCIDR }} \
        --allow-privileged=true \
        --pod-manifest-path=/opt/kubernetes/manifests \
        --tls-cert-file=/opt/kubernetes/certs/tls.pem \
        --tls-private-key-file=/opt/kubernetes/certs/tls-key.pem \
        --client-ca-file=/opt/kubernetes/certs/ca.pem {{ if .IsMaster }}\
        --register-with-taints='node-role.kubernetes.io/master=:NoSchedule'{{ end }}

        [Install]
        WantedBy=multi-user.target

    - name: load-k8s-images.service
      command: start
      content: |
        [Unit]
        Description=Load Kubernetes images to Docker
        After=opt-kubernetes-bin.mount docker.service
        Requires=opt-kubernetes-bin.mount docker.service

        [Service]
        WorkingDirectory=/opt/kubernetes/bin
        Type=forking
        KillMode=process
{{ if .IsMaster }}
        ExecStart=/bin/bash -c 'docker load -i kube-apiserver.tar && docker load -i kube-controller-manager.tar &&docker load -i kube-scheduler.tar'
{{ else }}
        ExecStart=/usr/bin/docker load -i kube-proxy.tar
{{ end }}
`
