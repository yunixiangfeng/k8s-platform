package model

import "time"

// 定义结构体，属性与mysql表字段对齐
type Workflow struct {
	// gorm:"primarykey"用于声明主键
	ID       uint       `json:"id" gorm:"primaryKey"`
	CreateAt *time.Time `json:"created_at"`
	UpdateAt *time.Time `json:"update_at"`
	DeleteAt *time.Time `json:"deleted_at"`

	Name       string `json:"name"`
	Namespace  string `json:"namespace"`
	Replicas   int32  `json:"replicas"`
	Deployment string `json:"deployment"`
	Service    string `json:"service"`
	Ingress    string `json:"ingress"`
	// gorm:"column:type"用于声明mysql中表的字段名
	Type string `json:"type" gorm:"column:type"`
}

// 定义TableName方法，返回mysql表名，以次定义mysql中的表名
func (*Workflow) TableName() string {
	return "workflow"
}
