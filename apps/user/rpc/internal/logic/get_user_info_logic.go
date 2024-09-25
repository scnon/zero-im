package logic

import (
	"context"
	"github.com/pkg/errors"
	"zero-im/apps/user/models"
	"zero-im/pkg/xerr"

	"zero-im/apps/user/rpc/internal/svc"
	"zero-im/apps/user/rpc/user"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrUserNotFound = errors.New("用户不存在")
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoReq) (*user.GetUserInfoResp, error) {
	model, err := l.svcCtx.UsersModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrUserNotFound)
		}

		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by id err %v, req %v", err, in.Id)
	}

	var entity user.UserEntity
	if err = copier.Copy(&entity, model); err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "copy user entity err %v", err)
	}

	return &user.GetUserInfoResp{
		User: &entity,
	}, nil
}
