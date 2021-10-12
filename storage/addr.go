package storage

import (
	"context"
	"crab/cluster"
	"crab/utils"
	"github.com/gin-gonic/gin"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"strings"
)

type Data struct{
	Name string `json:"name"`
	Addrs []string `json:"addrs"`
}

func GetAddrsHandlerFunc(c *gin.Context)  {
	conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("获取ROOK的键值对失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取地址信息失败"))
		return
	}
	var addrs []Data
	for i:=0; i< len(conf.Items); i++ {
		if strings.Contains(conf.Items[i].Name, "island-addrs-") {
			var cm AddrData
			err := json.Unmarshal([]byte(conf.Items[i].Data["value"]), &cm)
			if err != nil {
				klog.Errorln()
				continue
			}
			addrs = append(addrs, Data{
				Name: strings.TrimPrefix(conf.Items[i].Name, "island-addrs-"),
				Addrs: cm.Value,
			})
		}
	}
	c.JSON(200, utils.SuccessResponse(addrs))
}
