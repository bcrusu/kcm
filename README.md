## kcm: Kubernetes cluster manager

### Prerequisites

[libvirt](https://libvirt.org)

[QEMU](http://www.qemu.org)/[KVM](https://www.linux-kvm.org/page/Main_Page)

### Installation

TODO

### Usage

#### Create a cluster:
```
kcm create mykube
```
Creates a cluster named 'mykube', with 3 minion nodes and one master

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

### Items left to do:

- [ ] Bug: add socat binary to CoreOS image (kubectl port-forward does not work without it)
- [ ] Support clusters with multiple master nodes (via nginx/HAProxy)
- [ ] More netorking options (e.g. weave, calico, etc.)
- [ ] Allow users to pass configuration settings to newly-created clusters (e.g. all vars with prefix 'KCM_' should be made available to Kubernetes)
- [ ] Better status command output: fetch and display k8s cluster status using client-go library
- [ ] Add 'verbose' flag
- [ ] Add 'yes' flag for remove command
- [ ] Check if KSM is not enabled & warn user
- [ ] Better logging and output error messages

### Inspiration

[CCM (Cassandra Cluster Manager)](https://github.com/pcmanus/ccm): A script to easily create and destroy an Apache Cassandra cluster on localhost

[kube-up (deprecated)](https://github.com/kubernetes/kubernetes/tree/master/cluster)
