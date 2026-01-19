package public

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/public"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Check if initial setup is required
func GetPublicSetupStatusHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := public.NewGetPublicSetupStatusLogic(r.Context(), svcCtx)
		resp, err := l.GetPublicSetupStatus()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
