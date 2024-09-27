package svc

import (
	"zero-im/apps/social/api/internal/config"
	"zero-im/apps/social/rpc/social_client"
	"zero-im/apps/user/rpc/user_client"

	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	user_client.User
	social_client.Social
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User:   user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
		Social: social_client.NewSocial(zrpc.MustNewClient(c.SocialRpc)),
	}
}
