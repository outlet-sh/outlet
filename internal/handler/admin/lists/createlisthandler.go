package lists

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"outlet/internal/logic/admin/lists"
	"outlet/internal/svc"
	"outlet/internal/types"
)

func CreateListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := lists.NewCreateListLogic(r.Context(), svcCtx)
		resp, err := l.CreateList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
