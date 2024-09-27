package handler

import (
	"zero-im/apps/im/ws/internal/handler/user"
	"zero-im/apps/im/ws/internal/svc"
	"zero-im/apps/im/ws/websocket"
)

func RegisterHandlers(srv *websocket.Server, svc *svc.ServiceContext) {
	srv.AddRouters([]websocket.Route{
		{
			Method:  "user.online",
			Handler: user.OnLine(svc),
		},
	})
}
