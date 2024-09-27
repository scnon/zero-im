package group

import (
	"context"

	"zero-im/apps/social/api/internal/svc"
	"zero-im/apps/social/api/internal/types"
	"zero-im/apps/social/rpc/social_client"
	"zero-im/pkg/ctxdata"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateGroupLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创群
func NewCreateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateGroupLogic {
	return &CreateGroupLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateGroupLogic) CreateGroup(req *types.GroupCreateReq) (resp *types.GroupCreateResp, err error) {
	uid := ctxdata.GetUId(l.ctx)

	// 创建群
	_, err = l.svcCtx.Social.GroupCreate(l.ctx, &social_client.GroupCreateReq{
		Name:       req.Name,
		Icon:       req.Icon,
		CreatorUid: uid,
	})
	if err != nil {
		return nil, err
	}

	return
}
