package logic

import (
	"context"

	"zero-im/apps/social/api/social/internal/svc"
	"zero-im/apps/social/api/social/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SocialLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSocialLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SocialLogic {
	return &SocialLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SocialLogic) Social(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
