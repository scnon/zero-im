package websocket

type ServerOptions func(opt *serverOption)

type serverOption struct {
	Authentiation
	patten string
}

func newServerOptions(opts ...ServerOptions) serverOption {
	o := serverOption{
		patten:        "/ws",
		Authentiation: new(authentiation),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return o
}

func WithServerAuthentiation(auth Authentiation) ServerOptions {
	return func(opt *serverOption) {
		opt.Authentiation = auth
	}
}

func WithServerPatten(patten string) ServerOptions {
	return func(opt *serverOption) {
		opt.patten = patten
	}
}
