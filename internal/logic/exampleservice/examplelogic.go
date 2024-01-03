package exampleservicelogic

import (
	"context"

	"simple/internal/svc"
	"simple/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExampleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExampleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExampleLogic {
	return &ExampleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ExampleLogic) Example(in *pb.NoContentReq) (*pb.NoContentResp, error) {
	// todo: add your logic here and delete this line

	return &pb.NoContentResp{}, nil
}
