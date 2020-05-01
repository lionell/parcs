package parcs

import (
	"encoding/binary"
	"encoding/json"
	"log"
)

var ByteOrder = binary.BigEndian

const BytesInInt = 8

func encodeUint64(v uint64) []byte {
	b := make([]byte, BytesInInt)
	ByteOrder.PutUint64(b, v)
	return b
}

func decodeUint64(bytes []byte) uint64 {
	return ByteOrder.Uint64(bytes)
}

func marshal(v interface{}) []byte {
	bytes, err := json.Marshal(v)
	if err != nil {
		log.Fatalf("Error while serializing %v: %v", v, err)
	}
	return bytes
}

func unmarshal(bytes []byte, v interface{}) error {
	return json.Unmarshal(bytes, v)
}
