package handler

import (
	"net/http"

	"disko/internal/logic"
	"disko/internal/svc"
	"disko/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetOneFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetOneFileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetOneFileLogic(r.Context(), svcCtx)
		resp, err := l.GetOneFile(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
