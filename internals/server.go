package internals

import (
	"fmt"
	"log/slog"
	"net"
)

type Config struct {
	Addr string
}

type Server struct {
	Config
	Listener   net.Listener
	peers      map[*Peer]bool
	addPeer    chan *Peer
	deletePeer chan *Peer
	quitCh     chan struct{}
}

func NewServer(config Config) *Server {
	return &Server{
		Config:     config,
		peers:      make(map[*Peer]bool),
		addPeer:    make(chan *Peer),
		deletePeer: make(chan *Peer),
		quitCh:     make(chan struct{}),
	}
}

func (s *Server) loop() {
	for {
		select {
		case p := <-s.addPeer:
			s.peers[p] = true
		case <-s.quitCh:
			return
		}
	}
}

func (s *Server) ListenAndServer() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.Listener = listener

	go s.loop()

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
	peer := NewPeer(conn)

	s.addPeer <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remote Addr", conn.RemoteAddr().String())
	}
}
