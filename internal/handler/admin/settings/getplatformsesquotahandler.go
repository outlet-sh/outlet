package settings

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/settings"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetPlatformSESQuotaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := settings.NewGetPlatformSESQuotaLogic(r.Context(), svcCtx)
		resp, err := l.GetPlatformSESQuota()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
