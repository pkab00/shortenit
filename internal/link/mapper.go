package link

func (l *Link) toResponse() *LinkResponse {
	return &LinkResponse{
		Code:      l.Code,
		URL:       l.URL,
		CreatedAt: l.CreatedAt,
	}
}
