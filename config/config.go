package config

import "time"

const (
	ListenAddr = "0.0.0.0:9090"
	KubeConfig = "C:\\Users\\Administrator\\.kube\\config"
	// tail的日志行数
	// tail -n 2000
	PodLogTailLine = 2000

	// DB Config
	DbType = "mysql"
	DbHost = "192.168.204.130"
	DbPort = 3306
	DbName = "k8s_dashboard"
	DbUser = "root"
	DbPass = "1234"
	// 打印mysql debug的sql日志
	LogMode = false
	// 连接池配置
	MaxIdleConns = 10               // 最大空闲连接
	MaxOpenConns = 100              // 最大连接数
	MaxLifeTime  = 30 * time.Second // 会话时间
)
