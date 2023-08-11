package main

import (
	"disko/model"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type mysqlConfig struct {
	Username string
	Password string
	Host     string
	Port     int
	Database string
	Timeout  string
}

func main() {
	var (
		cfg mysqlConfig
		err error
		db  *gorm.DB
	)
	conf.MustLoad("./etc/mysql.yaml", &cfg)
	fmt.Printf("Read mysql config: %+v\n", cfg)
	dsnPattern := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s"
	dsn := fmt.Sprintf(dsnPattern, cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Timeout)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("连接数据库失败, error=" + err.Error())
	}
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./repository/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.User{}, model.File{})

	// Execute the generator
	g.Execute()
}
