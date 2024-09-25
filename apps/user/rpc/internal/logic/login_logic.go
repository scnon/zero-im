package logic

import (
	"context"
	"github.com/pkg/errors"
	"time"
	"zero-im/apps/user/models"
	"zero-im/pkg/ctxdata"
	"zero-im/pkg/encrypt"
	"zero-im/pkg/xerr"

	"zero-im/apps/user/rpc/internal/svc"
	"zero-im/apps/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	ErrPhoneNotRegistered = xerr.New(xerr.SERVER_COMMON_ERROR, "手机号未注册")
	ErrPasswordError      = xerr.New(xerr.SERVER_COMMON_ERROR, "密码错误")
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *user.LoginReq) (*user.LoginResp, error) {
	// 检查手机号是否注册
	entity, err := l.svcCtx.UsersModel.FindOneByPhone(l.ctx, &in.Phone)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.WithStack(ErrPhoneNotRegistered)
		}
		return nil, errors.Wrapf(xerr.NewDBErr(), "find user by phone err %v, req %v", err, in.Phone)
	}

	// 检查密码
	if !encrypt.ValidatePasswordHash(in.Password, entity.Password) {
		return nil, errors.WithStack(ErrPasswordError)
	}

	// 生成 token
	now := time.Now().Unix()
	token, err := ctxdata.GetJwtToken(l.svcCtx.Config.Jwt.AccessSecret, now, l.svcCtx.Config.Jwt.AccessExpire, entity.Id)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewInternalErr(), "generate jwt token err %v", err)
	}

	return &user.LoginResp{
		Token:   token,
		Expired: now + l.svcCtx.Config.Jwt.AccessExpire,
	}, nil
}
