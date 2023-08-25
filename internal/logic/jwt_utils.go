package logic

import (
	"context"
	"github.com/spf13/cast"
)

func GetUserId(ctx context.Context) uint {
	return cast.ToUint(ctx.Value("id"))
}
