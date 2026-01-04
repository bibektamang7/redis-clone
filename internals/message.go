package internals

type Message struct {
	cmd  Command
	peer *Peer
}
