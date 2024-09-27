package group

import (
	"context"

	"zero-im/apps/social/api/internal/svc"
	"zero-im/apps/social/api/internal/types"
	"zero-im/apps/social/rpc/social_client"
	"zero-im/apps/user/rpc/user_client"

	"github.com/zeromicro/go-zero/core/logx"
)

type GroupUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 成员列表列表
func NewGroupUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GroupUserListLogic {
	return &GroupUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GroupUserListLogic) GroupUserList(req *types.GroupUserListReq) (resp *types.GroupUserListResp, err error) {
	groupUsers, err := l.svcCtx.Social.GroupUsers(l.ctx, &social_client.GroupUsersReq{
		GroupId: req.GroupId,
	})

	if err != nil {
		return nil, err
	}

	// 还需要获取用户的信息
	uids := make([]string, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {
		uids = append(uids, v.UserId)
	}

	// 获取用户信息
	userList, err := l.svcCtx.User.FindUser(l.ctx, &user_client.FindUserReq{Ids: uids})
	if err != nil {
		return nil, err
	}

	userRecords := make(map[string]*user_client.UserEntity, len(userList.Users))
	for i := range userList.Users {
		userRecords[userList.Users[i].Id] = userList.Users[i]
	}

	respList := make([]*types.GroupMembers, 0, len(groupUsers.List))
	for _, v := range groupUsers.List {

		member := &types.GroupMembers{
			Id:        int64(v.Id),
			GroupId:   v.GroupId,
			UserId:    v.UserId,
			RoleLevel: int(v.RoleLevel),
		}
		if u, ok := userRecords[v.UserId]; ok {
			member.Nickname = u.Nickname
			member.UserAvatarUrl = u.Avatar
		}
		respList = append(respList, member)
	}

	return &types.GroupUserListResp{List: respList}, err
}
