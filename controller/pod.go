package controller

import (
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Pod pod

type pod struct{}

// Controller中的方法入参是gin.Context，用于从上下文中获取请求参数及定义响应内容
// 流程:绑定参数->调用service代码->根据调用结果响应具体内容
// 获取Pod列表，支持分页、过滤、排序
func (p *pod) GetPods(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		FilterName string `form:"filter_name"`
		Namespace  string `form:"namespace"`
		Limit      int    `form:"limit"`
		Page       int    `form:"page"`
	})
	// form格式使用Bind方法，json格式使用ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPods(params.FilterName, params.Namespace, params.Limit, params.Page)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}

// 获取pod详情
func (p *pod) GetPodDetail(ctx *gin.Context) {
	// 处理入参
	// 匿名结构体，用于定义入参，get请求为form格式，其他请求为json格式
	params := new(struct {
		PodName   string `form:"pod_name"`
		Namespace string `form:"namespace"`
	})
	// form格式使用Bind方法，json格式使用ShouldBindJSON方法
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind绑定参数失败" + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  "Bind绑定参数失败" + err.Error(),
			"data": nil,
		})
		return
	}
	data, err := service.Pod.GetPodDetail(params.PodName, params.Namespace)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Pod列表成功",
		"data": data,
	})
}
