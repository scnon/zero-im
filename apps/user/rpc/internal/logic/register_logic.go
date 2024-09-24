package logic

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"zero-im/apps/user/models"
	"zero-im/pkg/ctxdata"
	"zero-im/pkg/xerr"

	"zero-im/apps/user/rpc/internal/svc"
	"zero-im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneIsRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号已注册")
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *user.RegisterReq) (*user.RegisterResp, error) {
	userEntity, err := l.svcCtx.UsersModel.FindOneByPhone(l.ctx, &in.Phone)
	if err != nil && !errors.Is(err, models.ErrNotFound) {
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	if userEntity != nil {
		return nil, errors.WithStack(ErrPhoneIsRegistered)
	}

	userEntity = &models.Users{
		Phone: in.Phone,
	}

	if len(in.Password) > 0 {
		// todo)) 加密密码
		userEntity.Password = in.Password
	}

	// 插入用户数据
	_, err = l.svcCtx.UsersModel.Insert(l.ctx, userEntity)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewDBErr(), "insert user err %v, req %v", err, in)
	}

	// 生成 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, userEntity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "generate jwt token err %v", err)
	}

	return &user.RegisterResp{
		Token:   token,
		Expired: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
