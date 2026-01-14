package apikeys

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/apikeys"
	"outlet/internal/svc"
)

func ListMCPAPIKeysHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := apikeys.NewListMCPAPIKeysLogic(r.Context(), svcCtx)
		resp, err := l.ListMCPAPIKeys()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
