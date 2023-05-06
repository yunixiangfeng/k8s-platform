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
func (r *router) InitApiRouter(router *gin.Engine) {
	router.GET("/api/k8s/pods", Pod.GetPods)
}
