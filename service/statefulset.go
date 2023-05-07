package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var StatefulSet statefulSet

type statefulSet struct{}

type StatusfulSetsResp struct {
	Items []appsv1.StatefulSet `json:"items"`
	Total int                  `json:"total"`
}

// 获取statefulset列表，支持过滤、排序、分页
func (s *statefulSet) GetStatefulSets(filterName, namespace string, limit, page int) (statusfulSetsResp *StatusfulSetsResp, err error) {
	//获取statefulSetList类型的statefulSet列表
	statefulSetList, err := K8s.ClientSet.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet列表失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet列表失败, " + err.Error())
	}
	//将statefulSetList中的StatefulSet列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: s.toCells(statefulSetList.Items),
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

	//将[]DataCell类型的statefulset列表转为v1.statefulset列表
	statefulSets := s.fromCells(data.GenericDataList)

	return &StatusfulSetsResp{
		Items: statefulSets,
		Total: total,
	}, nil
}

// 获取statefulset详情
func (s *statefulSet) GetStatefulSetDetail(statefulSetName, namespace string) (statefulSet *appsv1.StatefulSet, err error) {
	statefulSet, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Get(context.TODO(), statefulSetName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取StatefulSet详情失败, " + err.Error()))
		return nil, errors.New("获取StatefulSet详情失败, " + err.Error())
	}

	return statefulSet, nil
}

// 删除statefulset
func (s *statefulSet) DeleteStatefulSet(statefulSetName, namespace string) (err error) {
	err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Delete(context.TODO(), statefulSetName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除StatefulSet失败, " + err.Error()))
		return errors.New("删除StatefulSet失败, " + err.Error())
	}

	return nil
}

// 更新statefulset
func (s *statefulSet) UpdateStatefulSet(namespace, content string) (err error) {
	var statefulSet = &appsv1.StatefulSet{}

	err = json.Unmarshal([]byte(content), statefulSet)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.AppsV1().StatefulSets(namespace).Update(context.TODO(), statefulSet, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新StatefulSet失败, " + err.Error()))
		return errors.New("更新StatefulSet失败, " + err.Error())
	}
	return nil
}

func (s *statefulSet) toCells(std []appsv1.StatefulSet) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = statefulSetCell(std[i])
	}
	return cells
}

func (s *statefulSet) fromCells(cells []DataCell) []appsv1.StatefulSet {
	statefulSets := make([]appsv1.StatefulSet, len(cells))
	for i := range cells {
		statefulSets[i] = appsv1.StatefulSet(cells[i].(statefulSetCell))
	}

	return statefulSets
}
