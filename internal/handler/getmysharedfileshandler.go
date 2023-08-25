package handler

import (
	"net/http"

	"disko/internal/logic"
	"disko/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetMySharedFilesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetMySharedFilesLogic(r.Context(), svcCtx)
		resp, err := l.GetMySharedFiles()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
