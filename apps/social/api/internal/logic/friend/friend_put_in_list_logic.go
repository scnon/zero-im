package friend

import (
	"context"

	"zero-im/apps/social/api/internal/svc"
	"zero-im/apps/social/api/internal/types"
	"zero-im/apps/social/rpc/social_client"
	"zero-im/pkg/ctxdata"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

type FriendPutInListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 好友申请列表
func NewFriendPutInListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendPutInListLogic {
	return &FriendPutInListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FriendPutInListLogic) FriendPutInList(req *types.FriendPutInListReq) (resp *types.FriendPutInListResp, err error) {
	list, err := l.svcCtx.Social.FriendPutInList(l.ctx, &social_client.FriendPutInListReq{
		UserId: ctxdata.GetUId(l.ctx),
	})
	if err != nil {
		return nil, err
	}

	var respList []*types.FriendRequests
	copier.Copy(&respList, list.List)

	return &types.FriendPutInListResp{List: respList}, nil
}
