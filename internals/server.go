package internals

import (
	"fmt"
	"net"
)

type Server struct {
	Addr string
}

func NewServer(addr string) *Server {
	return &Server{
		Addr: addr,
	}
}

func (s *Server) ListenAndServer() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	fmt.Println("Server is running at PORT ", s.Addr)
	conn, err := listener.Accept()
	if err != nil {
		return err
	}

	go s.handleConnection(conn)

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("ERROR IN SERVER: ", err)
			return
		}
		fmt.Println("READ FROM CONNECTION: ", string(buf[:n]))
		conn.Write([]byte("+OK\r\n"))
	}
}
