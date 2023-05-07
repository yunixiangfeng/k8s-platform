package service

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/wonderivan/logger"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Secret secret

type secret struct{}

type SecretsResp struct {
	Items []corev1.Secret `json:"items"`
	Total int             `json:"total"`
}

// 获取secret列表，支持过滤、排序、分页
func (s *secret) GetSecrets(filterName, namespace string, limit, page int) (secretsResp *SecretsResp, err error) {
	//获取secretList类型的secret列表
	secretList, err := K8s.ClientSet.CoreV1().Secrets(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error(errors.New("获取Secret列表失败, " + err.Error()))
		return nil, errors.New("获取Secret列表失败, " + err.Error())
	}
	//将secretList中的secret列表(Items)，放进dataselector对象中，进行排序
	selectableData := &dataSelector{
		GenericDataList: s.toCells(secretList.Items),
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

	//将[]DataCell类型的secret列表转为v1.secret列表
	secrets := s.fromCells(data.GenericDataList)

	return &SecretsResp{
		Items: secrets,
		Total: total,
	}, nil
}

// 获取secret详情
func (s *secret) GetSecretDetail(secretName, namespace string) (secret *corev1.Secret, err error) {
	secret, err = K8s.ClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err != nil {
		logger.Error(errors.New("获取Secret详情失败, " + err.Error()))
		return nil, errors.New("获取Secret详情失败, " + err.Error())
	}

	return secret, nil
}

// 删除secret
func (s *secret) DeleteSecret(secretName, namespace string) (err error) {
	err = K8s.ClientSet.CoreV1().Secrets(namespace).Delete(context.TODO(), secretName, metav1.DeleteOptions{})
	if err != nil {
		logger.Error(errors.New("删除Secret失败, " + err.Error()))
		return errors.New("删除Secret失败, " + err.Error())
	}

	return nil
}

// 更新secret
func (s *secret) UpdateSecret(namespace, content string) (err error) {
	var secret = &corev1.Secret{}

	err = json.Unmarshal([]byte(content), secret)
	if err != nil {
		logger.Error(errors.New("反序列化失败, " + err.Error()))
		return errors.New("反序列化失败, " + err.Error())
	}

	_, err = K8s.ClientSet.CoreV1().Secrets(namespace).Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		logger.Error(errors.New("更新Secret失败, " + err.Error()))
		return errors.New("更新Secret失败, " + err.Error())
	}
	return nil
}

func (s *secret) toCells(std []corev1.Secret) []DataCell {
	cells := make([]DataCell, len(std))
	for i := range std {
		cells[i] = secretCell(std[i])
	}
	return cells
}

func (s *secret) fromCells(cells []DataCell) []corev1.Secret {
	secrets := make([]corev1.Secret, len(cells))
	for i := range cells {
		secrets[i] = corev1.Secret(cells[i].(secretCell))
	}

	return secrets
}
