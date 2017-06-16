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

### Items left to do:

- [ ] Export cluster details to kubectl (i.e. kubectl config set-cluster)
- [ ] Support clusters with multiple master nodes (via nginx/HAProxy)
- [ ] Add kube-dns addon
- [ ] More netorking options (e.g. weave, calico, etc.)
- [ ] Allow users to pass configuration settings to newly-created clusters (e.g. all vars with previx 'KCM_' should be made available to Kubernetes)
- [ ] Better status command output: fetch and display k8s cluster status using client-go library
- [ ] Use the Libvirt StorageVolume.Upload(...) API call to add the base CoreOS image to the pool
- [ ] Add 'verbose' flag
- [ ] Add 'yes' flag for remove command
- [ ] Check if KSM is not enabled & warn user
- [ ] Better logging and output error messages

### Inspiration

[CCM (Cassandra Cluster Manager)](https://github.com/pcmanus/ccm): A script to easily create and destroy an Apache Cassandra cluster on localhost

[kube-up (deprecated)](https://github.com/kubernetes/kubernetes/tree/master/cluster)
