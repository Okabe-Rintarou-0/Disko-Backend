package dao

import (
	"cloud_disk/model"
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"os"
)

var (
	db  *gorm.DB
	ctx = context.TODO()
)

type mysqlConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Timeout  string `yaml:"timeout"`
}

func genDao() {
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./repository/query",
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
	file, err := os.ReadFile("./etc/mysql.yaml")
	if err != nil {
		panic(err)
	}

	cfg := mysqlConfig{}
	if err = yaml.Unmarshal(file, &cfg); err != nil {
		panic(err)
	}
	fmt.Printf("Read mysql config: %+v\n", cfg)

	dsnPattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s"
	dsn := fmt.Sprintf(dsnPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Timeout)
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
