package healthz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"simple/internal/api/logic/healthz"
	"simple/internal/svc"
	"simple/internal/api/types"
)

func HealthzHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HealthReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := healthz.NewHealthzLogic(r.Context(), svcCtx)
		resp, err := l.Healthz(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
