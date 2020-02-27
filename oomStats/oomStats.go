package oomStats

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
	"oom-stats/sftp"
	"oom-stats/ssh"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
)

const (
	dstDirPath  = "/tmp"
	shScript = `
#!/bin/bash
res=$(dmesg -LT|grep "killed as a result of limit"|awk -F ' ' '{print $8}'|awk -F '/' '{print $(NF-1)}'|cut -c 4- | awk '{a[$0]++}END{for(i in a){print i,a[i]}}')
echo -e "${res}\n"
`
)

var (
	shFileName = "oom-stats.sh"
	cmdList = []string{
		"chmod 777 /tmp/"+shFileName,
		"/tmp/"+shFileName,
		"rm -f /tmp/"+shFileName,
	}
	Wg sync.WaitGroup
)

type ConnectInfo struct {
	user     string
	host     string
	port     int
	password string
}

func NewConnectInfo(user, host string, port int, password string) ConnectInfo {
	return ConnectInfo{
		user:     user,
		host:     host,
		port:     port,
		password: password,
	}
}

func genTmpBashScript() (*string, error) {
	tDir, err := ioutil.TempDir("", "oom-stats")
	if err != nil {
		return nil, err
	}
	tFile := path.Join(tDir, shFileName)
	if err = ioutil.WriteFile(tFile, []byte(shScript), 0777); err != nil {
		return nil, err
	}

	return &tFile, nil
}

func GetOOMPodUid(c ConnectInfo, oomPodChan chan<- map[string][]map[string]string) {
	defer Wg.Done()
	// generate bash script file tmp
	shFile, err := genTmpBashScript()
	if err != nil {
		logrus.Fatal(err)
	}
	defer func() {
		os.Remove(*shFile)
	}()
	// sftp put bash
	sftpClient, err := sftp.Connect(c.user, c.password, c.host, c.port)
	if err != nil {
		logrus.Fatal(err)
	}
	defer sftpClient.Close()
	fmt.Println(*shFile)
	if err := sftp.Put(sftpClient, *shFile, dstDirPath); err != nil {
		logrus.Fatal(err)
	}
	// ssh run bash script
	returnList := make([]string, 0)

	for _, cmd := range cmdList {
		session, err := ssh.Connect(c.user, c.password, c.host, c.port)
		if err != nil {
			logrus.Fatal(err)
		}
		defer session.Close()

		output, err := ssh.RunCmd(session, cmd)
		logrus.Debugf("%s", output)
		returnList = append(returnList, output)
	}
	logrus.Debugf("CMD output:\n%s", returnList[1])

	// return oom pod's uid and oom times
	logrus.Debugf("goroutine GetOOMPodUid: %v not done\n", c)

	oomPodChan <- get_UID_OOMTimes(returnList[1], c.host)

	logrus.Debugf("goroutine GetOOMPodUid: %v done\n", c)

}

func get_UID_OOMTimes(s string, host string) map[string][]map[string]string {
	r := make(map[string][]map[string]string)
	logrus.Debugf("func: get_UID_OOMTimes: s is: %v \n", s)
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		if line != "" {
			tuple := strings.Fields(line)
			r[host] = append(r[host], map[string]string{
				"uid":       tuple[0],
				"oom_times": tuple[1],
			})
		}
	}
	logrus.Debugf("UID_OOMTimes map:\n%v", r)

	return r
}

func GetOOMPod(uidPodMap map[string]map[string]string, oomPodChan chan map[string][]map[string]string,
	oomStatsChan chan<- map[string][]map[string]string) {
	logrus.Debugf("func: GetOOMPod: oomPodChan %v\n", oomPodChan)
	for oomPod := range oomPodChan { // 会一直从 oomPodChan 中读取数据 deadlock，只要 oomPodChan 没有关闭
		logrus.Debugf("func: GetOOMPod: oomPod %v\n", oomPod)
		for _, mapList := range oomPod {
			for _, m := range mapList {
				uid := m["uid"]
				if pod, ok := uidPodMap[uid]; ok {
					m["name"] = pod["name"]
					m["namespace"] = pod["namespace"]
				} else {
					m["name"] = "N/A"
					m["namespace"] = "N/A"
				}
			}
		}
		oomStatsChan <- oomPod
	}

	logrus.Debugln("close oomStatsChan")
	close(oomStatsChan)

}

func GetSSHConfig(cluster string, authChan chan ConnectInfo) {
	viper.SetConfigName("ssh-config") // without suffix
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	authList := viper.GetStringSlice(cluster + ".auth")

	for _, auth := range authList {
		p, err := strconv.Atoi(strings.Fields(auth)[1])
		if err != nil {
			logrus.Fatalf("config error: port is not a number ", err)
		}
		authChan <- ConnectInfo{
			host:     strings.Fields(auth)[0],
			port:     p,
			user:     strings.Fields(auth)[2],
			password: strings.Fields(auth)[3],
		}
		logrus.Debugf("send auth: %v\n", ConnectInfo{})
	}
	close(authChan)

}

func GetAuthLen(cluster string) int {
	viper.SetConfigName("ssh-config") // without suffix
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	authList := viper.GetStringSlice(cluster + ".auth")
	return len(authList)
}
