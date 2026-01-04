package internals

import (
	"bytes"
	"fmt"
	"net"
)

type Peer struct {
	conn  net.Conn
	msgCh chan Message
}

func NewPeer(conn net.Conn) *Peer {
	return &Peer{conn: conn}
}

func (p *Peer) readLoop() error {
	resp := NewResp(p.conn)
	for {
		v, err := resp.Read()
		if err != nil {
			return err
		}
		var cmd Command
		if v.typ == "array" {
			requestCmd := v.array[0].bulk
			switch requestCmd {
			case CommandClient:
				cmd = ClientCommand{
					value: v.array[1].bulk,
				}
			case CommandHello:
				cmd = HelloCommand{
					value: v.array[1].bulk,
				}
			case CommandSet:
				cmd = SetCommand{
					key:   []byte(v.array[1].bulk),
					value: []byte(v.array[2].bulk),
				}
			case CommandGet:
				cmd = GetCommand{
					key: []byte(v.array[1].bulk),
				}
			default:
				fmt.Println("got unhandled command", requestCmd)
			}
			p.msgCh <- Message{
				peer: p,
				cmd:  cmd,
			}
		}
	}
}

func (p *Peer) WriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))

	for k, v := range m {
		buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(k)))
		buf.WriteString(fmt.Sprintf("%s\r\n", v))
	}

	return buf.Bytes()
}
