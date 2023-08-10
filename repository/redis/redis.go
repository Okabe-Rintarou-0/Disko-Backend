package redis

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	ctx = context.Background()
	rds *redis.Redis
)

func init() {
	file, err := os.ReadFile("./etc/redis.yaml")
	if err != nil {
		panic(err)
	}

	conf := redis.RedisConf{}
	if err = yaml.Unmarshal(file, &conf); err != nil {
		panic(err)
	}

	fmt.Printf("Read redis config: %+v\n", conf)
	rds = redis.MustNewRedis(conf)
}

func Set(key, value string) error {
	return rds.SetCtx(ctx, key, value)
}

func Exists(key string) (bool, error) {
	return rds.ExistsCtx(ctx, key)
}

func ExpireAt(key string, expireAt int64) error {
	return rds.ExpireatCtx(ctx, key, expireAt)
}

func TTL(key string) (int, error) {
	return rds.TtlCtx(ctx, key)
}
