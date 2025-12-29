package server

import "net"

// Option - настройки HTTP-сервера.
type Option func(*Server)

// Port - настройки порта HTTP-сервера.
func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}
