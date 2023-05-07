package db

import (
	"fmt"
	"k8s-platform/config"
	"time"

	"github.com/wonderivan/logger"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	isInit bool
	GORM   *gorm.DB
	err    error
)

// DB的初始化函数，与数据库建立连接
func Init() {
	// 判断是否已经初始化
	if isInit {
		return
	}
	// 组装连接配置
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.DbUser,
		config.DbPass,
		config.DbHost,
		config.DbPort,
		config.DbName)
	GORM, err := gorm.Open(config.DbType, dsn)
	if err != nil {
		panic("数据库连接失败," + err.Error())
	}
	// 打印sql语句
	GORM.LogMode(config.LogMode)
	// 开启连接池
	GORM.DB().SetMaxIdleConns(config.MaxIdleConns)
	GORM.DB().SetMaxOpenConns(config.MaxOpenConns)
	GORM.DB().SetConnMaxLifetime(time.Duration(config.MaxLifeTime))
	isInit = true
	logger.Info("数据库初始化成功")
}

// 关闭数据库连接
func Close() error {
	return GORM.Close()
}
