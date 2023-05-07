package controller

import (
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wonderivan/logger"
)

var Workflow workflow

type workflow struct{}

// 获取列表分页查询
func (w *workflow) GetList(ctx *gin.Context) {
	params := new(struct {
		Name      string `form:"name"`
		Namespace string `form:"namespace"`
		Page      int    `form:"page"`
		Limit     int    `form:"limit"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Workflow.GetList(params.Name, params.Namespace, params.Limit, params.Page)
	if err != nil {
		logger.Error("获取Workflow列表失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "获取Workflow列表成功",
		"data": data,
	})
}

// 查询workflow单条数据
func (w *workflow) GetById(ctx *gin.Context) {
	params := new(struct {
		ID int `form:"id"`
	})
	if err := ctx.Bind(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	data, err := service.Workflow.GetById(params.ID)
	if err != nil {
		logger.Error("查询Workflow单条数据失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "查询Workflow单条数据成功",
		"data": data,
	})
}

// 创建workflow
func (w *workflow) Create(ctx *gin.Context) {
	var (
		wc  = &service.WorkflowCreate{}
		err error
	)

	if err = ctx.ShouldBindJSON(wc); err != nil {
		logger.Error("Bind请求参数dc失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}
	if err = service.Workflow.CreateWorkFlow(wc); err != nil {
		logger.Error("创建Workflow失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "创建Workflow成功",
		"data": nil,
	})

}

// 删除workflow
func (w *workflow) DelById(ctx *gin.Context) {
	params := new(struct {
		ID int `json:"id"`
	})
	if err := ctx.ShouldBindJSON(params); err != nil {
		logger.Error("Bind请求参数失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	if err := service.Workflow.DelById(params.ID); err != nil {
		logger.Error("删除Workflow失败, " + err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg":  err.Error(),
			"data": nil,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg":  "删除Workflow成功",
		"data": nil,
	})
}
