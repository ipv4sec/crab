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

	routers.GET("/app/:id/logs", app.GetPodLogsHandlerFunc)

	routers.GET("/cluster/domain", domain.GetDomainHandlerFunc)
	routers.PUT("/cluster/domain", domain.PutDomainHandlerFunc)

	routers.PUT("/deployment/:id", deployment.PutDeploymentHandlerFunc)

	routers.POST("/trait", trait.PostTraitHandlerFunc)
	routers.GET("/trait", trait.GetTraitsHandlerFunc)
	routers.GET("/trait/:id", trait.GetTraitHandlerFunc)
	routers.PUT("/trait/:id", trait.PutTraitHandlerFunc)
	routers.DELETE("/trait/:id", trait.DeleteTraitHandlerFunc)

	routers.POST("/workload/type", workloadType.PostTypeHandlerFunc)
	routers.GET("/workload/type", workloadType.GetTypesHandlerFunc)
	routers.GET("/workload/type/:id", workloadType.GetTypeHandlerFunc)
	routers.PUT("/workload/type/:id", workloadType.PutTypeHandlerFunc)
	routers.DELETE("/workload/type/:id", workloadType.DeleteTypeHandlerFunc)

	routers.POST("/workload/vendor", workloadVendor.PostVendorHandlerFunc)
	routers.GET("/workload/vendor", workloadVendor.GetVendorsHandlerFunc)
	routers.GET("/workload/vendor/:id", workloadVendor.GetVendorHandlerFunc)
	routers.PUT("/workload/vendor/:id", workloadVendor.PutVendorHandlerFunc)
	routers.DELETE("/workload/vendor/:id", workloadVendor.DeleteVendorHandlerFunc)

	routers.POST("/tool/spelling", tool.PostSpellingHandlerFunc)
	routers.POST("/tool/convertion", tool.PostConvertionHandlerFunc)
	routers.GET("/tool/systemTemplate", tool.GetTemplateHandlerFunc)

	routers.GET("/resource/:namespace/:resourceType/:resourceName", cluster.GetResourceHandlerFunc)

	err = routers.Run("0.0.0.0:3000")
	if err != nil {
		panic(fmt.Errorf("监听端口失败: %w", err))
	}
}