package dao

import (
	"errors"
	"k8s-platform/db"
	"k8s-platform/model"

	"github.com/wonderivan/logger"
)

var Workflow workflow

type workflow struct{}

//定义列表的返回内容，Items是workflow元素列表，Total为workflow元素数量
type WorkflowResp struct {
	Items []*model.Workflow `json:"items"`
	Total int               `json:"total"`
}

// 获取workflow列表
func (w *workflow) GetWorkflows(filterName, namespace string, limit, page int) (data *WorkflowResp, err error) {
	//定义分页的起始位置
	startSet := (page - 1) * limit
	//定义数据库查询返回的内容
	var (
		workflowList []*model.Workflow
		total        int
	)
	//数据库查询，Limit方法用于限制条数，Offset方法用于设置起始位置
	tx := db.GORM.
		Model(&model.Workflow{}).
		Where("name like ?", "%"+filterName+"%").
		Count(&total).
		Limit(limit).
		Offset(startSet).
		Order("id desc").
		Find(&workflowList)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow列表失败, " + tx.Error.Error())
		return nil, errors.New("获取Workflow列表失败, " + tx.Error.Error())
	}
	return &WorkflowResp{
		Items: workflowList,
		Total: total,
	}, nil
}

// 获取详情
func (w *workflow) GetById(id int) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	tx := db.GORM.Where("id = ?", id).First(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow详情失败, " + tx.Error.Error())
		return nil, errors.New("获取Workflow详情失败, " + tx.Error.Error())
	}
	return workflow, nil
}

// 创建
func (w *workflow) Add(workflow *model.Workflow) (err error) {
	tx := db.GORM.Create(&workflow)
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("创建Workflow失败, " + tx.Error.Error())
		return errors.New("创建Workflow失败, " + tx.Error.Error())
	}
	return nil
}

// 删除
func (w *workflow) DelById(id int) (err error) {
	tx := db.GORM.Where("id = ?", id).Delete(&model.Workflow{})
	if tx.Error != nil && tx.Error.Error() != "record not found" {
		logger.Error("获取Workflow详情失败, " + tx.Error.Error())
		return errors.New("获取Workflow详情失败, " + tx.Error.Error())
	}
	return nil
}
