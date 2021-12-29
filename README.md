# Container Storage Interface (CSI) Driver for Distributed Asynchronous Object Storage (DAOS)



## Usage

The followings show how root user in a Pod accesses to a DAOS server with dfuse (tested on Ubuntu 20.04).

**Note**

* It is assumed that host machine has IP address: 192.168.122.1 and it is accessible from minikube Pods.
* It is recommended that Host machine has more than 32 GiB memory.


### Build CSI driver

```bash
make
```

### Start minikube

```bash
minikube start --cpus=8 --memory=16G
```

### Build CSI driver image and put it on minikube's registry

```bash
eval $(minikube -p minikube docker-env)
docker build https://github.com/daos-stack/daos.git#release/1.2 \
        -f utils/docker/Dockerfile.centos.7 -t daos
docker build . -t dfuse-csi-plugin
```

### Build DAOS server image and put it on localhost

```bash
sudo docker build https://github.com/daos-stack/daos.git#release/1.2 \
        -f utils/docker/Dockerfile.centos.7 -t daos
cd server
sudo docker build . -t daos_server
cd ../
```

### Start DAOS server

```bash
sudo docker run -it -d --network=host --privileged --cap-add=ALL --name server -v /dev:/dev daos_server
```

### Create Storage Pool

```
sudo docker exec -it server /bin/bash
sudo su
/opt/daos/bin/dmg pool create -i --size=8G
#[OUTPUT]
#Creating DAOS pool with automatic storage allocation: 8.0 GB NVMe + 6.00% SCM
#Pool created with 100.00% SCM/NVMe ratio
#-----------------------------------------
#  UUID          : 7170ec68-52d5-4f39-98fb-27494cabb47c
#  Service Ranks : 0
#  Storage Ranks : 0
#  Total Size    : 8.0 GB
#  SCM           : 8.0 GB (8.0 GB / rank)
#  NVMe          : 0 B (0 B / rank)
/opt/daos/bin/dmg pool -i get-acl --pool=7170ec68-52d5-4f39-98fb-27494cabb47c
#[OUTPUT]
## Owner: root@
# Owner Group: root@
# Entries:
#A::OWNER@:rw
#A:G:GROUP@:rw
exit
```

### Deploy CSI driver

```bash
cd deploy/dfuse
kubectl create -f csi-dfuse-driverinfo.yaml 
# csidriver.storage.k8s.io/dfuse.csi.k8s.io created
kubectl get csidriver
# NAME               ATTACHREQUIRED   PODINFOONMOUNT   MODES                  AGE
# dfuse.csi.k8s.io   false            false            Persistent,Ephemeral   15s
kubectl create -f csi-dfuse-nodeplugin.yaml 
# daemonset.apps/csi-dfuse-nodeplugin created
kubectl -n kube-system get pod
# NAME                               READY   STATUS    RESTARTS   AGE
# coredns-74ff55c5b-9mg49            1/1     Running   2          40h
# csi-dfuse-nodeplugin-t7lq9         2/2     Running   0          12s
# ...
```

### Deploy Pod

Copy the pool UUID above and modify `deploy/dfuse/example/app.yaml`.

```yaml
# app.yaml
# ...
        volumeAttributes:
            uid: "0"
            poolid: "7170ec68-52d5-4f39-98fb-27494cabb47c" # Paste your UUID
```

```bash
cd example
kubectl create -f app.yaml
# persistentvolume/pv-dfuseplugin created
# persistentvolumeclaim/pvc-dfuseplugin created
# pod/test-app created
kubectl get po
# NAME       READY   STATUS    RESTARTS   AGE
# test-app   1/1     Running   0          6s
```

### Login to the app container and check if DAOS mounted successfully

```bash
kubectl exec -it test-app -- /bin/sh
/ # df
Filesystem           1K-blocks      Used Available Use% Mounted on
overlay              959863856 174735984 736299664  19% /
tmpfs                    65536         0     65536   0% /dev
tmpfs                 49402308         0  49402308   0% /sys/fs/cgroup
dfuse                  7812500       281   7812219   0% /data <<<<<
/dev/nvme0n1p2       959863856 174735984 736299664  19% /dev/termination-log
/dev/nvme0n1p2       959863856 174735984 736299664  19% /etc/resolv.conf
/dev/nvme0n1p2       959863856 174735984 736299664  19% /etc/hostname
/dev/nvme0n1p2       959863856 174735984 736299664  19% /etc/hosts
shm                      65536         0     65536   0% /dev/shm
tmpfs                 49402308        12  49402296   0% /var/run/secrets/kubernetes.io/serviceaccount
tmpfs                 49402308         0  49402308   0% /proc/acpi
tmpfs                    65536         0     65536   0% /proc/kcore
tmpfs                    65536         0     65536   0% /proc/keys
tmpfs                    65536         0     65536   0% /proc/timer_list
tmpfs                    65536         0     65536   0% /proc/sched_debug
tmpfs                 49402308         0  49402308   0% /proc/scsi
tmpfs                 49402308         0  49402308   0% /sys/firmware
/ # dd if=/dev/zero of=/data/file bs=1M count=1024
1024+0 records in
1024+0 records out
1073741824 bytes (1.0GB) copied, 0.867793 seconds, 1.2GB/s
```
