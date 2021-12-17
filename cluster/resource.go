package cluster

import (
	"context"
	"crab/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func GetResourceHandlerFunc(c *gin.Context) {
	namespace := c.Param("namespace")
	resourceType := c.Param("resourceType")
	resourceName := c.Param("resourceName")
	var v interface{}
	var err error
	switch resourceType {
	case "cronJob":
		v, err = Client.Clientset.BatchV1().CronJobs(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "daemonSet":
		v, err = Client.Clientset.AppsV1().DaemonSets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "deployment":
		v, err = Client.Clientset.AppsV1().Deployments(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "job":
		v, err = Client.Clientset.BatchV1().Jobs(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "pod":
		v, err = Client.Clientset.CoreV1().Pods(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "replicaSet":
		v, err = Client.Clientset.AppsV1().ReplicaSets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "replicationController":
		v, err = Client.Clientset.CoreV1().ReplicationControllers(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "statefulSet":
		v, err = Client.Clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "service":
		v, err = Client.Clientset.CoreV1().Services(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "configMap":
		v, err = Client.Clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "pvc":
		v, err = Client.Clientset.CoreV1().PersistentVolumeClaims(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "secret":
		v, err = Client.Clientset.CoreV1().Secrets(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "roleBinding":
		v, err = Client.Clientset.RbacV1().RoleBindings(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "role":
		v, err = Client.Clientset.RbacV1().Roles(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	case "serviceAccount":
		v, err = Client.Clientset.CoreV1().ServiceAccounts(namespace).Get(context.Background(), resourceName, metav1.GetOptions{})
	}
	if err != nil {
		klog.Errorln("获取Resource错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequest, "获取数据错误"))
		return
	}
	c.JSON(200, utils.SuccessResponse(v))
}
