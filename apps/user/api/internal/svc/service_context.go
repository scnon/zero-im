package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"zero-im/apps/user/api/internal/config"
	"zero-im/apps/user/rpc/user_client"
)

type ServiceContext struct {
	Config config.Config

	user_client.User
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		User: user_client.NewUser(zrpc.MustNewClient(c.UserRpc)),
	}
}
