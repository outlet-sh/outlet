package backup

import (
	"net/http"

	"github.com/outlet-sh/outlet/internal/logic/admin/backup"
	"github.com/outlet-sh/outlet/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBackupSettingsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := backup.NewGetBackupSettingsLogic(r.Context(), svcCtx)
		resp, err := l.GetBackupSettings()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
