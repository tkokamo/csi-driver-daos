# csi-driver-daos

```
sudo ./csc --endpoint unix:///tmp/csi.sock node publish volid --cap=5,1 --target-path=`pwd`/dfuse --vol-context uid=$UID,size=8G
```

```
sudo ./csc --endpoint unix:///tmp/csi.sock node unpublish volid --target-path=`pwd`/dfuse
```


start server

```
sudo docker run -it -d --network=host --privileged --cap-add=ALL --name server -v /dev:/dev daos_server:v1.1
```

start plugin container

```
sudo docker run -it -d --privileged --cap-add=ALL --name dfusep -v /dev:/dev dfuse-csi-plugin
```

# Build container images

```
eval $(minikube -p minikube docker-env)

# Build DAOS base image
docker build https://github.com/daos-stack/daos.git#release/1.2 \
        -f utils/docker/Dockerfile.centos.7 -t daos

# Build CSI driver
docker build . -t dfuse-csi-plugin

# Build DAOS server
```

# Usage:

```
(topse2021) (base) takuya@takuya-MS-7C35:~/git/csi-driver-daos/csi-driver-daos/deploy/dfuse/example$ kubectl -n kube-system logs csi-dfuse-nodeplugin-8nkkn dfuse
I1228 15:34:43.883486       9 server.go:117] Listening for connections on address: &net.UnixAddr{Name:"//csi/csi.sock", Net:"unix"}
unix:///csi/csi.sockDAOS Agent v1.2 (pid 10) listening on /var/run/daos_agent/daos_agent.sock
I1228 15:34:46.076447       9 nodeserver.go:28] NodePublishVolume: volume_id:"csi-4dc79a0eb614edba3f4ef498448e96b70e5470c054a67f590df3bd964a3fa47a" target_path:"/var/lib/kubelet/pods/c53958eb-b849-4d1a-b8a0-4aa1c2be63ac/volumes/kubernetes.io~csi/dfuse/mount" volume_capability:<mount:<> access_mode:<mode:SINGLE_NODE_WRITER > > 
I1228 15:34:48.358497       9 nodeserver.go:48] cmd: /opt/daos/bin/dmg pool create -i --ranks=0 --scm-size=8G | grep UUID | awk '{print $3}'
I1228 15:34:48.358524       9 nodeserver.go:49] poolid: 55cf4674-5241-40b0-a071-d08046b3c85e
I1228 15:34:49.176584       9 nodeserver.go:59] cmd: /opt/daos/bin/daos cont create --pool=55cf4674-5241-40b0-a071-d08046b3c85e --type=POSIX | awk '{print $4}'
I1228 15:34:49.176598       9 nodeserver.go:60] contid: cd8c00eb-9842-4b46-ada1-54fbe6f26cf9
I1228 15:34:49.721925       9 nodeserver.go:73] output: 
(topse2021) (base) takuya@takuya-MS-7C35:~/git/csi-driver-daos/csi-driver-daos/deploy/dfuse/example$ kubectl get po
NAME                   READY   STATUS    RESTARTS   AGE
csi-hostpathplugin-0   8/8     Running   8          26h
test-app               1/1     Running   0          57s
(topse2021) (base) takuya@takuya-MS-7C35:~/git/csi-driver-daos/csi-driver-daos/deploy/dfuse/example$ kubectl exec -it test-app -- /bin/sh
/ # df
Filesystem           1K-blocks      Used Available Use% Mounted on
overlay              959863856 174062164 736973484  19% /
tmpfs                    65536         0     65536   0% /dev
tmpfs                 49402308         0  49402308   0% /sys/fs/cgroup
dfuse                  7812500       144   7812356   0% /data
/dev/nvme0n1p2       959863856 174062164 736973484  19% /dev/termination-log
/dev/nvme0n1p2       959863856 174062164 736973484  19% /etc/resolv.conf
/dev/nvme0n1p2       959863856 174062164 736973484  19% /etc/hostname
/dev/nvme0n1p2       959863856 174062164 736973484  19% /etc/hosts
shm                      65536         0     65536   0% /dev/shm
tmpfs                 49402308        12  49402296   0% /var/run/secrets/kubernetes.io/serviceaccount
tmpfs                 49402308         0  49402308   0% /proc/acpi
tmpfs                    65536         0     65536   0% /proc/kcore
tmpfs                    65536         0     65536   0% /proc/keys
tmpfs                    65536         0     65536   0% /proc/timer_list
tmpfs                    65536         0     65536   0% /proc/sched_debug
tmpfs                 49402308         0  49402308   0% /proc/scsi
tmpfs                 49402308         0  49402308   0% /sys/firmware
/ # dd if=/dev/zero of=/data bs=1M count=1024
dd: can't open '/data': Is a directory
/ # dd if=/dev/zero of=/data/file bs=1M count=1024
1024+0 records in
1024+0 records out
1073741824 bytes (1.0GB) copied, 0.887934 seconds, 1.1GB/s

``` 
