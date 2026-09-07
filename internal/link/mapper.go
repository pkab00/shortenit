package link

import "github.com/pkab00/shortenit/pkg/encoder"

func (l *Link) toResponse() *LinkResponse {
	encoder := encoder.NewEncoder()
	return &LinkResponse{
		Code:      encoder.Encode(uint64(l.ID)),
		Body:      l.Body,
		CreatedAt: l.CreatedAt,
	}
}
