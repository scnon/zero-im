package user

import (
	"context"
	"github.com/jinzhu/copier"
	"zero-im/apps/user/api/internal/svc"
	"zero-im/apps/user/api/internal/types"
	"zero-im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRegisterLogic 用户注册
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	registerResp, err := l.svcCtx.User.Register(l.ctx, &user.RegisterReq{
		Phone:    req.Phone,
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Sex:      int32(req.Sex),
	})
	if err != nil {
		return nil, err
	}
	var res = types.RegisterResp{}
	err = copier.Copy(&res, registerResp)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
