package config

import "github.com/zeromicro/go-zero/rest"

type FileStorageConf struct {
	// in GB
	MaxFileSize  float64
	LocalStorage *LocalStorage
}

type LocalStorage struct {
	Dir string
}

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	FileStorage FileStorageConf
}
