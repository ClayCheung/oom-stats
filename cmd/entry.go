package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"oom-stats/k8s"
	"oom-stats/oomStats"
	"sort"
)

const (
	WORKER_NUM int = 10
)

func init() {
	//set loglevel
	logrus.SetLevel(logrus.InfoLevel)
}

func worker(id int, auths <-chan oomStats.ConnectInfo, oomPodChan chan<- map[string][]map[string]string) {
	for auth := range auths {
		logrus.Debugf("worker:%d start, job: %v\n", id, auth)

		oomStats.GetOOMPodUid(auth, oomPodChan)
	}
}

func main() {
	// set flag
	k8sconfig := flag.String("k8sconfig", "./k8sconfig", "kubernetes config file path")
	cluster := flag.String("cluster", "compass-stack", "cluster load from ssh-config")
	flag.Parse()

	// get pod uid map
	clientset := k8s.NewClientSet(k8sconfig)
	uidMap, _ := k8s.GetAllPodUID(clientset)

	// load ssh-config
	authChan := make(chan oomStats.ConnectInfo)
	go oomStats.GetSSHConfig(*cluster, authChan)

	// every node to get oomPod
	oomPodChan := make(chan map[string][]map[string]string)

	oomStats.Wg.Add(oomStats.GetAuthLen(*cluster)) // wait group for close chan oomPodChan
	for w := 0; w < WORKER_NUM; w++ {
		go worker(w, authChan, oomPodChan)
	}

	//logrus.Debugln("close oomPodChan")
	//close(oomPodChan)

	// get result in Chan
	oomStatsChan := make(chan map[string][]map[string]string)

	go oomStats.GetOOMPod(uidMap, oomPodChan, oomStatsChan)

	// get result
	oomStatsResult := make(map[string][]map[string]string)
	go func() {
		for oom_stats := range oomStatsChan {
			for k, v := range oom_stats {
				logrus.Debugf("oomStatsChan :\n%v", oom_stats)
				oomStatsResult[k] = v
			}
		}
	}()

	oomStats.Wg.Wait() // waiting for stopping write data to oomPodChan
	close(oomPodChan)

	logrus.Debugf("oomStatsResult is : %v\n\n", oomStatsResult)

	// print result on stdout
	fmt.Printf("UID\t\t\t\t\tNAMESPACE\tPOD_NAME\tOOM_TIMES\n")
	hosts := []string{}
	for host := range oomStatsResult {
		hosts = append(hosts, host)
	}
	sort.Strings(hosts)
	for _, host := range hosts {
		fmt.Printf("%s:\n", host)
		for _, r := range oomStatsResult[host] {
			fmt.Printf("%s\t%s\t%s\t%s\n", r["uid"], r["namespace"], r["name"], r["oom_times"])
		}
	}
}
