package encode

import (
	"encoding/base64"
	"encoding/binary"
	"log"
)

type Decoder interface {
	Decode(code string) (*int, error)
}

type Base64Decoder struct{}

func NewDecoder() *Base64Decoder {
	return &Base64Decoder{}
}

func (d *Base64Decoder) Decode(code string) (*int, error) {
	data, err := base64.RawURLEncoding.DecodeString(code)
	if err != nil {
		log.Println("code: ", code)
		return nil, err
	}
	res := int(binary.BigEndian.Uint64(data))
	return &res, nil
}
