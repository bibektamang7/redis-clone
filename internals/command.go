package internals

const (
	CommandHello  = "hello"
	CommandClient = "client"
	CommandSet    = "set"
	CommandGet    = "get"
)

type Command interface{}

type SetCommand struct {
	key, value []byte
}
type GetCommand struct {
	key []byte
}

type ClientCommand struct {
	value string
}
type HelloCommand struct {
	value string
}
