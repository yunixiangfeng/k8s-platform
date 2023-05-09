package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化k8s client
	service.K8s.Init() // 可以使用service.K8s.clientset 进行跨包调用

	// 初始化数据库
	db.Init()
	// 初始化gin对象/路由配置
	r := gin.Default()
	// 加载jwt中间件
	//r.Use(middle.JWTAuth())
	// 加载跨域中间件
	r.Use(middle.Cors())
	// 初始化路由规则
	controller.Router.InitApiRouter(r)
	// 启动websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":8081", nil)
	}()
	// gin程序启动
	r.Run(config.ListenAddr)

	// 关闭数据库
	db.Close()
}
