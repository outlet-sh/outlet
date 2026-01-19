package webhooks

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/sdk/webhooks"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListWebhooksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := webhooks.NewListWebhooksLogic(r.Context(), svcCtx)
		resp, err := l.ListWebhooks()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
