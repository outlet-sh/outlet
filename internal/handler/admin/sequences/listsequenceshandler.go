package sequences

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/sequences"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func ListSequencesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := sequences.NewListSequencesLogic(r.Context(), svcCtx)
		resp, err := l.ListSequences()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
