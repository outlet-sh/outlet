package transactional

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/transactional"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func GetTransactionalEmailStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetTransactionalEmailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := transactional.NewGetTransactionalEmailStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetTransactionalEmailStats(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
