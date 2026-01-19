package system

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/system"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func CheckForUpdatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := system.NewCheckForUpdatesLogic(r.Context(), svcCtx)
		resp, err := l.CheckForUpdates()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
