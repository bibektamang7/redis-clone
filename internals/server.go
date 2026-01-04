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
	Listener     net.Listener
	peers        map[*Peer]bool
	addPeerCh    chan *Peer
	deletePeerCh chan *Peer
	msgCh        chan Message
	quitCh       chan struct{}
	kv           *KeyValue
}

func NewServer(config Config) *Server {
	return &Server{
		Config:       config,
		peers:        make(map[*Peer]bool),
		addPeerCh:    make(chan *Peer),
		msgCh:        make(chan Message),
		deletePeerCh: make(chan *Peer),
		quitCh:       make(chan struct{}),
		kv:           NewKeyValue(),
	}
}

func (s *Server) handleMessage(msg Message) error {
	switch v := msg.cmd.(type) {
	case HelloCommand:
		// spec := map[string]string{
		// 	"server": "redis",
		// }
		// TODO: HANDLE IT PROPERLY
		val := "%3\r\n$6\r\nserver\r\n$15\r\nmy-custom-redis\r\n$7\r\nversion\r\n$5\r\n1.0.0\r\n$5\r\nproto\r\n:3\r\n"

		// b := WriteMap(spec)
		// w := NewWriter(msg.peer.conn)
		if _, err := msg.peer.conn.Write([]byte(val)); err != nil {
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
			}
		case p := <-s.addPeerCh:
			slog.Info("peer connected:", "remoteAddr", p.conn.RemoteAddr())
			s.peers[p] = true
		case p := <-s.deletePeerCh:
			slog.Info("peer disconnected: ", "remoteAddr", p.conn.RemoteAddr())
			delete(s.peers, p)
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
			continue
		}
		go s.handleConnection(conn)

	}

}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	peer := NewPeer(conn, s.msgCh, s.deletePeerCh)

	s.addPeerCh <- peer
	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", "err", err, "remote Addr", conn.RemoteAddr().String())
	}
}
