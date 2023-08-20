package handler

import (
	"net/http"

	"disko/internal/logic"
	"disko/internal/svc"
	"disko/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func DownloadSharedFileHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadSharedFileRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewDownloadSharedFileLogic(r.Context(), svcCtx)
		resp, err := l.DownloadSharedFile(&req, w)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
