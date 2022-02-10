package main

import (
	"crab/app"
	"crab/cluster"
	"crab/config"
	"crab/db"
	"crab/deployment"
	"crab/domain"
	"crab/tool"
	"crab/trait"
	"crab/user"
	"crab/workloadType"
	"crab/workloadVendor"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
)

func main() {
	var conf string
	flag.StringVar(&conf,"config", "config.yaml", "配置文件")
	flag.Parse()

	klog.Infoln("读取配置文件")
	bytes, err := ioutil.ReadFile(conf)
	if err != nil {
		panic(err)
	}

	klog.Infoln("解析配置文件")
	var cfg config.Config
	err = yaml.Unmarshal(bytes, &cfg)
	if err != nil {
		panic(err)
	}

	klog.Infoln("连接到数据库")
	err = db.Init(&cfg.Mysql)
	if err != nil {
		panic(err)
	}

	klog.Infoln("开始集群认证")
	err = cluster.Init()
	if err != nil {
		if os.Getenv("CRAB_DEBUG") == "" {
			panic(fmt.Errorf("获取集群认证失败: %w", err))
		}
	}

	klog.Infoln("开始提供服务")
	gin.SetMode(gin.ReleaseMode)
	routers := gin.Default()

	routers.GET("/user/:username", user.GetUserHandlerFunc)
	routers.PUT("/user/:username", user.PutUserHandlerFunc)

	routers.POST("/app", app.PostAppHandlerFunc)
	routers.GET("/app", app.GetAppsHandlerFunc)
	routers.GET("/app/:id", app.GetAppHandlerFunc)
	routers.PUT("/app/:id", app.PutAppHandlerFunc)
	routers.DELETE("/app/:id", app.DeleteAppHandlerFunc)

	routers.GET("/app/:id/:pod/logs", app.GetPodLogsHandlerFunc)

	routers.GET("/cluster/domain", domain.GetDomainHandlerFunc)
	routers.PUT("/cluster/domain", domain.PutDomainHandlerFunc)

	routers.PUT("/deployment/:id", deployment.PutDeploymentHandlerFunc)

	routers.POST("/trait", trait.PostTraitHandlerFunc)
	routers.GET("/trait", trait.GetTraitsHandlerFunc)
	routers.GET("/trait/:id", trait.GetTraitHandlerFunc)
	routers.PUT("/trait/:id", trait.PutTraitHandlerFunc)
	routers.DELETE("/trait/:id", trait.DeleteTraitHandlerFunc)

	routers.POST("/workloadType", workloadType.PostTypeHandlerFunc)
	routers.GET("/workloadType", workloadType.GetTypesHandlerFunc)
	routers.GET("/workloadType/:id", workloadType.GetTypeHandlerFunc)
	routers.PUT("/workloadType/:id", workloadType.PutTypeHandlerFunc)
	routers.DELETE("/workloadType/:id", workloadType.DeleteTypeHandlerFunc)

	routers.POST("/workloadVendor", workloadVendor.PostVendorHandlerFunc)
	routers.GET("/workloadVendor", workloadVendor.GetVendorsHandlerFunc)
	routers.GET("/workloadVendor/:id", workloadVendor.GetVendorHandlerFunc)
	routers.PUT("/workloadVendor/:id", workloadVendor.PutVendorHandlerFunc)
	routers.DELETE("/workloadVendor/:id", workloadVendor.DeleteVendorHandlerFunc)

	routers.POST("/tool/spelling", tool.PostSpellingHandlerFunc)
	routers.POST("/tool/convertion", tool.PostConvertionHandlerFunc)
	routers.GET("/tool/systemTemplate", tool.GetTemplateHandlerFunc)

	routers.GET("/resource/:namespace/:resourceType/:resourceName", cluster.GetResourceHandlerFunc)
	routers.GET("/metrics/:namespace/:resourceName", cluster.GetMetricsHandlerFunc)
	routers.GET("/plugin/:resourceName", cluster.GetConfigMapResourceHandlerFunc)

	err = routers.Run("0.0.0.0:3000")
	if err != nil {
		panic(fmt.Errorf("监听端口失败: %w", err))
	}
}