package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var DaemonSet daemonSet

type daemonSet struct{}

type DaemonSetsResp struct {
	Items []appsv1.DaemonSet `json:"items"`
	Total int                `json:"total"`
}

// 获取daemonset列表，支持过滤、排序、分页
func (d *daemonSet) GetDaemonSets(filterName, namespace string, limit, page int) (daemonSetsResp *DaemonSetsResp, err error) {
	//获取daemonSetList类型的daemonSet列表
	daemonSetList, err := K8s.ClientSet.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取DaemonSet列表失败, " + err.Error()))
		return nil, errors.New("获取DaemonSet列表失败, " + err.Error())
	}
	//将daemonSetList中的daemonSet列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: d.toCells(daemonSetList.Items),
		DataSelect: &DataSelectQuery{
			Filter: &FilterQuery{Name: filterName},
			Paginate: &PaginateQuery{
				Limit: limit,
				Page:  page,
			},
		},
	}

	filtered := selectableData.Filter()
	total := len(filtered.GenericDataList)
	data := filtered.Sort().Paginate()

	//将[]DataCell类型的daemonset列表转为v1.daemonset列表
	daemonSets := d.fromCells(data.GenericDataList)

	return &DaemonSetsResp{
		Items: daemonSets,
		Total: total,
	}, nil
}

// 获取daemonset详情
func (d *daemonSet) GetDaemonSetDetail(daemonSetName, namespace string) (daemonSet *appsv1.DaemonSet, err error) {
	daemonSet, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Get(context.TODO(), daemonSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取DaemonSet详情失败, " + err.Error()))
		return nil, errors.New("获取DaemonSet详情失败, " + err.Error())
	}

	return daemonSet, nil
}

// 删除daemonset
func (d *daemonSet) DeleteDaemonSet(daemonSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Delete(context.TODO(), daemonSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除DaemonSet失败, " + err.Error()))
		return errors.New("删除DaemonSet失败, " + err.Error())
	}
	return nil
}

// 更新daemonset
func (d *daemonSet) UpdateDaemonSet(namespace, content string) (err error) {
	var daemonSet = &appsv1.DaemonSet{}

	err = json.Unmarshal([]byte(content), daemonSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().DaemonSets(namespace).Update(context.TODO(), daemonSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新DaemonSet失败, " + err.Error()))
		return errors.New("更新DaemonSet失败, " + err.Error())
	}
	return nil
}

func (d *daemonSet) toCells(std []appsv1.DaemonSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = daemonSetCell(std[i])
	}
	return cells
}

func (d *daemonSet) fromCells(cells []DataCell) []appsv1.DaemonSet {
	daemonSets := make([]appsv1.DaemonSet, len(cells))
	for i := range cells {
		daemonSets[i] = appsv1.DaemonSet(cells[i].(daemonSetCell))
	}

	return daemonSets
}
