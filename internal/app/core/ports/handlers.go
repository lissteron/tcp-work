package ports

import "net"

type TCPHandler interface {
	Handle(conn net.Conn)
}
