package handler

import (
	"net/http"

	"cloud_disk/internal/logic"
	"cloud_disk/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func MeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewMeLogic(r.Context(), svcCtx)
		resp, err := l.Me()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
