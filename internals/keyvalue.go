package internals

import "sync"

type KeyValue struct {
	mu   sync.RWMutex
	Data map[string][]byte
}

func NewKeyValue() *KeyValue {
	return &KeyValue{
		Data: make(map[string][]byte),
	}
}

func (kv *KeyValue) Set(key, value []byte) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.Data[string(key)] = value
	return nil
}

func (kv *KeyValue) Get(key []byte) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.Unlock()

	val, ok := kv.Data[string(key)]
	return val, ok

}
