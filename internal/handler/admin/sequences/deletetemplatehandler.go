package sequences

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/sequences"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func DeleteTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteTemplateRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := sequences.NewDeleteTemplateLogic(r.Context(), svcCtx)
		resp, err := l.DeleteTemplate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
