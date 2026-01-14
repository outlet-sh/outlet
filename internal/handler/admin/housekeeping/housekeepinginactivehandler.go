package housekeeping

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/housekeeping"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func HousekeepingInactiveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HousekeepingInactiveRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := housekeeping.NewHousekeepingInactiveLogic(r.Context(), svcCtx)
		resp, err := l.HousekeepingInactive(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
