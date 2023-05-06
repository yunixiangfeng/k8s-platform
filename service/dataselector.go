package service

import (
	"sort"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// dataselector用于排序，过滤，分页的数据类型
type dataSelector struct {
	GenericDataList []DataCell
	DataSelect      *DataSelectQuery
}

// DataCell接口，用于各种资源List的类型转换，转换后可以使用dataselector的排序，过滤，分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}

// DataSelectQuery 定义过滤和分页的结构体，过滤:Name 分页:Limit和Page
type DataSelectQuery struct {
	Filter   *FilterQuery
	Paginate *PaginateQuery
}

// FilterQuery用于查询 过滤:Name
type FilterQuery struct {
	Name string
}

// 分页:Limit和Page Limit是单页的数据条数，Page是第几页
type PaginateQuery struct {
	Page  int
	Limit int
}

// 实现自定义结构的排序，需要重写Len、Swap、Less方法
// Len方法用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.GenericDataList)
}

// Swap用于数据比较大小后的位置变更
func (d *dataSelector) Swap(i, j int) {
	d.GenericDataList[i], d.GenericDataList[j] = d.GenericDataList[j], d.GenericDataList[i]
}

// Less用于比较大小
func (d *dataSelector) Less(i, j int) bool {
	return d.GenericDataList[i].GetCreation().Before(d.GenericDataList[j].GetCreation())
}

// 重写以上三个方法,用sort.Sort 方法触发排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(d)
	return d
}

// Filter方法用于过滤,比较数据Name属性,若包含则返回
func (d *dataSelector) Filter() *dataSelector {
	if d.DataSelect.Filter.Name == "" {
		return d
	}

	filtered := []DataCell{}
	for _, value := range d.GenericDataList {
		// 定义是否匹配的标签变量，默认是匹配的
		matches := true
		objName := value.GetName()
		if !strings.Contains(objName, d.DataSelect.Filter.Name) {
			matches = false
			continue
		}
		if matches {
			filtered = append(filtered, value)
		}
	}
	d.GenericDataList = filtered
	return d
}

// Paginate分页，根据Limit和Page的传参，取一定范围内的数据返回
func (d *dataSelector) Paginate() *dataSelector {
	limit := d.DataSelect.Paginate.Limit
	page := d.DataSelect.Paginate.Page
	// 验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	// 举例：25个元素的数组，limit是10, page是3，startIndex是20，endIndex是29（实际上endIndex是24）
	startIndex := limit * (page - 1)
	endIndex := limit*page - 1

	// 处理最后一页，这时候就把endIndex由30改为25了
	if len(d.GenericDataList) < endIndex {
		endIndex = len(d.GenericDataList)
	}
	d.GenericDataList = d.GenericDataList[startIndex:endIndex]
	return d
}

// 定义podCell, 重写GetCreation和GetName方法后，可以进行数据转换
// covev1.Pod --> podCell  --> DataCell
// appsv1.Deployment --> deployCell --> DataCell
type podCell corev1.Pod

// 重写DataCell接口的两个方法
func (p podCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p podCell) GetName() string {
	return p.Name
}
