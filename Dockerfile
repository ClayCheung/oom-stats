FROM busybox:1.31
MAINTAINER clay.zhang@outlook.com
WORKDIR /opt/oom-stats
ADD oom-stats oom-stats.sh ssh-config.yaml.temp ./

