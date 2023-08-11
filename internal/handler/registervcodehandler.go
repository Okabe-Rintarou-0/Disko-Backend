package handler

import (
	"net/http"

	"disko/internal/logic"
	"disko/internal/svc"
	"disko/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterVcodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterVcodeRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRegisterVcodeLogic(r.Context(), svcCtx)
		resp, err := l.RegisterVcode(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
