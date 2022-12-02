package server

import "net"

const (
	CONN_HOST = "localhost"
	CONN_PORT = "4545"
	CONN_TYPE = "tcp"
)

type Server struct {
	Listener net.Listener
}

func (s *Server) NewConnection() (err error) {
	s.Listener, err = net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		return
	}
	return
}
