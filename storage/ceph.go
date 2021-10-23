package storage

import (
	"bytes"
	"context"
	"crab/cluster"
	"crab/provider"
	"crab/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/klog/v2"
	"strings"
	"time"
)

type Device struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

type Volume struct {
	Name        string `json:"name"`
	Size        string `json:"size"`
	Hostname    string `json:"hostname"`
	HasChildren bool   `json:"hasChildren"`
}

type Node struct {
	Name    string `json:"name"`
	Devices []struct {
		Name string `json:"name"`
	} `json:"devices"`
}

type Param struct {
	Hostname string `json:"hostname"`
	Volume   string `json:"volume"`
}



func fetchStorageStatus() int {
	gvr := schema.GroupVersionResource{
		Group:    "ceph.rook.io",
		Version:  "v1",
		Resource: "cephclusters",
	}
	rook, err := cluster.Client.DynamicClient.Resource(gvr).Namespace("rook-ceph").
		Get(context.Background(), "rook-ceph", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取ROOK资源错误", err.Error())
		return 0
	}
	value, found, err := unstructured.NestedString(rook.Object, "status", "phase")
	if err != nil {
		klog.Errorln("解析ROOK资源错误", err.Error())
		return 0
	}
	if !found {
		klog.Errorln("未找到ROOK资源")
		return 0
	}
	if value == "Ready" {
		return 2
	}
	return 1
}
func GetStorageHandlerFunc(c *gin.Context) {
	conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		klog.Errorln("获取ROOK的键值对失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "获取磁盘信息失败"))
		return
	}

	volumes := []Volume{}
	for i := 0; i < len(conf.Items); i++ {
		if strings.Contains(conf.Items[i].Name, "island-devices-") {
			var cm DiskData
			err := json.Unmarshal([]byte(conf.Items[i].Data["value"]), &cm)
			if err != nil {
				klog.Errorln("解析键值对失败", err.Error())
				continue
			}
			for j := 0; j < len(cm.Value); j++ {
				volumes = append(volumes, Volume{
					Name:        cm.Value[j].Name,
					Size:        fmt.Sprintf("%vG", cm.Value[j].Size/1024/1024/1024),
					Hostname:    strings.TrimPrefix(conf.Items[i].Name, "island-devices-"),
					HasChildren: cm.Value[j].HasChildren,
				})
			}
		}
	}
	deploy, err := cluster.Client.Clientset.AppsV1().Deployments("rook-ceph").
		Get(context.Background(), "rook-ceph-tools", metav1.GetOptions{})
	if err != nil {
		klog.Errorln("获取存储状态失败", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetDeployment, "获取存储状态失败"))
		return
	}
	status := struct {
		Status int    `json:"status"`
		Message string `json:"message"`
	}{}
	if deploy.Status.AvailableReplicas == 0 {
		status.Status = 0
		status.Message = "存储集群还未初始化"
	}
	if deploy.Status.AvailableReplicas > 0 {
		s := fetchStorageStatus()
		if s == 0 {
			status.Status = 0
			status.Message = "存储集群还未初始化"
		}
		if s == 1 {
			status.Status = 1
			status.Message = "存储集群正在构建中"
		}
		if s == 2 {
			status.Status = 2
			status.Message = "存储集群平稳运行中"
		}
	}
	c.JSON(200, utils.SuccessResponse(struct {
		Store struct {
			Status int    `json:"status"`
			Message string `json:"message"`
		} `json:"store"`
		Volumes []Volume `json:"volumes"`
	}{
		Store:   status,
		Volumes: volumes,
	}))
}

func PostStorageHandlerFunc(c *gin.Context) {
	var param struct {
		Default bool     `json:"default"`
		Volumes []Volume `json:"volumes"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	a, err := ioutil.ReadFile("assets/ceph/cluster.yaml")
	if err != nil {
		klog.Errorln("读取文件错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "服务器内部错误"))
		return
	}
	err = cluster.Client.Apply(context.Background(), a)
	if err != nil {
		klog.Errorln("初始化存储错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrClusterClientApply, "服务器内部错误"))
		return
	}
	for {
		if 0 != fetchStorageStatus() {
			break
		}
		time.Sleep(time.Second)
	}
	if param.Default {
		gvr := schema.GroupVersionResource{
			Group:    "ceph.rook.io",
			Version:  "v1",
			Resource: "cephclusters", // oc api-resources
		}
		p := []byte(`{"spec": {"storage": {"useAllNodes": true, "useAllDevices": true}}}`)
		_, err = cluster.Client.DynamicClient.Resource(gvr).Namespace("rook-ceph").
			Patch(context.Background(), "rook-ceph", types.MergePatchType, p, metav1.PatchOptions{})
		if err != nil {
			klog.Errorln("使用全部存储错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrClusterClientPatch, "服务器内部错误"))
			return
		}
	} else {
		gvr := schema.GroupVersionResource{
			Group:    "ceph.rook.io",
			Version:  "v1",
			Resource: "cephclusters", // oc api-resources
		}
		m := map[string][]struct{ Name string }{}
		for i := 0; i < len(param.Volumes); i++ {
			if _, ok := m[param.Volumes[i].Hostname]; !ok {
				m[param.Volumes[i].Hostname] = []struct{ Name string }{
					{Name: param.Volumes[i].Name},
				}
				continue
			}
			m[param.Volumes[i].Hostname] = append(m[param.Volumes[i].Hostname], struct{ Name string }{Name: param.Volumes[i].Name})
			conf, err := cluster.Client.Clientset.CoreV1().ConfigMaps("island-system").
				Get(context.Background(), fmt.Sprintf("island-addrs-%s", param.Volumes[i].Hostname), metav1.GetOptions{})
			if err != nil {
				klog.Errorln("通过集群获取某节点地址失败", err.Error())
				c.JSON(200, utils.ErrorResponse(utils.ErrClusterGetConfigMap, "通过集群获取某节点地址失败"))
				return
			}
			// cm.Data["value"]
			var cm AddrData
			err = json.Unmarshal([]byte(conf.Data["value"]), &cm)
			if err != nil {
				klog.Errorln("解析键值对失败", err.Error())
				continue
			}
			for j := 0; j < len(cm.Value); j++ {
				p := &Param{
					Hostname: param.Volumes[i].Hostname,
					Volume:   param.Volumes[i].Name,
				}
				s , err := json.Marshal(p)
				if err != nil {
					klog.Errorln("序列化参数错误", err.Error())
					continue
				}
				res, err := provider.HTTPClient.Post(fmt.Sprintf("http://%s:3000/disk", cm.Value[j]),
					bytes.NewBuffer(s), nil)
				if err != nil {
					klog.Errorln("发送格式化请求错误", err.Error())
					continue
				}
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					klog.Errorln("读取格式化返回结果错误", err.Error())
					continue
				}
				var reply utils.Reply
				err = json.Unmarshal(body, &reply)
				if err != nil {
					klog.Errorln("反序列化格式化返回结果错误", err.Error())
					continue
				}
				if reply.Code != 0 {
					klog.Errorln("接口返回错误", reply.Result)
					continue
				}
				klog.Errorln("格式化磁盘成功", p.Hostname, p.Volume)
			}
		}
		var nodes []Node
		for k, v := range m {
			nodes = append(nodes, Node{
				Name: k,
				Devices: []struct {
					Name string `json:"name"`
				}(v),
			})
		}
		nodeBytes, err := json.Marshal(nodes)
		if err != nil {
			klog.Errorln("序列化错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrJSONMarshal, "服务器内部错误"))
			return
		}
		p := fmt.Sprintf(`{"spec":{"storage":{"useAllNodes": false,"useAllDevices": false, "nodes":%s}}}`,
			string(nodeBytes))
		klog.Errorln("修改内容为:", p)
		_, err = cluster.Client.DynamicClient.Resource(gvr).Namespace("rook-ceph").
			Patch(context.Background(), "rook-ceph", types.MergePatchType, []byte(p), metav1.PatchOptions{})
		if err != nil {
			klog.Errorln("使用指定存储错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrClusterClientPatch, "服务器内部错误"))
			return
		}
	}
	c.JSON(200, utils.SuccessResponse("初始化存储成功, 等待内部操作完成"))
}

func PutStorageHandlerFunc(c *gin.Context) {
	var param struct {
		Default bool     `json:"default"`
		Volumes []Volume `json:"volumes"`
	}
	err := c.ShouldBindJSON(&param)
	if err != nil {
		c.JSON(200, utils.ErrorResponse(utils.ErrBadRequestParam, "参数错误"))
		return
	}
	for {
		if 0 != fetchStorageStatus() {
			break
		}
		time.Sleep(time.Second)
	}
	if param.Default {
		gvr := schema.GroupVersionResource{
			Group:    "ceph.rook.io",
			Version:  "v1",
			Resource: "cephclusters", // oc api-resources
		}
		p := []byte(`{"spec": {"storage": {"useAllNodes": true, "useAllDevices": true}}}`)
		_, err = cluster.Client.DynamicClient.Resource(gvr).Namespace("rook-ceph").
			Patch(context.Background(), "rook-ceph", types.MergePatchType, p, metav1.PatchOptions{})
		if err != nil {
			klog.Errorln("使用全部存储错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrClusterClientPatch, "服务器内部错误"))
			return
		}
	} else {
		gvr := schema.GroupVersionResource{
			Group:    "ceph.rook.io",
			Version:  "v1",
			Resource: "cephclusters", // oc api-resources
		}
		m := map[string][]struct{ Name string }{}
		for i := 0; i < len(param.Volumes); i++ {
			if _, ok := m[param.Volumes[i].Hostname]; !ok {
				m[param.Volumes[i].Hostname] = []struct{ Name string }{
					{Name: param.Volumes[i].Name},
				}
				continue
			}
			m[param.Volumes[i].Hostname] = append(m[param.Volumes[i].Hostname], struct{ Name string }{Name: param.Volumes[i].Name})
		}
		var nodes []Node
		for k, v := range m {
			nodes = append(nodes, Node{
				Name: k,
				Devices: []struct {
					Name string `json:"name"`
				}(v),
			})
		}
		bytes, err := json.Marshal(nodes)
		if err != nil {
			klog.Errorln("序列Git化错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrJSONMarshal, "服务器内部错误"))
			return
		}
		p := fmt.Sprintf(`{"spec":{"storage":{"useAllNodes": false,"useAllDevices": false, "nodes":%s}}}`, string(bytes))
		klog.Errorln("修改内容为:", p)
		_, err = cluster.Client.DynamicClient.Resource(gvr).Namespace("rook-ceph").
			Patch(context.Background(), "rook-ceph", types.MergePatchType, []byte(p), metav1.PatchOptions{})
		if err != nil {
			klog.Errorln("使用指定存储错误", err.Error())
			c.JSON(200, utils.ErrorResponse(utils.ErrClusterClientPatch, "服务器内部错误"))
			return
		}
	}
	c.JSON(200, utils.SuccessResponse("修改存储成功, 等待内部操作完成"))
}