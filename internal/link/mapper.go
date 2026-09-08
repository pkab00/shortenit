package link

import "github.com/pkab00/shortenit/pkg/encode"

func (l *Link) toResponse() *LinkResponse {
	encoder := encode.NewEncoder()
	return &LinkResponse{
		Code:      encoder.Encode(uint64(l.ID)),
		Body:      l.Body,
		CreatedAt: l.CreatedAt,
	}
}
