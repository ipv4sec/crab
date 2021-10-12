package cluster

import (
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

type ComponentCode int
const (
	ComponentError ComponentCode = iota
	ComponentNotReady
	ComponentAvailable
)

type ComponentStatus struct {
	Name string `json:"name"`
	Status int `json:"status"`
}

func (c ComponentCode) ByName(name string) ComponentStatus {
	if c == ComponentError {
		return ComponentStatus{
			Name:   name,
			Status: 0,
		}
	}
	if c == ComponentNotReady {
		return ComponentStatus{
			Name:   name,
			Status: 1,
		}
	}
	if c == ComponentAvailable {
		return ComponentStatus{
			Name:   name,
			Status: 2,
		}
	}
	return ComponentStatus{
		Name:   name,
		Status: 0,
	}
}

func DiscoveryDeploymentsStatus(namespace string, deployments []string) ComponentCode {
	v := 0
	for i := 0; i < len(deployments); i++ {
		deploy, err := Client.Clientset.AppsV1().Deployments(namespace).
			Get(context.Background(), deployments[i], metav1.GetOptions{})
		if err != nil {
			klog.Errorln(deployments[i], err.Error())
			return ComponentError
		}
		if deploy.Status.AvailableReplicas == 0 {
			v += 1
			continue
		}
		v += 2
	}
	if v == len(deployments) * 2 {
		return ComponentAvailable
	}
	return ComponentNotReady

}