package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化k8s client
	service.K8s.Init() // 可以使用service.K8s.clientset 进行跨包调用

	// 初始化gin对象
	r := gin.Default()
	// 初始化路由规则
	controller.Router.InitApiRouter(r)
	// gin程序启动
	r.Run(config.ListenAddr)
}
