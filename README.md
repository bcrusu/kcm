## kcm: Kubernetes cluster manager

### Prerequisites

[libvirt](https://libvirt.org)

[QEMU](http://www.qemu.org)/[KVM](https://www.linux-kvm.org/page/Main_Page)

### Installation

TODO

### Usage

TODO

### Items left to do:

- [ ] Export cluster details to kubectl (i.e. kubectl config set-cluster)
- [ ] Support clusters with multiple master nodes (via nginx/HAProxy)
- [ ] Allow users to pass configuration settings to newly-created clusters (e.g. all vars with previx 'KCM_' should be made available to Kubernetes)
- [ ] Use the Libvirt StorageVolume.Upload(...) API call to add the base CoreOS image to the pool
- [ ] Add 'verbose' flag
- [ ] Add 'yes' flag for remove command
- [ ] Check if KSM is not enabled & warn user
- [ ] Better logging and output error messages

### Inspiration

[CCM (Cassandra Cluster Manager)](https://github.com/pcmanus/ccm): A script to easily create and destroy an Apache Cassandra cluster on localhost
[kube-up (deprecated)](https://github.com/kubernetes/kubernetes/tree/master/cluster)
