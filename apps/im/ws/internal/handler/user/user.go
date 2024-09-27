package user

import (
	"zero-im/apps/im/ws/internal/svc"
	websocketx "zero-im/apps/im/ws/websocket"

	"github.com/gorilla/websocket"
)

func OnLine(svc *svc.ServiceContext) websocketx.HandlerFunc {
	return func(c *websocketx.Server, conn *websocket.Conn, msg *websocketx.Message) {
		uids := c.GetUids()
		u := c.GetUids(conn)
		err := c.Send(websocketx.NewMessage(u[0], uids), conn)
		c.Info("err ", err)
	}
}
