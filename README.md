## kcm: Kubernetes cluster manager

### Prerequisites

TODO

libvirt
KVM/QEMU

### Installation

TODO

### Usage

TODO

### Items left to do:

- [ ] Export cluster details to kubectl (i.e. kubectl config set-cluster)
- [ ] Use the Libvirt StorageVolume.Upload(...) API call to add the base CoreOS image to the pool
- [ ] Add 'verbose' flag
- [ ] Add 'yes' flag for remove command
- [ ] Check if KSM is not enabled & warn user

### Inspiration

[CCM (Cassandra Cluster Manager)](https://github.com/pcmanus/ccm): A script to easily create and destroy an Apache Cassandra cluster on localhost

