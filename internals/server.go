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
	msgCh      chan Message
	quitCh     chan struct{}
	kv         *KeyValue
}

func NewServer(config Config) *Server {
	return &Server{
		Config:     config,
		peers:      make(map[*Peer]bool),
		addPeer:    make(chan *Peer),
		deletePeer: make(chan *Peer),
		quitCh:     make(chan struct{}),
		kv:         NewKeyValue(),
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case HelloCommand:
		spec := map[string]string{
			"server": "redis",
		}

		b := WriteMap(spec)
		w := NewWriter(msg.peer.conn)

		if _, err := w.writer.Write(b); err != nil {
			return err
		}

	case ClientCommand:
		w := NewWriter(msg.peer.conn)

		if err := w.Write(Value{typ: "string", rawValue: "OK"}); err != nil {
			return err
		}

	case SetCommand:
		if err := s.kv.Set(v.key, v.value); err != nil {
			return err
		}
		w := NewWriter(msg.peer.conn)

		if err := w.Write(Value{typ: "string", rawValue: "OK"}); err != nil {
			return err
		}

	case GetCommand:
		val, ok := s.kv.Get(v.key)
		if !ok {
			return fmt.Errorf("key not found")
		}
		w := NewWriter(msg.peer.conn)
		if err := w.Write(Value{typ: "string", rawValue: string(val)}); err != nil {
			return err
		}

	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case msg := <-s.msgCh:
			if err := s.handleMessage(msg); err != nil {
				fmt.Println("failed to handle message")
				return
			}
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
