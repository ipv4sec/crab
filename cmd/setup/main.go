package main

import (
	"bytes"
	"context"
	"crab/cluster"
	"crab/exec"
	"crab/utils"
	"errors"
	"fmt"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

func main() {
	var err error
	executor := exec.CommandExecutor{}

	klog.Infoln("开始集群认证")
	err = cluster.Init()
	if err != nil {
		panic(fmt.Errorf("获取集群认证失败: %w", err))
	}
	klog.Infoln("集群认证成功")

	klog.Infoln("现有集群检查")
	klog.Infoln("集群版本检查")
	ver, err := cluster.Client.Clientset.ServerVersion()
	if err != nil {
		panic(fmt.Errorf("获取集群版本失败: %w", err))
	}
	if ver.Major != "1" && ver.Minor != "20" {
		panic(errors.New(fmt.Sprintf("当前集群版本错误: %s", ver.String())))
	}

	klog.Infoln("集群组件检查")
	status := cluster.DiscoveryDeploymentsStatus("kube-system", []string{"coredns"})
	if status != cluster.ComponentAvailable {
		panic(errors.New("集群组件错误"))
	}

	klog.Infoln("网格组件检查")
	_, err = cluster.Client.Clientset.CoreV1().Namespaces().Get(context.Background(), "istio-system",
		metav1.GetOptions{})
	if err != nil {
		klog.Errorln(fmt.Errorf("获取网格命名空间错误: %w", err).Error())
		_, err = cluster.Client.Clientset.CoreV1().Namespaces().Create(context.Background(), &v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: "istio-system",
			},
		}, metav1.CreateOptions{})
		if err != nil {
			klog.Errorln(fmt.Errorf("创建网格命名空间错误: %w", err).Error())
		}
	}
	// klog.Infoln("ns:", ns)

	svcs, err := cluster.Client.Clientset.CoreV1().Services("istio-system").List(context.Background(),
		metav1.ListOptions{})
	if err != nil {
		panic(fmt.Errorf("列出集群资源错误: %w",err))
	}
	// klog.Infoln("svc:", svcs)
	var n = 0
	var components = []string{"istio-egressgateway", "istio-ingressgateway", "istiod"}
	for i := 0; i < len(components); i++ {
		for j := 0; j < len(svcs.Items); j++ {
			// klog.Infoln(svcs.Items[j].ObjectMeta.Name)
			if utils.Contains(components, svcs.Items[j].ObjectMeta.Name) {
				n++
			}
		}
	}
	if n ==0 {
		output, err := executor.ExecuteCommandWithCombinedOutput("scripts/istio.sh")
		if err != nil {
			panic(fmt.Errorf("初始化网格失败: %w", err))
		}
		klog.Infoln("初始化网格: ", output)
		yaml, err := ioutil.ReadFile("assets/istio/operator.yaml")
		if err != nil {
			panic(fmt.Errorf("读取yaml错误: %w", err))
		}
		err = cluster.Client.Apply(context.Background(), yaml)
		if err != nil {
			klog.Errorln("安装网格失败: ", err.Error())
		}
		for {
			if cluster.ComponentAvailable == cluster.DiscoveryDeploymentsStatus("istio-operator", []string{"istio-operator"}) {
				break
			}
			time.Sleep(time.Second * 5)
		}
		for {
			if cluster.ComponentAvailable == cluster.DiscoveryDeploymentsStatus("istio-system",
				[]string{"istio-egressgateway", "istio-ingressgateway", "istiod"}) {
				break
			}
			time.Sleep(time.Second * 5)
		}
	}
	if n != len(components) * len(svcs.Items) {
		panic(errors.New("网格中必备组件缺失"))
	}
	pods, err := cluster.Client.Clientset.CoreV1().Pods("istio-system").List(context.Background(), metav1.ListOptions{
		LabelSelector: "app=istio-ingressgateway",
	})
	if err != nil {
		panic(fmt.Errorf("列出资源错误: %w", err))
	}
	if len(pods.Items) == 0 {
		panic(errors.New("列出资源错误: 数量为空"))
	}

	req := cluster.Client.Clientset.CoreV1().RESTClient().Post().Resource("pods").Name(pods.Items[0].ObjectMeta.Name).
		Namespace("istio-system").SubResource("exec")
	req.VersionedParams(
		&v1.PodExecOptions{
			Command: []string{
				"sh",
				"-c",
				"env |grep ISTIO_META_ISTIO_VERSION",
			},
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
		},
		scheme.ParameterCodec,
	)

	var stdout, stderr bytes.Buffer
	conf, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("集群认证错误 :%w", err))
	}
	e, err := remotecommand.NewSPDYExecutor(conf, "POST", req.URL())
	if err != nil {
		panic(fmt.Errorf("执行POST错误 :%w", err))
	}
	err = e.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: &stdout,
		Stderr: &stderr,
	})
	if err != nil {
		panic(fmt.Errorf("执行命令错误 :%w %v %v", err, stdout.String(), stderr.String()))
	}
	v := strings.Trim(strings.TrimSpace(stdout.String()), "ISTIO_META_ISTIO_VERSION=")
	if !strings.HasPrefix(v, "1.9") {
		panic(errors.New(fmt.Sprintf("网格版本错误: %s", v)))
	}

	klog.Infoln("开始设置根域")
	yaml := `
apiVersion: v1
kind: ConfigMap
metadata:
  name: island-info
  namespace: island-system
data:
  root-domain: ""
  mirror: https://github.com/GlobalSphare/workloads
`
	err = cluster.Client.Apply(context.Background(), []byte(yaml))
	if err != nil {
		klog.Errorln("设置根域失败: ", err.Error())
	}
	klog.Infoln("设置根域完成")

	klog.Infoln("开始设置密码")
	_, err = cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		Create(context.Background(), &v1.ConfigMap{
			TypeMeta: metav1.TypeMeta{},
			ObjectMeta: metav1.ObjectMeta{
				Name: "island-administrator",
			},
			Data: map[string]string{"root": "toor"},
		}, metav1.CreateOptions{})
	if err != nil {
		klog.Errorln("设置密码失败: ", err.Error())
	}
	klog.Infoln("设置密码完成")

	klog.Infoln("开始部署应用")
	files, err := ioutil.ReadDir("assets/island/")
	if err != nil {
		panic(fmt.Errorf("读取应用列表错误 :%w", err))
	}
	for i := 0; i < len(files); i++ {
		klog.Infoln("要安装的应用为: ", files[i].Name())
		yaml, err := ioutil.ReadFile("assets/island/"+files[i].Name())
		if err != nil {
			panic(fmt.Errorf("读取yaml错误: %w", err))
		}
		err = cluster.Client.Apply(context.Background(), yaml)
		if err != nil {
			panic(fmt.Errorf("安装应用失败: %s %w", files[i].Name(), err))
		}
	}
	klog.Infoln("部署应用完成")
	klog.Info("结束退出程序")
}
