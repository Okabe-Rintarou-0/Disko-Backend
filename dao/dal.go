package dao

import (
	"context"
	"disko/model"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	ctx = context.TODO()
)

type mysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Timeout  string
}

func init() {
	var (
		cfg mysqlConfig
		err error
	)
	conf.MustLoad("./etc/mysql.yaml", &cfg)
	fmt.Printf("Read mysql config: %+v\n", cfg)
	dsnPattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s"
	dsn := fmt.Sprintf(dsnPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Timeout)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	err = db.AutoMigrate(&model.User{}, &model.File{})
	if err != nil {
		panic(err)
	}

	// 延时关闭数据库连接
	//defer func() {
	//	if sql, err := db.DB(); err == nil {
	//		_ = sql.Close()
	//	}
	//}()
}