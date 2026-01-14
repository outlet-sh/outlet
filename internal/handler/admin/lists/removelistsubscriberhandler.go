package lists

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/lists"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func RemoveListSubscriberHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RemoveSubscriberRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lists.NewRemoveListSubscriberLogic(r.Context(), svcCtx)
		resp, err := l.RemoveListSubscriber(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
