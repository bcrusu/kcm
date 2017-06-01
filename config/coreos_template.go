package config

const coreOSCloudConfigTemplate = `#cloud-config

hostname: {{ .Name }}

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
    advertise-client-urls: http://${public_ip}:2379
    initial-advertise-peer-urls: http://${public_ip}:2380
    listen-client-urls: http://0.0.0.0:2379
    listen-peer-urls: http://${public_ip}:2380
    initial-cluster-state: new
    initial-cluster: ${etcd2_initial_cluster}
{{ end }}
  units:
    - name: etcd2.service
      command: start
      drop-ins:
        - name: 10-override-name.conf
          content: |
            [Service]
            Environment=ETCD_NAME=%H
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
        What=kubernetes
        Where=/opt/kubernetes
        Options=ro,trans=virtio,version=9p2000.L
        Type=9p
  update:
    group: {{ .CoreOSChannel }}
reboot-strategy: off
`
