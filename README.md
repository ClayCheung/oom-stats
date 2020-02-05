# kubernets pod OOM statistics
For some reason， sometimes events of k8s can not record pod oom, but log in linux os dmesg
this tool help to get pod OOM statistics
## Usage

### start by executeable file

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

### start as docker container

1. Prepare configuration file directory and set the configuration file
```bash
mkdir /tmp/oom-stats && cp ~/.kube/config /tmp/oom-stats/
vim /tmp/oom-stats/ssh-config.yaml
```
2. run container
```bash
# run a container
[root@kube-master-1 ~]$ docker run -it --rm --network=host -v=/tmp/oom-stats:/opt/oom-stats/config --name=oom-stats clayz95/oom-stats
# cmd in container
/opt/oom-stats # ./oom-stats -cluster=compass-stack -k8sconfig=config/config
```

## output
- output like this:
```bash
UID					NAMESPACE	POD_NAME	OOM_TIMES
192.168.21.197:
678e557a-30f7-11ea-a3e6-525400810531	N/A	N/A	1
192.168.21.199:
b5173eb8-30f7-11ea-a3e6-525400810531	N/A	N/A	34
192.168.21.224:
2c895f40-3292-11ea-af76-525400c705b2	N/A	N/A	1
192.168.21.225:
eef10ef9-30f7-11ea-a3e6-525400810531	N/A	N/A	13
192.168.21.227:
5c7beb48-30f8-11ea-a3e6-525400810531	N/A	N/A	3
192.168.21.229:
830a85b7-2e0a-11ea-a3e6-525400810531	kube-system	config-reference-reference-v1-0-688c446d86-c2plw	2
```
PS: if the `NAMESPACE` and `POD_NAME` is N/A, that means we couldn't find the pod with the `UID` in cluster currently.  
we find oom pod's UID by dmesg log, so we could still get the UID.  