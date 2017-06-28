## kcm: Kubernetes cluster manager

### Prerequisites

[libvirt](https://libvirt.org)

[QEMU](http://www.qemu.org)/[KVM](https://www.linux-kvm.org/page/Main_Page)

### Installation

```
go install github.com/bcrusu/kcm
```

### Usage

#### Create a cluster:
```
kcm create mykube
```
Creates a cluster named 'mykube', with 2 minion nodes and one master

```
kcm create mykube --node-count=5 --start
```
Creates a cluster with 5 nodes and one master and starts it immediately

#### Start/stop a cluster:
```
kcm start mykube
kcm stop mykube
```

#### Remove a cluster:
```
kcm remove cluster mykube
```
Removes the cluster named 'mykube' and its artefacts (i.e. libvirt objects and files on disk)

#### Use kubectl to interact with the cluster:
```
kcm ctl get pods
kcm ctl apply -f FILENAME
...
```
The 'ctl' command calls the right version of kubectl binary and sets the "--kubeconfig" argument. It uses the following files:
* kubectl: ~/.kcm/cache/kubernetes/KUBE_VERSION/kubernetes/server/bin/kubectl 
* kubeconfig: ~/.kcm/config/CLUSTER_NAME/kubeconfig/kubectl

#### Get cluster status:
```
kcm status
```
Outputs information similar to:
```
CLUSTER   STATUS    DNS DOMAIN      KUBE VERSION   COREOS VERSION
mykube    Active    mykube.kube     1.7.0-beta.2   stable/1353.8.0

NETWORK      STATUS    CIDR          DNS SERVER
kcm.mykube   Active    10.1.0.0/16   10.1.0.1

NODE      STATUS    DNS NAME             DNS LOOKUP   IP
master    Active    master.mykube.kube   OK           10.1.238.138
node1     Active    node1.mykube.kube    OK           10.1.199.19
node2     Active    node2.mykube.kube    OK           10.1.97.155
```

### Items left to do:

- [ ] Add Kubernetes Dashboard
- [ ] Support clusters with multiple master nodes (via nginx/HAProxy)
- [ ] More netorking options (e.g. weave, calico, etc.)
- [ ] Allow users to pass configuration settings to newly-created clusters (e.g. all vars with prefix 'KCM_' should be made available to Kubernetes)

### Inspiration

[CCM (Cassandra Cluster Manager)](https://github.com/pcmanus/ccm): A script to easily create and destroy an Apache Cassandra cluster on localhost

[kube-up (deprecated)](https://github.com/kubernetes/kubernetes/tree/master/cluster)
