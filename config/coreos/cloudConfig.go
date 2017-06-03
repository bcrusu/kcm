package coreos

type CloudConfigParams struct {
	Hostname        string
	IsMaster        bool
	SSHPublicKey    string
	PodsNetworkCIDR string
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
#TODO: etcd config
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
{{ end }}

    - name: docker.service
      command: start
      drop-ins:
        - name: 50-opts.conf
          content: |
            [Service]
            Environment='DOCKER_OPTS=--bridge=cbr0 --iptables=false'
            
    - name: opt-kubernetes.mount
      command: start
      content: |
        [Unit]
        ConditionVirtualization=|vm
        [Mount]
        What=kubernetesConfig
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
        What=kubernetesBin
        Where=/opt/kubernetes/bin
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p
`
