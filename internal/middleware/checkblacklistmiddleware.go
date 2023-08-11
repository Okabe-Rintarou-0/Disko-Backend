package middleware

import (
	"disko/repository/redis"
	"fmt"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"strings"
)

type CheckBlackListMiddleware struct {
}

func NewCheckBlackListMiddleware() *CheckBlackListMiddleware {
	return &CheckBlackListMiddleware{}
}

func (m *CheckBlackListMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if strings.HasPrefix(token, "Bearer ") && len(token) > 7 {
			token = token[7:]
		}
		key := fmt.Sprintf("blacklist:%s", token)
		// check if token is blocked by blacklist in redis
		if exists, err := redis.Exists(key); err == nil {
			if exists {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		} else {
			httpx.ErrorCtx(r.Context(), w, err)
		}
		// Pass through to next handler if needed
		next(w, r)
	}
}
