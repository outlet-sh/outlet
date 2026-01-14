package subscribers

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/subscribers"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func ResubscribeSubscriberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubscriberActionRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := subscribers.NewResubscribeSubscriberLogic(r.Context(), svcCtx)
		resp, err := l.ResubscribeSubscriber(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
