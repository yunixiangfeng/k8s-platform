package service

import (
	"context"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Pod pod

type pod struct{}

// 定义列表的返回内容，Items是pod元素列表，Total是元素数量
type PodsResp struct {
	Total int          `json:"total"`
	Items []corev1.Pod `json:"items`
}

// 获取pod列表，支持过滤、排序、分页
func (p *pod) GetPods(filterName, namespace string, limit, page int) (podsResp *PodsResp, err error) {
	//context.TODO()  用于声明一个空的context上下文,用于List方法内设置这个请求超时
	//metav1.ListOptions{} 用于过滤List数据,如label,field等
	podList, err := K8s.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Info("获取Pod列表失败，" + err.Error())
		// 返回给上一层,最终返回给前端,前端捕获到后打印出来
		return nil, errors.New("获取Pod列表失败，" + err.Error())
	}
	// 实例化dataselector结构体，组装数据
	selectableData := &dataSelector{
		GenericDataList: p.toCells(podList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}
	// 先过滤
	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	// 排序和分页
	data := filtered.Sort().Paginate()

	// 将DataCell类型转成Pod
	pods := p.fromCells(data.GenericDataList)
	return &PodsResp{
		Total: total,
		Items: pods,
	}, nil
}

// 类型转换方法corev1.Pod --> DataCell,DataCell-->corev1.Pod
func (p *pod) toCells(pods []corev1.Pod) []DataCell {
	cells := make([]DataCell, len(pods))
	for i := range pods {
		cells[i] = podCell(pods[i])
	}
	return cells
}

func (p *pod) fromCells(cells []DataCell) []corev1.Pod {
	pods := make([]corev1.Pod, len(cells))
	for i := range cells {
		// 是将DataCell类型转成podCell
		pods[i] = corev1.Pod(cells[i].(podCell))
	}
	return pods
}
