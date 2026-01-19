package blocklist

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/blocklist"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ClearSuppressionListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := blocklist.NewClearSuppressionListLogic(r.Context(), svcCtx)
		resp, err := l.ClearSuppressionList()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
