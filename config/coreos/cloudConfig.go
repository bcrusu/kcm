package coreos

import "github.com/bcrusu/kcm/util"

type CloudConfigParams struct {
	Hostname          string
	IsMaster          bool
	SSHPublicKey      string
	NonMasqueradeCIDR string
	Network           util.NetworkInfo
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
            Environment='DOCKER_OPTS=--bridge=cbr0 --iptables=false'
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
    - name: opt-kubernetes-metadata.mount
      command: start
      content: |
        [Unit]
        ConditionVirtualization=|vm
        After=opt-kubernetes.mount
        Requires=opt-kubernetes.mount
        [Mount]
        What=k8sConfigMetadata
        Where=/opt/kubernetes/metadata
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p
`
