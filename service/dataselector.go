package service

import "time"

// dataselector用于排序，过滤，分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelect *DataSelectQuery
}

// DataCell接口，用于各种资源List的类型转换，转换后可以使用dataselector的排序，过滤，分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体，过滤:Name 分页:Limit和Page
type DataSelectQuery struct {
	Filter *FilterQuery
	Paginate *PaginateQuery
}

// FilterQuery用于查询 过滤:Name
type FilterQuery struct {
	Name string
}

// 分页:Limit和Page Limit是单页的数据条数，Page是第几页
type PaginateQuery struct {
	Page int
	Limit int
}

// 实现自定义结构的排序，需要重写Len、Swap、Less方法
// Len方法用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

/