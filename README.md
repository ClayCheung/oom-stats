# kubernets pod OOM statistics
For some reason， sometimes events of k8s can not record pod oom, but log in linux os dmesg
this tool help to get pod OOM statistics
## Usage
1. Put the `ssh-config.yaml` files in the same directory of executable file `oom-stats`
> example `ssh-config.yaml`
```
cluster01:
  auth:
  - 192.168.21.224 22 root Your_Password 
  - 192.168.21.225 22 root Your_Password
  - 192.168.21.226 22 root Your_Password
  - 192.168.21.227 22 root Your_Password
  - 192.168.21.228 22 root Your_Password
  - 192.168.21.229 22 root Your_Password
  - 192.168.21.197 22 root Your_Password 
  - 192.168.21.199 22 root Your_Password
cluster02:
  auth:
  - 192.168.1.100 22 root Your_Password
  - 192.168.1.101 22 root Your_Password

```

2. Put the `oom-stats.sh` files in the same directory of executable file `oom-stats`
> `oom-stats.sh` files is order to get oom pod UID and times

3. Run cmd with your parameters
> k8sconfig: kubeconfig file path
> cluster: cluster in ssh-config.yaml
```
oom-stats --k8sconfig="kubeconfig_path" --cluster="k8s_cluster_name"
```