package redis

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var (
	ctx = context.Background()
	rds *redis.Redis
	cfg redis.RedisConf
)

func init() {
	conf.MustLoad("./etc/redis.yaml", &cfg)
	fmt.Printf("Read redis config: %+v\n", cfg)
	rds = redis.MustNewRedis(cfg)
}

func Host() string {
	return cfg.Host
}

func Set(key, value string) error {
	return rds.SetCtx(ctx, key, value)
}

func Get(key string) (string, error) {
	return rds.GetCtx(ctx, key)
}

func Exists(key string) (bool, error) {
	return rds.ExistsCtx(ctx, key)
}

func Del(key string) (int, error) {
	return rds.DelCtx(ctx, key)
}

func ExpireAt(key string, expireAt int64) error {
	return rds.ExpireatCtx(ctx, key, expireAt)
}

func Expire(key string, seconds int) error {
	return rds.ExpireCtx(ctx, key, seconds)
}

func TTL(key string) (int, error) {
	return rds.TtlCtx(ctx, key)
}
