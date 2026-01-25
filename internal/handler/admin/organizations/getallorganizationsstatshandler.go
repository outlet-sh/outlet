package organizations

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/organizations"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetAllOrganizationsStatsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := organizations.NewGetAllOrganizationsStatsLogic(r.Context(), svcCtx)
		resp, err := l.GetAllOrganizationsStats()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
