package controller

import (
	"github.com/gin-gonic/gin"
)

// // 初始化router类型对象，首字母大写，用于跨包调用
// var Router router

// // 声明一个router的结构体
// type router struct{}

// func (r *router) InitApiRouter(router *gin.Engine) {
// 	router.GET("/", Index)
// }

// func Index(ctx *gin.Context) {
// 	ctx.JSON(200, gin.H{
// 		"code": 200,
// 		"msg":  "In index",
// 	})
// }

// 实例化router结构体，可使用该对象点出首字母大写的方法（包外调用）
var Router router

// 创建router的结构体
type router struct{}

// // 初始化路由规则，创建测试api接口
// func (r *router) InitApiRouter(router *gin.Engine) {
// 	router.GET("/testapi", func(ctx *gin.Context) {
// 		ctx.JSON(http.StatusOK, gin.H{
// 			"msg":  "testapi success!",
// 			"data": nil,
// 		})
// 	})
// }
// 初始化路由规则
// func (r *router) InitApiRouter(router *gin.Engine) {
// 	router.
// 		GET("/api/k8s/pods", Pod.GetPods).
// 		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
// 		POST("/api/k8s/pods", Pod.DeletePod).
func (r *router) InitApiRouter(router *gin.Engine) {
	router.
		// Pods
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNp).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment).
		// workflows
		GET("/api/k8s/workflows", Workflow.GetList).
		GET("/api/k8s/workflow/detail", Workflow.GetById).
		POST("/api/k8s/workflow/create", Workflow.Create).
		DELETE("/api/k8s/workflow/del", Workflow.DelById)

}
