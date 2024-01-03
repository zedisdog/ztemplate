package healthz

import (
	"context"

	"simple/internal/svc"
	"simple/internal/api/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HealthzLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHealthzLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HealthzLogic {
	return &HealthzLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HealthzLogic) Healthz(req *types.HealthReq) (resp *types.HealthResp, err error) {
	// todo: add your logic here and delete this line

	return
}
