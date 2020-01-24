package k8s

import (
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetAllPodUID(clientset *kubernetes.Clientset) (map[string]map[string]string, error) {
	//var uidMap map[string]map[string]string
	// eg:
	//	{
	//		"00ae3665-11a1-11ea-931c-5254005453a1": {		// UID
	//				"name":"kube-apiserver-8fvh2" ,	// pod name
	//				"namespace": "kube-system",				// namespace
	//		},
	//	}
	uidMap := make(map[string]map[string]string)

	//get POD
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		logrus.Errorln(err)
		return nil, err
	}

	for _, pod := range pods.Items {
		logrus.Debugf("%s\t%s\t%s\n", pod.Namespace, pod.Name, pod.UID)
		uidMap[string(pod.UID)] = map[string]string{
			"name":      pod.Name,
			"namespace": pod.Namespace,
		}
	}

	logrus.Infoln("get uid pod map success")
	logrus.Debugf("uid map is: %v\n\n", uidMap)

	return uidMap, nil
}
