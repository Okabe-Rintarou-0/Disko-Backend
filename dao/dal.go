package dao

import (
	"cloud_disk/dao/model"
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	ctx = context.TODO()
)

func genDao() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./dao/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.User{})

	// Execute the generator
	g.Execute()
}

func init() {
	username := "root"     // 账号
	password := "123"      // 密码
	host := "127.0.0.1"    // 数据库地址，可以是Ip或者域名
	port := 3306           // 数据库端口
	Dbname := "cloud_disk" // 数据库名
	timeout := "10s"       // 连接超时，10秒

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)

	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}

	//genDao()

	err = db.AutoMigrate(&model.User{})
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
