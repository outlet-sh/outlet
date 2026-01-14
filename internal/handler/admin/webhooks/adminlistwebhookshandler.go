package webhooks

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/webhooks"
	"outlet/internal/svc"
)

func AdminListWebhooksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := webhooks.NewAdminListWebhooksLogic(r.Context(), svcCtx)
		resp, err := l.AdminListWebhooks()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
