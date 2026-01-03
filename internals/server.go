package internals

import (
	"fmt"
	"net"
	"strings"
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
	for {

		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)

	}

}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}

		data := string(buf[:n])
		commands := strings.Split(data, "*")

		for _, cmd := range commands {
			if cmd == "" {
				continue
			}

			fmt.Printf("Processing Command: *%s\n", strings.ReplaceAll(cmd, "\r\n", " "))

			if strings.Contains(cmd, "hello") {
				conn.Write([]byte("%3\r\n$6\r\nserver\r\n$15\r\nmy-custom-redis\r\n$7\r\nversion\r\n$5\r\n1.0.0\r\n$5\r\nproto\r\n:3\r\n"))
			} else if strings.Contains(cmd, "client") {
				conn.Write([]byte("+OK\r\n"))
			} else if strings.Contains(cmd, "set") {
				fmt.Println("--- RECEIVED THE SET COMMAND! ---")
				fmt.Println("data: ", cmd)
				conn.Write([]byte("+OK\r\n"))
			}
		}
	}
}
