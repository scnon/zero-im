package logic

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-im/apps/user/models"

	"zero-im/apps/user/rpc/internal/svc"
	"zero-im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type FindUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFindUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FindUserLogic {
	return &FindUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *FindUserLogic) FindUser(in *user.FindUserReq) (*user.FindUserResp, error) {
	var (
		list []*models.Users
		err  error
	)

	if in.Phone != nil {
		model, err := l.svcCtx.UsersModel.FindOneByPhone(l.ctx, in.Phone)
		if err == nil {
			list = append(list, model)
		}
	} else if in.Nickname == nil {
		list, err = l.svcCtx.UsersModel.ListByName(l.ctx, in.Nickname)
	} else if len(in.Ids) > 0 {
		list, err = l.svcCtx.UsersModel.ListByIds(l.ctx, in.Ids)
	}

	if err != nil {
		return nil, err
	}

	var entities []*user.UserEntity
	if err = copier.Copy(&entities, list); err != nil {
		return nil, err
	}

	return &user.FindUserResp{
		Users: entities,
	}, nil
}
