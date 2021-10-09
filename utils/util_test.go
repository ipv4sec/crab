package utils

import (
	"errors"
	"fmt"
	"k8s.io/klog/v2"
	"testing"
)

func TestContains(t *testing.T) {
	if !Contains([]string{"aaa","bbb", "ccc"}, "aaa") {
		panic(errors.New("aaa"))
	}
	var n = 0
	var components = []string{"istio-egressgateway", "istio-ingressgateway", "istiod"}
	var svc = []string{"istio-egressgateway", "istio-ingressgateway", "sss"}
	for i := 0; i < len(components); i++ {
		for j := 0; j < len(svc); j++ {
			klog.Infoln(svc[j])
			if Contains(components, svc[j]) {
				n++
			}
		}
	}
	fmt.Println(n)
}
