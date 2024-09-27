package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"zero-im/apps/social/api/social/internal/logic"
	"zero-im/apps/social/api/social/internal/svc"
	"zero-im/apps/social/api/social/internal/types"
)

func SocialHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSocialLogic(r.Context(), svcCtx)
		resp, err := l.Social(&req)
		response.Response(w, resp, err)
	}
}
