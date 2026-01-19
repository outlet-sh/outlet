package apikeys

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/apikeys"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
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
