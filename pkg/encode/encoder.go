package encode

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/speps/go-hashids"
)

type Encoder interface {
	Encode(val uint64) string
}

type Base64Encoder struct{}

type HashEncoder struct{}

func NewBase64Encoder() *Base64Encoder {
	return &Base64Encoder{}
}

func NewHashEncoder() *HashEncoder {
	return &HashEncoder{}
}

func (e *Base64Encoder) Encode(val uint64) string {
	buf := make([]byte, 8)

	binary.BigEndian.PutUint64(buf, val)
	return base64.RawURLEncoding.EncodeToString(buf)
}

func (e *HashEncoder) Encode(val uint64) string {
	hd := hashids.NewData()
	hd.MinLength = 10
	hd.Salt = "AHHH SHHHH ALMOST FORGOT TO ADD SOME SALT EHEHE"

	h, _ := hashids.NewWithData(hd)
	out, _ := h.Encode([]int{int(val)})
	return out
}
