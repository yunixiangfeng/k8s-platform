package service

import (
	"k8s-platform/config"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// 用于初始化k8s client
var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

// 初始化k8s
func (k *k8s) Init() {
	// 将kubeconfig格式化为rest.config类型
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config, err := clientcmd.BuildConfigFromFlags("", config.KubeConfig)
	if err != nil {
		panic("获取K8s配置失败:" + err.Error())
	} else {
		logger.Info("获取K8s配置成功！")
	}
	// 通过config创建clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic("创建K8s client失败:" + err.Error())
	} else {
		logger.Info("创建K8s client 成功!")
	}
	k.ClientSet = clientset
}
