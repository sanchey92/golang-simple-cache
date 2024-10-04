package protocol

import (
	"bytes"
	"fmt"
	"github.com/tidwall/resp"
)

const (
	CommandSET = "set"
	CommandGET = "get"
)

type Command interface{}

type SetCommand struct {
	Key, Value []byte
}

type GetCommand struct {
	Key []byte
}

// respWriteMap serializes a map of string key-value pairs to the RESP (REdis Serialization Protocol) format.
//
// Params:
// - m (map[string]string): The map containing key-value pairs to serialize.
//
// Returns:
// - []byte: A byte slice containing the serialized RESP format.
func respWriteMap(m map[string]string) []byte {
	buf := &bytes.Buffer{}
	buf.WriteString("%" + fmt.Sprintf("%d\r\n", len(m)))
	rw := resp.NewWriter(buf)

	for key, val := range m {
		rw.WriteString(key)
		rw.WriteString(":" + val)
	}

	return buf.Bytes()
}
