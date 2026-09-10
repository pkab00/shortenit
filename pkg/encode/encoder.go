package encode

import (
	"encoding/base64"
	"encoding/binary"
)

type Encoder interface {
	Encode(val uint64) string
}

type Base64Encoder struct{}

func NewEncoder() *Base64Encoder {
	return &Base64Encoder{}
}

func (e *Base64Encoder) Encode(val uint64) string {
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, val)
	return base64.RawURLEncoding.EncodeToString(buf)
}
