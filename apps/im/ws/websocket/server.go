package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

type Server struct {
	sync.RWMutex   // 读写锁 保护 connToUser 和 userToConn
	addr           string
	partten        string
	upgrader       websocket.Upgrader
	routers        map[string]Route
	anthentication Authentiation // 鉴权

	connToUser map[*websocket.Conn]string // connToUser 保存连接对应的用户id
	userToConn map[string]*websocket.Conn // userToConn 保存用户id对应的连接
	logx.Logger
}

func NewServer(addr string, opts ...ServerOptions) *Server {
	opt := newServerOptions(opts...)

	return &Server{
		addr:           addr,
		partten:        opt.patten,
		anthentication: opt.Authentiation,
		upgrader:       websocket.Upgrader{},
		routers:        make(map[string]Route),
		connToUser:     make(map[*websocket.Conn]string),
		userToConn:     make(map[string]*websocket.Conn),
		Logger:         logx.WithContext(context.Background()),
	}
}

func (s *Server) ServeWs(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			s.Errorf("server handler ws recover: %v", err)
		}
	}()

	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		s.Errorf("upgrade err: %v", err)
		return
	}

	// 鉴权
	if !s.anthentication.Auth(w, r) {
		conn.WriteMessage(websocket.TextMessage, []byte("auth failed"))
		s.Close(conn)
		return
	}

	// 存储连接
	s.addConn(conn, r)

	go s.handlerConn(conn)
}

// 处理连接
func (s *Server) handlerConn(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			s.Errorf("websocket read message err: %v", err)
			s.Close(conn)
			return
		}

		var message Message
		if err = json.Unmarshal(msg, &message); err != nil {
			s.Errorf("websocket unmarshal message err: %v", err)
			s.Close(conn)
			return
		}

		// 路由分发
		if route, ok := s.routers[message.Method]; ok {
			route.Handler(s, conn, &message)
		} else {
			conn.WriteMessage(websocket.TextMessage, []byte("route not found"))
			s.Errorf("websocket route not found: %s", message.Method)
		}
	}
}

func (s *Server) addConn(conn *websocket.Conn, req *http.Request) {
	uid := s.anthentication.UserId(req)

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	s.connToUser[conn] = uid
	s.userToConn[uid] = conn
}

func (s *Server) GetConn(uid string) *websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	return s.userToConn[uid]
}

func (s *Server) GetConns(uids ...string) []*websocket.Conn {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	conns := make([]*websocket.Conn, 0, len(uids))
	for _, uid := range uids {
		if conn, ok := s.userToConn[uid]; ok {
			conns = append(conns, conn)
		}
	}

	return conns
}

func (s *Server) GetUids(conns ...*websocket.Conn) []string {
	s.RWMutex.RLock()
	defer s.RWMutex.RUnlock()

	var uids []string
	if len(conns) == 0 {
		// 获取所有在线用户
		uids := make([]string, 0, len(s.connToUser))
		for _, uid := range s.connToUser {
			uids = append(uids, uid)
		}
	} else {
		// 获取部分连接的用户
		uids := make([]string, 0, len(conns))
		for _, conn := range conns {
			uids = append(uids, s.connToUser[conn])
		}
	}

	return uids
}

func (s *Server) Close(conn *websocket.Conn) {
	conn.Close()

	s.RWMutex.Lock()
	defer s.RWMutex.Unlock()

	uid := s.connToUser[conn]
	delete(s.connToUser, conn)
	delete(s.userToConn, uid)
}

func (s *Server) Send(msg interface{}, conns ...*websocket.Conn) error {
	if len(conns) == 0 {
		return nil
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, conn := range conns {
		if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) SendByUserId(msg interface{}, uids ...string) error {
	if len(uids) == 0 {
		return nil
	}

	return s.Send(msg, s.GetConns(uids...)...)
}

func (s *Server) AddRouters(routes []Route) {
	for _, route := range routes {
		s.routers[route.Method] = route
	}
}

func (s *Server) Start() {
	http.HandleFunc(s.partten, s.ServeWs)
	s.Info(http.ListenAndServe(s.addr, nil))
}

func (s *Server) Stop() {
	s.Info("stop")
}
