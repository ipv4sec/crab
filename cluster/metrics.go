package cluster

import (
	"crab/utils"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gojek/heimdall/v7/httpclient"
	"github.com/pkg/errors"
	"io/ioutil"
	"k8s.io/klog/v2"
	"time"
)

var (
	HTTPClient = httpclient.NewClient(httpclient.WithHTTPTimeout(time.Second * 30))
)

func GetMetricsHandlerFunc(c *gin.Context) {
	namespace := c.Param("namespace")
	resourceName := c.Param("resourceName")

	cpuMetricsPoints, err := GetDataPointsFromMetricsService(namespace, "cpu", resourceName)
	if err != nil {
		klog.Errorln("获取CPU图表数据错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "获取CPU图表数据错误"))
		return
	}
	memoryMetricsPoints, err := GetDataPointsFromMetricsService(namespace, "mem", resourceName)
	if err != nil {
		klog.Errorln("获取内存图表数据错误", err.Error())
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "获取内存图表数据错误"))
		return
	}
	if len(cpuMetricsPoints) == 0 && len(memoryMetricsPoints)==0 {
		c.JSON(200, utils.ErrorResponse(utils.ErrInternalServer, "获取数据为空"))
		return
	}
	c.JSON(200, utils.SuccessResponse(map[string] DataPoints{
		"cpu": cpuMetricsPoints,
		"mem": memoryMetricsPoints,
	}))
}

type SourceMetrics struct {
	Items []struct{
		MetricPoints []struct {
			Timestamp time.Time `json:"timestamp"`
			Value     int64     `json:"value"`
		} `json:"metricPoints"`
	} `json:"items"`
}

type DataPoint struct {
	X int64 `json:"x"`
	Y int64 `json:"y"`
}
type DataPoints []DataPoint

func GetDataPointsFromMetricsService(namespace, resourceType, name string) (DataPoints, error) {
	endpoint := fmt.Sprintf(
		"http://island-metrics/api/v1/dashboard/namespaces/%s/pod-list/%s/metrics/%s/%v",
		namespace, name, resourceType,time.Now().UnixMicro())
	resp, err := HTTPClient.Get(endpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "请求MetricsService失败")
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "解析MetricsService返回值失败")
	}
	klog.Info("MetricsService返回值:", string(bodyBytes))
	var val SourceMetrics
	err = json.Unmarshal(bodyBytes, &val)
	if err != nil {
		return nil, errors.Wrap(err, "反序列化MetricsService返回值失败")
	}
	var dataPoints DataPoints
	for _, raw := range val.Items[0].MetricPoints {
		converted := DataPoint{
			X: raw.Timestamp.Unix(),
			Y: raw.Value,
		}

		if converted.Y < 0 {
			converted.Y = 0
		}
		dataPoints = append(dataPoints, converted)
	}
	return dataPoints, nil
}